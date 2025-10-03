const { GoogleGenerativeAI } = require('@google/generative-ai');

class GeminiService {
  constructor() {
    const apiKey = process.env.GEMINI_API_KEY;
    if (!apiKey) {
      throw new Error('GEMINI_API_KEY is not set in environment variables');
    }
    this.genAI = new GoogleGenerativeAI(apiKey);
    this.model = this.genAI.getGenerativeModel({ model: 'gemini-pro' });
  }

  /**
   * Generate a coding exercise using Gemini AI
   * @param {string} difficulty - The difficulty level (beginner, intermediate, advanced)
   * @param {string} topic - The programming topic (e.g., arrays, loops, functions)
   * @param {string} language - The programming language (e.g., javascript, python, java)
   * @returns {Promise<object>} The generated exercise
   */
  async generateExercise(difficulty = 'beginner', topic = 'general', language = 'javascript') {
    try {
      const prompt = `Generate a coding exercise with the following requirements:
- Difficulty: ${difficulty}
- Topic: ${topic}
- Programming Language: ${language}

Please provide the response in the following JSON format:
{
  "title": "Exercise title",
  "description": "Detailed description of the problem",
  "difficulty": "${difficulty}",
  "language": "${language}",
  "starterCode": "Initial code template for the student",
  "hints": ["hint1", "hint2", "hint3"],
  "expectedOutput": "Description of expected output or behavior",
  "testCases": [
    {"input": "test input 1", "expectedOutput": "expected output 1"},
    {"input": "test input 2", "expectedOutput": "expected output 2"}
  ]
}

Make the exercise educational and engaging. Include clear instructions and appropriate test cases.`;

      const result = await this.model.generateContent(prompt);
      const response = await result.response;
      const text = response.text();
      
      // Try to parse JSON from the response
      const jsonMatch = text.match(/\{[\s\S]*\}/);
      if (jsonMatch) {
        const exercise = JSON.parse(jsonMatch[0]);
        return exercise;
      }
      
      // Fallback if JSON parsing fails
      return {
        title: 'Generated Exercise',
        description: text,
        difficulty,
        language,
        starterCode: `// Write your ${language} code here\n`,
        hints: [],
        expectedOutput: 'See description for details',
        testCases: []
      };
    } catch (error) {
      console.error('Error generating exercise:', error);
      throw new Error('Failed to generate exercise');
    }
  }

  /**
   * Check submitted code using Gemini AI
   * @param {string} exercise - The exercise description
   * @param {string} submittedCode - The student's submitted code
   * @param {string} language - The programming language
   * @returns {Promise<object>} Feedback on the submission
   */
  async checkSubmission(exercise, submittedCode, language = 'javascript') {
    try {
      const prompt = `You are a code reviewer and teacher. Review the following code submission:

EXERCISE:
${exercise}

PROGRAMMING LANGUAGE: ${language}

SUBMITTED CODE:
\`\`\`${language}
${submittedCode}
\`\`\`

Please provide feedback in the following JSON format:
{
  "passed": true/false,
  "score": 0-100,
  "feedback": "Overall feedback on the submission",
  "strengths": ["strength1", "strength2"],
  "improvements": ["improvement1", "improvement2"],
  "correctness": "Assessment of correctness",
  "codeQuality": "Assessment of code quality and style",
  "suggestions": "Specific suggestions for improvement"
}

Be constructive, educational, and encouraging in your feedback.`;

      const result = await this.model.generateContent(prompt);
      const response = await result.response;
      const text = response.text();
      
      // Try to parse JSON from the response
      const jsonMatch = text.match(/\{[\s\S]*\}/);
      if (jsonMatch) {
        const feedback = JSON.parse(jsonMatch[0]);
        return feedback;
      }
      
      // Fallback if JSON parsing fails
      return {
        passed: false,
        score: 50,
        feedback: text,
        strengths: [],
        improvements: [],
        correctness: 'Unable to determine',
        codeQuality: 'Unable to determine',
        suggestions: 'Please review the feedback above'
      };
    } catch (error) {
      console.error('Error checking submission:', error);
      throw new Error('Failed to check submission');
    }
  }
}

module.exports = new GeminiService();

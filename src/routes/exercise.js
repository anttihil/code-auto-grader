const express = require('express');
const router = express.Router();
const geminiService = require('../services/gemini');

/**
 * POST /api/exercise/generate
 * Generate a new coding exercise
 */
router.post('/generate', async (req, res) => {
  try {
    const { difficulty, topic, language } = req.body;
    
    const exercise = await geminiService.generateExercise(
      difficulty || 'beginner',
      topic || 'general',
      language || 'javascript'
    );
    
    res.json({
      success: true,
      exercise
    });
  } catch (error) {
    console.error('Error generating exercise:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to generate exercise',
      message: error.message
    });
  }
});

/**
 * GET /api/exercise/sample
 * Get a sample exercise for testing
 */
router.get('/sample', (req, res) => {
  const sampleExercise = {
    title: 'Sum Two Numbers',
    description: 'Write a function that takes two numbers as parameters and returns their sum.',
    difficulty: 'beginner',
    language: 'javascript',
    starterCode: '// Write a function called sum that takes two parameters: a and b\n// Return the sum of a and b\n\nfunction sum(a, b) {\n  // Your code here\n}\n\n// Test your function\nconsole.log(sum(2, 3)); // Should output: 5\nconsole.log(sum(10, 20)); // Should output: 30',
    hints: [
      'Use the + operator to add two numbers',
      'Remember to return the result',
      'The function should work with any two numbers'
    ],
    expectedOutput: 'The function should return the sum of two numbers',
    testCases: [
      { input: 'sum(2, 3)', expectedOutput: '5' },
      { input: 'sum(10, 20)', expectedOutput: '30' },
      { input: 'sum(-5, 5)', expectedOutput: '0' }
    ]
  };
  
  res.json({
    success: true,
    exercise: sampleExercise
  });
});

module.exports = router;

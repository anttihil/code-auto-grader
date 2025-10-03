const express = require('express');
const router = express.Router();
const geminiService = require('../services/gemini');
const codeExecutionService = require('../services/codeExecution');

/**
 * POST /api/submission/submit
 * Submit code for evaluation
 */
router.post('/submit', async (req, res) => {
  try {
    const { code, exercise, language } = req.body;
    
    if (!code || !exercise) {
      return res.status(400).json({
        success: false,
        error: 'Code and exercise are required'
      });
    }
    
    // Execute the code to check if it runs
    const executionResult = await codeExecutionService.executeCode(
      code,
      language || 'javascript'
    );
    
    // Get AI feedback on the submission
    const feedback = await geminiService.checkSubmission(
      exercise,
      code,
      language || 'javascript'
    );
    
    res.json({
      success: true,
      executionResult,
      feedback
    });
  } catch (error) {
    console.error('Error processing submission:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to process submission',
      message: error.message
    });
  }
});

/**
 * POST /api/submission/run
 * Run code without evaluation
 */
router.post('/run', async (req, res) => {
  try {
    const { code, language } = req.body;
    
    if (!code) {
      return res.status(400).json({
        success: false,
        error: 'Code is required'
      });
    }
    
    const executionResult = await codeExecutionService.executeCode(
      code,
      language || 'javascript'
    );
    
    res.json({
      success: true,
      executionResult
    });
  } catch (error) {
    console.error('Error running code:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to run code',
      message: error.message
    });
  }
});

module.exports = router;

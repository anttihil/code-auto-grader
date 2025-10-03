const { exec } = require('child_process');
const { promisify } = require('util');
const fs = require('fs').promises;
const path = require('path');
const crypto = require('crypto');

const execAsync = promisify(exec);

class CodeExecutionService {
  constructor() {
    this.tempDir = '/tmp/code-executions';
    this.timeout = 5000; // 5 seconds timeout
  }

  /**
   * Initialize the temporary directory for code execution
   */
  async init() {
    try {
      await fs.mkdir(this.tempDir, { recursive: true });
    } catch (error) {
      console.error('Error creating temp directory:', error);
    }
  }

  /**
   * Execute JavaScript code securely in an isolated environment
   * @param {string} code - The code to execute
   * @param {string} input - Optional input for the code
   * @returns {Promise<object>} Execution result
   */
  async executeJavaScript(code, input = '') {
    const executionId = crypto.randomBytes(16).toString('hex');
    const filePath = path.join(this.tempDir, `${executionId}.js`);
    
    try {
      // Write code to a temporary file
      await fs.writeFile(filePath, code);
      
      // Execute with Node.js in a restricted way
      // Using timeout and limiting resources
      const nodePath = process.execPath || 'node';
      const command = `timeout ${this.timeout / 1000}s ${nodePath} --max-old-space-size=128 "${filePath}"`;
      
      const { stdout, stderr } = await execAsync(command, {
        timeout: this.timeout,
        maxBuffer: 1024 * 1024, // 1MB max output
        env: { ...process.env, NODE_ENV: 'sandbox' }
      });
      
      return {
        success: true,
        output: stdout || '',
        error: stderr || null,
        executionTime: 'completed'
      };
    } catch (error) {
      return {
        success: false,
        output: '',
        error: error.message || 'Execution failed',
        executionTime: error.killed ? 'timeout' : 'error'
      };
    } finally {
      // Clean up the temporary file
      try {
        await fs.unlink(filePath);
      } catch (err) {
        // Ignore cleanup errors
      }
    }
  }

  /**
   * Execute Python code securely
   * @param {string} code - The code to execute
   * @param {string} input - Optional input for the code
   * @returns {Promise<object>} Execution result
   */
  async executePython(code, input = '') {
    const executionId = crypto.randomBytes(16).toString('hex');
    const filePath = path.join(this.tempDir, `${executionId}.py`);
    
    try {
      await fs.writeFile(filePath, code);
      
      const pythonPath = 'python3';
      const command = `timeout ${this.timeout / 1000}s ${pythonPath} -u "${filePath}"`;
      
      const { stdout, stderr } = await execAsync(command, {
        timeout: this.timeout,
        maxBuffer: 1024 * 1024,
        env: { ...process.env }
      });
      
      return {
        success: true,
        output: stdout || '',
        error: stderr || null,
        executionTime: 'completed'
      };
    } catch (error) {
      return {
        success: false,
        output: '',
        error: error.message || 'Execution failed',
        executionTime: error.killed ? 'timeout' : 'error'
      };
    } finally {
      try {
        await fs.unlink(filePath);
      } catch (err) {
        // Ignore cleanup errors
      }
    }
  }

  /**
   * Execute code based on language
   * @param {string} code - The code to execute
   * @param {string} language - The programming language
   * @param {string} input - Optional input for the code
   * @returns {Promise<object>} Execution result
   */
  async executeCode(code, language = 'javascript', input = '') {
    await this.init();
    
    // Sanitize and validate code
    if (!code || code.trim().length === 0) {
      return {
        success: false,
        output: '',
        error: 'No code provided',
        executionTime: 'error'
      };
    }
    
    // Check for dangerous patterns (basic security)
    const dangerousPatterns = [
      /require\s*\(\s*['"]child_process['"]\s*\)/,
      /require\s*\(\s*['"]fs['"]\s*\)/,
      /require\s*\(\s*['"]net['"]\s*\)/,
      /import\s+os/,
      /import\s+subprocess/,
      /import\s+socket/,
      /eval\s*\(/,
      /exec\s*\(/,
      /system\s*\(/,
    ];
    
    for (const pattern of dangerousPatterns) {
      if (pattern.test(code)) {
        return {
          success: false,
          output: '',
          error: 'Code contains potentially dangerous operations',
          executionTime: 'blocked'
        };
      }
    }
    
    switch (language.toLowerCase()) {
      case 'javascript':
      case 'js':
        return await this.executeJavaScript(code, input);
      case 'python':
      case 'py':
        return await this.executePython(code, input);
      default:
        return {
          success: false,
          output: '',
          error: `Language ${language} is not supported yet`,
          executionTime: 'error'
        };
    }
  }
}

module.exports = new CodeExecutionService();

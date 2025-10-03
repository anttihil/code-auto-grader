// Initialize CodeMirror editor
let editor;
let currentExercise = null;

document.addEventListener('DOMContentLoaded', () => {
    // Initialize CodeMirror
    editor = CodeMirror.fromTextArea(document.getElementById('codeEditor'), {
        mode: 'javascript',
        theme: 'dracula',
        lineNumbers: true,
        autoCloseBrackets: true,
        matchBrackets: true,
        indentUnit: 2,
        tabSize: 2,
        lineWrapping: true
    });

    // Set initial content
    editor.setValue('// Write your code here\n');

    // Event listeners
    document.getElementById('generateBtn').addEventListener('click', generateExercise);
    document.getElementById('sampleBtn').addEventListener('click', loadSampleExercise);
    document.getElementById('runBtn').addEventListener('click', runCode);
    document.getElementById('submitBtn').addEventListener('click', submitCode);
    document.getElementById('clearBtn').addEventListener('click', clearEditor);
    document.getElementById('language').addEventListener('change', onLanguageChange);
});

function showLoading(show) {
    document.getElementById('loadingOverlay').style.display = show ? 'flex' : 'none';
}

function onLanguageChange() {
    const language = document.getElementById('language').value;
    const mode = language === 'python' ? 'python' : 'javascript';
    editor.setOption('mode', mode);
}

async function generateExercise() {
    const difficulty = document.getElementById('difficulty').value;
    const topic = document.getElementById('topic').value;
    const language = document.getElementById('language').value;

    showLoading(true);

    try {
        const response = await fetch('/api/exercise/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ difficulty, topic, language })
        });

        const data = await response.json();

        if (data.success) {
            currentExercise = data.exercise;
            displayExercise(data.exercise);
            if (data.exercise.starterCode) {
                editor.setValue(data.exercise.starterCode);
            }
        } else {
            showError('Failed to generate exercise. Please try again.');
        }
    } catch (error) {
        console.error('Error generating exercise:', error);
        showError('Error connecting to server. Make sure the server is running and try again.');
    } finally {
        showLoading(false);
    }
}

async function loadSampleExercise() {
    showLoading(true);

    try {
        const response = await fetch('/api/exercise/sample');
        const data = await response.json();

        if (data.success) {
            currentExercise = data.exercise;
            displayExercise(data.exercise);
            if (data.exercise.starterCode) {
                editor.setValue(data.exercise.starterCode);
            }
        } else {
            showError('Failed to load sample exercise.');
        }
    } catch (error) {
        console.error('Error loading sample:', error);
        showError('Error connecting to server. Make sure the server is running and try again.');
    } finally {
        showLoading(false);
    }
}

function displayExercise(exercise) {
    const content = document.getElementById('exerciseContent');
    
    let hintsHtml = '';
    if (exercise.hints && exercise.hints.length > 0) {
        hintsHtml = `
            <div class="hints">
                <details>
                    <summary>💡 Show Hints</summary>
                    <ul>
                        ${exercise.hints.map(hint => `<li>${hint}</li>`).join('')}
                    </ul>
                </details>
            </div>
        `;
    }

    let testCasesHtml = '';
    if (exercise.testCases && exercise.testCases.length > 0) {
        testCasesHtml = `
            <div class="test-cases">
                <h4>Test Cases:</h4>
                <ul>
                    ${exercise.testCases.map(tc => `<li><code>${tc.input}</code> → <code>${tc.expectedOutput}</code></li>`).join('')}
                </ul>
            </div>
        `;
    }

    content.innerHTML = `
        <div class="difficulty ${exercise.difficulty}">${exercise.difficulty.toUpperCase()}</div>
        <h3>${exercise.title}</h3>
        <p><strong>Language:</strong> ${exercise.language}</p>
        <p>${exercise.description}</p>
        ${testCasesHtml}
        ${hintsHtml}
    `;
}

async function runCode() {
    const code = editor.getValue();
    const language = document.getElementById('language').value;
    const output = document.getElementById('output');

    if (!code.trim()) {
        output.className = 'output error';
        output.textContent = 'Please write some code first.';
        return;
    }

    showLoading(true);
    output.className = 'output';
    output.textContent = 'Running code...';

    try {
        const response = await fetch('/api/submission/run', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ code, language })
        });

        const data = await response.json();

        if (data.success && data.executionResult) {
            const result = data.executionResult;
            if (result.success) {
                output.className = 'output success';
                output.textContent = result.output || 'Code executed successfully (no output)';
            } else {
                output.className = 'output error';
                output.textContent = `Error: ${result.error}`;
            }
        } else {
            output.className = 'output error';
            output.textContent = 'Failed to run code. Please try again.';
        }
    } catch (error) {
        console.error('Error running code:', error);
        output.className = 'output error';
        output.textContent = 'Error connecting to server. Make sure the server is running and try again.';
    } finally {
        showLoading(false);
    }
}

async function submitCode() {
    if (!currentExercise) {
        alert('Please load an exercise first.');
        return;
    }

    const code = editor.getValue();
    const language = document.getElementById('language').value;

    if (!code.trim()) {
        alert('Please write some code first.');
        return;
    }

    showLoading(true);

    try {
        const response = await fetch('/api/submission/submit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                code,
                exercise: currentExercise.description,
                language
            })
        });

        const data = await response.json();

        if (data.success) {
            displayFeedback(data.feedback, data.executionResult);
        } else {
            showError('Failed to submit code. Please try again.');
        }
    } catch (error) {
        console.error('Error submitting code:', error);
        showError('Error connecting to server. Make sure the server is running and try again.');
    } finally {
        showLoading(false);
    }
}

function displayFeedback(feedback, executionResult) {
    const feedbackDiv = document.getElementById('feedback');
    
    const scoreClass = feedback.score >= 80 ? 'high' : feedback.score >= 60 ? 'medium' : 'low';
    
    let strengthsHtml = '';
    if (feedback.strengths && feedback.strengths.length > 0) {
        strengthsHtml = `
            <h4>✓ Strengths:</h4>
            <ul>
                ${feedback.strengths.map(s => `<li>${s}</li>`).join('')}
            </ul>
        `;
    }

    let improvementsHtml = '';
    if (feedback.improvements && feedback.improvements.length > 0) {
        improvementsHtml = `
            <h4>→ Areas for Improvement:</h4>
            <ul>
                ${feedback.improvements.map(i => `<li>${i}</li>`).join('')}
            </ul>
        `;
    }

    let executionHtml = '';
    if (executionResult) {
        const execStatus = executionResult.success ? '✓ Code executed successfully' : '✗ Execution failed';
        const execClass = executionResult.success ? 'success' : 'error';
        executionHtml = `
            <h4>Execution Result:</h4>
            <p class="${execClass}">${execStatus}</p>
            ${executionResult.output ? `<pre>${executionResult.output}</pre>` : ''}
            ${executionResult.error ? `<p class="error">Error: ${executionResult.error}</p>` : ''}
        `;
    }

    feedbackDiv.innerHTML = `
        <div class="score ${scoreClass}">Score: ${feedback.score}/100</div>
        <h4>Overall Feedback:</h4>
        <p>${feedback.feedback}</p>
        ${executionHtml}
        ${strengthsHtml}
        ${improvementsHtml}
        <h4>Code Quality:</h4>
        <p>${feedback.codeQuality}</p>
        <h4>Correctness:</h4>
        <p>${feedback.correctness}</p>
        ${feedback.suggestions ? `
            <h4>Suggestions:</h4>
            <p>${feedback.suggestions}</p>
        ` : ''}
    `;
}

function clearEditor() {
    if (confirm('Are you sure you want to clear the editor?')) {
        editor.setValue('// Write your code here\n');
        document.getElementById('output').textContent = '';
    }
}

function showError(message) {
    alert(message);
}

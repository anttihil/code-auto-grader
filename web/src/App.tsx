import { createSignal, Show, For } from "solid-js";

interface Exercise {
  title: string;
  description: string;
  difficulty: string;
  language: string;
  starterCode?: string;
  testCases?: TestCase[];
  hints?: string[];
}

interface TestCase {
  input: string;
  expectedOutput: string;
}

interface Feedback {
  score: number;
  feedback: string;
  strengths?: string[];
  improvements?: string[];
  codeQuality: string;
  correctness: string;
  suggestions?: string;
}

function App() {
  const [difficulty, setDifficulty] = createSignal("beginner");
  const [topic, setTopic] = createSignal("general");
  const [language, setLanguage] = createSignal("javascript");
  const [code, setCode] = createSignal("// Write your code here\n");
  const [output, setOutput] = createSignal("");
  const [exercise, setExercise] = createSignal<Exercise | null>(null);
  const [feedback, setFeedback] = createSignal<Feedback | null>(null);
  const [loading, setLoading] = createSignal(false);
  const [outputClass, setOutputClass] = createSignal("output");

  // Handle language change
  const handleLanguageChange = (lang: string) => {
    setLanguage(lang);
    if (lang === "python") {
      if (code() === "// Write your code here\n") {
        setCode("# Write your code here\n");
      }
    } else {
      if (code() === "# Write your code here\n") {
        setCode("// Write your code here\n");
      }
    }
  };

  // Handle tab key in textarea
  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const target = e.target as HTMLTextAreaElement;
      const start = target.selectionStart;
      const end = target.selectionEnd;
      const value = code();
      setCode(value.substring(0, start) + "  " + value.substring(end));
      setTimeout(() => {
        target.selectionStart = target.selectionEnd = start + 2;
      }, 0);
    }
  };

  // Load sample exercise
  const loadSampleExercise = async () => {
    setLoading(true);
    try {
      const response = await fetch("/api/exercise/sample");
      const data = await response.json();
      if (data.success) {
        setExercise(data.exercise);
        if (data.exercise?.starterCode) {
          setCode(data.exercise.starterCode);
        }
      } else {
        alert("Failed to load sample exercise.");
      }
    } catch (error) {
      console.error("Error loading sample:", error);
      alert(
        "Error connecting to server. Make sure the server is running and try again."
      );
    } finally {
      setLoading(false);
    }
  };

  // Generate new exercise
  const generateExercise = async () => {
    setLoading(true);
    try {
      const response = await fetch("/api/exercise/generate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          difficulty: difficulty(),
          topic: topic(),
          language: language(),
        }),
      });
      const data = await response.json();
      if (data.success) {
        setExercise(data.exercise);
        if (data.exercise?.starterCode) {
          setCode(data.exercise.starterCode);
        }
      } else {
        alert("Failed to generate exercise. Please try again.");
      }
    } catch (error) {
      console.error("Error generating exercise:", error);
      alert(
        "Error connecting to server. Make sure the server is running and try again."
      );
    } finally {
      setLoading(false);
    }
  };

  // Run code
  const runCode = async () => {
    if (!code().trim()) {
      setOutputClass("output error");
      setOutput("Please write some code first.");
      return;
    }

    setLoading(true);
    setOutputClass("output");
    setOutput("Running code...");

    try {
      const response = await fetch("/api/submission/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code: code(), language: language() }),
      });
      const data = await response.json();

      if (data.success && data.executionResult) {
        const result = data.executionResult;
        if (result.success) {
          setOutputClass("output success");
          setOutput(result.output || "Code executed successfully (no output)");
        } else {
          setOutputClass("output error");
          setOutput(`Error: ${result.error}`);
        }
      } else {
        setOutputClass("output error");
        setOutput("Failed to run code. Please try again.");
      }
    } catch (error) {
      console.error("Error running code:", error);
      setOutputClass("output error");
      setOutput(
        "Error connecting to server. Make sure the server is running and try again."
      );
    } finally {
      setLoading(false);
    }
  };

  // Submit code for grading
  const submitCode = async () => {
    const loadedExercise = exercise();
    const currentCode = code();
    if (!loadedExercise) {
      alert("Please load an exercise first.");
      return;
    }

    if (!currentCode.trim()) {
      alert("Please write some code first.");
      return;
    }

    setLoading(true);

    try {
      const response = await fetch("/api/submission/submit", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          code: currentCode,
          exercise: loadedExercise.description,
          language: language(),
        }),
      });
      const data = await response.json();

      if (data.success) {
        setFeedback(data.feedback);
      } else {
        alert("Failed to submit code. Please try again.");
      }
    } catch (error) {
      console.error("Error submitting code:", error);
      alert(
        "Error connecting to server. Make sure the server is running and try again."
      );
    } finally {
      setLoading(false);
    }
  };

  // Clear editor
  const clearEditor = () => {
    if (confirm("Are you sure you want to clear the editor?")) {
      setCode("// Write your code here\n");
      setOutput("");
    }
  };

  return (
    <div class="container">
      <header>
        <h1>🎓 Code Auto-Grader</h1>
        <p>LLM-Powered Code Learning Platform</p>
      </header>

      <div class="content">
        {/* Exercise Panel */}
        <div class="exercise-panel">
          <div class="panel-header">
            <h2>Exercise</h2>
            <div class="exercise-controls">
              <select
                value={difficulty()}
                onInput={(e) => setDifficulty(e.target.value)}
              >
                <option value="beginner">Beginner</option>
                <option value="intermediate">Intermediate</option>
                <option value="advanced">Advanced</option>
              </select>
              <select value={topic()} onInput={(e) => setTopic(e.target.value)}>
                <option value="general">General</option>
                <option value="arrays">Arrays</option>
                <option value="loops">Loops</option>
                <option value="functions">Functions</option>
                <option value="objects">Objects</option>
                <option value="algorithms">Algorithms</option>
              </select>
              <select
                value={language()}
                onInput={(e) => handleLanguageChange(e.target.value)}
              >
                <option value="javascript">JavaScript</option>
                <option value="python">Python</option>
              </select>
              <button onClick={generateExercise} class="btn btn-primary">
                Generate New Exercise
              </button>
              <button onClick={loadSampleExercise} class="btn btn-secondary">
                Load Sample
              </button>
            </div>
          </div>

          <div class="exercise-content">
            <Show
              when={exercise()}
              fallback={
                <div class="loading-message">
                  <p>
                    👋 Welcome! Click "Load Sample" for a demo or "Generate New
                    Exercise" to create a custom problem.
                  </p>
                </div>
              }
            >
              <div class={`difficulty ${exercise()!.difficulty}`}>
                {exercise()!.difficulty.toUpperCase()}
              </div>
              <h3>{exercise()!.title}</h3>
              <p>
                <strong>Language:</strong> {exercise()!.language}
              </p>
              <p>{exercise()!.description}</p>
              <Show
                when={
                  exercise()!.testCases && exercise()!.testCases!.length > 0
                }
              >
                <div class="test-cases">
                  <h4>Test Cases:</h4>
                  <ul>
                    <For each={exercise()!.testCases}>
                      {(tc) => (
                        <li>
                          <code>{tc.input}</code> →{" "}
                          <code>{tc.expectedOutput}</code>
                        </li>
                      )}
                    </For>
                  </ul>
                </div>
              </Show>
              <Show when={exercise()!.hints && exercise()!.hints!.length > 0}>
                <div class="hints">
                  <details>
                    <summary>💡 Show Hints</summary>
                    <ul>
                      <For each={exercise()!.hints}>
                        {(hint) => <li>{hint}</li>}
                      </For>
                    </ul>
                  </details>
                </div>
              </Show>
            </Show>
          </div>
        </div>

        {/* Code Editor Panel */}
        <div class="editor-panel">
          <div class="panel-header">
            <h2>Code Editor</h2>
            <div class="editor-controls">
              <button onClick={runCode} class="btn btn-success">
                ▶ Run Code
              </button>
              <button onClick={submitCode} class="btn btn-primary">
                ✓ Submit for Grading
              </button>
              <button onClick={clearEditor} class="btn btn-secondary">
                Clear
              </button>
            </div>
          </div>

          <textarea
            class="code-textarea"
            spellcheck={false}
            value={code()}
            onInput={(e) => setCode(e.target.value)}
            onKeyDown={handleKeyDown}
          />

          <div class="output-section">
            <h3>Output</h3>
            <div class={outputClass()}>{output()}</div>
          </div>
        </div>

        {/* Feedback Panel */}
        <div class="feedback-panel">
          <div class="panel-header">
            <h2>AI Feedback</h2>
          </div>
          <div class="feedback-content">
            <Show
              when={feedback()}
              fallback={
                <p class="info-message">
                  Submit your code to receive AI-powered feedback and grading.
                </p>
              }
            >
              <div
                class={`score ${
                  feedback()!.score >= 80
                    ? "high"
                    : feedback()!.score >= 60
                    ? "medium"
                    : "low"
                }`}
              >
                Score: {feedback()!.score}/100
              </div>
              <h4>Overall Feedback:</h4>
              <p>{feedback()!.feedback}</p>
              <Show
                when={
                  feedback()!.strengths && feedback()!.strengths!.length > 0
                }
              >
                <h4>✓ Strengths:</h4>
                <ul>
                  <For each={feedback()!.strengths}>{(s) => <li>{s}</li>}</For>
                </ul>
              </Show>
              <Show
                when={
                  feedback()!.improvements &&
                  feedback()!.improvements!.length > 0
                }
              >
                <h4>→ Areas for Improvement:</h4>
                <ul>
                  <For each={feedback()!.improvements}>
                    {(i) => <li>{i}</li>}
                  </For>
                </ul>
              </Show>
              <h4>Code Quality:</h4>
              <p>{feedback()!.codeQuality}</p>
              <h4>Correctness:</h4>
              <p>{feedback()!.correctness}</p>
              <Show when={feedback()!.suggestions}>
                <h4>Suggestions:</h4>
                <p>{feedback()!.suggestions}</p>
              </Show>
            </Show>
          </div>
        </div>
      </div>

      {/* Loading Overlay */}
      <Show when={loading()}>
        <div class="loading-overlay">
          <div class="spinner"></div>
          <p>Processing...</p>
        </div>
      </Show>
    </div>
  );
}

export default App;

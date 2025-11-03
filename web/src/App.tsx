import { createSignal, Show, For, onMount } from "solid-js";

// ============================================================================
// TYPE DEFINITIONS
// ============================================================================

interface TechTree {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

interface Node {
  id: string;
  tech_tree_id: string;
  name: string;
  description: string;
  order_index: number;
  created_at: string;
}

interface Exercise {
  id: string;
  node_id: string;
  title: string;
  description: string;
  difficulty: string;
  order_index: number;
  content: {
    starter_code?: string;
    language?: string;
    test_cases?: TestCase[];
    hints?: string[];
    learning_objectives?: string[];
    evaluation_prompt?: string;
  };
  created_at: string;
}

interface TestCase {
  name: string;
  code: string;
}

interface Submission {
  id: string;
  exercise_id: string;
  user_id: string;
  code: string;
  language: string;
  test_results: {
    total: number;
    passed: number;
    failed: number;
    details: TestResult[];
  };
  ai_feedback: {
    score: number;
    strengths: string[];
    improvements: string[];
    suggestions: string;
  };
  total_score: number;
  created_at: string;
}

interface TestResult {
  name: string;
  passed: boolean;
  output?: string;
  error?: string;
}

interface UserProgress {
  node_id: string;
  exercises_completed: number;
  total_exercises: number;
  best_scores: { [exerciseId: string]: number };
  mastery_achieved: boolean;
  is_unlocked: boolean;
}

interface NodeWithProgress extends Node {
  progress?: UserProgress;
  exercises?: Exercise[];
}

interface TechTreeWithNodes extends TechTree {
  nodes: NodeWithProgress[];
}

interface ExerciseWithHistory extends Exercise {
  submissions: Submission[];
  best_score: number;
  attempt_count: number;
}

// ============================================================================
// MAIN APP COMPONENT
// ============================================================================

type ViewType = "tech-tree" | "node-detail" | "exercise";

function App() {
  // Navigation State
  const [currentView, setCurrentView] = createSignal<ViewType>("tech-tree");
  const [currentTechTreeId, setCurrentTechTreeId] = createSignal<string>("");
  const [currentNodeId, setCurrentNodeId] = createSignal<string>("");

  // Data State
  const [techTree, setTechTree] = createSignal<TechTreeWithNodes | null>(null);
  const [currentNode, setCurrentNode] = createSignal<NodeWithProgress | null>(null);
  const [currentExercise, setCurrentExercise] = createSignal<ExerciseWithHistory | null>(null);

  // Editor State
  const [code, setCode] = createSignal("// Write your code here\n");
  const [output, setOutput] = createSignal("");
  const [outputClass, setOutputClass] = createSignal("output");

  // UI State
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  // ============================================================================
  // NAVIGATION FUNCTIONS
  // ============================================================================

  const navigateToTechTree = (treeId: string) => {
    setCurrentTechTreeId(treeId);
    setCurrentView("tech-tree");
    loadTechTree(treeId);
  };

  const navigateToNode = (nodeId: string) => {
    setCurrentNodeId(nodeId);
    setCurrentView("node-detail");
    loadNode(nodeId);
  };

  const navigateToExercise = (exerciseId: string) => {
    setCurrentView("exercise");
    loadExercise(exerciseId);
  };

  const navigateToNextExercise = () => {
    const node = currentNode();
    const currentEx = currentExercise();
    if (!node || !node.exercises || !currentEx) return;

    const currentIndex = node.exercises.findIndex(ex => ex.id === currentEx.id);
    if (currentIndex >= 0 && currentIndex < node.exercises.length - 1) {
      navigateToExercise(node.exercises[currentIndex + 1].id);
    } else {
      // Last exercise in node, go back to node view
      navigateToNode(currentNodeId());
    }
  };

  // ============================================================================
  // DATA LOADING FUNCTIONS
  // ============================================================================

  const loadTechTree = async (treeId: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/tech-tree/${treeId}`);
      if (!response.ok) throw new Error("Failed to load tech tree");
      const data = await response.json();

      if (data.success) {
        setTechTree(data.tree);
      } else {
        setError(data.error || "Failed to load tech tree");
      }
    } catch (err) {
      console.error("Error loading tech tree:", err);
      setError("Error connecting to server. Make sure the server is running.");
    } finally {
      setLoading(false);
    }
  };

  const loadNode = async (nodeId: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/node/${nodeId}`);
      if (!response.ok) throw new Error("Failed to load node");
      const data = await response.json();

      if (data.success) {
        setCurrentNode(data.node);
      } else {
        setError(data.error || "Failed to load node");
      }
    } catch (err) {
      console.error("Error loading node:", err);
      setError("Error connecting to server. Make sure the server is running.");
    } finally {
      setLoading(false);
    }
  };

  const loadExercise = async (exerciseId: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/exercise/${exerciseId}`);
      if (!response.ok) throw new Error("Failed to load exercise");
      const data = await response.json();

      if (data.success) {
        setCurrentExercise(data.exercise);
        // Set starter code
        if (data.exercise.content?.starter_code) {
          setCode(data.exercise.content.starter_code);
        } else {
          setCode("// Write your code here\n");
        }
        setOutput("");
      } else {
        setError(data.error || "Failed to load exercise");
      }
    } catch (err) {
      console.error("Error loading exercise:", err);
      setError("Error connecting to server. Make sure the server is running.");
    } finally {
      setLoading(false);
    }
  };

  // ============================================================================
  // CODE EXECUTION & SUBMISSION
  // ============================================================================

  const runCode = async () => {
    const exercise = currentExercise();
    if (!exercise || !code().trim()) {
      setOutputClass("output error");
      setOutput("Please write some code first.");
      return;
    }

    setLoading(true);
    setOutputClass("output");
    setOutput("Running code...");

    try {
      const language = exercise.content?.language || "rust";
      const response = await fetch("/api/submission/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code: code(), language }),
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
      setOutput("Error connecting to server.");
    } finally {
      setLoading(false);
    }
  };

  const submitCode = async () => {
    const exercise = currentExercise();
    if (!exercise || !code().trim()) {
      alert("Please write some code first.");
      return;
    }

    setLoading(true);
    try {
      const language = exercise.content?.language || "rust";
      const response = await fetch(`/api/submission/${exercise.id}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code: code(), language }),
      });
      const data = await response.json();

      if (data.success) {
        // Reload exercise to get updated submission history
        await loadExercise(exercise.id);
        // Scroll to feedback section
        setTimeout(() => {
          const feedbackSection = document.querySelector('.submission-history');
          feedbackSection?.scrollIntoView({ behavior: 'smooth' });
        }, 100);
      } else {
        alert(data.error || "Failed to submit code. Please try again.");
      }
    } catch (error) {
      console.error("Error submitting code:", error);
      alert("Error connecting to server.");
    } finally {
      setLoading(false);
    }
  };

  // ============================================================================
  // UTILITY FUNCTIONS
  // ============================================================================

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

  const clearEditor = () => {
    if (confirm("Are you sure you want to clear the editor?")) {
      const exercise = currentExercise();
      if (exercise?.content?.starter_code) {
        setCode(exercise.content.starter_code);
      } else {
        setCode("// Write your code here\n");
      }
      setOutput("");
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  // ============================================================================
  // LIFECYCLE
  // ============================================================================

  onMount(() => {
    // Load the Rust Fundamentals tech tree
    // Tech tree ID from content/tech-trees.json
    navigateToTechTree("10000000-0000-0000-0000-000000000001");
  });

  // ============================================================================
  // RENDER
  // ============================================================================

  return (
    <div class="container">
      <header>
        <h1>🎓 Code Auto-Grader</h1>
        <p>LLM-Powered Code Learning Platform</p>
      </header>

      {/* Navigation Breadcrumbs */}
      <Show when={currentView() !== "tech-tree"}>
        <div class="breadcrumbs">
          <button class="breadcrumb-link" onClick={() => navigateToTechTree(currentTechTreeId())}>
            {techTree()?.name || "Tech Tree"}
          </button>
          <Show when={currentView() === "exercise"}>
            <span class="breadcrumb-separator">→</span>
            <button class="breadcrumb-link" onClick={() => navigateToNode(currentNodeId())}>
              {currentNode()?.name || "Node"}
            </button>
            <span class="breadcrumb-separator">→</span>
            <span class="breadcrumb-current">{currentExercise()?.title || "Exercise"}</span>
          </Show>
          <Show when={currentView() === "node-detail"}>
            <span class="breadcrumb-separator">→</span>
            <span class="breadcrumb-current">{currentNode()?.name || "Node"}</span>
          </Show>
        </div>
      </Show>

      {/* Error Display */}
      <Show when={error()}>
        <div class="error-banner">
          <strong>Error:</strong> {error()}
          <button onClick={() => setError(null)}>✕</button>
        </div>
      </Show>

      {/* TECH TREE VIEW */}
      <Show when={currentView() === "tech-tree"}>
        <div class="content">
          <div class="tech-tree-view">
            <Show when={techTree()}>
              <div class="tech-tree-header">
                <h2>{techTree()!.name}</h2>
                <p>{techTree()!.description}</p>
              </div>

              <div class="nodes-list">
                <For each={techTree()!.nodes}>
                  {(node, index) => {
                    const progress = node.progress;
                    const isLocked = progress ? !progress.is_unlocked : index() > 0;
                    const isComplete = progress?.mastery_achieved || false;
                    const exercisesCompleted = progress?.exercises_completed || 0;
                    const totalExercises = progress?.total_exercises || 0;

                    return (
                      <div
                        class={`node-card ${isLocked ? 'locked' : ''} ${isComplete ? 'complete' : ''}`}
                        onClick={() => !isLocked && navigateToNode(node.id)}
                      >
                        <div class="node-card-header">
                          <h3>
                            <Show when={isLocked}>🔒 </Show>
                            <Show when={isComplete}>✅ </Show>
                            {node.name}
                          </h3>
                          <Show when={!isLocked}>
                            <div class="node-progress-badge">
                              {exercisesCompleted}/{totalExercises}
                            </div>
                          </Show>
                        </div>
                        <p>{node.description}</p>
                        <Show when={!isLocked && !isComplete}>
                          <div class="progress-bar">
                            <div
                              class="progress-fill"
                              style={{ width: `${totalExercises > 0 ? (exercisesCompleted / totalExercises) * 100 : 0}%` }}
                            />
                          </div>
                        </Show>
                        <Show when={isLocked}>
                          <div class="locked-message">
                            Complete previous nodes to unlock
                          </div>
                        </Show>
                      </div>
                    );
                  }}
                </For>
              </div>
            </Show>
          </div>
        </div>
      </Show>

      {/* NODE DETAIL VIEW */}
      <Show when={currentView() === "node-detail"}>
        <div class="content">
          <div class="node-detail-view">
            <Show when={currentNode()}>
              <div class="node-detail-header">
                <h2>{currentNode()!.name}</h2>
                <p>{currentNode()!.description}</p>
                <Show when={currentNode()!.progress}>
                  <div class="node-progress-summary">
                    <Show when={currentNode()!.progress!.mastery_achieved}>
                      <div class="mastery-badge">✅ Mastery Achieved!</div>
                    </Show>
                    <div class="progress-stats">
                      Exercises Completed: {currentNode()!.progress!.exercises_completed} / {currentNode()!.progress!.total_exercises}
                    </div>
                  </div>
                </Show>
              </div>

              <div class="exercises-list">
                <h3>Exercises</h3>
                <For each={currentNode()!.exercises}>
                  {(exercise) => {
                    const progress = currentNode()!.progress;
                    const bestScore = progress?.best_scores?.[exercise.id] || 0;
                    const isComplete = bestScore >= 80;

                    return (
                      <div
                        class={`exercise-card ${isComplete ? 'complete' : ''}`}
                        onClick={() => navigateToExercise(exercise.id)}
                      >
                        <div class="exercise-card-header">
                          <h4>
                            <Show when={isComplete}>✅ </Show>
                            {exercise.title}
                          </h4>
                          <div class="exercise-meta">
                            <span class={`difficulty ${exercise.difficulty}`}>
                              {exercise.difficulty}
                            </span>
                            <Show when={bestScore > 0}>
                              <span class={`score-badge ${bestScore >= 80 ? 'high' : bestScore >= 60 ? 'medium' : 'low'}`}>
                                Best: {bestScore}%
                              </span>
                            </Show>
                          </div>
                        </div>
                        <p>{exercise.description}</p>
                        <button class="btn btn-primary">
                          {bestScore > 0 ? 'Continue Exercise' : 'Start Exercise'}
                        </button>
                      </div>
                    );
                  }}
                </For>
              </div>
            </Show>
          </div>
        </div>
      </Show>

      {/* EXERCISE VIEW */}
      <Show when={currentView() === "exercise"}>
        <div class="content">
          <div class="exercise-view">
            <Show when={currentExercise()}>
              {/* Exercise Description Panel */}
              <div class="exercise-panel">
                <div class="panel-header">
                  <h2>{currentExercise()!.title}</h2>
                  <div class="exercise-meta">
                    <span class={`difficulty ${currentExercise()!.difficulty}`}>
                      {currentExercise()!.difficulty}
                    </span>
                    <span class="language-badge">
                      {currentExercise()!.content?.language || "rust"}
                    </span>
                    <Show when={currentExercise()!.best_score > 0}>
                      <span class={`score-badge ${currentExercise()!.best_score >= 80 ? 'high' : currentExercise()!.best_score >= 60 ? 'medium' : 'low'}`}>
                        Best: {currentExercise()!.best_score}%
                      </span>
                    </Show>
                    <span class="attempt-badge">
                      Attempts: {currentExercise()!.attempt_count}
                    </span>
                  </div>
                </div>

                <div class="exercise-content">
                  <p>{currentExercise()!.description}</p>

                  <Show when={currentExercise()!.content?.learning_objectives && currentExercise()!.content.learning_objectives!.length > 0}>
                    <div class="learning-objectives">
                      <h4>Learning Objectives:</h4>
                      <ul>
                        <For each={currentExercise()!.content.learning_objectives}>
                          {(objective) => <li>{objective}</li>}
                        </For>
                      </ul>
                    </div>
                  </Show>

                  <Show when={currentExercise()!.content?.hints && currentExercise()!.content.hints!.length > 0}>
                    <div class="hints">
                      <details>
                        <summary>💡 Show Hints</summary>
                        <ul>
                          <For each={currentExercise()!.content.hints}>
                            {(hint) => <li>{hint}</li>}
                          </For>
                        </ul>
                      </details>
                    </div>
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
                      Reset
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
                  <div class={outputClass()}>{output() || "Run your code to see output here..."}</div>
                </div>
              </div>

              {/* Submission History Panel */}
              <div class="submission-history">
                <h2>Submission History</h2>
                <Show
                  when={currentExercise()!.submissions && currentExercise()!.submissions.length > 0}
                  fallback={
                    <p class="info-message">
                      No submissions yet. Submit your code to receive AI-powered feedback and grading.
                    </p>
                  }
                >
                  <div class="submissions-list">
                    <For each={currentExercise()!.submissions}>
                      {(submission, index) => (
                        <div class="submission-card">
                          <div class="submission-header">
                            <h3>Attempt #{currentExercise()!.submissions.length - index()}</h3>
                            <div class="submission-meta">
                              <span class="submission-date">{formatDate(submission.created_at)}</span>
                              <span class={`submission-score ${submission.total_score >= 80 ? 'high' : submission.total_score >= 60 ? 'medium' : 'low'}`}>
                                Score: {submission.total_score}%
                              </span>
                            </div>
                          </div>

                          <div class="submission-details">
                            {/* Test Results */}
                            <div class="test-results">
                              <h4>Test Results: {submission.test_results.passed}/{submission.test_results.total} passed</h4>
                              <div class="test-results-list">
                                <For each={submission.test_results.details}>
                                  {(test) => (
                                    <div class={`test-result ${test.passed ? 'passed' : 'failed'}`}>
                                      <span class="test-icon">{test.passed ? '✅' : '❌'}</span>
                                      <span class="test-name">{test.name}</span>
                                      <Show when={test.error}>
                                        <div class="test-error">{test.error}</div>
                                      </Show>
                                    </div>
                                  )}
                                </For>
                              </div>
                            </div>

                            {/* AI Feedback */}
                            <div class="ai-feedback">
                              <h4>AI Feedback (Score: {submission.ai_feedback.score}%)</h4>

                              <Show when={submission.ai_feedback.strengths && submission.ai_feedback.strengths.length > 0}>
                                <div class="feedback-section">
                                  <h5>✓ Strengths:</h5>
                                  <ul>
                                    <For each={submission.ai_feedback.strengths}>
                                      {(strength) => <li>{strength}</li>}
                                    </For>
                                  </ul>
                                </div>
                              </Show>

                              <Show when={submission.ai_feedback.improvements && submission.ai_feedback.improvements.length > 0}>
                                <div class="feedback-section">
                                  <h5>→ Areas for Improvement:</h5>
                                  <ul>
                                    <For each={submission.ai_feedback.improvements}>
                                      {(improvement) => <li>{improvement}</li>}
                                    </For>
                                  </ul>
                                </div>
                              </Show>

                              <Show when={submission.ai_feedback.suggestions}>
                                <div class="feedback-section">
                                  <h5>💡 Suggestions:</h5>
                                  <p>{submission.ai_feedback.suggestions}</p>
                                </div>
                              </Show>
                            </div>

                            {/* Code Display */}
                            <details class="submission-code">
                              <summary>View Code</summary>
                              <pre><code>{submission.code}</code></pre>
                            </details>
                          </div>
                        </div>
                      )}
                    </For>
                  </div>
                </Show>
              </div>

              {/* Next Exercise Button */}
              <div class="exercise-navigation">
                <button class="btn btn-primary btn-large" onClick={navigateToNextExercise}>
                  Next Exercise →
                </button>
              </div>
            </Show>
          </div>
        </div>
      </Show>

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

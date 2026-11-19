# Development Plan

Code Auto-Grader helps self-taught developers build programming skills through lessons, hands-on practice, and feedback tied to their code. The single-user Rust prototype has three seeded nodes and nine exercises, a browser editor, Docker execution, submission history, and progress tracking. The next milestone is a reliable Rust Basics learning loop.

## Phase 1 — Validate the Rust learning loop

1. **Fix grading and progression.** Give each exercise one function to implement or repair, with a clear expected result. Test the learner’s implementation and verify that the starter, an incorrect solution, and the reference solution produce the expected results. Use test-only grading for deterministic tasks. Reserve automatic AI review for exercises where implementation choices, design, or code quality warrant it. Store raw test and AI percentages and the weights used for each submission. Distinguish a skipped or unavailable AI review from a score of zero.
2. **Make mastery reliable.** Count each mastered exercise once, using its best qualifying result and the node’s mastery threshold. Unlock nodes when their prerequisites are met. Make submission and progress updates consistent across retries and partial failures. Verify the behavior with existing submission history and concurrent attempts.
3. **Connect lessons, exercises, and help.** Seed a lesson for each node and show it in the learner UI. Keep exercise instructions focused on what to implement and put concept explanations in the lesson. Offer a hint or **Explain this error** after a failed run. Use the actual compiler diagnostic or test failure to guide the learner. Track attempts and help requests per exercise, and improve lessons when learners repeatedly get stuck on the same step.
4. **Improve the editor workflow.** Preserve code after a submission and make the editor the main workspace. Use collapsible navigation for the tree, lessons, and exercise instructions. Add a resizable results panel with separate tabs for execution output and AI feedback. Use accessible light and dark terminal themes. [Terminal design reference](dogfooding/90_STYLE_TERMINAL.png)

**Exit criteria.** From a fresh seed, a learner can read the Basics lesson, recover from a failed run using an explanation, submit a passing solution, retry without inflating progress, and unlock Ownership & Borrowing after completing Basics. Apply the same pattern to the other two seeded nodes. Check that learners can explain and fix an ownership error and solve a related exercise independently after receiving feedback.

### Rust curriculum

| Node | Practice | Prerequisites | Scope |
| --- | --- | --- | --- |
| Rust Basics | Variables and mutability, functions and return values, control flow | None | Seeded |
| Ownership & Borrowing | Moves, references, borrowing rules | Rust Basics | Seeded |
| Structs & Enums | Structs, methods, enum variants and pattern matching | Ownership & Borrowing | Seeded |
| Collections & Strings | Vectors, strings, counting with hash maps | Ownership & Borrowing | Next content after validation |
| Error Handling | `Option`, `Result`, error propagation | Ownership & Borrowing | Next content after validation |
| Traits & Generics | Traits, generic functions and types, trait bounds | Structs & Enums, Collections & Strings, Error Handling | Later |
| Iterators & Closures | Capturing variables, iterator methods, implementing `Iterator` | Traits & Generics | Optional |

## Phase 2 — Private multi-user beta

- Replace the shared learner identity with authenticated accounts and isolate submission history and progress per learner. Protect content editing and manual unlocks with admin permissions. Verify runner isolation, execution quotas, and AI request limits before accepting external submissions.
- Add conceptual text exercises and focused debugging tasks. Give admins the ability to preview and edit content as beta feedback comes in.
- Start with a small cohort and expand toward 50 learners through developer communities. Collect feedback in the app. Baseline completion, repeat practice, AI helpfulness, and issues that block learning. Initial targets are 20% of beta users completing all seeded nodes and 70% of respondents finding AI feedback helpful.
- Add automated checks for grading, progression, and account isolation. Track AI usage and cost per learner and set budget alerts.

## Phase 3 — Expand learning and content

- Recommend the next node using mastery, time spent, and recurring errors. Explain the recommendation when the learner finishes a node.
- Add multi-step projects with repository evaluation and system design exercises with diagrams and written explanations.
- Add exercise voting and community submissions with moderator review. Pilot generated lessons and exercises from uploaded learning materials with human review before publishing.
- Use learning data to tune difficulty, recommend help when learners get stuck, and measure where they disengage. Evaluate streaks and other motivation features against repeat practice.

## Phase 4 — Grow after beta

Start with self-taught developers, followed by students and teams. Recruit through developer communities, demos, and educational content. Track exercises completed per week, nodes mastered per month, and retention. Test demand for paid learning paths and contextual feedback before investing in paid acquisition. If team pilots demonstrate demand, expand into custom learning paths and team reporting.

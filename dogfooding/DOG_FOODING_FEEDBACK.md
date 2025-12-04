## Learnings from dogfooding

IMPORTANT: Ask for clarification if anything is clear.

1. Each exercise should have a single function to fix / debug. It's bad practice to have a main function with many unrelated assignments inside.

2. I'd expect that the function does something concrete, like returns value.

3. Split exercises when needed to fix the above issues.

4. Design a better AI feedback logic. Currently every submission gets AI feedback even if trivial. Some conditions when feedback is warranted:
   a. exercise has multiple ways of implementing
   b. user has made a clear error even though the answer is correct
   c. exercise has design elements
   d. exercise is a project

5. hints should be provided only on the basis of failed attempts. If instructions and learning materials are insufficient, they need to be improved. (Evaluate this suggestion)

6. There should be a separate feature of asking AI assistance & interpretation on the output received from the code execution output (e.g. Rust compiler error)

7. The database schema and business logic for AI vs code execution based scoring needs to be changed. There should be one column for AI feedback percentage (out of 100) and another for test based scoring (out of 100). Then there should be an AI feedback weight column (or some other way of indicating the weights). This allows for ignoring AI feedback for some exercises.

8. There needs to be a separate database level for learning material content. Each node in the tech tree should have a learning material section that explains the core concepts in the node. The exercises should focus on that material only and have their own instructions. The learning materials should provide the necessary materials for the exercises. The exercise instructions can offer refreshers and remind about specific syntax without revealing the solution that the student needs to struggle with.

9. After learning material change from item 8 is ready, the db seed content and content/ .md's need to be updated.

10. In addition to counting attempts, the exercise data model should keep track of hints that the user has requested from the API when asking to explain a test outcome or other feedback.

11. The UI layout needs an overhaul.
    a. the editor has to take the whole screen space, no large header
    b. there should be a hideable left hand side navigation where the user can go back to the tree view. It should also have the ability to navigate to visited nodes and exercises, and see upcoming exercises.
    c. the left hand navigation should also allow the user go back to the learning material
    d. the left hand panel should have the instructions for the current exercise
    e. there should be a right hand panel that can be hidden and resized
    f. the right hand panel should have tabs for: code execution outcomes and AI feedback for the current submission.
    g. the right hand panel code execution outcome should have buttons that send a request to the API to explain the outcome (e.g. Rust compile error)
    f. the overall look of the UI should be a stylized terminal. In dark mode: black/dark for backgrounds, and white for outlines, contrast, and text.
    g. devise a palette for accent colors
    h. use accent colors for errors, success and decoration, but sparingly. Make sure the palette is accessibility compliant.
    i. elements, such as buttons, should be angular and outlined, with sharp, hatched box shadows offset to the lower-right direction. See the provided image "90_STYLE_TERMINAL.png" as inspiration.
    j. add dark and light modes to the UI

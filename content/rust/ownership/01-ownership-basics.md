{
  "id": "30000000-0000-0000-0000-000000000004",
  "node_id": "20000000-0000-0000-0000-000000000002",
  "title": "Exercise 2.1: Ownership and Moves",
  "slug": "ownership-basics",
  "difficulty": "medium",
  "estimated_minutes": 30,
  "order_index": 1
}
---

# Exercise 2.1: Ownership and Moves

## Learning Objectives
- Understand Rust's ownership model (every value has one owner)
- Recognize when ownership is transferred (moved)
- Distinguish between Copy types (stack) and Move types (heap)
- Use `.clone()` when you need deep copies
- Understand why Rust prevents use-after-move

## Problem Description

Welcome to Rust's **ownership system** - the concept that makes Rust unique!

### The Rules of Ownership

1. Each value in Rust has a single **owner** (the variable that holds it)
2. When the owner goes out of scope, the value is **dropped** (memory freed)
3. When you assign or pass a value to a function, ownership is **moved** (transferred)

### Copy vs Move

- **Copy types** (stored on stack): integers, booleans, floats, char
  - These are copied automatically when assigned
  - `let x = 5; let y = x;` → Both `x` and `y` are valid

- **Move types** (stored on heap): String, Vec, custom structs without Copy
  - These are moved when assigned
  - `let s1 = String::from("hello"); let s2 = s1;` → `s1` is now invalid!

### Your Task

You'll fix several ownership errors in code that tries to use values after they've been moved. You'll learn to:
- Recognize move semantics
- Use references to borrow instead of taking ownership
- Clone data when you need multiple owners

## Starting Code

```rust
fn main() {
    // Part 1: Basic move semantics
    println!("=== Part 1: String Moves ===");
    let s1 = String::from("hello");
    let s2 = s1;  // s1's ownership moves to s2

    println!("s1: {}", s1);  // ERROR: s1 is no longer valid
    println!("s2: {}", s2);

    // Part 2: Function call moves
    println!("\n=== Part 2: Function Moves ===");
    let message = String::from("Hello, Rust!");
    print_message(message);  // Ownership moves into function
    println!("Message again: {}", message);  // ERROR: message was moved

    // Part 3: Return values and ownership
    println!("\n=== Part 3: Returning Ownership ===");
    let greeting = String::from("Hi");
    let greeting = add_exclamation(greeting);  // Ownership moves in and back out
    println!("{}", greeting);  // This should work

    // Part 4: Multiple uses with clone
    println!("\n=== Part 4: Cloning ===");
    let original = String::from("original");
    let copy = original;  // Move happens here
    println!("Original: {}", original);  // ERROR: original was moved
    println!("Copy: {}", copy);

    // Part 5: Copy types vs Move types
    println!("\n=== Part 5: Copy Types ===");
    let x = 42;
    let y = x;  // Copy happens (integers implement Copy trait)
    println!("x: {}, y: {}", x, y);  // Both are valid!

    println!("\n✅ All parts working!");
}

fn print_message(msg: String) {
    println!("Message: {}", msg);
    // msg is dropped here when function ends
}

fn add_exclamation(mut s: String) -> String {
    s.push_str("!");
    s  // Return ownership back to caller
}
```

## Expected Behavior

After fixing the code, it should:
1. Compile without ownership errors
2. Print all messages correctly
3. Demonstrate understanding of moves, clones, and borrows

## Hints

<details>
<summary>Hint 1: Understanding the Errors</summary>

The compiler error messages are your friend! They'll say things like:
```
error[E0382]: borrow of moved value: `s1`
```

This means `s1` was moved somewhere else and you can't use it anymore.

When you see this error, you have 3 options:
1. **Don't use the moved value** (use the new owner instead)
2. **Clone before moving** (`let s2 = s1.clone();`)
3. **Borrow instead of move** (use `&s1` - covered in next exercise)
</details>

<details>
<summary>Hint 2: Part 1 - String Moves</summary>

In Part 1, `s1` is moved to `s2`. You can't use both.

Options:
- Remove the line that uses `s1`
- Or: Clone before moving: `let s2 = s1.clone();`
- Or: Only use `s2` and remove `s1` reference
</details>

<details>
<summary>Hint 3: Part 2 - Function Moves</summary>

When you call `print_message(message)`, ownership of `message` moves into the function. After the function ends, `message` is dropped.

Options:
- Remove the second `println!` that uses `message`
- Or: Clone before passing: `print_message(message.clone());`
- Or: Have the function return the String back
- Or: Pass a reference instead (we'll learn this next!)
</details>

<details>
<summary>Hint 4: Part 3 - This One Works!</summary>

Part 3 actually works correctly! The function takes ownership and returns it back. This pattern is common in Rust.

No fix needed here.
</details>

<details>
<summary>Hint 5: Part 4 - Clone When You Need Multiple Owners</summary>

If you need both `original` and `copy` to be valid, use `.clone()`:
```rust
let copy = original.clone();  // Deep copy the data
```

This creates a new String with the same contents, so both are valid owners.
</details>

<details>
<summary>Hint 6: Part 5 - This One Works Too!</summary>

Part 5 works because integers implement the `Copy` trait. They're stored on the stack and copied automatically.

No fix needed here.
</details>

## Reference Solution

```rust
fn main() {
    // Part 1: Basic move semantics
    println!("=== Part 1: String Moves ===");
    let s1 = String::from("hello");
    let s2 = s1.clone();  // FIX: Clone instead of move

    println!("s1: {}", s1);  // Now valid because we cloned
    println!("s2: {}", s2);

    // Part 2: Function call moves
    println!("\n=== Part 2: Function Moves ===");
    let message = String::from("Hello, Rust!");
    print_message(message.clone());  // FIX: Clone before passing
    println!("Message again: {}", message);  // Now valid

    // Part 3: Return values and ownership (already works!)
    println!("\n=== Part 3: Returning Ownership ===");
    let greeting = String::from("Hi");
    let greeting = add_exclamation(greeting);
    println!("{}", greeting);

    // Part 4: Multiple uses with clone
    println!("\n=== Part 4: Cloning ===");
    let original = String::from("original");
    let copy = original.clone();  // FIX: Clone to create independent copy
    println!("Original: {}", original);  // Now valid
    println!("Copy: {}", copy);

    // Part 5: Copy types vs Move types (already works!)
    println!("\n=== Part 5: Copy Types ===");
    let x = 42;
    let y = x;  // Copy happens automatically for integers
    println!("x: {}, y: {}", x, y);

    println!("\n✅ All parts working!");
}

fn print_message(msg: String) {
    println!("Message: {}", msg);
}

fn add_exclamation(mut s: String) -> String {
    s.push_str("!");
    s
}
```

### Alternative Solution (More Efficient)

```rust
// This solution uses borrowing (which we'll learn next) to avoid cloning
// For now, the clone-based solution above is fine!

fn main() {
    let s1 = String::from("hello");
    let s2 = s1.clone();  // Need clone here to use both
    println!("s1: {}, s2: {}", s1, s2);

    let message = String::from("Hello, Rust!");
    print_message_borrowed(&message);  // Borrow instead of move
    println!("Message again: {}", message);

    // ... rest is same
}

fn print_message_borrowed(msg: &String) {
    println!("Message: {}", msg);
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_string_cloning() {
        let s1 = String::from("hello");
        let s2 = s1.clone();
        assert_eq!(s1, "hello");
        assert_eq!(s2, "hello");
        // Both should be valid
    }

    #[test]
    fn test_function_with_clone() {
        let message = String::from("test");
        let cloned = message.clone();
        print_message(cloned);
        // Original message should still be valid
        assert_eq!(message, "test");
    }

    #[test]
    fn test_ownership_transfer_and_return() {
        let s = String::from("hello");
        let s = add_exclamation(s);
        assert_eq!(s, "hello!");
    }

    #[test]
    fn test_copy_types_dont_move() {
        let x = 42;
        let y = x;
        assert_eq!(x, 42);  // x is still valid
        assert_eq!(y, 42);  // y is also valid
    }

    #[test]
    fn test_multiple_clones() {
        let original = String::from("data");
        let copy1 = original.clone();
        let copy2 = original.clone();
        let copy3 = copy1.clone();

        assert_eq!(original, "data");
        assert_eq!(copy1, "data");
        assert_eq!(copy2, "data");
        assert_eq!(copy3, "data");
    }

    #[test]
    fn test_ownership_transfer_chain() {
        let s1 = String::from("start");
        let s2 = s1;  // Move
        let s3 = s2;  // Move again
        // Only s3 is valid now
        assert_eq!(s3, "start");
    }

    #[test]
    fn test_function_takes_ownership() {
        let s = String::from("test");
        consume(s);
        // s is no longer valid here
    }

    fn consume(_s: String) {
        // _s is dropped when function ends
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All 7 tests pass: 60 points
- 5-6 tests pass: 48 points
- 3-4 tests pass: 30 points
- <3 tests pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Understanding of move semantics (15 points)**
   - Recognizes when moves happen
   - Understands that moved values become invalid
   - Correctly identifies heap vs stack types

2. **Appropriate use of `.clone()` (10 points)**
   - Uses clone when truly needing multiple owners
   - Doesn't clone unnecessarily (e.g., in Part 5 with integers)
   - Understands performance implications of cloning

3. **Code organization (10 points)**
   - Clear, readable fixes
   - Good comments explaining ownership decisions
   - Follows Rust conventions

4. **Conceptual understanding (5 points)**
   - Comments or code show grasp of ownership rules
   - Can explain *why* certain fixes work

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on ownership and move semantics.

Exercise: Fix ownership errors in code that tries to use values after they've been moved.

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Understanding of when moves happen (assignments, function calls)
2. Appropriate use of .clone() (not overusing or underusing)
3. Recognition of Copy types vs Move types
4. Code clarity and comments

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. Focus on:\n- If student clones everything: explain performance implications and that we'll learn borrowing next\n- If student doesn't understand moves: explain ownership transfer concept\n- If errors remain: guide toward understanding which values are moved where\n- Praise correct understanding of Copy vs Move types",
  "keyTakeaways": ["Ownership transfer", "Move vs Copy", "When to use .clone()"],
  "encouragement": "Motivational message. Ownership is hard - acknowledge the challenge!"
}

Common patterns:
- Student clones everywhere → Works but inefficient, prepare for borrowing
- Student removes functionality → Point out they can use clone or borrowing
- Student confused about Copy types → Explain stack vs heap allocation
- Student gets it right → Praise understanding, preview that borrowing is more efficient
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Understand the ownership rules (one owner, transfer on move)
- ✅ Recognize when ownership moves happen
- ✅ Know the difference between Copy and Move types
- ✅ Use `.clone()` to create independent copies
- ✅ Read ownership-related compiler errors
- ✅ Feel ready to learn borrowing (more efficient than cloning)

## Common Mistakes to Watch For

1. **Cloning everything**: Students might add `.clone()` everywhere without understanding why
   - Guide: Cloning works but is inefficient; borrowing (next exercise) is better

2. **Not understanding Copy types**: Thinking integers also move
   - Guide: Explain stack vs heap allocation

3. **Removing functionality**: Instead of fixing moves, students delete the problematic code
   - Guide: Explain the goal is to make both uses work, not remove one

4. **Confusion about ownership return**: Not understanding Part 3 works because ownership returns
   - Guide: Explain the "ownership chain" pattern

## Why This Exercise Works

This exercise is pedagogically effective because:

1. **Immediate feedback**: Compiler errors show exactly what's wrong
2. **Progressive difficulty**: 5 parts, from simple to nuanced
3. **Multiple solutions**: Students can clone OR restructure
4. **Prepares for next**: Sets up motivation for borrowing (more efficient)
5. **Practical context**: Real scenarios, not abstract examples

The "aha moment" comes when students realize:
> "Oh! I can't use a value after it's moved. I need to either clone it or... wait, isn't there a better way? (Yes - borrowing!)"

## Extension Ideas (Post-MVP)

- Add visualization showing heap vs stack
- Interactive diagram of ownership transfers
- Challenge: "Minimize clones" mode (find solution with fewest clones)
- Compare memory usage with/without clones

## Notes for Implementation

When storing this exercise in the database:
- Mark Parts 3 & 5 as "already correct" in hints (avoid confusion)
- The starting code should NOT compile (ownership errors are intentional)
- Reference solution should be well-commented for student review after completion

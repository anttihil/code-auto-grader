{
  "id": "30000000-0000-0000-0000-000000000001",
  "node_id": "20000000-0000-0000-0000-000000000001",
  "title": "Exercise 1.1: Variables and Mutability",
  "slug": "variables-mutability",
  "difficulty": "easy",
  "estimated_minutes": 15,
  "order_index": 1
}
---

# Exercise 1.1: Variables and Mutability

## Learning Objectives
- Understand Rust's variable binding with `let`
- Distinguish between mutable (`mut`) and immutable variables
- Recognize compile-time enforcement of mutability rules
- Apply type inference and explicit type annotations

## Problem Description

Welcome to your first Rust exercise! In Rust, variables are **immutable by default** - once you bind a value to a name, you can't change it unless you explicitly mark it as mutable with the `mut` keyword.

This is different from most programming languages and is a key part of Rust's safety guarantees. It helps prevent bugs by making it clear which values are meant to change.

Your task is to fix the compilation errors in the code below. The code tries to reassign variables, but some lack the `mut` keyword.

## Starting Code

```rust
fn main() {
    let x = 5;
    println!("The value of x is: {}", x);
    x = 6;
    println!("The value of x is now: {}", x);

    let y = 10;
    let y = y + 5;
    println!("The value of y is: {}", y);

    let spaces = "   ";
    let spaces = spaces.len();
    println!("Number of spaces: {}", spaces);

    let mut z = 1;
    println!("The value of z is: {}", z);
    z = 2;
    println!("The value of z is now: {}", z);
    z = "three";
    println!("The value of z is now: {}", z);
}
```

## Expected Behavior

After fixing the code, it should:
1. Compile without errors
2. Print the values correctly
3. Demonstrate proper use of mutability and shadowing

## Hints

<details>
<summary>Hint 1: Mutability vs Shadowing</summary>

There are two ways to "change" a value in Rust:
- **Mutability**: Mark a variable as `mut` to allow reassignment
- **Shadowing**: Use `let` again to create a new variable with the same name

Shadowing is useful when you want to change the type or transform a value while keeping the same name.
</details>

<details>
<summary>Hint 2: Type Consistency</summary>

Even with `mut`, you can't change a variable's type. Once `z` is an integer, it must stay an integer.

If you want to change types, you need shadowing (another `let` statement).
</details>

<details>
<summary>Hint 3: Which Lines Have Errors?</summary>

Compilation errors occur on these lines:
- Line 4: `x = 6` (x is immutable)
- Line 18: `z = "three"` (type mismatch)

The `y` and `spaces` examples use **shadowing**, which is perfectly valid!
</details>

## Reference Solution

```rust
fn main() {
    // Fix: Add `mut` to allow reassignment
    let mut x = 5;
    println!("The value of x is: {}", x);
    x = 6;
    println!("The value of x is now: {}", x);

    // This is shadowing (creating a new variable), not mutation
    // No `mut` needed
    let y = 10;
    let y = y + 5;
    println!("The value of y is: {}", y);

    // Another shadowing example - changing type from &str to usize
    let spaces = "   ";
    let spaces = spaces.len();
    println!("Number of spaces: {}", spaces);

    // This is mutation - z was declared with `mut`
    let mut z = 1;
    println!("The value of z is: {}", z);
    z = 2;
    println!("The value of z is now: {}", z);

    // Fix: Use shadowing to change type
    let z = "three";
    println!("The value of z is now: {}", z);
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    // Note: This exercise is primarily about compilation
    // The tests verify behavior after fixing compilation errors

    #[test]
    fn test_compiles() {
        // If this test runs, the code compiles successfully
        assert!(true);
    }

    #[test]
    fn test_x_mutability() {
        let mut x = 5;
        x = 6;
        assert_eq!(x, 6);
    }

    #[test]
    fn test_y_shadowing() {
        let y = 10;
        let y = y + 5;
        assert_eq!(y, 15);
    }

    #[test]
    fn test_spaces_shadowing_type_change() {
        let spaces = "   ";
        let spaces = spaces.len();
        assert_eq!(spaces, 3);
    }

    #[test]
    fn test_z_operations() {
        let mut z = 1;
        z = 2;
        assert_eq!(z, 2);

        let z = "three";
        assert_eq!(z, "three");
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All tests pass: 60 points
- 4/5 tests pass: 48 points
- 3/5 tests pass: 36 points
- <3 tests pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Correct use of `mut` (10 points)**
   - Uses `mut` only when reassigning variables
   - Doesn't use `mut` for shadowed variables

2. **Understanding of shadowing (10 points)**
   - Recognizes when shadowing is appropriate (type changes)
   - Doesn't unnecessarily shadow instead of using `mut`

3. **Code style (10 points)**
   - Follows Rust naming conventions
   - No unnecessary modifications to working code
   - Clean, readable solution

4. **Demonstrates understanding (10 points)**
   - Comments show grasp of mutability vs shadowing
   - Solution isn't just trial-and-error

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on variables and mutability.

Exercise: Fix compilation errors related to variable mutability and shadowing.

Student's submission:
{student_code}

Test results:
{test_results}

Evaluate the solution on:
1. Correct use of `mut` keyword
2. Understanding of shadowing vs mutation
3. Code style and clarity

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback explaining what was done well and what could improve",
  "keyTakeaways": ["Main lesson 1", "Main lesson 2"],
  "encouragement": "Motivational message"
}
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Understand that Rust variables are immutable by default
- ✅ Know when to use `mut` vs shadowing
- ✅ Recognize that shadowing can change types, mutation cannot
- ✅ Be comfortable with Rust's strict compile-time checks

## Common Mistakes to Watch For

1. **Adding `mut` everywhere**: Students might add `mut` to `y` and `spaces`, not recognizing shadowing
2. **Trying to change types with `mut`**: Attempting `mut z` and reassigning to a string
3. **Confusion about shadowing**: Not understanding that `let y = y + 5` creates a new variable

The AI feedback should gently correct these misconceptions while encouraging experimentation.

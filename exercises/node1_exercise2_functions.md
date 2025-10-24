# Exercise 1.2: Functions and Return Values

## Metadata
- **Node**: 1 (Rust Basics)
- **Exercise**: 2 of 3
- **Difficulty**: Easy
- **Estimated Time**: 20-25 minutes
- **Prerequisites**: Exercise 1.1 (Variables and Mutability)

## Learning Objectives
- Define functions with parameters and return types
- Understand Rust's expression-based return syntax
- Use type annotations for function signatures
- Apply functions to solve simple problems
- Distinguish between statements and expressions

## Problem Description

Functions are the building blocks of Rust programs. In Rust, functions:
- Are declared with the `fn` keyword
- Use snake_case naming convention
- Require explicit type annotations for parameters
- Can return values using `->` syntax
- Use **implicit returns** (the last expression without a semicolon)

Your task is to implement three functions with different purposes:

1. **`add`**: Add two integers
2. **`rectangle_area`**: Calculate the area of a rectangle
3. **`is_even`**: Check if a number is even

Pay attention to:
- Parameter types
- Return types
- Expression-based returns (no `return` keyword needed at the end)

## Starting Code

```rust
// TODO: Write a function called `add` that takes two i32 parameters
// and returns their sum as an i32

// TODO: Write a function called `rectangle_area` that takes width and height
// as f64 parameters and returns the area as f64

// TODO: Write a function called `is_even` that takes an i32
// and returns true if it's even, false otherwise
// Hint: Use the modulo operator %

fn main() {
    // Test add function
    let sum = add(5, 7);
    println!("5 + 7 = {}", sum);
    assert_eq!(sum, 12);

    let sum2 = add(-3, 3);
    println!("-3 + 3 = {}", sum2);
    assert_eq!(sum2, 0);

    // Test rectangle_area function
    let area = rectangle_area(10.5, 20.0);
    println!("Rectangle area: {}", area);
    assert_eq!(area, 210.0);

    let area2 = rectangle_area(0.0, 10.0);
    println!("Rectangle area (zero width): {}", area2);
    assert_eq!(area2, 0.0);

    // Test is_even function
    println!("Is 4 even? {}", is_even(4));
    assert_eq!(is_even(4), true);

    println!("Is 7 even? {}", is_even(7));
    assert_eq!(is_even(7), false);

    println!("Is 0 even? {}", is_even(0));
    assert_eq!(is_even(0), true);

    println!("All tests passed!");
}
```

## Expected Behavior

Your implementation should:
1. Define all three functions with correct signatures
2. Pass all assertions in main()
3. Use idiomatic Rust style (expression-based returns)

## Hints

<details>
<summary>Hint 1: Function Syntax</summary>

Rust function syntax:
```rust
fn function_name(param1: Type1, param2: Type2) -> ReturnType {
    // function body
    expression_without_semicolon  // implicit return
}
```

The last line without a semicolon is automatically returned.
</details>

<details>
<summary>Hint 2: Modulo Operator</summary>

To check if a number is even, use the modulo operator `%`:
- `n % 2 == 0` means n is even
- `n % 2 != 0` (or `n % 2 == 1`) means n is odd
</details>

<details>
<summary>Hint 3: Return Syntax</summary>

In Rust, you can return in two ways:
```rust
// Explicit return (with semicolon on next line)
fn add_explicit(a: i32, b: i32) -> i32 {
    return a + b;
}

// Implicit return (no semicolon, idiomatic)
fn add_implicit(a: i32, b: i32) -> i32 {
    a + b  // No semicolon!
}
```

The implicit return is preferred in Rust.
</details>

## Reference Solution

```rust
fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn rectangle_area(width: f64, height: f64) -> f64 {
    width * height
}

fn is_even(n: i32) -> bool {
    n % 2 == 0
}

fn main() {
    // Test add function
    let sum = add(5, 7);
    println!("5 + 7 = {}", sum);
    assert_eq!(sum, 12);

    let sum2 = add(-3, 3);
    println!("-3 + 3 = {}", sum2);
    assert_eq!(sum2, 0);

    // Test rectangle_area function
    let area = rectangle_area(10.5, 20.0);
    println!("Rectangle area: {}", area);
    assert_eq!(area, 210.0);

    let area2 = rectangle_area(0.0, 10.0);
    println!("Rectangle area (zero width): {}", area2);
    assert_eq!(area2, 0.0);

    // Test is_even function
    println!("Is 4 even? {}", is_even(4));
    assert_eq!(is_even(4), true);

    println!("Is 7 even? {}", is_even(7));
    assert_eq!(is_even(7), false);

    println!("Is 0 even? {}", is_even(0));
    assert_eq!(is_even(0), true);

    println!("All tests passed!");
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add_positive_numbers() {
        assert_eq!(add(5, 7), 12);
        assert_eq!(add(100, 200), 300);
    }

    #[test]
    fn test_add_with_zero() {
        assert_eq!(add(5, 0), 5);
        assert_eq!(add(0, 10), 10);
    }

    #[test]
    fn test_add_negative_numbers() {
        assert_eq!(add(-3, 3), 0);
        assert_eq!(add(-5, -10), -15);
        assert_eq!(add(-20, 10), -10);
    }

    #[test]
    fn test_rectangle_area_positive() {
        assert_eq!(rectangle_area(10.5, 20.0), 210.0);
        assert_eq!(rectangle_area(5.0, 5.0), 25.0);
    }

    #[test]
    fn test_rectangle_area_zero() {
        assert_eq!(rectangle_area(0.0, 10.0), 0.0);
        assert_eq!(rectangle_area(10.0, 0.0), 0.0);
    }

    #[test]
    fn test_rectangle_area_decimal() {
        assert_eq!(rectangle_area(2.5, 4.0), 10.0);
        assert_eq!(rectangle_area(1.1, 2.2), 2.42);
    }

    #[test]
    fn test_is_even_true() {
        assert_eq!(is_even(0), true);
        assert_eq!(is_even(2), true);
        assert_eq!(is_even(4), true);
        assert_eq!(is_even(100), true);
        assert_eq!(is_even(-8), true);
    }

    #[test]
    fn test_is_even_false() {
        assert_eq!(is_even(1), false);
        assert_eq!(is_even(3), false);
        assert_eq!(is_even(7), false);
        assert_eq!(is_even(99), false);
        assert_eq!(is_even(-5), false);
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All 8 test suites pass: 60 points
- 6-7 test suites pass: 48 points
- 4-5 test suites pass: 36 points
- <4 test suites pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Correct function signatures (15 points)**
   - Parameter types match requirements
   - Return types are correct
   - Function names follow snake_case convention

2. **Idiomatic return syntax (10 points)**
   - Uses expression-based returns (no unnecessary `return` keyword)
   - No semicolons on the final expression
   - Clean, readable function bodies

3. **Correctness (10 points)**
   - Functions implement the correct logic
   - Handle edge cases (zero, negative numbers)

4. **Code quality (5 points)**
   - No unnecessary complexity
   - Efficient implementations
   - Good variable naming

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on functions and return values.

Exercise: Implement three functions (add, rectangle_area, is_even).

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Correct function signatures and types
2. Use of idiomatic expression-based returns
3. Correctness and edge case handling
4. Code quality and style

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. If tests pass but code isn't idiomatic (e.g., uses explicit 'return' statements), acknowledge correctness but suggest improvements.",
  "keyTakeaways": ["Understanding of function syntax", "Expression-based returns"],
  "encouragement": "Motivational message"
}

Focus on:
- If student uses `return` keyword unnecessarily, explain expression-based returns
- If function signatures are wrong, explain type annotations
- If logic is incorrect, provide hints without giving away the answer
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Write function signatures with type annotations
- ✅ Understand implicit returns (expression-based)
- ✅ Use basic operators (arithmetic, modulo, comparison)
- ✅ Read and write idiomatic Rust functions
- ✅ Test their code with assertions

## Common Mistakes to Watch For

1. **Missing return types**: Forgetting the `-> Type` syntax
2. **Adding semicolons**: Putting `;` after the final expression turns it into a statement
3. **Wrong types**: Using `i32` instead of `f64` for rectangle_area
4. **Explicit returns everywhere**: Overusing `return` keyword (not idiomatic)
5. **Wrong comparison**: Using `n % 2 == 1` instead of `n % 2 == 0` for is_even

The AI should identify these patterns and provide targeted feedback.

## Extension Ideas (Post-MVP)

For advanced students who finish early:
- Add a `triangle_area` function using Heron's formula
- Implement a `power` function (without using `pow()`)
- Write a `fizz_buzz_number` function that returns "Fizz", "Buzz", "FizzBuzz", or the number as a string

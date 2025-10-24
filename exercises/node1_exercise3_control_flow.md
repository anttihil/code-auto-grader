# Exercise 1.3: Control Flow and Pattern Matching

## Metadata
- **Node**: 1 (Rust Basics)
- **Exercise**: 3 of 3
- **Difficulty**: Medium
- **Estimated Time**: 30-40 minutes
- **Prerequisites**: Exercises 1.1 & 1.2

## Learning Objectives
- Use `if/else` expressions for conditional logic
- Apply `match` for pattern matching
- Understand that control flow constructs are expressions
- Combine multiple conditions and patterns
- Work with ranges and multiple match arms

## Problem Description

Rust's control flow is **expression-based**, meaning `if` and `match` return values. This makes the code more concise and functional in style.

In this exercise, you'll implement a function that categorizes numbers with detailed descriptions. It will combine:
- **`match`** expressions to check the sign (negative, zero, positive)
- **`if/else`** to determine magnitude (small, medium, large)
- **String formatting** to build descriptive output

### Requirements

Implement the `categorize_number` function that:
1. Returns `"zero"` for 0
2. For non-zero numbers, returns a string like:
   - `"negative (small)"` for -5
   - `"positive (medium)"` for 50
   - `"negative (large)"` for -150

### Magnitude Categories
- **Small**: Absolute value 1-10
- **Medium**: Absolute value 11-100
- **Large**: Absolute value > 100

## Starting Code

```rust
fn categorize_number(n: i32) -> String {
    // TODO: Implement this function
    // 1. Use `match` to check if n is negative, zero, or positive
    // 2. For non-zero, use if/else to check magnitude
    // 3. Return a formatted string

    todo!() // Replace this with your implementation
}

fn main() {
    // Test cases
    println!("{}", categorize_number(0));      // "zero"
    println!("{}", categorize_number(5));      // "positive (small)"
    println!("{}", categorize_number(-5));     // "negative (small)"
    println!("{}", categorize_number(50));     // "positive (medium)"
    println!("{}", categorize_number(-75));    // "negative (medium)"
    println!("{}", categorize_number(150));    // "positive (large)"
    println!("{}", categorize_number(-200));   // "negative (large)"

    // Edge cases
    println!("{}", categorize_number(10));     // "positive (small)"
    println!("{}", categorize_number(11));     // "positive (medium)"
    println!("{}", categorize_number(100));    // "positive (medium)"
    println!("{}", categorize_number(101));    // "positive (large)"
}
```

## Expected Output

```
zero
positive (small)
negative (small)
positive (medium)
negative (medium)
positive (large)
negative (large)
positive (small)
positive (medium)
positive (medium)
positive (large)
```

## Hints

<details>
<summary>Hint 1: Match Structure</summary>

Use `match` to handle the three cases:
```rust
match n {
    0 => "zero".to_string(),
    n if n > 0 => {
        // positive case - check magnitude
    },
    _ => {
        // negative case - check magnitude
    }
}
```

Note: `match` arms can use guards (`if` conditions).
</details>

<details>
<summary>Hint 2: Calculating Magnitude</summary>

To check magnitude regardless of sign, use absolute value:
```rust
let abs_value = n.abs();

if abs_value <= 10 {
    "small"
} else if abs_value <= 100 {
    "medium"
} else {
    "large"
}
```
</details>

<details>
<summary>Hint 3: Building the String</summary>

You can use `format!` macro to build strings:
```rust
format!("positive ({})", magnitude)
```

Or use `to_string()` for simple conversions.
</details>

<details>
<summary>Hint 4: Match as an Expression</summary>

Remember, `match` returns a value! You can do:
```rust
let result = match n {
    0 => "zero".to_string(),
    // ...
};
result
```

Or even more concisely:
```rust
match n {
    // ...
}
// No need for explicit return if this is the last expression
```
</details>

## Reference Solution

```rust
fn categorize_number(n: i32) -> String {
    match n {
        0 => "zero".to_string(),
        n if n > 0 => {
            let magnitude = get_magnitude(n);
            format!("positive ({})", magnitude)
        }
        n => {
            // n is negative
            let magnitude = get_magnitude(n.abs());
            format!("negative ({})", magnitude)
        }
    }
}

fn get_magnitude(abs_value: i32) -> &'static str {
    if abs_value <= 10 {
        "small"
    } else if abs_value <= 100 {
        "medium"
    } else {
        "large"
    }
}

fn main() {
    // Test cases
    println!("{}", categorize_number(0));      // "zero"
    println!("{}", categorize_number(5));      // "positive (small)"
    println!("{}", categorize_number(-5));     // "negative (small)"
    println!("{}", categorize_number(50));     // "positive (medium)"
    println!("{}", categorize_number(-75));    // "negative (medium)"
    println!("{}", categorize_number(150));    // "positive (large)"
    println!("{}", categorize_number(-200));   // "negative (large)"

    // Edge cases
    println!("{}", categorize_number(10));     // "positive (small)"
    println!("{}", categorize_number(11));     // "positive (medium)"
    println!("{}", categorize_number(100));    // "positive (medium)"
    println!("{}", categorize_number(101));    // "positive (large)"
}
```

### Alternative Solution (Without Helper Function)

```rust
fn categorize_number(n: i32) -> String {
    match n {
        0 => "zero".to_string(),
        n if n > 0 => {
            let abs_value = n;
            let magnitude = if abs_value <= 10 {
                "small"
            } else if abs_value <= 100 {
                "medium"
            } else {
                "large"
            };
            format!("positive ({})", magnitude)
        }
        n => {
            let abs_value = n.abs();
            let magnitude = if abs_value <= 10 {
                "small"
            } else if abs_value <= 100 {
                "medium"
            } else {
                "large"
            };
            format!("negative ({})", magnitude)
        }
    }
}

fn main() {
    // Same main as above...
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_zero() {
        assert_eq!(categorize_number(0), "zero");
    }

    #[test]
    fn test_positive_small() {
        assert_eq!(categorize_number(1), "positive (small)");
        assert_eq!(categorize_number(5), "positive (small)");
        assert_eq!(categorize_number(10), "positive (small)");
    }

    #[test]
    fn test_positive_medium() {
        assert_eq!(categorize_number(11), "positive (medium)");
        assert_eq!(categorize_number(50), "positive (medium)");
        assert_eq!(categorize_number(100), "positive (medium)");
    }

    #[test]
    fn test_positive_large() {
        assert_eq!(categorize_number(101), "positive (large)");
        assert_eq!(categorize_number(150), "positive (large)");
        assert_eq!(categorize_number(1000), "positive (large)");
    }

    #[test]
    fn test_negative_small() {
        assert_eq!(categorize_number(-1), "negative (small)");
        assert_eq!(categorize_number(-5), "negative (small)");
        assert_eq!(categorize_number(-10), "negative (small)");
    }

    #[test]
    fn test_negative_medium() {
        assert_eq!(categorize_number(-11), "negative (medium)");
        assert_eq!(categorize_number(-75), "negative (medium)");
        assert_eq!(categorize_number(-100), "negative (medium)");
    }

    #[test]
    fn test_negative_large() {
        assert_eq!(categorize_number(-101), "negative (large)");
        assert_eq!(categorize_number(-200), "negative (large)");
        assert_eq!(categorize_number(-5000), "negative (large)");
    }

    #[test]
    fn test_boundary_values() {
        // Test exact boundaries
        assert_eq!(categorize_number(10), "positive (small)");
        assert_eq!(categorize_number(11), "positive (medium)");
        assert_eq!(categorize_number(100), "positive (medium)");
        assert_eq!(categorize_number(101), "positive (large)");

        assert_eq!(categorize_number(-10), "negative (small)");
        assert_eq!(categorize_number(-11), "negative (medium)");
        assert_eq!(categorize_number(-100), "negative (medium)");
        assert_eq!(categorize_number(-101), "negative (large)");
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All 8 test suites pass: 60 points
- 6-7 test suites pass: 48 points
- 4-5 test suites pass: 30 points
- <4 test suites pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Use of `match` expression (15 points)**
   - Uses `match` for sign checking (not just if/else)
   - Handles all cases (zero, positive, negative)
   - Match is used as an expression (returns value)

2. **Control flow correctness (10 points)**
   - Magnitude categorization logic is correct
   - Boundary values handled properly (10, 11, 100, 101)
   - Uses `.abs()` for negative numbers

3. **Code organization (10 points)**
   - Clean, readable structure
   - No code duplication (or justified duplication)
   - Good use of intermediate variables

4. **Rust idioms (5 points)**
   - Expression-based returns
   - Appropriate use of `format!` or string methods
   - Follows naming conventions

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on control flow and pattern matching.

Exercise: Implement categorize_number(n) that returns descriptive strings like "positive (small)", "negative (large)", or "zero".

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Proper use of `match` expression (not just nested if/else)
2. Correctness of magnitude logic and boundary cases
3. Code organization and readability
4. Use of Rust idioms (expressions, no unnecessary returns)

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. If student uses only if/else instead of match, suggest match is more idiomatic. If boundary cases are wrong, explain the issue.",
  "keyTakeaways": ["Match expressions", "Expression-based control flow"],
  "encouragement": "Motivational message"
}

Common issues to check:
- Not using `match` (only if/else) - suggest match for clarity
- Wrong boundary values (e.g., 10 categorized as medium instead of small)
- Not using `.abs()` for negative numbers, leading to incorrect magnitude
- Code duplication between positive and negative cases
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Use `match` expressions for pattern matching
- ✅ Understand match guards (`n if n > 0`)
- ✅ Combine `match` and `if/else` effectively
- ✅ Write expression-based control flow
- ✅ Handle boundary cases correctly
- ✅ Use `.abs()` to get absolute values
- ✅ Format strings with `format!` macro

## Common Mistakes to Watch For

1. **Using only if/else, no match**: While it works, `match` is more idiomatic for this case
2. **Wrong boundaries**:
   - Treating 10 as medium (should be small)
   - Treating 100 as large (should be medium)
3. **Forgetting `.abs()` for negative numbers**: Leads to all negatives being "large"
4. **Adding semicolons**: Making `match` a statement instead of an expression
5. **String ownership issues**: Forgetting `.to_string()` or using `&str` incorrectly

The AI should gently guide students toward more idiomatic solutions while valuing working code.

## Extension Ideas (Post-MVP)

For advanced students:
- Add more categories (tiny: 0-5, huge: >1000)
- Return a struct instead of a string with `sign` and `magnitude` fields
- Handle i64 or accept generic numeric types
- Add a "describe" mode that includes the actual number: "50 is positive (medium)"

## Node 1 Completion

After finishing this exercise, students have completed **Node 1: Rust Basics**!

They should now understand:
- Variables and mutability
- Functions and type annotations
- Control flow with match and if/else
- Basic Rust syntax and idioms

**Next**: Node 2 (Ownership & Borrowing) - the heart of Rust's unique features!

# Exercise 2.2: References and Borrowing

## Metadata
- **Node**: 2 (Ownership & Borrowing)
- **Exercise**: 2 of 3
- **Difficulty**: Medium
- **Estimated Time**: 35-45 minutes
- **Prerequisites**: Exercise 2.1 (Ownership and Moves)

## Learning Objectives
- Understand references (`&T`) as "borrowing" without taking ownership
- Use immutable references to read data without moving it
- Apply mutable references (`&mut T`) to modify data in place
- Recognize that borrowing is more efficient than cloning
- Understand reference lifetimes (basic)

## Problem Description

In the previous exercise, you learned that moving ownership can be inconvenient. If you want to use a value in multiple places, you had to `.clone()` it, which copies data and can be expensive.

**Borrowing** is Rust's solution: you can create a **reference** to a value that lets you read (or modify) it without taking ownership.

### References and Borrowing

- **Immutable reference** (`&T`): Read-only access, no ownership transfer
  - `let len = calculate_length(&s);` - `s` is still valid after!
  - You can have multiple immutable references simultaneously

- **Mutable reference** (`&mut T`): Read-write access, no ownership transfer
  - `modify_string(&mut s);` - `s` is modified in place and still valid after
  - You can have only ONE mutable reference at a time

### Your Task

Implement several functions that use references instead of taking ownership:

1. **`calculate_length`**: Take a reference to a String, return its length
2. **`first_word_length`**: Take a reference to a String, return length of first word
3. **`append_exclamation`**: Take a mutable reference and add "!" to the end
4. **`capitalize_first_char`**: Take a mutable reference and capitalize the first letter

## Starting Code

```rust
// TODO: Implement this function to take a REFERENCE to String
// and return its length without taking ownership
fn calculate_length(/* add parameter */) -> usize {
    todo!()
}

// TODO: Implement this function to take a reference to String
// and return the length of the first word (up to first space or end of string)
fn first_word_length(/* add parameter */) -> usize {
    todo!()
}

// TODO: Implement this function to take a MUTABLE reference to String
// and append "!" to the end
fn append_exclamation(/* add parameter */) {
    todo!()
}

// TODO: Implement this function to take a mutable reference to String
// and capitalize the first character
// Hint: Use .remove(0) to get first char, to_uppercase(), and .insert()
fn capitalize_first_char(/* add parameter */) {
    todo!()
}

fn main() {
    // Test calculate_length
    println!("=== Test calculate_length ===");
    let s1 = String::from("hello world");
    let len = calculate_length(&s1);  // Borrowing s1
    println!("The length of '{}' is {}.", s1, len);
    // s1 is still valid here!

    // Test first_word_length
    println!("\n=== Test first_word_length ===");
    let s2 = String::from("hello world");
    let len = first_word_length(&s2);
    println!("First word length of '{}' is {}", s2, len);

    // Test append_exclamation
    println!("\n=== Test append_exclamation ===");
    let mut s3 = String::from("hello");
    append_exclamation(&mut s3);  // Mutable borrow
    println!("After append: {}", s3);

    // Test capitalize_first_char
    println!("\n=== Test capitalize_first_char ===");
    let mut s4 = String::from("rust");
    capitalize_first_char(&mut s4);
    println!("After capitalize: {}", s4);

    // Demonstrate multiple immutable borrows
    println!("\n=== Multiple Immutable Borrows ===");
    let s5 = String::from("multiple borrows");
    let r1 = &s5;
    let r2 = &s5;
    let r3 = &s5;
    println!("r1: {}, r2: {}, r3: {}", r1, r2, r3);
    // All valid! Multiple immutable borrows are allowed

    println!("\n✅ All tests passed!");
}
```

## Expected Output

```
=== Test calculate_length ===
The length of 'hello world' is 11.

=== Test first_word_length ===
First word length of 'hello world' is 5

=== Test append_exclamation ===
After append: hello!

=== Test capitalize_first_char ===
After capitalize: Rust

=== Multiple Immutable Borrows ===
r1: multiple borrows, r2: multiple borrows, r3: multiple borrows

✅ All tests passed!
```

## Hints

<details>
<summary>Hint 1: Reference Syntax</summary>

**Immutable reference**:
```rust
fn function_name(parameter: &String) -> ReturnType {
    // Can read parameter but not modify it
    parameter.len()  // Example: reading length
}
```

**Mutable reference**:
```rust
fn function_name(parameter: &mut String) {
    // Can modify parameter
    parameter.push_str("more");  // Example: appending
}
```

When calling:
- Immutable: `function_name(&my_string)`
- Mutable: `function_name(&mut my_string)`
</details>

<details>
<summary>Hint 2: calculate_length</summary>

```rust
fn calculate_length(s: &String) -> usize {
    s.len()  // Just return the length!
}
```

The `&String` means "borrow a String" - you can read it but not take ownership.
</details>

<details>
<summary>Hint 3: first_word_length</summary>

You need to:
1. Iterate through characters
2. Count until you hit a space or reach the end

One approach:
```rust
fn first_word_length(s: &String) -> usize {
    for (i, ch) in s.chars().enumerate() {
        if ch == ' ' {
            return i;
        }
    }
    s.len()  // No space found, return full length
}
```
</details>

<details>
<summary>Hint 4: append_exclamation</summary>

```rust
fn append_exclamation(s: &mut String) {
    s.push_str("!");
}
```

The `&mut` means you can modify the String in place.
</details>

<details>
<summary>Hint 5: capitalize_first_char</summary>

This is tricky! You need to:
1. Remove the first character
2. Convert it to uppercase
3. Insert it back at the beginning

```rust
fn capitalize_first_char(s: &mut String) {
    if !s.is_empty() {
        let first = s.remove(0);  // Remove and get first char
        let upper = first.to_uppercase().to_string();
        s.insert_str(0, &upper);  // Insert at beginning
    }
}
```
</details>

<details>
<summary>Hint 6: Understanding Borrowing Benefits</summary>

Compare these two approaches:

**Without borrowing (cloning)**:
```rust
let s = String::from("hello");
let len = calculate_length_clone(s.clone());  // Expensive copy!
// s is still valid but we copied the data
```

**With borrowing**:
```rust
let s = String::from("hello");
let len = calculate_length(&s);  // No copy! Just a reference
// s is still valid and no extra memory used
```

Borrowing is both safer AND more efficient!
</details>

## Reference Solution

```rust
fn calculate_length(s: &String) -> usize {
    s.len()
}

fn first_word_length(s: &String) -> usize {
    for (i, ch) in s.chars().enumerate() {
        if ch == ' ' {
            return i;
        }
    }
    s.len()  // No space found
}

fn append_exclamation(s: &mut String) {
    s.push_str("!");
}

fn capitalize_first_char(s: &mut String) {
    if !s.is_empty() {
        let first = s.remove(0);
        let upper = first.to_uppercase().to_string();
        s.insert_str(0, &upper);
    }
}

fn main() {
    // Test calculate_length
    println!("=== Test calculate_length ===");
    let s1 = String::from("hello world");
    let len = calculate_length(&s1);
    println!("The length of '{}' is {}.", s1, len);

    // Test first_word_length
    println!("\n=== Test first_word_length ===");
    let s2 = String::from("hello world");
    let len = first_word_length(&s2);
    println!("First word length of '{}' is {}", s2, len);

    // Test append_exclamation
    println!("\n=== Test append_exclamation ===");
    let mut s3 = String::from("hello");
    append_exclamation(&mut s3);
    println!("After append: {}", s3);

    // Test capitalize_first_char
    println!("\n=== Test capitalize_first_char ===");
    let mut s4 = String::from("rust");
    capitalize_first_char(&mut s4);
    println!("After capitalize: {}", s4);

    // Demonstrate multiple immutable borrows
    println!("\n=== Multiple Immutable Borrows ===");
    let s5 = String::from("multiple borrows");
    let r1 = &s5;
    let r2 = &s5;
    let r3 = &s5;
    println!("r1: {}, r2: {}, r3: {}", r1, r2, r3);

    println!("\n✅ All tests passed!");
}
```

### Alternative Solution (Using &str)

```rust
// More idiomatic: accept &str instead of &String
// &str works with both String and string literals

fn calculate_length_better(s: &str) -> usize {
    s.len()
}

fn first_word_length_better(s: &str) -> usize {
    s.split_whitespace()
        .next()
        .map(|word| word.len())
        .unwrap_or(0)
}

// Note: Mutable functions still need &mut String since &str is immutable
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calculate_length() {
        let s = String::from("hello");
        assert_eq!(calculate_length(&s), 5);

        let s2 = String::from("");
        assert_eq!(calculate_length(&s2), 0);

        let s3 = String::from("hello world!");
        assert_eq!(calculate_length(&s3), 12);
    }

    #[test]
    fn test_original_string_still_valid_after_borrow() {
        let s = String::from("test");
        let len = calculate_length(&s);
        assert_eq!(len, 4);
        assert_eq!(s, "test");  // s is still valid!
    }

    #[test]
    fn test_first_word_length() {
        assert_eq!(first_word_length(&String::from("hello world")), 5);
        assert_eq!(first_word_length(&String::from("hello")), 5);
        assert_eq!(first_word_length(&String::from("")), 0);
        assert_eq!(first_word_length(&String::from("a b c")), 1);
        assert_eq!(first_word_length(&String::from("   ")), 0);
    }

    #[test]
    fn test_append_exclamation() {
        let mut s = String::from("hello");
        append_exclamation(&mut s);
        assert_eq!(s, "hello!");

        append_exclamation(&mut s);
        assert_eq!(s, "hello!!");
    }

    #[test]
    fn test_append_empty_string() {
        let mut s = String::from("");
        append_exclamation(&mut s);
        assert_eq!(s, "!");
    }

    #[test]
    fn test_capitalize_first_char() {
        let mut s = String::from("hello");
        capitalize_first_char(&mut s);
        assert_eq!(s, "Hello");

        let mut s2 = String::from("rust");
        capitalize_first_char(&mut s2);
        assert_eq!(s2, "Rust");
    }

    #[test]
    fn test_capitalize_already_capital() {
        let mut s = String::from("Hello");
        capitalize_first_char(&mut s);
        assert_eq!(s, "Hello");  // Should stay the same
    }

    #[test]
    fn test_capitalize_empty_string() {
        let mut s = String::from("");
        capitalize_first_char(&mut s);
        assert_eq!(s, "");  // Should handle gracefully
    }

    #[test]
    fn test_capitalize_single_char() {
        let mut s = String::from("a");
        capitalize_first_char(&mut s);
        assert_eq!(s, "A");
    }

    #[test]
    fn test_multiple_immutable_borrows() {
        let s = String::from("test");
        let r1 = &s;
        let r2 = &s;
        let r3 = &s;
        assert_eq!(r1, "test");
        assert_eq!(r2, "test");
        assert_eq!(r3, "test");
        assert_eq!(s, "test");  // Original still valid
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All 10 test suites pass: 60 points
- 8-9 tests pass: 50 points
- 6-7 tests pass: 40 points
- 4-5 tests pass: 25 points
- <4 tests pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Correct use of references (20 points)**
   - Uses `&String` for immutable borrows
   - Uses `&mut String` for mutable borrows
   - Doesn't take ownership when not needed
   - Demonstrates understanding of when to use each

2. **Implementation correctness (10 points)**
   - Functions work correctly for all edge cases
   - Handles empty strings gracefully
   - Efficient implementations (no unnecessary operations)

3. **Code quality (10 points)**
   - Clean, readable code
   - Good variable names
   - Appropriate error handling (if any)
   - Follows Rust conventions

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on references and borrowing.

Exercise: Implement functions using immutable and mutable references to avoid taking ownership.

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Correct use of & and &mut in function signatures
2. Understanding of when to use immutable vs mutable references
3. Implementation correctness and edge case handling
4. Code quality and idioms

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. Focus on:\n- If student takes ownership instead of borrowing: explain the benefit of references\n- If student uses wrong reference type: explain immutable vs mutable\n- If edge cases fail: guide toward handling empty strings, etc.\n- If using advanced patterns like &str: acknowledge and praise\n- Compare efficiency to cloning from previous exercise",
  "keyTakeaways": ["References don't take ownership", "Immutable vs mutable borrows", "Borrowing is efficient"],
  "encouragement": "Building on ownership knowledge!"
}

Common patterns:
- Takes ownership (fn calculate_length(s: String)) → Should borrow instead
- Uses &mut when & would work → Explain principle of least privilege
- Doesn't handle empty strings → Point to edge case testing
- Correctly implements everything → Praise understanding, preview borrowing rules next
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Use `&T` references to borrow data without taking ownership
- ✅ Use `&mut T` to modify data in place
- ✅ Understand that borrowing is more efficient than cloning
- ✅ Create multiple immutable references to the same data
- ✅ Call functions with `&` and `&mut` syntax
- ✅ Recognize that references don't live beyond their scope
- ✅ Feel prepared for borrowing rules (next exercise)

## Common Mistakes to Watch For

1. **Taking ownership instead of borrowing**: `fn calculate_length(s: String)` instead of `fn calculate_length(s: &String)`
   - Guide: Explain that borrowing lets the caller keep using the value

2. **Using `&mut` when `&` would work**: Over-privileged access
   - Guide: Explain principle of least privilege - use immutable unless you need to modify

3. **Forgetting `&` or `&mut` in function calls**: `calculate_length(s)` instead of `calculate_length(&s)`
   - Guide: Show the compiler error and explain reference syntax

4. **Not handling empty strings**: `capitalize_first_char` crashes on empty string
   - Guide: Point to the need for bounds checking

5. **Confusion about dereferencing**: Trying to use `*` explicitly when not needed
   - Guide: Explain Rust's auto-deref feature for references

## Connection to Previous Exercise

Exercise 2.1 taught that moving ownership can be inconvenient (need to clone). This exercise shows the solution: **borrowing lets you use values without taking ownership**.

Key insight:
```rust
// Exercise 2.1 approach: Clone (works but inefficient)
let s = String::from("hello");
print_message(s.clone());
println!("{}", s);  // s still valid

// Exercise 2.2 approach: Borrow (efficient!)
let s = String::from("hello");
print_message_ref(&s);
println!("{}", s);  // s still valid, no clone needed!
```

## Preview of Next Exercise

This exercise allows multiple immutable borrows, but what about mutable borrows? The next exercise explores **the borrowing rules**: why you can't have mutable and immutable references simultaneously.

Teaser:
```rust
let mut s = String::from("hello");
let r1 = &s;       // OK
let r2 = &mut s;   // ERROR! Can't borrow as mutable while immutable ref exists
```

Why? We'll find out in Exercise 2.3!

## Extension Ideas (Post-MVP)

- Add `reverse_in_place` function (more complex mutable borrow)
- Implement `swap_words` that swaps first and second word
- Challenge: Implement functions that return references (introduces explicit lifetimes)
- Comparison tool showing memory usage: borrowing vs cloning

## Notes for Implementation

- This exercise should compile successfully once completed (unlike 2.1 which starts with errors)
- The `capitalize_first_char` function is intentionally the hardest (involves remove + insert)
- Consider showing memory diagrams in hints (stack vs heap, pointers)
- The alternative solution using `&str` is more idiomatic but more advanced (good for AI to mention)

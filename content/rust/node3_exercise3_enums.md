# Exercise 3.3: Enums and Pattern Matching

## Metadata

- **Node**: 3 (Structs & Enums)
- **Exercise**: 3 of 3
- **Difficulty**: Medium
- **Estimated Time**: 40-50 minutes
- **Prerequisites**: Exercises 3.1 & 3.2

## Learning Objectives

- Define enums with multiple variants
- Create enums with data (tuple and struct variants)
- Use `match` expressions to handle enum variants
- Apply pattern matching with destructuring
- Work with `Option<T>` (Rust's built-in enum)
- Understand exhaustive matching requirements

## Problem Description

**Enums** (enumerations) let you define a type that can be one of several variants. Unlike structs (which hold multiple values at once), enums hold one of several possible values.

### Enum Basics

```rust
enum IpAddrKind {
    V4,
    V6,
}

let ipv4 = IpAddrKind::V4;
let ipv6 = IpAddrKind::V6;
```

### Enums with Data

Enums can hold data in each variant:

```rust
enum IpAddr {
    V4(u8, u8, u8, u8),         // Tuple variant
    V6(String),                  // Holds a String
}

let home = IpAddr::V4(127, 0, 0, 1);
let loopback = IpAddr::V6(String::from("::1"));
```

### Pattern Matching with match

Use `match` to handle different variants:

```rust
match ip_addr {
    IpAddr::V4(a, b, c, d) => {
        println!("IPv4: {}.{}.{}.{}", a, b, c, d)
    }
    IpAddr::V6(addr) => {
        println!("IPv6: {}", addr)
    }
}
```

### Your Task

Implement a message processing system with an enum representing different message types:

1. **Define `Message` enum** with variants:

   - `Quit` (no data)
   - `Move { x: i32, y: i32 }` (struct variant)
   - `Write(String)` (tuple variant)
   - `ChangeColor(u8, u8, u8)` (RGB color)

2. **Implement functions**:

   - `process_message` - returns a description string
   - `is_move_message` - checks if it's a Move variant
   - `extract_text` - returns `Option<String>` with text from Write messages

3. **Work with `Option<T>`**:
   - Use Rust's built-in `Option` enum
   - Handle Some and None cases

## Starting Code

```rust
// TODO: Define the Message enum with these variants:
// - Quit (no data)
// - Move { x: i32, y: i32 }
// - Write(String)
// - ChangeColor(u8, u8, u8)

impl Message {
    // TODO: Implement this function to return a descriptive string
    // based on the message variant
    fn process_message(&self) -> String {
        todo!()
    }

    // TODO: Implement this function to check if the message is Move variant
    fn is_move_message(&self) -> bool {
        todo!()
    }

    // TODO: Implement this function to extract text from Write variant
    // Returns Some(text) if Write, None otherwise
    fn extract_text(&self) -> Option<String> {
        todo!()
    }
}

// TODO: Implement this function to calculate area of different shapes
// Use match to handle each variant
fn calculate_area(shape: &Shape) -> f64 {
    todo!()
}

// TODO: Define a Shape enum with these variants:
// - Circle { radius: f64 }
// - Rectangle { width: f64, height: f64 }
// - Triangle { base: f64, height: f64 }

fn main() {
    println!("=== Testing Message Enum ===");

    // Test Quit variant
    let msg1 = Message::Quit;
    assert_eq!(msg1.process_message(), "Quit message");
    assert_eq!(msg1.is_move_message(), false);
    assert_eq!(msg1.extract_text(), None);
    println!("✓ Quit message works");

    // Test Move variant
    let msg2 = Message::Move { x: 10, y: 20 };
    assert_eq!(msg2.process_message(), "Move to (10, 20)");
    assert_eq!(msg2.is_move_message(), true);
    assert_eq!(msg2.extract_text(), None);
    println!("✓ Move message works");

    // Test Write variant
    let msg3 = Message::Write(String::from("Hello, Rust!"));
    assert_eq!(msg3.process_message(), "Write: Hello, Rust!");
    assert_eq!(msg3.is_move_message(), false);
    assert_eq!(msg3.extract_text(), Some(String::from("Hello, Rust!")));
    println!("✓ Write message works");

    // Test ChangeColor variant
    let msg4 = Message::ChangeColor(255, 0, 128);
    assert_eq!(msg4.process_message(), "Change color to RGB(255, 0, 128)");
    assert_eq!(msg4.is_move_message(), false);
    assert_eq!(msg4.extract_text(), None);
    println!("✓ ChangeColor message works");

    println!("\n=== Testing Shape Enum ===");

    // Test shapes
    let circle = Shape::Circle { radius: 5.0 };
    let rect = Shape::Rectangle { width: 10.0, height: 20.0 };
    let triangle = Shape::Triangle { base: 6.0, height: 8.0 };

    assert_eq!(calculate_area(&circle), 78.5);  // π * 5^2 ≈ 78.5
    assert_eq!(calculate_area(&rect), 200.0);   // 10 * 20
    assert_eq!(calculate_area(&triangle), 24.0); // (6 * 8) / 2

    println!("✓ Shape area calculations work");

    println!("\n✅ All tests passed!");
}
```

## Expected Output

```
=== Testing Message Enum ===
✓ Quit message works
✓ Move message works
✓ Write message works
✓ ChangeColor message works

=== Testing Shape Enum ===
✓ Shape area calculations work

✅ All tests passed!
```

## Hints

<details>
<summary>Hint 1: Defining Enums</summary>

```rust
enum Message {
    Quit,                             // Unit variant (no data)
    Move { x: i32, y: i32 },          // Struct variant
    Write(String),                    // Tuple variant
    ChangeColor(u8, u8, u8),          // Tuple variant with 3 values
}
```

Each variant can hold different types and amounts of data!

</details>

<details>
<summary>Hint 2: Pattern Matching with match</summary>

Use `match` to handle each variant:

```rust
fn process_message(&self) -> String {
    match self {
        Message::Quit => String::from("Quit message"),
        Message::Move { x, y } => format!("Move to ({}, {})", x, y),
        Message::Write(text) => format!("Write: {}", text),
        Message::ChangeColor(r, g, b) => format!("Change color to RGB({}, {}, {})", r, g, b),
    }
}
```

Note: `match` must be **exhaustive** - handle all variants!

</details>

<details>
<summary>Hint 3: Checking for Specific Variant</summary>

Use `matches!` macro or match with boolean:

```rust
fn is_move_message(&self) -> bool {
    matches!(self, Message::Move { .. })
}
```

Or:

```rust
fn is_move_message(&self) -> bool {
    match self {
        Message::Move { .. } => true,
        _ => false,
    }
}
```

The `..` ignores the fields inside the variant.

</details>

<details>
<summary>Hint 4: Extracting Data with Option</summary>

Return `Some` if the variant matches, `None` otherwise:

```rust
fn extract_text(&self) -> Option<String> {
    match self {
        Message::Write(text) => Some(text.clone()),
        _ => None,
    }
}
```

`_` is a catch-all pattern for "all other cases".

</details>

<details>
<summary>Hint 5: Shape Enum and Area Calculation</summary>

```rust
enum Shape {
    Circle { radius: f64 },
    Rectangle { width: f64, height: f64 },
    Triangle { base: f64, height: f64 },
}

fn calculate_area(shape: &Shape) -> f64 {
    match shape {
        Shape::Circle { radius } => 3.14159 * radius * radius,
        Shape::Rectangle { width, height } => width * height,
        Shape::Triangle { base, height } => (base * height) / 2.0,
    }
}
```

Pattern matching automatically destructures the variant fields!

</details>

<details>
<summary>Hint 6: Understanding Option<T></summary>

`Option<T>` is Rust's built-in enum for nullable values:

```rust
enum Option<T> {
    Some(T),
    None,
}
```

It forces you to handle the "no value" case explicitly, preventing null pointer errors!

</details>

## Reference Solution

```rust
enum Message {
    Quit,
    Move { x: i32, y: i32 },
    Write(String),
    ChangeColor(u8, u8, u8),
}

impl Message {
    fn process_message(&self) -> String {
        match self {
            Message::Quit => String::from("Quit message"),
            Message::Move { x, y } => format!("Move to ({}, {})", x, y),
            Message::Write(text) => format!("Write: {}", text),
            Message::ChangeColor(r, g, b) => {
                format!("Change color to RGB({}, {}, {})", r, g, b)
            }
        }
    }

    fn is_move_message(&self) -> bool {
        matches!(self, Message::Move { .. })
    }

    fn extract_text(&self) -> Option<String> {
        match self {
            Message::Write(text) => Some(text.clone()),
            _ => None,
        }
    }
}

enum Shape {
    Circle { radius: f64 },
    Rectangle { width: f64, height: f64 },
    Triangle { base: f64, height: f64 },
}

fn calculate_area(shape: &Shape) -> f64 {
    match shape {
        Shape::Circle { radius } => 3.14159 * radius * radius,
        Shape::Rectangle { width, height } => width * height,
        Shape::Triangle { base, height } => (base * height) / 2.0,
    }
}

fn main() {
    println!("=== Testing Message Enum ===");

    let msg1 = Message::Quit;
    assert_eq!(msg1.process_message(), "Quit message");
    assert_eq!(msg1.is_move_message(), false);
    assert_eq!(msg1.extract_text(), None);
    println!("✓ Quit message works");

    let msg2 = Message::Move { x: 10, y: 20 };
    assert_eq!(msg2.process_message(), "Move to (10, 20)");
    assert_eq!(msg2.is_move_message(), true);
    assert_eq!(msg2.extract_text(), None);
    println!("✓ Move message works");

    let msg3 = Message::Write(String::from("Hello, Rust!"));
    assert_eq!(msg3.process_message(), "Write: Hello, Rust!");
    assert_eq!(msg3.is_move_message(), false);
    assert_eq!(msg3.extract_text(), Some(String::from("Hello, Rust!")));
    println!("✓ Write message works");

    let msg4 = Message::ChangeColor(255, 0, 128);
    assert_eq!(msg4.process_message(), "Change color to RGB(255, 0, 128)");
    assert_eq!(msg4.is_move_message(), false);
    assert_eq!(msg4.extract_text(), None);
    println!("✓ ChangeColor message works");

    println!("\n=== Testing Shape Enum ===");

    let circle = Shape::Circle { radius: 5.0 };
    let rect = Shape::Rectangle { width: 10.0, height: 20.0 };
    let triangle = Shape::Triangle { base: 6.0, height: 8.0 };

    assert_eq!(calculate_area(&circle), 78.5);
    assert_eq!(calculate_area(&rect), 200.0);
    assert_eq!(calculate_area(&triangle), 24.0);

    println!("✓ Shape area calculations work");

    println!("\n✅ All tests passed!");
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_quit_message() {
        let msg = Message::Quit;
        assert_eq!(msg.process_message(), "Quit message");
        assert_eq!(msg.is_move_message(), false);
        assert_eq!(msg.extract_text(), None);
    }

    #[test]
    fn test_move_message() {
        let msg = Message::Move { x: 10, y: 20 };
        assert_eq!(msg.process_message(), "Move to (10, 20)");
        assert_eq!(msg.is_move_message(), true);
        assert_eq!(msg.extract_text(), None);
    }

    #[test]
    fn test_move_negative_coordinates() {
        let msg = Message::Move { x: -5, y: -10 };
        assert_eq!(msg.process_message(), "Move to (-5, -10)");
        assert_eq!(msg.is_move_message(), true);
    }

    #[test]
    fn test_write_message() {
        let msg = Message::Write(String::from("Hello, Rust!"));
        assert_eq!(msg.process_message(), "Write: Hello, Rust!");
        assert_eq!(msg.is_move_message(), false);
        assert_eq!(msg.extract_text(), Some(String::from("Hello, Rust!")));
    }

    #[test]
    fn test_write_empty_string() {
        let msg = Message::Write(String::from(""));
        assert_eq!(msg.process_message(), "Write: ");
        assert_eq!(msg.extract_text(), Some(String::from("")));
    }

    #[test]
    fn test_change_color() {
        let msg = Message::ChangeColor(255, 0, 128);
        assert_eq!(msg.process_message(), "Change color to RGB(255, 0, 128)");
        assert_eq!(msg.is_move_message(), false);
        assert_eq!(msg.extract_text(), None);
    }

    #[test]
    fn test_circle_area() {
        let circle = Shape::Circle { radius: 5.0 };
        assert_eq!(calculate_area(&circle), 78.5);
    }

    #[test]
    fn test_rectangle_area() {
        let rect = Shape::Rectangle {
            width: 10.0,
            height: 20.0,
        };
        assert_eq!(calculate_area(&rect), 200.0);
    }

    #[test]
    fn test_triangle_area() {
        let triangle = Shape::Triangle {
            base: 6.0,
            height: 8.0,
        };
        assert_eq!(calculate_area(&triangle), 24.0);
    }

    #[test]
    fn test_zero_radius_circle() {
        let circle = Shape::Circle { radius: 0.0 };
        assert_eq!(calculate_area(&circle), 0.0);
    }

    #[test]
    fn test_square_as_rectangle() {
        let square = Shape::Rectangle {
            width: 5.0,
            height: 5.0,
        };
        assert_eq!(calculate_area(&square), 25.0);
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)

- All 11 tests pass: 60 points
- 8-10 tests pass: 48 points
- 5-7 tests pass: 30 points
- <5 tests pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Correct enum definitions (15 points)**

   - All variants defined with correct data types
   - Follows naming conventions (PascalCase for variants)

2. **Pattern matching (15 points)**

   - Exhaustive matching (handles all variants)
   - Correct destructuring of variant data
   - Appropriate use of `_` for catch-all

3. **Option handling (5 points)**

   - Returns Some/None appropriately
   - Understands Option semantics

4. **Code quality (5 points)**
   - Clean match expressions
   - Uses `format!` or string interpolation
   - Idiomatic Rust patterns

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on enums and pattern matching.

Exercise: Define Message and Shape enums with methods using pattern matching.

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Enum definitions (correct variants with appropriate data)
2. Pattern matching (exhaustive, correct destructuring)
3. Option<T> usage (understanding Some/None)
4. Code quality and idioms

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. Focus on:\n- If enum variants wrong: explain syntax for different variant types\n- If match not exhaustive: explain Rust requires handling all cases\n- If destructuring wrong: show correct syntax for extracting fields\n- If Option confused: explain Some/None semantics\n- If using matches! macro: praise idiomatic approach\n- If verbose match for boolean: suggest matches! as alternative",
  "keyTakeaways": ["Enums can hold different data per variant", "Match must be exhaustive", "Option prevents null errors"],
  "encouragement": "Pattern matching is one of Rust's superpowers!"
}

Common patterns:
- Non-exhaustive match → Explain compiler enforces handling all cases
- Wrong destructuring syntax → Show { field } for struct variants, (value) for tuple
- Not using Option correctly → Explain Some/None pattern
- Overly verbose is_move_message → Suggest matches! macro
```

## Expected Student Outcomes

After completing this exercise, students should:

- ✅ Define enums with various variant types
- ✅ Use pattern matching to handle enum variants
- ✅ Destructure data from enum variants
- ✅ Work with `Option<T>` for nullable values
- ✅ Understand exhaustive matching requirements
- ✅ Apply enums to model real-world scenarios

## Common Mistakes to Watch For

1. **Non-exhaustive match**: Forgetting to handle all variants

   - Guide: Rust's compiler enforces exhaustiveness - this prevents bugs!

2. **Wrong destructuring syntax**:

   - `Message::Move(x, y)` instead of `Message::Move { x, y }`
   - Guide: Struct variants use `{}`, tuple variants use `()`

3. **Not cloning in extract_text**: Trying to move String out of borrowed &self

   - Guide: Need `.clone()` to return owned String from borrowed reference

4. **Using if/else instead of match**: Multiple if statements instead of match

   - Guide: Match is more idiomatic and compiler-checked for completeness

5. **Confused about Option**: Not understanding Some/None pattern
   - Guide: Option is just an enum! `Some(value)` or `None`

## Connection to Previous Exercises

**Node 3** has shown two ways to group data:

- **Structs** (Exercise 3.1-3.2): Hold multiple values simultaneously
- **Enums** (Exercise 3.3): Hold ONE of several possible values

Both use methods, both integrate with ownership:

```rust
// Struct: always has all fields
struct Rectangle { width: u32, height: u32 }

// Enum: ONE variant at a time
enum Shape {
    Rectangle { width: u32, height: u32 },
    Circle { radius: u32 },
}
```

## Why Enums Matter

Enums solve a common problem: "This value could be one of several different things."

Without enums, you'd need:

- Multiple structs
- Lots of `Option<T>` fields (most of which are None)
- Runtime checks for which type it is

With enums:

- One type that clearly represents all possibilities
- Compile-time checking that you handle all cases
- Pattern matching makes it elegant

## Real-World Rust Enums

Rust's standard library uses enums extensively:

```rust
// Result<T, E> - for operations that might fail
enum Result<T, E> {
    Ok(T),
    Err(E),
}

// Option<T> - for nullable values
enum Option<T> {
    Some(T),
    None,
}
```

You'll use these constantly in Rust!

## Extension Ideas (Post-MVP)

- Add methods to Shape enum (area, perimeter)
- Implement `From` trait to convert between types
- Add Result<T, E> exercise for error handling
- Create a more complex enum (e.g., JSON value representation)

## Node 3 Completion

**Congratulations!** After finishing this exercise, students have completed Node 3: Structs & Enums.

They now understand:

- Custom structs and data modeling
- Instance methods and associated functions
- Enums and pattern matching
- Working with Option<T>

Combined with Nodes 1-2, students now have:

- ✅ Rust fundamentals (syntax, control flow)
- ✅ Ownership and borrowing (memory safety)
- ✅ Custom types (structs and enums)

**This is a solid foundation for learning Rust!** These 9 exercises (3 nodes) cover the core concepts needed to write real Rust programs.

**Next**: Node 4 (Collections & Strings), Node 5 (Error Handling), or Node 6 (Traits & Generics) - each builds on this foundation but is conceptually easier than Node 2's ownership challenges.

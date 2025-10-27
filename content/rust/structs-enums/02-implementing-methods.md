{
  "id": "30000000-0000-0000-0000-000000000008",
  "node_id": "20000000-0000-0000-0000-000000000003",
  "title": "Exercise 3.2: Methods and Self",
  "slug": "implementing-methods",
  "difficulty": "medium",
  "estimated_minutes": 30,
  "order_index": 2
}
---

# Exercise 3.2: Methods and Self

## Learning Objectives
- Implement instance methods that operate on struct data
- Understand `&self`, `&mut self`, and `self` in method signatures
- Use method chaining patterns
- Apply the builder pattern (basic)
- Distinguish between associated functions and instance methods

## Problem Description

In the previous exercise, you created structs and associated functions (constructors). Now you'll add **instance methods** - functions that operate on a specific instance of a struct.

### Method Syntax

Methods are defined in `impl` blocks and take some form of `self`:

```rust
impl Rectangle {
    // Borrows self (read-only access)
    fn area(&self) -> u32 {
        self.width * self.height
    }

    // Borrows self mutably (can modify)
    fn resize(&mut self, new_width: u32) {
        self.width = new_width;
    }

    // Takes ownership of self (consumes it)
    fn into_square(self) -> Square {
        Square { size: self.width }
    }
}
```

### When to Use Each?

- **`&self`**: Most common - read data without modifying
- **`&mut self`**: When you need to modify the struct
- **`self`**: Rare - when the method consumes the struct (transforms it)

### Your Task

Implement a `Rectangle` struct with various methods:

1. **Constructor**: `new(width, height)`
2. **Read methods**: `area()`, `perimeter()`, `is_square()`
3. **Comparison**: `can_hold(&other_rect)` - checks if this can contain another
4. **Mutation**: `scale(factor)` - multiply width and height
5. **Builder**: `with_width(w)` and `with_height(h)` for method chaining

## Starting Code

```rust
// TODO: Define a Rectangle struct with:
// - width: u32
// - height: u32

impl Rectangle {
    // TODO: Implement constructor
    // fn new(width: u32, height: u32) -> Rectangle

    // TODO: Implement area calculation
    // Returns width * height
    // fn area(&self) -> u32

    // TODO: Implement perimeter calculation
    // Returns 2 * (width + height)
    // fn perimeter(&self) -> u32

    // TODO: Implement is_square check
    // Returns true if width == height
    // fn is_square(&self) -> bool

    // TODO: Implement can_hold
    // Returns true if self can contain other (both dimensions larger)
    // fn can_hold(&self, other: &Rectangle) -> bool

    // TODO: Implement scale (mutable method)
    // Multiplies both width and height by factor
    // fn scale(&mut self, factor: u32)

    // TODO: Implement with_width (builder pattern)
    // Returns self with modified width (for method chaining)
    // fn with_width(mut self, width: u32) -> Self

    // TODO: Implement with_height (builder pattern)
    // Returns self with modified height
    // fn with_height(mut self, height: u32) -> Self
}

fn main() {
    // Test basic creation and methods
    let rect1 = Rectangle::new(30, 50);
    assert_eq!(rect1.area(), 1500);
    assert_eq!(rect1.perimeter(), 160);
    assert_eq!(rect1.is_square(), false);
    println!("✓ Basic methods work");

    // Test can_hold
    let rect2 = Rectangle::new(10, 40);
    let rect3 = Rectangle::new(60, 45);

    assert_eq!(rect1.can_hold(&rect2), true);
    assert_eq!(rect1.can_hold(&rect3), false);
    println!("✓ can_hold works");

    // Test is_square
    let square = Rectangle::new(25, 25);
    assert_eq!(square.is_square(), true);
    println!("✓ is_square works");

    // Test scale (mutable method)
    let mut rect4 = Rectangle::new(10, 20);
    rect4.scale(3);
    assert_eq!(rect4.width, 30);
    assert_eq!(rect4.height, 60);
    assert_eq!(rect4.area(), 1800);
    println!("✓ scale works");

    // Test builder pattern (method chaining)
    let rect5 = Rectangle::new(10, 10)
        .with_width(50)
        .with_height(100);
    assert_eq!(rect5.width, 50);
    assert_eq!(rect5.height, 100);
    assert_eq!(rect5.area(), 5000);
    println!("✓ Builder pattern works");

    println!("\n✅ All tests passed!");
}
```

## Expected Output

```
✓ Basic methods work
✓ can_hold works
✓ is_square works
✓ scale works
✓ Builder pattern works

✅ All tests passed!
```

## Hints

<details>
<summary>Hint 1: Struct Definition</summary>

```rust
struct Rectangle {
    width: u32,
    height: u32,
}
```
</details>

<details>
<summary>Hint 2: Read-Only Methods (&self)</summary>

Methods that only read data use `&self`:

```rust
impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }

    fn perimeter(&self) -> u32 {
        2 * (self.width + self.height)
    }

    fn is_square(&self) -> bool {
        self.width == self.height
    }
}
```
</details>

<details>
<summary>Hint 3: can_hold Method</summary>

To hold another rectangle, both dimensions must be larger:

```rust
fn can_hold(&self, other: &Rectangle) -> bool {
    self.width > other.width && self.height > other.height
}
```

Note: We take `&Rectangle` because we just need to read `other`, not own it.
</details>

<details>
<summary>Hint 4: Mutable Method (&mut self)</summary>

The `scale` method modifies the rectangle, so it needs `&mut self`:

```rust
fn scale(&mut self, factor: u32) {
    self.width *= factor;
    self.height *= factor;
}
```

Call it on a `mut` variable: `rect.scale(2);`
</details>

<details>
<summary>Hint 5: Builder Pattern (self)</summary>

Builder methods take ownership of `self` and return the modified value:

```rust
fn with_width(mut self, width: u32) -> Self {
    self.width = width;
    self
}

fn with_height(mut self, height: u32) -> Self {
    self.height = height;
    self
}
```

This allows chaining: `Rectangle::new(10, 10).with_width(20).with_height(30)`

Key points:
- Takes `mut self` (ownership, but can modify)
- Returns `Self` (the struct type)
- Enables method chaining
</details>

<details>
<summary>Hint 6: Understanding Self Types</summary>

| Signature | Meaning | When to use |
|-----------|---------|-------------|
| `&self` | Borrow (read-only) | Most methods - just reading data |
| `&mut self` | Mutable borrow | Modifying the struct in place |
| `self` | Take ownership | Transforming/consuming the struct |

99% of methods use `&self` or `&mut self`!
</details>

## Reference Solution

```rust
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn new(width: u32, height: u32) -> Rectangle {
        Rectangle { width, height }
    }

    fn area(&self) -> u32 {
        self.width * self.height
    }

    fn perimeter(&self) -> u32 {
        2 * (self.width + self.height)
    }

    fn is_square(&self) -> bool {
        self.width == self.height
    }

    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }

    fn scale(&mut self, factor: u32) {
        self.width *= factor;
        self.height *= factor;
    }

    fn with_width(mut self, width: u32) -> Self {
        self.width = width;
        self
    }

    fn with_height(mut self, height: u32) -> Self {
        self.height = height;
        self
    }
}

fn main() {
    // Test basic creation and methods
    let rect1 = Rectangle::new(30, 50);
    assert_eq!(rect1.area(), 1500);
    assert_eq!(rect1.perimeter(), 160);
    assert_eq!(rect1.is_square(), false);
    println!("✓ Basic methods work");

    // Test can_hold
    let rect2 = Rectangle::new(10, 40);
    let rect3 = Rectangle::new(60, 45);

    assert_eq!(rect1.can_hold(&rect2), true);
    assert_eq!(rect1.can_hold(&rect3), false);
    println!("✓ can_hold works");

    // Test is_square
    let square = Rectangle::new(25, 25);
    assert_eq!(square.is_square(), true);
    println!("✓ is_square works");

    // Test scale (mutable method)
    let mut rect4 = Rectangle::new(10, 20);
    rect4.scale(3);
    assert_eq!(rect4.width, 30);
    assert_eq!(rect4.height, 60);
    assert_eq!(rect4.area(), 1800);
    println!("✓ scale works");

    // Test builder pattern (method chaining)
    let rect5 = Rectangle::new(10, 10)
        .with_width(50)
        .with_height(100);
    assert_eq!(rect5.width, 50);
    assert_eq!(rect5.height, 100);
    assert_eq!(rect5.area(), 5000);
    println!("✓ Builder pattern works");

    println!("\n✅ All tests passed!");
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_area() {
        let rect = Rectangle::new(30, 50);
        assert_eq!(rect.area(), 1500);

        let square = Rectangle::new(10, 10);
        assert_eq!(square.area(), 100);
    }

    #[test]
    fn test_perimeter() {
        let rect = Rectangle::new(30, 50);
        assert_eq!(rect.perimeter(), 160);

        let square = Rectangle::new(10, 10);
        assert_eq!(square.perimeter(), 40);
    }

    #[test]
    fn test_is_square() {
        let rect = Rectangle::new(30, 50);
        assert_eq!(rect.is_square(), false);

        let square = Rectangle::new(25, 25);
        assert_eq!(square.is_square(), true);
    }

    #[test]
    fn test_can_hold_true() {
        let rect1 = Rectangle::new(30, 50);
        let rect2 = Rectangle::new(10, 40);
        assert_eq!(rect1.can_hold(&rect2), true);
    }

    #[test]
    fn test_can_hold_false() {
        let rect1 = Rectangle::new(30, 50);
        let rect2 = Rectangle::new(60, 45);
        assert_eq!(rect1.can_hold(&rect2), false);
    }

    #[test]
    fn test_can_hold_equal_size() {
        let rect1 = Rectangle::new(30, 50);
        let rect2 = Rectangle::new(30, 50);
        assert_eq!(rect1.can_hold(&rect2), false);  // Equal size doesn't count
    }

    #[test]
    fn test_scale() {
        let mut rect = Rectangle::new(10, 20);
        rect.scale(3);
        assert_eq!(rect.width, 30);
        assert_eq!(rect.height, 60);
        assert_eq!(rect.area(), 1800);
    }

    #[test]
    fn test_scale_by_one() {
        let mut rect = Rectangle::new(10, 20);
        rect.scale(1);
        assert_eq!(rect.width, 10);
        assert_eq!(rect.height, 20);
    }

    #[test]
    fn test_builder_width() {
        let rect = Rectangle::new(10, 20).with_width(50);
        assert_eq!(rect.width, 50);
        assert_eq!(rect.height, 20);
    }

    #[test]
    fn test_builder_height() {
        let rect = Rectangle::new(10, 20).with_height(100);
        assert_eq!(rect.width, 10);
        assert_eq!(rect.height, 100);
    }

    #[test]
    fn test_builder_chaining() {
        let rect = Rectangle::new(10, 10)
            .with_width(50)
            .with_height(100);
        assert_eq!(rect.width, 50);
        assert_eq!(rect.height, 100);
    }

    #[test]
    fn test_builder_multiple_chains() {
        let rect = Rectangle::new(1, 1)
            .with_width(10)
            .with_height(20)
            .with_width(30)  // Overwrite previous
            .with_height(40);
        assert_eq!(rect.width, 30);
        assert_eq!(rect.height, 40);
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All 12 tests pass: 60 points
- 9-11 tests pass: 50 points
- 6-8 tests pass: 40 points
- <6 tests pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Correct use of self parameters (20 points)**
   - Read-only methods use `&self`
   - Mutating method uses `&mut self`
   - Builder methods use `self` (ownership)

2. **Implementation correctness (10 points)**
   - All calculations correct
   - Edge cases handled (equal rectangles, scale by 1)

3. **Code quality (10 points)**
   - Clean, readable implementations
   - No unnecessary complexity
   - Idiomatic Rust

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on struct methods and self.

Exercise: Implement Rectangle methods with different self types (&self, &mut self, self).

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Correct self parameter types (&self for reading, &mut self for mutation, self for builders)
2. Implementation correctness (calculations, logic)
3. Understanding of when to use each self type
4. Code quality and style

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. Focus on:\n- If wrong self types: explain when to use &self vs &mut self vs self\n- If builder pattern incorrect: explain ownership and method chaining\n- If can_hold wrong: check for > vs >= (equal doesn't hold)\n- If calculations wrong: guide toward correct formulas\n- Praise correct use of different self types",
  "keyTakeaways": ["&self for reading", "&mut self for mutation", "self for transformation/building"],
  "encouragement": "Methods bring structs to life!"
}

Common patterns:
- Uses &self for scale → Should be &mut self
- Uses &mut self for area → Unnecessary, &self is enough
- Forgets to return self in builders → Method chaining breaks
- Equal rectangles in can_hold → Should be strictly greater
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Implement methods with appropriate `self` types
- ✅ Understand when to use `&self`, `&mut self`, and `self`
- ✅ Use the builder pattern for method chaining
- ✅ Distinguish associated functions from instance methods
- ✅ Apply methods to solve practical problems

## Common Mistakes to Watch For

1. **Wrong self type for scale**: Using `&self` instead of `&mut self`
   - Guide: Modifying methods need `&mut self`

2. **Over-privileged access**: Using `&mut self` when `&self` would work
   - Guide: Principle of least privilege - use immutable when possible

3. **Builder pattern errors**: Not returning `self` or not taking ownership
   - Guide: Builders take `self`, modify, and return `self`

4. **Forgetting `mut` keyword**: Calling `rect.scale(2)` on non-mut rectangle
   - Guide: Need `let mut rect = ...` to call mutating methods

5. **can_hold logic**: Using `>=` instead of `>` (equal rectangles don't "hold" each other)
   - Guide: Clarify the requirement - strictly larger

## Connection to Previous Exercises

**Exercise 3.1** introduced structs and associated functions (constructors). This exercise adds instance methods, showing the full picture:

```rust
impl Rectangle {
    // Associated function (no self) - creates instances
    fn new(width: u32, height: u32) -> Rectangle { ... }

    // Instance method (&self) - reads data
    fn area(&self) -> u32 { ... }

    // Instance method (&mut self) - modifies data
    fn scale(&mut self, factor: u32) { ... }

    // Instance method (self) - transforms/consumes
    fn with_width(self, width: u32) -> Self { ... }
}
```

## Preview of Next Exercise

This exercise focused on methods for a single struct. The next exercise introduces **enums** - types that can be one of several variants - and powerful **pattern matching** to handle them.

Teaser:
```rust
enum Shape {
    Rectangle { width: u32, height: u32 },
    Circle { radius: u32 },
    Triangle { base: u32, height: u32 },
}

// Pattern match to handle each variant differently
match shape {
    Shape::Rectangle { width, height } => width * height,
    Shape::Circle { radius } => 3 * radius * radius,
    Shape::Triangle { base, height } => (base * height) / 2,
}
```

## Extension Ideas (Post-MVP)

- Add `rotate()` method (swaps width and height)
- Implement `intersection` method (returns overlapping area)
- Add validation (reject zero/negative dimensions)
- Implement `Debug` and `Display` traits

## Notes for Implementation

- This exercise has more test cases (12) because it covers more concepts
- The builder pattern is a preview of more advanced Rust patterns
- Students often confuse when to use each `self` type - AI feedback should focus on this
- The `can_hold` edge case (equal sizes) is intentionally subtle - good for testing comprehension

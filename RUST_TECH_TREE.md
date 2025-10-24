# Rust Fundamentals Tech Tree

This document outlines the learning path structure for the Rust fundamentals tech tree in the Code Auto-Grader platform.

## Design Philosophy

**Goal**: Learn Rust's core concepts through progressive, hands-on practice.

**Key Principles**:
1. **Ownership Early**: Introduce Rust's unique ownership model early (can't write real Rust without it)
2. **Build Confidence**: Start with familiar concepts (variables, functions) before diving into complexity
3. **Practical Focus**: Exercises produce working programs, not just isolated snippets
4. **Progressive Difficulty**: Each node builds on previous knowledge with increasing complexity

## Tech Tree Structure (DAG)

```
[Node 1: Rust Basics]
         ↓
[Node 2: Ownership & Borrowing] ←──────┐
         ↓                              │
[Node 3: Structs & Enums] ──────────────┤
         ↓                              │
[Node 4: Collections & Strings] ────────┤
         ↓                              │
[Node 5: Error Handling] ───────────────┤
         ↓                              │
[Node 6: Traits & Generics] ←───────────┘
         ↓
[Node 7: Iterators & Closures] (Optional)
```

**Dependencies**:
- Node 1 → Node 2 (must understand basics before ownership)
- Node 2 → Node 3 (ownership needed for struct methods)
- Node 2 → Node 4 (ownership critical for collections)
- Node 2 → Node 5 (ownership affects Result/Option handling)
- Nodes 3, 4, 5 → Node 6 (traits build on structs, collections, error handling)
- Node 6 → Node 7 (generics needed for iterator trait bounds)

---

## Node 1: Rust Basics

**Learning Objectives**:
- Understand Rust syntax fundamentals
- Declare and use variables (mut vs immut)
- Write functions with parameters and return values
- Use basic types (integers, floats, booleans, tuples)
- Apply control flow (if/else, loops, match on simple types)

**Prerequisites**: None (starting point)

**Estimated Time**: 2-3 hours

### Exercise 1.1: Variables and Mutability
**Difficulty**: Easy
**Concept**: Variable binding, mutability, type inference

**Description**:
Fix the compilation errors in the provided code that deals with variable declarations and reassignments.

**Example Problem**:
```rust
fn main() {
    let x = 5;
    println!("The value of x is: {}", x);
    x = 6; // Error: cannot assign twice to immutable variable
    println!("The value of x is: {}", x);
}
```

**Learning Goal**: Understand that Rust variables are immutable by default and require `mut` keyword.

**Test Cases**:
- Code compiles successfully
- Output shows correct values
- Uses `mut` keyword appropriately

---

### Exercise 1.2: Functions and Return Values
**Difficulty**: Easy
**Concept**: Function signatures, parameters, return values

**Description**:
Write a function that takes two integers and returns their sum. Then write another function that calculates the area of a rectangle.

**Starting Code**:
```rust
// TODO: Write a function called `add` that takes two i32 parameters
// and returns their sum

// TODO: Write a function called `rectangle_area` that takes width and height
// as f64 parameters and returns the area

fn main() {
    let sum = add(5, 7);
    println!("Sum: {}", sum);

    let area = rectangle_area(10.5, 20.0);
    println!("Area: {}", area);
}
```

**Test Cases**:
- `add(5, 7)` returns 12
- `add(-3, 3)` returns 0
- `rectangle_area(10.5, 20.0)` returns 210.0
- `rectangle_area(0.0, 10.0)` returns 0.0

---

### Exercise 1.3: Control Flow and Pattern Matching
**Difficulty**: Medium
**Concept**: if/else, loops, match expressions

**Description**:
Write a function that categorizes a number as "negative", "zero", or "positive", then extend it to describe the magnitude ("small", "medium", "large") using nested conditions.

**Requirements**:
- Use `match` for the sign check
- Use `if/else` for magnitude (small: 0-10, medium: 11-100, large: >100)
- Handle negative numbers' absolute values for magnitude

**Starting Code**:
```rust
fn categorize_number(n: i32) -> String {
    // TODO: Implement using match for sign and if/else for magnitude
    todo!()
}

fn main() {
    println!("{}", categorize_number(-5));   // Should print: "negative (small)"
    println!("{}", categorize_number(0));    // Should print: "zero"
    println!("{}", categorize_number(50));   // Should print: "positive (medium)"
    println!("{}", categorize_number(150));  // Should print: "positive (large)"
}
```

**Test Cases**:
- categorize_number(-5) returns "negative (small)"
- categorize_number(0) returns "zero"
- categorize_number(50) returns "positive (medium)"
- categorize_number(150) returns "positive (large)"
- categorize_number(-75) returns "negative (medium)"

---

## Node 2: Ownership & Borrowing

**Learning Objectives**:
- Understand Rust's ownership model (one owner, no GC)
- Apply borrowing rules (immutable vs mutable references)
- Recognize when data is moved vs borrowed
- Use references to avoid unnecessary cloning
- Understand lifetimes in simple cases

**Prerequisites**: Node 1 (Rust Basics)

**Estimated Time**: 4-6 hours

**Why This Matters**: Ownership is Rust's superpower - it prevents memory leaks, data races, and null pointer errors at compile time. This is the hardest but most important concept.

### Exercise 2.1: Ownership and Moves
**Difficulty**: Medium
**Concept**: Ownership transfer, move semantics

**Description**:
Fix the ownership errors in code that tries to use a String after it's been moved.

**Example Problem**:
```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1;

    println!("{}, world!", s1); // Error: value borrowed here after move
}
```

**Learning Goal**: Understand that ownership transfers when assigning String (heap-allocated data).

**Test Cases**:
- Code compiles after fixing ownership issues
- Demonstrates correct use of cloning or borrowing
- No unnecessary copies

---

### Exercise 2.2: References and Borrowing
**Difficulty**: Medium
**Concept**: Immutable and mutable references

**Description**:
Write a function that calculates the length of a string without taking ownership. Then write another that modifies a string in place.

**Starting Code**:
```rust
// TODO: Implement this function to take a reference to String
// and return its length without taking ownership
fn calculate_length(s: /* fill in type */) -> usize {
    todo!()
}

// TODO: Implement this function to append " world" to the string
// It should modify the string in place using a mutable reference
fn append_world(s: /* fill in type */) {
    todo!()
}

fn main() {
    let s1 = String::from("hello");
    let len = calculate_length(&s1);
    println!("The length of '{}' is {}.", s1, len);

    let mut s2 = String::from("hello");
    append_world(&mut s2);
    println!("{}", s2); // Should print: "hello world"
}
```

**Test Cases**:
- calculate_length works without taking ownership
- append_world modifies the original string
- Original variables remain usable after function calls

---

### Exercise 2.3: The Borrowing Rules
**Difficulty**: Hard
**Concept**: Multiple immutable refs OR one mutable ref (not both)

**Description**:
You're given code with borrowing violations. Fix them by restructuring the code to follow Rust's borrowing rules.

**Example Problem**:
```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    let r3 = &mut s; // Error: cannot borrow as mutable while immutable refs exist

    println!("{}, {}, and {}", r1, r2, r3);
}
```

**Learning Goal**: Understand that you can't have mutable and immutable references to the same data simultaneously.

**Test Cases**:
- Code compiles after restructuring
- Demonstrates understanding of reference scopes
- Maintains intended functionality

---

## Node 3: Structs & Enums

**Learning Objectives**:
- Define custom structs with named fields
- Implement methods on structs
- Use associated functions (constructors)
- Create and match on enums
- Apply Option<T> for nullable values

**Prerequisites**: Node 2 (Ownership & Borrowing)

**Estimated Time**: 3-4 hours

### Exercise 3.1: Defining and Using Structs
**Difficulty**: Easy
**Concept**: Struct definition, instantiation, field access

**Description**:
Create a `User` struct with username, email, and active status. Implement a constructor function.

**Starting Code**:
```rust
// TODO: Define a User struct with fields:
// - username: String
// - email: String
// - active: bool

// TODO: Implement an associated function `new` that creates a User
impl User {
    // Your code here
}

fn main() {
    let user1 = User::new(
        String::from("alice"),
        String::from("alice@example.com"),
    );

    println!("User: {}, Active: {}", user1.username, user1.active);
}
```

**Test Cases**:
- User struct has correct fields
- Constructor creates user with active=true by default
- Fields are accessible

---

### Exercise 3.2: Methods and Self
**Difficulty**: Medium
**Concept**: Instance methods, &self, &mut self

**Description**:
Add methods to a `Rectangle` struct to calculate area and check if it can hold another rectangle.

**Starting Code**:
```rust
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    // TODO: Implement `new(width, height) -> Rectangle`

    // TODO: Implement `area(&self) -> u32`

    // TODO: Implement `can_hold(&self, other: &Rectangle) -> bool`
    // Returns true if self can hold other rectangle inside it
}

fn main() {
    let rect1 = Rectangle::new(30, 50);
    let rect2 = Rectangle::new(10, 40);

    println!("Area: {}", rect1.area());
    println!("Can rect1 hold rect2? {}", rect1.can_hold(&rect2));
}
```

**Test Cases**:
- area() returns correct calculation
- can_hold() returns true when both dimensions are larger
- can_hold() returns false otherwise

---

### Exercise 3.3: Enums and Pattern Matching
**Difficulty**: Medium
**Concept**: Enum variants, match expressions, Option<T>

**Description**:
Create a `Message` enum with different variants (Quit, Move, Write, ChangeColor). Implement a function that processes messages.

**Starting Code**:
```rust
// TODO: Define enum Message with variants:
// - Quit (no data)
// - Move { x: i32, y: i32 }
// - Write(String)
// - ChangeColor(u8, u8, u8)

// TODO: Implement this function to return a description of the message
fn process_message(msg: Message) -> String {
    // Use match to handle each variant
    todo!()
}

fn main() {
    let msg1 = Message::Quit;
    let msg2 = Message::Move { x: 10, y: 20 };
    let msg3 = Message::Write(String::from("Hello"));

    println!("{}", process_message(msg1));
    println!("{}", process_message(msg2));
    println!("{}", process_message(msg3));
}
```

**Test Cases**:
- All variants compile correctly
- match handles each variant appropriately
- Returns sensible descriptions

---

## Node 4: Collections & Strings

**Learning Objectives**:
- Use Vec<T> for dynamic arrays
- Manipulate String and &str types
- Work with HashMap<K, V>
- Iterate over collections safely
- Understand String vs &str differences

**Prerequisites**: Node 2 (Ownership & Borrowing)

**Estimated Time**: 3-4 hours

### Exercise 4.1: Working with Vectors
**Difficulty**: Easy
**Concept**: Vec creation, push, pop, indexing, iteration

**Description**:
Write functions to find the median and mode of a list of integers.

**Starting Code**:
```rust
fn median(numbers: &Vec<i32>) -> f64 {
    // TODO: Sort the numbers and return the middle value
    // If even length, return average of two middle values
    todo!()
}

fn mode(numbers: &Vec<i32>) -> i32 {
    // TODO: Return the most frequently occurring number
    // You may use a HashMap to count occurrences
    todo!()
}

fn main() {
    let numbers = vec![1, 2, 2, 3, 4, 4, 4, 5];
    println!("Median: {}", median(&numbers));
    println!("Mode: {}", mode(&numbers));
}
```

**Test Cases**:
- median([1,2,3,4,5]) returns 3.0
- median([1,2,3,4]) returns 2.5
- mode([1,2,2,3]) returns 2

---

### Exercise 4.2: String Manipulation
**Difficulty**: Medium
**Concept**: String vs &str, concatenation, ownership with strings

**Description**:
Write a function that takes a sentence and returns it with each word capitalized (title case).

**Starting Code**:
```rust
fn to_title_case(s: &str) -> String {
    // TODO: Split the string into words, capitalize first letter of each,
    // and join them back together
    todo!()
}

fn main() {
    let sentence = "hello world from rust";
    let title = to_title_case(sentence);
    println!("{}", title); // Should print: "Hello World From Rust"
}
```

**Test Cases**:
- to_title_case("hello world") returns "Hello World"
- to_title_case("a") returns "A"
- to_title_case("") returns ""

---

### Exercise 4.3: HashMap for Counting
**Difficulty**: Medium
**Concept**: HashMap insertion, retrieval, entry API

**Description**:
Write a program that counts the frequency of words in a text.

**Starting Code**:
```rust
use std::collections::HashMap;

fn word_frequency(text: &str) -> HashMap<String, u32> {
    // TODO: Split text into words (by whitespace)
    // Count each word's frequency in a HashMap
    // Convert words to lowercase for case-insensitive counting
    todo!()
}

fn main() {
    let text = "hello world hello rust world";
    let freq = word_frequency(text);

    for (word, count) in &freq {
        println!("{}: {}", word, count);
    }
}
```

**Test Cases**:
- word_frequency("hello world hello") returns {"hello": 2, "world": 1}
- Case insensitive: "Hello hello" counts as 2
- Empty string returns empty HashMap

---

## Node 5: Error Handling

**Learning Objectives**:
- Use Option<T> for values that might be absent
- Apply Result<T, E> for operations that can fail
- Propagate errors with the ? operator
- Choose between panic! and Result
- Create custom error types (basic)

**Prerequisites**: Node 2 (Ownership) and Node 3 (Enums)

**Estimated Time**: 3-4 hours

**Why This Matters**: Rust forces you to handle errors explicitly, preventing unexpected crashes.

### Exercise 5.1: Option and None
**Difficulty**: Easy
**Concept**: Option<T>, Some(T), None, unwrap, pattern matching

**Description**:
Write a function that finds the first even number in a vector.

**Starting Code**:
```rust
fn first_even(numbers: &Vec<i32>) -> Option<i32> {
    // TODO: Return Some(n) for the first even number found,
    // or None if no even numbers exist
    todo!()
}

fn main() {
    let numbers = vec![1, 3, 5, 8, 9];

    match first_even(&numbers) {
        Some(n) => println!("First even: {}", n),
        None => println!("No even numbers found"),
    }
}
```

**Test Cases**:
- first_even([1,3,5,8,9]) returns Some(8)
- first_even([1,3,5]) returns None
- first_even([2,4,6]) returns Some(2)

---

### Exercise 5.2: Result and Error Propagation
**Difficulty**: Medium
**Concept**: Result<T, E>, ? operator, error handling

**Description**:
Write a function that parses a string to an integer and doubles it. Handle parse errors gracefully.

**Starting Code**:
```rust
use std::num::ParseIntError;

fn double_number(s: &str) -> Result<i32, ParseIntError> {
    // TODO: Parse the string to i32, then double it
    // Use ? to propagate errors
    todo!()
}

fn main() {
    match double_number("21") {
        Ok(n) => println!("Result: {}", n),
        Err(e) => println!("Error: {}", e),
    }

    match double_number("invalid") {
        Ok(n) => println!("Result: {}", n),
        Err(e) => println!("Error: {}", e),
    }
}
```

**Test Cases**:
- double_number("21") returns Ok(42)
- double_number("0") returns Ok(0)
- double_number("invalid") returns Err(ParseIntError)

---

### Exercise 5.3: Combining Option and Result
**Difficulty**: Hard
**Concept**: Nested error handling, combinators (map, and_then)

**Description**:
Write a function that reads a "database" (HashMap) and retrieves a user's age, then determines if they're an adult.

**Starting Code**:
```rust
use std::collections::HashMap;

struct User {
    name: String,
    age: u8,
}

fn is_adult(users: &HashMap<String, User>, name: &str) -> Result<bool, String> {
    // TODO:
    // 1. Look up the user by name (returns Option)
    // 2. Check if age >= 18
    // 3. Return Ok(true/false) if found, Err("User not found") if not
    todo!()
}

fn main() {
    let mut users = HashMap::new();
    users.insert(
        String::from("Alice"),
        User { name: String::from("Alice"), age: 25 }
    );
    users.insert(
        String::from("Bob"),
        User { name: String::from("Bob"), age: 16 }
    );

    println!("{:?}", is_adult(&users, "Alice")); // Ok(true)
    println!("{:?}", is_adult(&users, "Bob"));   // Ok(false)
    println!("{:?}", is_adult(&users, "Charlie")); // Err("User not found")
}
```

**Test Cases**:
- is_adult with age >= 18 returns Ok(true)
- is_adult with age < 18 returns Ok(false)
- is_adult with non-existent user returns Err

---

## Node 6: Traits & Generics

**Learning Objectives**:
- Define traits to share behavior
- Implement traits for custom types
- Use generic types with trait bounds
- Understand trait objects (basic)
- Apply common traits (Debug, Clone, PartialEq)

**Prerequisites**: Nodes 3 (Structs & Enums), 4 (Collections), 5 (Error Handling)

**Estimated Time**: 4-5 hours

**Why This Matters**: Traits are Rust's interface system, enabling polymorphism and code reuse.

### Exercise 6.1: Defining and Implementing Traits
**Difficulty**: Medium
**Concept**: Trait definition, implementation for structs

**Description**:
Create a `Summary` trait and implement it for `Article` and `Tweet` structs.

**Starting Code**:
```rust
// TODO: Define a trait called Summary with a method:
// fn summarize(&self) -> String;

struct Article {
    title: String,
    author: String,
    content: String,
}

struct Tweet {
    username: String,
    content: String,
}

// TODO: Implement Summary for Article
// Format: "Title by Author"

// TODO: Implement Summary for Tweet
// Format: "@username: content"

fn main() {
    let article = Article {
        title: String::from("Rust is great"),
        author: String::from("Alice"),
        content: String::from("..."),
    };

    let tweet = Tweet {
        username: String::from("bob"),
        content: String::from("Learning Rust!"),
    };

    println!("{}", article.summarize());
    println!("{}", tweet.summarize());
}
```

**Test Cases**:
- Article summarize returns "Title by Author" format
- Tweet summarize returns "@username: content" format

---

### Exercise 6.2: Generic Functions
**Difficulty**: Medium
**Concept**: Generic type parameters, trait bounds

**Description**:
Write a generic function that finds the largest element in a slice.

**Starting Code**:
```rust
// TODO: Implement this function with generic type T
// T must implement PartialOrd and Clone traits
fn largest<T>(list: &[T]) -> T {
    todo!()
}

fn main() {
    let numbers = vec![34, 50, 25, 100, 65];
    let result = largest(&numbers);
    println!("The largest number is {}", result);

    let chars = vec!['y', 'm', 'a', 'q'];
    let result = largest(&chars);
    println!("The largest char is {}", result);
}
```

**Test Cases**:
- largest([34, 50, 25, 100, 65]) returns 100
- largest(['y', 'm', 'a', 'q']) returns 'y'
- largest([1.0, 3.5, 2.1]) returns 3.5

---

### Exercise 6.3: Generic Structs with Traits
**Difficulty**: Hard
**Concept**: Generic structs, trait bounds on implementations

**Description**:
Create a generic `Pair<T>` struct and implement methods with trait bounds.

**Starting Code**:
```rust
struct Pair<T> {
    first: T,
    second: T,
}

impl<T> Pair<T> {
    fn new(first: T, second: T) -> Self {
        // TODO: Implement constructor
        todo!()
    }
}

// TODO: Implement a method `cmp_display` that only works when T implements
// Display + PartialOrd. It should print which element is larger.
impl<T: /* add trait bounds */> Pair<T> {
    fn cmp_display(&self) {
        // TODO: Compare first and second, print the larger one
        todo!()
    }
}

fn main() {
    let pair = Pair::new(5, 10);
    pair.cmp_display(); // Should print: "The largest member is 10"
}
```

**Test Cases**:
- cmp_display correctly identifies larger element
- Works with i32, f64, etc.
- Trait bounds prevent compilation with non-comparable types

---

## Node 7: Iterators & Closures (Optional/Advanced)

**Learning Objectives**:
- Create and use closures (anonymous functions)
- Understand closure capture modes (by reference, by mutable reference, by value)
- Use iterator methods (map, filter, fold, etc.)
- Implement the Iterator trait
- Chain iterator adaptors for functional-style code

**Prerequisites**: Node 6 (Traits & Generics)

**Estimated Time**: 4-5 hours

**Note**: This node is optional for MVP. Include only if you want to reach intermediate Rust proficiency.

### Exercise 7.1: Closures and Capturing Environment
**Difficulty**: Medium
**Concept**: Closure syntax, environment capture

**Description**:
Use closures to filter and transform a list of numbers based on runtime conditions.

**Starting Code**:
```rust
fn main() {
    let numbers = vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    let threshold = 5;

    // TODO: Use a closure to filter numbers greater than threshold
    // Then use another closure to square each number
    // Collect into a new Vec

    // let result = numbers.iter()...

    println!("{:?}", result); // Should print: [36, 49, 64, 81, 100]
}
```

**Test Cases**:
- Correctly filters numbers > 5
- Correctly squares each number
- Uses closures (not separate functions)

---

### Exercise 7.2: Iterator Methods
**Difficulty**: Medium
**Concept**: map, filter, fold, collect

**Description**:
Given a list of strings, filter out short strings, convert to uppercase, and concatenate.

**Starting Code**:
```rust
fn process_strings(strings: Vec<String>, min_length: usize) -> String {
    // TODO: Use iterator methods to:
    // 1. Filter strings with length >= min_length
    // 2. Convert each to uppercase
    // 3. Concatenate with ", " separator
    todo!()
}

fn main() {
    let words = vec![
        String::from("hi"),
        String::from("hello"),
        String::from("world"),
        String::from("rust"),
    ];

    let result = process_strings(words, 4);
    println!("{}", result); // Should print: "HELLO, WORLD, RUST"
}
```

**Test Cases**:
- Filters correctly by length
- Converts to uppercase
- Joins with ", "

---

### Exercise 7.3: Custom Iterator
**Difficulty**: Hard
**Concept**: Implementing the Iterator trait

**Description**:
Implement a counter that generates a sequence of numbers.

**Starting Code**:
```rust
struct Counter {
    count: u32,
    max: u32,
}

impl Counter {
    fn new(max: u32) -> Counter {
        Counter { count: 0, max }
    }
}

// TODO: Implement the Iterator trait for Counter
// It should yield numbers from 1 to max (inclusive)

fn main() {
    let counter = Counter::new(5);

    for num in counter {
        println!("{}", num);
    }
    // Should print: 1, 2, 3, 4, 5
}
```

**Test Cases**:
- Iterates from 1 to max
- Stops after max
- Works with for loops and iterator methods

---

## Implementation Notes

### Exercise Format
Each exercise should be stored in the database with:
- **Title**: Short name (e.g., "Variables and Mutability")
- **Description**: Problem statement in markdown
- **Starting Code**: Template code with TODOs
- **Reference Solution**: Correct implementation (hidden from students)
- **Test Cases**: Automated tests to verify correctness
- **Hints**: Optional graduated hints (defer to post-MVP)

### Test Case Structure
```rust
// Test format for automated grading
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_case_1() {
        assert_eq!(function_name(input), expected_output);
    }
}
```

### Scoring Criteria
- **Test Pass Rate**: 60% of score (how many test cases pass)
- **AI Evaluation**: 40% of score (code quality, style, understanding)
- **Mastery Threshold**: 80% overall score to count toward node completion

### Content Generation Priority
For dogfooding MVP, create exercises in this order:
1. **Node 1 (Basics)** - 3 exercises (start here)
2. **Node 2 (Ownership)** - 3 exercises (most critical)
3. **Node 3 (Structs/Enums)** - 3 exercises
4. **Node 4 (Collections)** - 3 exercises
5. **Node 5 (Error Handling)** - 3 exercises
6. **Node 6 (Traits/Generics)** - Defer to post-dogfooding
7. **Node 7 (Iterators)** - Defer to post-dogfooding

Start with Nodes 1-3 for initial testing (9 exercises total).

---

## Success Metrics for Dogfooding

After completing Nodes 1-3:
- [ ] Can write simple Rust programs confidently
- [ ] Understand ownership errors and how to fix them
- [ ] Feel comfortable with struct and enum patterns
- [ ] AI feedback helped improve code quality
- [ ] Want to continue learning (intrinsic motivation)

If these criteria are met, the learning loop works. Proceed to build remaining nodes and multi-user features.

{
  "id": "30000000-0000-0000-0000-000000000007",
  "node_id": "20000000-0000-0000-0000-000000000003",
  "title": "Exercise 3.1: Defining and Using Structs",
  "slug": "defining-structs",
  "difficulty": "easy",
  "estimated_minutes": 25,
  "order_index": 1
}
---

# Exercise 3.1: Defining and Using Structs

## Learning Objectives
- Define custom structs with named fields
- Instantiate structs and access fields
- Implement associated functions (constructors)
- Apply struct update syntax
- Understand ownership with structs (structs own their data)

## Problem Description

So far, you've worked with Rust's built-in types (integers, strings, etc.). But real programs need **custom data types** to model domain concepts.

**Structs** (short for "structures") let you group related data together with named fields. They're similar to classes in other languages, but without methods (yet - that's the next exercise!).

### Struct Basics

```rust
// Define a struct
struct User {
    username: String,
    email: String,
    active: bool,
    sign_in_count: u64,
}

// Create an instance
let user = User {
    username: String::from("alice"),
    email: String::from("alice@example.com"),
    active: true,
    sign_in_count: 1,
};

// Access fields
println!("Username: {}", user.username);
```

### Associated Functions (Constructors)

Instead of manually creating instances, we usually write **associated functions** (like static methods in other languages) to construct them:

```rust
impl User {
    fn new(username: String, email: String) -> User {
        User {
            username,
            email,
            active: true,
            sign_in_count: 0,
        }
    }
}

// Usage
let user = User::new(
    String::from("alice"),
    String::from("alice@example.com")
);
```

### Your Task

Implement several structs to model a simple book management system:

1. **`Book`**: Title, author, pages, available (bool)
2. **`Library`**: Name, books (Vec<Book>)

You'll implement:
- Struct definitions
- Constructor functions
- Field access and modification
- A function to count available books in a library

## Starting Code

```rust
// TODO: Define a Book struct with the following fields:
// - title: String
// - author: String
// - pages: u32
// - available: bool

// TODO: Implement an associated function `new` for Book
// that takes title, author, and pages and returns a Book
// with available set to true by default
impl Book {
    // Your code here
}

// TODO: Define a Library struct with:
// - name: String
// - books: Vec<Book>

// TODO: Implement an associated function `new` for Library
// that takes a name and returns a Library with an empty books vector
impl Library {
    // Your code here
}

// TODO: Implement a method `count_available_books` on Library
// that returns how many books have available = true
impl Library {
    fn count_available_books(&self) -> usize {
        todo!()
    }
}

fn main() {
    // Test Book creation
    let book1 = Book::new(
        String::from("The Rust Programming Language"),
        String::from("Steve Klabnik and Carol Nichols"),
        552,
    );

    assert_eq!(book1.title, "The Rust Programming Language");
    assert_eq!(book1.author, "Steve Klabnik and Carol Nichols");
    assert_eq!(book1.pages, 552);
    assert_eq!(book1.available, true);
    println!("✓ Book created successfully");

    // Test Library creation
    let mut library = Library::new(String::from("City Library"));
    assert_eq!(library.name, "City Library");
    assert_eq!(library.books.len(), 0);
    println!("✓ Library created successfully");

    // Add books to library
    library.books.push(book1);

    let book2 = Book::new(
        String::from("Programming Rust"),
        String::from("Jim Blandy and Jason Orendorff"),
        622,
    );
    library.books.push(book2);

    let mut book3 = Book::new(
        String::from("Rust for Rustaceans"),
        String::from("Jon Gjengset"),
        282,
    );
    book3.available = false;  // This book is checked out
    library.books.push(book3);

    assert_eq!(library.books.len(), 3);
    println!("✓ Books added to library");

    // Test count_available_books
    let available_count = library.count_available_books();
    assert_eq!(available_count, 2);
    println!("✓ Available books counted correctly");

    println!("\n✅ All tests passed!");
}
```

## Expected Output

```
✓ Book created successfully
✓ Library created successfully
✓ Books added to library
✓ Available books counted correctly

✅ All tests passed!
```

## Hints

<details>
<summary>Hint 1: Struct Definition Syntax</summary>

```rust
struct StructName {
    field1: Type1,
    field2: Type2,
    field3: Type3,
}
```

Example:
```rust
struct Book {
    title: String,
    author: String,
    pages: u32,
    available: bool,
}
```
</details>

<details>
<summary>Hint 2: Associated Function (Constructor)</summary>

Associated functions are defined in an `impl` block and don't take `self`:

```rust
impl Book {
    fn new(title: String, author: String, pages: u32) -> Book {
        Book {
            title,        // Shorthand when variable name matches field name
            author,
            pages,
            available: true,  // Default value
        }
    }
}
```

Call it with: `Book::new(...)`
</details>

<details>
<summary>Hint 3: Vec<Book> and Library Constructor</summary>

For the Library constructor:
```rust
impl Library {
    fn new(name: String) -> Library {
        Library {
            name,
            books: Vec::new(),  // Empty vector
        }
    }
}
```
</details>

<details>
<summary>Hint 4: Counting Available Books</summary>

You need to iterate through the books and count how many have `available == true`:

```rust
fn count_available_books(&self) -> usize {
    let mut count = 0;
    for book in &self.books {
        if book.available {
            count += 1;
        }
    }
    count
}
```

Or using iterators (more idiomatic):
```rust
fn count_available_books(&self) -> usize {
    self.books.iter().filter(|book| book.available).count()
}
```
</details>

<details>
<summary>Hint 5: Ownership and Structs</summary>

Important: Structs **own** their data. When you create a `Book`:
```rust
let book = Book::new(String::from("Title"), ...);
```

The `Book` instance owns the `String`s. When `book` goes out of scope, the strings are dropped.

This is why `count_available_books` takes `&self` (borrow) - it doesn't take ownership of the Library!
</details>

## Reference Solution

```rust
// Define Book struct
struct Book {
    title: String,
    author: String,
    pages: u32,
    available: bool,
}

// Implement constructor for Book
impl Book {
    fn new(title: String, author: String, pages: u32) -> Book {
        Book {
            title,
            author,
            pages,
            available: true,
        }
    }
}

// Define Library struct
struct Library {
    name: String,
    books: Vec<Book>,
}

// Implement constructor for Library
impl Library {
    fn new(name: String) -> Library {
        Library {
            name,
            books: Vec::new(),
        }
    }

    fn count_available_books(&self) -> usize {
        self.books.iter().filter(|book| book.available).count()
    }
}

fn main() {
    // Test Book creation
    let book1 = Book::new(
        String::from("The Rust Programming Language"),
        String::from("Steve Klabnik and Carol Nichols"),
        552,
    );

    assert_eq!(book1.title, "The Rust Programming Language");
    assert_eq!(book1.author, "Steve Klabnik and Carol Nichols");
    assert_eq!(book1.pages, 552);
    assert_eq!(book1.available, true);
    println!("✓ Book created successfully");

    // Test Library creation
    let mut library = Library::new(String::from("City Library"));
    assert_eq!(library.name, "City Library");
    assert_eq!(library.books.len(), 0);
    println!("✓ Library created successfully");

    // Add books to library
    library.books.push(book1);

    let book2 = Book::new(
        String::from("Programming Rust"),
        String::from("Jim Blandy and Jason Orendorff"),
        622,
    );
    library.books.push(book2);

    let mut book3 = Book::new(
        String::from("Rust for Rustaceans"),
        String::from("Jon Gjengset"),
        282,
    );
    book3.available = false;
    library.books.push(book3);

    assert_eq!(library.books.len(), 3);
    println!("✓ Books added to library");

    // Test count_available_books
    let available_count = library.count_available_books();
    assert_eq!(available_count, 2);
    println!("✓ Available books counted correctly");

    println!("\n✅ All tests passed!");
}
```

### Alternative Implementation (Verbose Loop)

```rust
// More explicit version of count_available_books
impl Library {
    fn count_available_books(&self) -> usize {
        let mut count = 0;
        for book in &self.books {
            if book.available {
                count += 1;
            }
        }
        count
    }
}
```

Both are correct! The iterator version is more idiomatic Rust.

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_book_creation() {
        let book = Book::new(
            String::from("Test Book"),
            String::from("Test Author"),
            100,
        );
        assert_eq!(book.title, "Test Book");
        assert_eq!(book.author, "Test Author");
        assert_eq!(book.pages, 100);
        assert_eq!(book.available, true);
    }

    #[test]
    fn test_library_creation() {
        let library = Library::new(String::from("Test Library"));
        assert_eq!(library.name, "Test Library");
        assert_eq!(library.books.len(), 0);
    }

    #[test]
    fn test_library_with_books() {
        let mut library = Library::new(String::from("Library"));

        let book1 = Book::new(String::from("Book 1"), String::from("Author 1"), 100);
        let book2 = Book::new(String::from("Book 2"), String::from("Author 2"), 200);

        library.books.push(book1);
        library.books.push(book2);

        assert_eq!(library.books.len(), 2);
    }

    #[test]
    fn test_count_all_available() {
        let mut library = Library::new(String::from("Library"));

        library.books.push(Book::new(String::from("Book 1"), String::from("A"), 100));
        library.books.push(Book::new(String::from("Book 2"), String::from("B"), 200));
        library.books.push(Book::new(String::from("Book 3"), String::from("C"), 300));

        assert_eq!(library.count_available_books(), 3);
    }

    #[test]
    fn test_count_some_available() {
        let mut library = Library::new(String::from("Library"));

        let mut book1 = Book::new(String::from("Book 1"), String::from("A"), 100);
        let book2 = Book::new(String::from("Book 2"), String::from("B"), 200);
        let mut book3 = Book::new(String::from("Book 3"), String::from("C"), 300);

        book1.available = false;
        book3.available = false;

        library.books.push(book1);
        library.books.push(book2);
        library.books.push(book3);

        assert_eq!(library.count_available_books(), 1);
    }

    #[test]
    fn test_count_none_available() {
        let mut library = Library::new(String::from("Library"));

        let mut book1 = Book::new(String::from("Book 1"), String::from("A"), 100);
        book1.available = false;
        library.books.push(book1);

        assert_eq!(library.count_available_books(), 0);
    }

    #[test]
    fn test_empty_library_count() {
        let library = Library::new(String::from("Empty Library"));
        assert_eq!(library.count_available_books(), 0);
    }

    #[test]
    fn test_book_field_modification() {
        let mut book = Book::new(String::from("Test"), String::from("Author"), 100);
        assert_eq!(book.available, true);

        book.available = false;
        assert_eq!(book.available, false);

        book.available = true;
        assert_eq!(book.available, true);
    }
}
```

## Automated Grading Rubric

### Test Pass Rate (60% of score)
- All 8 tests pass: 60 points
- 6-7 tests pass: 48 points
- 4-5 tests pass: 30 points
- <4 tests pass: 0 points

### AI Evaluation Criteria (40% of score)

The AI should evaluate based on:

1. **Correct struct definitions (15 points)**
   - Fields have correct names and types
   - Follows naming conventions (snake_case for fields, PascalCase for structs)

2. **Proper use of associated functions (10 points)**
   - Constructor functions follow `new` convention
   - Don't take `self` (they're associated, not instance methods)
   - Return instances correctly

3. **Implementation correctness (10 points)**
   - `count_available_books` correctly counts
   - Handles edge cases (empty library)
   - Uses appropriate iteration

4. **Code quality (5 points)**
   - Clean, readable code
   - Idiomatic Rust (iterator style preferred but not required)
   - Good use of ownership (borrowing `&self` in methods)

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on struct definitions and usage.

Exercise: Define Book and Library structs with constructors and a counting method.

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

Evaluate the solution on:
1. Struct definitions with correct field types
2. Associated functions (constructors) vs instance methods
3. Correct implementation of count_available_books
4. Code style and idioms (iterator vs loop)

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. Focus on:\n- If struct fields wrong: explain type requirements\n- If constructor takes self: explain associated functions vs methods\n- If count wrong: guide toward iteration and filtering\n- If using verbose loop: acknowledge correctness, suggest iterator style as more idiomatic\n- If using iterator: praise idiomatic Rust!",
  "keyTakeaways": ["Structs own their data", "Associated functions vs methods", "Borrowing self"],
  "encouragement": "Building on ownership knowledge!"
}

Common patterns:
- Wrong field types (e.g., &str instead of String) → Explain struct ownership
- Constructor takes self → Explain associated vs instance functions
- Count implementation broken → Guide toward filtering
- Verbose but correct → Suggest iterator as alternative
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Define custom structs with named fields
- ✅ Understand that structs own their data
- ✅ Create associated functions (constructors)
- ✅ Distinguish between associated functions and instance methods
- ✅ Access and modify struct fields
- ✅ Work with structs inside collections (Vec<Book>)
- ✅ Borrow structs with `&self` in methods

## Common Mistakes to Watch For

1. **Using &str instead of String in struct fields**: Causes lifetime issues
   - Guide: Structs should own their data; use String, not &str

2. **Making constructors take `self`**: `fn new(self, ...) -> Book`
   - Guide: Constructors don't have a self yet! They're associated functions

3. **Not borrowing `self` in `count_available_books`**: Taking ownership
   - Guide: Methods usually borrow `&self` or `&mut self`, not take ownership

4. **Incorrect counting logic**: Off-by-one, not filtering properly
   - Guide: Use `.iter()` and `.filter()` for cleaner counting

5. **Forgetting `mut` for library**: `let library = Library::new(...)` without `mut`
   - Guide: Need `mut` to push books into the vector

## Connection to Previous Nodes

**Node 2 (Ownership)** taught that values have owners and can be borrowed. This exercise shows:
- Structs **own** their fields (e.g., `Book` owns its `String`s)
- Methods **borrow** instances with `&self` (don't take ownership)
- Collections of structs work just like collections of basic types

Key insight:
```rust
let book = Book::new(...);           // book owns the Book
library.books.push(book);            // ownership moves into vector
// book is no longer valid here

library.count_available_books();     // Borrows library with &self
// library is still valid after the method call
```

## Preview of Next Exercise

This exercise introduced structs and associated functions. The next exercise adds **instance methods** - functions that can read or modify a struct's fields.

Teaser:
```rust
impl Book {
    // Associated function (no self)
    fn new(...) -> Book { ... }

    // Instance method (takes &self)
    fn summary(&self) -> String { ... }

    // Mutable instance method (takes &mut self)
    fn checkout(&mut self) { self.available = false; }
}
```

## Extension Ideas (Post-MVP)

- Add a `Genre` enum and a `genre` field to `Book`
- Implement `find_book_by_title` method on `Library`
- Add `checkout` and `return_book` methods
- Implement `Display` trait for pretty printing

## Notes for Implementation

- This is a gentler re-introduction after the intense Node 2
- Deliberately uses familiar concepts (ownership, borrowing) in a new context
- The next exercise (methods) builds directly on this foundation
- Consider showing memory diagrams (struct layout on stack/heap)

# Exercise 2.3: The Borrowing Rules

## Metadata
- **Node**: 2 (Ownership & Borrowing)
- **Exercise**: 3 of 3
- **Difficulty**: Hard
- **Estimated Time**: 45-60 minutes
- **Prerequisites**: Exercises 2.1 & 2.2

## Learning Objectives
- Understand Rust's borrowing rules (the core safety guarantees)
- Recognize when mutable and immutable references conflict
- Apply reference scoping to resolve borrowing conflicts
- Restructure code to satisfy the borrow checker
- Understand **why** these rules prevent data races and bugs

## Problem Description

Rust's borrowing system has strict rules that prevent data races and use-after-free bugs at compile time. These rules might seem frustrating at first, but they're catching real bugs before your code ever runs!

### The Borrowing Rules

1. **Rule 1**: You can have **either** (but not both):
   - ONE mutable reference (`&mut T`), **OR**
   - Any number of immutable references (`&T`)

2. **Rule 2**: References must always be valid (no dangling references)

3. **Rule 3**: A reference's scope starts where it's introduced and continues through its last use

### Why These Rules Exist

```rust
// WITHOUT these rules, this could crash:
let mut data = vec![1, 2, 3];
let reference = &data[0];  // Immutable borrow
data.push(4);              // Mutable operation - could reallocate memory!
println!("{}", reference); // reference might point to freed memory! 💥
```

Rust prevents this **at compile time**.

### Your Task

You'll fix several borrowing violations in code that tries to:
1. Use mutable and immutable references simultaneously
2. Modify data while it's being read
3. Hold references across operations that invalidate them

You'll learn to restructure code to satisfy the borrow checker.

## Starting Code

```rust
fn main() {
    // Part 1: Can't have mutable and immutable references simultaneously
    println!("=== Part 1: Mixed Borrows ===");
    let mut s = String::from("hello");

    let r1 = &s;       // Immutable borrow
    let r2 = &s;       // Another immutable borrow - OK!
    let r3 = &mut s;   // ERROR: Can't borrow as mutable while immutable borrows exist

    println!("r1: {}, r2: {}, r3: {}", r1, r2, r3);

    // Part 2: Can't have multiple mutable references simultaneously
    println!("\n=== Part 2: Multiple Mutable Borrows ===");
    let mut s = String::from("hello");

    let r1 = &mut s;
    let r2 = &mut s;  // ERROR: Can't borrow as mutable more than once

    println!("r1: {}, r2: {}", r1, r2);

    // Part 3: Scope-based borrowing
    println!("\n=== Part 3: Reference Scopes ===");
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);
    // r1 and r2 are no longer used after this point

    let r3 = &mut s;  // This should be OK - immutable refs are done
    r3.push_str(" world");
    println!("r3: {}", r3);

    // Part 4: Modifying while borrowed
    println!("\n=== Part 4: Modification During Borrow ===");
    let mut numbers = vec![1, 2, 3, 4, 5];

    let first = &numbers[0];  // Immutable borrow
    numbers.push(6);          // ERROR: Can't mutate while borrowed
    println!("First number: {}", first);

    // Part 5: Function calls and borrowing
    println!("\n=== Part 5: Functions and Borrows ===");
    let mut text = String::from("hello");

    let len = get_length(&text);
    append_exclamation(&mut text);  // ERROR: if immutable borrow still active
    println!("Text: {}, Length was: {}", text, len);

    println!("\n✅ All parts working!");
}

fn get_length(s: &String) -> usize {
    s.len()
}

fn append_exclamation(s: &mut String) {
    s.push('!');
}
```

## Expected Output (After Fixes)

```
=== Part 1: Mixed Borrows ===
r1: hello, r2: hello, r3: hello

(or with restructured code, might print differently but compile successfully)

=== Part 2: Multiple Mutable Borrows ===
r1: hello, r2: hello

=== Part 3: Reference Scopes ===
r1: hello, r2: hello
r3: hello world

=== Part 4: Modification During Borrow ===
First number: 1

=== Part 5: Functions and Borrows ===
Text: hello!, Length was: 5

✅ All parts working!
```

## Hints

<details>
<summary>Hint 1: Understanding the Errors</summary>

When you see:
```
error[E0502]: cannot borrow `s` as mutable because it is also borrowed as immutable
```

This means you have both types of borrows active at the same time. The compiler is preventing a potential bug!

**Solutions**:
1. **Restructure**: Use the immutable references first, then the mutable one
2. **Scope limiting**: Put borrows in explicit blocks `{ ... }`
3. **Remove conflicting borrows**: Only keep what you actually need
</details>

<details>
<summary>Hint 2: Part 1 - Mixed Borrows</summary>

The problem: `r1` and `r2` are still in use when you create `r3`.

**Solution 1**: Don't print `r1` and `r2` after creating `r3`
```rust
let r1 = &s;
let r2 = &s;
println!("r1: {}, r2: {}", r1, r2);  // Use them here
// They're no longer needed after this

let r3 = &mut s;  // Now OK!
println!("r3: {}", r3);
```

**Solution 2**: Clone the data for r3
```rust
let r3 = &mut s.clone();
```

**Solution 3**: Restructure to avoid needing all three
</details>

<details>
<summary>Hint 3: Part 2 - Multiple Mutable Borrows</summary>

You can't have two mutable references to the same data simultaneously. Period.

**Solution**: Use one, then the other (not simultaneously)
```rust
let r1 = &mut s;
println!("r1: {}", r1);
// r1 is done being used

let r2 = &mut s;
println!("r2: {}", r2);
```

Or restructure so you only need one mutable reference.
</details>

<details>
<summary>Hint 4: Part 3 - This One Works!</summary>

Part 3 actually demonstrates correct scoping! The immutable references `r1` and `r2` are no longer used after the first `println!`, so their scope ends. This allows `r3` to be a mutable borrow.

**This code should work as-is** (it's showing the solution pattern for other parts).

If it doesn't compile, check that you're not using `r1` or `r2` after creating `r3`.
</details>

<details>
<summary>Hint 5: Part 4 - Modifying Vector While Borrowed</summary>

The problem: `first` holds a reference into the vector, and `push` might reallocate memory.

**Solution 1**: Use `first` before modifying
```rust
let first = &numbers[0];
println!("First number: {}", first);
// first is no longer used

numbers.push(6);  // Now OK!
```

**Solution 2**: Copy the value instead of holding a reference
```rust
let first = numbers[0];  // Copy the value (integers are Copy)
numbers.push(6);         // OK, no reference exists
println!("First number: {}", first);
```
</details>

<details>
<summary>Hint 6: Part 5 - Functions and Borrows</summary>

The issue here is subtle: if you're using `len` after calling `append_exclamation`, you need to ensure the immutable borrow from `get_length` has ended.

**Solution**: The immutable borrow ends when `get_length` returns. As long as you don't use `len` before calling `append_exclamation`, you're fine. If needed, restructure:

```rust
let len = get_length(&text);
// Use len here if needed
append_exclamation(&mut text);  // OK, len's borrow is done
println!("Text: {}, Length was: {}", text, len);
```

Actually, this should work as-is! The borrow from `get_length` ends immediately.
</details>

<details>
<summary>Hint 7: Why These Rules Matter</summary>

These rules prevent:
- **Data races**: Two threads can't modify the same data simultaneously
- **Iterator invalidation**: Can't modify a collection while iterating
- **Use-after-free**: Can't hold a reference to freed memory

The compiler is your friend, catching bugs before runtime!
</details>

## Reference Solution

```rust
fn main() {
    // Part 1: Can't have mutable and immutable references simultaneously
    println!("=== Part 1: Mixed Borrows ===");
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);
    // r1 and r2 are no longer used after this point

    // FIX: Create r3 after r1 and r2 are done
    let r3 = &mut s;
    println!("r3: {}", r3);

    // Part 2: Can't have multiple mutable references simultaneously
    println!("\n=== Part 2: Multiple Mutable Borrows ===");
    let mut s = String::from("hello");

    // FIX: Use one mutable reference, then the other (not at the same time)
    let r1 = &mut s;
    println!("r1: {}", r1);
    // r1 is done

    let r2 = &mut s;
    println!("r2: {}", r2);

    // Part 3: Scope-based borrowing (already correct!)
    println!("\n=== Part 3: Reference Scopes ===");
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);
    // r1 and r2 are no longer used after this point

    let r3 = &mut s;  // OK - immutable refs are done
    r3.push_str(" world");
    println!("r3: {}", r3);

    // Part 4: Modifying while borrowed
    println!("\n=== Part 4: Modification During Borrow ===");
    let mut numbers = vec![1, 2, 3, 4, 5];

    // FIX: Copy the value instead of holding a reference
    let first = numbers[0];  // Copy (i32 is Copy)
    numbers.push(6);         // OK now
    println!("First number: {}", first);

    // Alternative fix: use the reference before modifying
    // let first = &numbers[0];
    // println!("First number: {}", first);
    // numbers.push(6);

    // Part 5: Functions and borrowing (already correct!)
    println!("\n=== Part 5: Functions and Borrows ===");
    let mut text = String::from("hello");

    let len = get_length(&text);
    append_exclamation(&mut text);  // OK - immutable borrow from get_length ended
    println!("Text: {}, Length was: {}", text, len);

    println!("\n✅ All parts working!");
}

fn get_length(s: &String) -> usize {
    s.len()
}

fn append_exclamation(s: &mut String) {
    s.push('!');
}
```

### Alternative Solution with Explicit Scopes

```rust
fn main() {
    // Part 1: Using explicit scopes to limit borrow lifetimes
    println!("=== Part 1: Mixed Borrows ===");
    let mut s = String::from("hello");

    {
        let r1 = &s;
        let r2 = &s;
        println!("r1: {}, r2: {}", r1, r2);
    }  // r1 and r2 go out of scope here

    {
        let r3 = &mut s;
        println!("r3: {}", r3);
    }  // r3 goes out of scope here

    // ... rest similar with explicit scopes
}
```

## Test Cases

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sequential_immutable_then_mutable() {
        let mut s = String::from("hello");

        let r1 = &s;
        let r2 = &s;
        assert_eq!(r1, "hello");
        assert_eq!(r2, "hello");
        // r1 and r2 no longer used

        let r3 = &mut s;
        r3.push_str(" world");
        assert_eq!(r3, "hello world");
    }

    #[test]
    fn test_sequential_mutable_borrows() {
        let mut s = String::from("hello");

        {
            let r1 = &mut s;
            r1.push('!');
            assert_eq!(r1, "hello!");
        }

        {
            let r2 = &mut s;
            r2.push('?');
            assert_eq!(r2, "hello!?");
        }

        assert_eq!(s, "hello!?");
    }

    #[test]
    fn test_copy_vs_borrow_vec() {
        let mut numbers = vec![1, 2, 3];

        // Copy the value
        let first = numbers[0];
        numbers.push(4);  // OK because we copied, not borrowed
        assert_eq!(first, 1);
        assert_eq!(numbers, vec![1, 2, 3, 4]);
    }

    #[test]
    fn test_function_borrow_ends() {
        let mut text = String::from("test");

        let len = get_length(&text);  // Borrow
        // Borrow ends when function returns

        append_exclamation(&mut text);  // Mutable borrow OK
        assert_eq!(text, "test!");
        assert_eq!(len, 4);
    }

    #[test]
    fn test_multiple_immutable_borrows_ok() {
        let s = String::from("hello");

        let r1 = &s;
        let r2 = &s;
        let r3 = &s;
        let r4 = &s;

        assert_eq!(r1, "hello");
        assert_eq!(r2, "hello");
        assert_eq!(r3, "hello");
        assert_eq!(r4, "hello");
    }

    #[test]
    fn test_reference_scope_ends_at_last_use() {
        let mut s = String::from("hello");

        let r1 = &s;
        println!("{}", r1);  // Last use of r1
        // r1's scope ends here (non-lexical lifetimes)

        let r2 = &mut s;  // OK
        r2.push('!');
        assert_eq!(r2, "hello!");
    }

    #[test]
    fn test_nested_scopes() {
        let mut s = String::from("hello");

        {
            let _r1 = &s;
            let _r2 = &s;
            // Immutable borrows end at scope exit
        }

        {
            let r3 = &mut s;
            r3.push('!');
            // Mutable borrow ends at scope exit
        }

        assert_eq!(s, "hello!");
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

1. **Understanding of borrowing rules (20 points)**
   - Demonstrates understanding of when borrows conflict
   - Recognizes reference scope boundaries
   - Knows the difference between Copy and borrow for vectors

2. **Solution approach (10 points)**
   - Uses appropriate fix for each situation
   - Doesn't just delete code to make it compile
   - Maintains intended functionality

3. **Code organization (10 points)**
   - Clean restructuring
   - Appropriate use of scopes (if needed)
   - Comments explaining borrow checker fixes

## Sample AI Feedback Prompt

```
You are evaluating a Rust exercise on the borrowing rules.

Exercise: Fix borrowing violations where mutable and immutable references conflict.

Student's submission:
{student_code}

Test results:
{test_results}

Previous attempts (if any):
{previous_submissions}

This is a HARD exercise. The borrowing rules are Rust's most challenging concept.

Evaluate the solution on:
1. Understanding of when mutable and immutable borrows conflict
2. Appropriate solutions (restructuring, scoping, copying vs borrowing)
3. Maintaining functionality (not just deleting code to compile)
4. Recognition of reference scope boundaries

Provide feedback as JSON:
{
  "score": 0-40,
  "feedback": "Detailed markdown feedback. Focus on:\n- If student removes functionality: suggest restructuring instead\n- If student doesn't understand scopes: explain non-lexical lifetimes\n- If student over-uses cloning: explain that copying (for Copy types) or restructuring is better\n- If student gets it right: Celebrate! This is hard! Explain why these rules prevent bugs\n- Emphasize that the borrow checker is helping, not hindering",
  "keyTakeaways": ["Mutable XOR immutable references", "Reference scope ends at last use", "Why these rules prevent bugs"],
  "encouragement": "The borrow checker is tough but it's preventing real bugs! This understanding is the hardest part of learning Rust."
}

Common patterns:
- Student deletes prints to avoid using references → Suggest restructuring
- Student clones everything → For Copy types, explain copying vs cloning
- Student doesn't understand NLL → Explain scope ends at last use, not block end
- Student frustrated → Acknowledge difficulty, explain this is Rust's core safety feature
- Student succeeds → Major milestone! The rest of Rust is easier than this.
```

## Expected Student Outcomes

After completing this exercise, students should:
- ✅ Understand the "mutable XOR immutable" borrowing rule
- ✅ Recognize when references conflict
- ✅ Know that reference scopes end at their last use (non-lexical lifetimes)
- ✅ Choose between copying (for Copy types) and restructuring
- ✅ Use explicit scopes `{}` to limit borrow lifetimes when needed
- ✅ Read and understand borrow checker error messages
- ✅ Appreciate WHY these rules exist (preventing data races and bugs)

## Common Mistakes to Watch For

1. **Deleting functionality to make it compile**: Instead of fixing borrows, removing prints
   - Guide: The goal is to keep functionality while satisfying the borrow checker

2. **Not understanding non-lexical lifetimes (NLL)**: Thinking references live until `}`
   - Guide: References' scopes end at their **last use**, not the closing brace

3. **Over-cloning**: Using `.clone()` instead of restructuring or copying
   - Guide: For Copy types like `i32`, just copy. For others, consider restructuring.

4. **Frustration**: "The compiler is being pedantic!"
   - Guide: Acknowledge frustration, explain these rules catch **real bugs** that would crash in other languages

5. **Not using explicit scopes**: Not knowing they can use `{}` to limit lifetimes
   - Guide: Show how explicit scopes can help organize borrows

## Why This Exercise Is Hard

This is the hardest concept in Rust because:
1. **Most languages don't have these rules** - it's a new mental model
2. **The compiler is strict** - code that "looks fine" doesn't compile
3. **Solutions aren't always obvious** - sometimes major restructuring is needed

But mastering this unlocks Rust's superpowers:
- No data races
- No use-after-free
- No iterator invalidation
- Memory safety **without garbage collection**

## Pedagogical Note

After students complete this exercise, they've conquered the hardest part of Rust! The rest of the language (traits, generics, etc.) is conceptually easier.

This is a **major milestone** worth celebrating. The AI should acknowledge this achievement.

## Connection to Real-World Bugs

These rules prevent actual crashes:
```rust
// In C++, this compiles but crashes:
vector<int> v = {1, 2, 3};
int& first = v[0];  // Reference to first element
v.push_back(4);     // Might reallocate - first now points to freed memory!
cout << first;      // 💥 Crash or garbage data

// Rust prevents this at compile time!
```

## Extension Ideas (Post-MVP)

- Interactive visualization of reference lifetimes
- "Debug the borrow checker error" mode with real error messages
- Challenge: Refactor badly-structured code to satisfy borrow checker
- Compare to equivalent C++ code that has undefined behavior

## Node 2 Completion

**Congratulations!** After finishing this exercise, students have completed Node 2: Ownership & Borrowing.

They now understand:
- Ownership and moves
- References and borrowing
- The borrowing rules

This is **the hardest content in the entire Rust learning path**. Everything else builds on this foundation but is conceptually easier.

**Next**: Node 3 (Structs & Enums) - applying ownership to custom types!

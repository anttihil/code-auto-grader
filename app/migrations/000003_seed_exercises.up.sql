-- Migration: 000003_seed_exercises
-- Description: Seed all exercises from Node 1, 2, and 3 with JSONB content

-- ============================================
-- NODE 1: Rust Basics (3 exercises)
-- ============================================

-- Exercise 1.1: Variables and Mutability
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    'Variables and Mutability',
    'variables-mutability',
    'Learn about Rust''s variable binding with `let`, understand the difference between mutable and immutable variables, and recognize compile-time enforcement of mutability rules.',
    'easy',
    20,
    1,
    jsonb_build_object(
        'starter_code', 'fn main() {
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
}',
        'reference_solution', 'fn main() {
    // Fix: Add `mut` to allow reassignment
    let mut x = 5;
    println!("The value of x is: {}", x);
    x = 6;
    println!("The value of x is now: {}", x);

    // This is shadowing (creating a new variable), not mutation
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
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object(
                'name', 'test_compiles',
                'code', '#[test]
fn test_compiles() {
    assert!(true);
}',
                'visible', true
            ),
            jsonb_build_object(
                'name', 'test_x_mutability',
                'code', '#[test]
fn test_x_mutability() {
    let mut x = 5;
    x = 6;
    assert_eq!(x, 6);
}',
                'visible', true
            ),
            jsonb_build_object(
                'name', 'test_y_shadowing',
                'code', '#[test]
fn test_y_shadowing() {
    let y = 10;
    let y = y + 5;
    assert_eq!(y, 15);
}',
                'visible', true
            ),
            jsonb_build_object(
                'name', 'test_spaces_shadowing_type_change',
                'code', '#[test]
fn test_spaces_shadowing_type_change() {
    let spaces = "   ";
    let spaces = spaces.len();
    assert_eq!(spaces, 3);
}',
                'visible', true
            ),
            jsonb_build_object(
                'name', 'test_z_operations',
                'code', '#[test]
fn test_z_operations() {
    let mut z = 1;
    z = 2;
    assert_eq!(z, 2);

    let z = "three";
    assert_eq!(z, "three");
}',
                'visible', true
            )
        ),
        'hints', jsonb_build_array(
            jsonb_build_object(
                'order', 1,
                'title', 'Mutability vs Shadowing',
                'content', 'There are two ways to "change" a value in Rust: **Mutability** with `mut` keyword allows reassignment, while **Shadowing** uses `let` again to create a new variable. Shadowing is useful when you want to change the type or transform a value.'
            ),
            jsonb_build_object(
                'order', 2,
                'title', 'Type Consistency',
                'content', 'Even with `mut`, you cannot change a variable''s type. Once `z` is an integer, it must stay an integer. If you want to change types, you need shadowing with another `let` statement.'
            ),
            jsonb_build_object(
                'order', 3,
                'title', 'Which Lines Have Errors?',
                'content', 'Compilation errors occur on: Line 4 (`x = 6` - x is immutable) and Line 18 (`z = "three"` - type mismatch). The `y` and `spaces` examples use shadowing, which is perfectly valid!'
            )
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object(
                    'name', 'Correct use of mut',
                    'points', 10,
                    'description', 'Uses `mut` only when reassigning variables, doesn''t use `mut` for shadowed variables'
                ),
                jsonb_build_object(
                    'name', 'Understanding of shadowing',
                    'points', 10,
                    'description', 'Recognizes when shadowing is appropriate (type changes) and doesn''t unnecessarily shadow instead of using `mut`'
                ),
                jsonb_build_object(
                    'name', 'Code style',
                    'points', 10,
                    'description', 'Follows Rust naming conventions, no unnecessary modifications to working code, clean readable solution'
                ),
                jsonb_build_object(
                    'name', 'Demonstrates understanding',
                    'points', 10,
                    'description', 'Comments show grasp of mutability vs shadowing, solution isn''t just trial-and-error'
                )
            )
        ),
        'ai_evaluation_prompt', 'You are evaluating a Rust exercise on variables and mutability. The student was asked to fix compilation errors related to variable mutability and shadowing. Evaluate the solution on: 1) Correct use of `mut` keyword, 2) Understanding of shadowing vs mutation, 3) Code style and clarity. Provide constructive feedback explaining what was done well and what could improve.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object(
                'pattern', 'Adding `mut` everywhere',
                'guidance', 'Students might add `mut` to `y` and `spaces`, not recognizing shadowing. Explain that shadowing creates a new variable and doesn''t need `mut`.'
            ),
            jsonb_build_object(
                'pattern', 'Trying to change types with `mut`',
                'guidance', 'Attempting `mut z` and reassigning to a string. Explain that mutation cannot change types - only shadowing can.'
            ),
            jsonb_build_object(
                'pattern', 'Confusion about shadowing',
                'guidance', 'Not understanding that `let y = y + 5` creates a new variable. Clarify the difference between mutation and shadowing.'
            )
        ),
        'expected_outcomes', jsonb_build_array(
            'Understand that Rust variables are immutable by default',
            'Know when to use `mut` vs shadowing',
            'Recognize that shadowing can change types, mutation cannot',
            'Be comfortable with Rust''s strict compile-time checks'
        )
    )
);

-- Exercise 1.2: Functions and Return Values
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000001',
    'Functions and Return Values',
    'functions-return-values',
    'Define functions with parameters and return types, understand Rust''s expression-based return syntax, and use type annotations for function signatures.',
    'easy',
    25,
    2,
    jsonb_build_object(
        'starter_code', '// TODO: Write a function called `add` that takes two i32 parameters
// and returns their sum as an i32

// TODO: Write a function called `rectangle_area` that takes width and height
// as f64 parameters and returns the area as f64

// TODO: Write a function called `is_even` that takes an i32
// and returns true if it''s even, false otherwise
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
}',
        'reference_solution', 'fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn rectangle_area(width: f64, height: f64) -> f64 {
    width * height
}

fn is_even(n: i32) -> bool {
    n % 2 == 0
}

fn main() {
    let sum = add(5, 7);
    println!("5 + 7 = {}", sum);
    assert_eq!(sum, 12);

    let sum2 = add(-3, 3);
    println!("-3 + 3 = {}", sum2);
    assert_eq!(sum2, 0);

    let area = rectangle_area(10.5, 20.0);
    println!("Rectangle area: {}", area);
    assert_eq!(area, 210.0);

    let area2 = rectangle_area(0.0, 10.0);
    println!("Rectangle area (zero width): {}", area2);
    assert_eq!(area2, 0.0);

    println!("Is 4 even? {}", is_even(4));
    assert_eq!(is_even(4), true);

    println!("Is 7 even? {}", is_even(7));
    assert_eq!(is_even(7), false);

    println!("Is 0 even? {}", is_even(0));
    assert_eq!(is_even(0), true);

    println!("All tests passed!");
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_add_positive_numbers', 'code', '#[test]\nfn test_add_positive_numbers() {\n    assert_eq!(add(5, 7), 12);\n    assert_eq!(add(100, 200), 300);\n}', 'visible', true),
            jsonb_build_object('name', 'test_add_with_zero', 'code', '#[test]\nfn test_add_with_zero() {\n    assert_eq!(add(5, 0), 5);\n    assert_eq!(add(0, 10), 10);\n}', 'visible', true),
            jsonb_build_object('name', 'test_add_negative_numbers', 'code', '#[test]\nfn test_add_negative_numbers() {\n    assert_eq!(add(-3, 3), 0);\n    assert_eq!(add(-5, -10), -15);\n}', 'visible', true),
            jsonb_build_object('name', 'test_rectangle_area', 'code', '#[test]\nfn test_rectangle_area() {\n    assert_eq!(rectangle_area(10.5, 20.0), 210.0);\n    assert_eq!(rectangle_area(5.0, 5.0), 25.0);\n}', 'visible', true),
            jsonb_build_object('name', 'test_is_even', 'code', '#[test]\nfn test_is_even() {\n    assert_eq!(is_even(0), true);\n    assert_eq!(is_even(2), true);\n    assert_eq!(is_even(1), false);\n    assert_eq!(is_even(7), false);\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Function Syntax', 'content', 'Rust function syntax: `fn function_name(param1: Type1, param2: Type2) -> ReturnType { expression_without_semicolon }`. The last line without a semicolon is automatically returned.'),
            jsonb_build_object('order', 2, 'title', 'Modulo Operator', 'content', 'To check if a number is even, use the modulo operator `%`: `n % 2 == 0` means n is even.'),
            jsonb_build_object('order', 3, 'title', 'Implicit Returns', 'content', 'In Rust, you can return implicitly (without `return` keyword) by omitting the semicolon on the last expression. This is preferred in idiomatic Rust.')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Correct function signatures', 'points', 15, 'description', 'Parameter types and return types are correct, function names follow snake_case'),
                jsonb_build_object('name', 'Idiomatic return syntax', 'points', 10, 'description', 'Uses expression-based returns without unnecessary return keyword'),
                jsonb_build_object('name', 'Correctness', 'points', 10, 'description', 'Functions implement correct logic and handle edge cases'),
                jsonb_build_object('name', 'Code quality', 'points', 5, 'description', 'No unnecessary complexity, efficient implementations')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this Rust exercise on functions. Check: 1) Correct function signatures and types, 2) Use of idiomatic expression-based returns, 3) Correctness and edge case handling. If student uses explicit `return` statements, acknowledge correctness but suggest implicit returns as more idiomatic.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Missing return types', 'guidance', 'Forgetting the `-> Type` syntax. Remind student that Rust requires explicit return types.'),
            jsonb_build_object('pattern', 'Adding semicolons', 'guidance', 'Putting `;` after the final expression turns it into a statement. Explain that implicit returns need no semicolon.'),
            jsonb_build_object('pattern', 'Wrong types', 'guidance', 'Using `i32` instead of `f64` for rectangle_area. Point out the parameter type requirements.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Write function signatures with type annotations',
            'Understand implicit returns (expression-based)',
            'Use basic operators (arithmetic, modulo, comparison)',
            'Read and write idiomatic Rust functions'
        )
    )
);

-- Exercise 1.3: Control Flow and Pattern Matching
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000003',
    '20000000-0000-0000-0000-000000000001',
    'Control Flow and Pattern Matching',
    'control-flow-pattern-matching',
    'Use if/else expressions for conditional logic, apply match for pattern matching, and understand that control flow constructs are expressions in Rust.',
    'medium',
    35,
    3,
    jsonb_build_object(
        'starter_code', 'fn categorize_number(n: i32) -> String {
    // TODO: Implement this function
    // 1. Use `match` to check if n is negative, zero, or positive
    // 2. For non-zero, use if/else to check magnitude
    // 3. Return a formatted string

    todo!()
}

fn main() {
    println!("{}", categorize_number(0));      // "zero"
    println!("{}", categorize_number(5));      // "positive (small)"
    println!("{}", categorize_number(-5));     // "negative (small)"
    println!("{}", categorize_number(50));     // "positive (medium)"
    println!("{}", categorize_number(-75));    // "negative (medium)"
    println!("{}", categorize_number(150));    // "positive (large)"
    println!("{}", categorize_number(-200));   // "negative (large)"
}',
        'reference_solution', 'fn categorize_number(n: i32) -> String {
    match n {
        0 => "zero".to_string(),
        n if n > 0 => {
            let magnitude = get_magnitude(n);
            format!("positive ({})", magnitude)
        }
        n => {
            let magnitude = get_magnitude(n.abs());
            format!("negative ({})", magnitude)
        }
    }
}

fn get_magnitude(abs_value: i32) -> &''static str {
    if abs_value <= 10 {
        "small"
    } else if abs_value <= 100 {
        "medium"
    } else {
        "large"
    }
}

fn main() {
    println!("{}", categorize_number(0));
    println!("{}", categorize_number(5));
    println!("{}", categorize_number(-5));
    println!("{}", categorize_number(50));
    println!("{}", categorize_number(-75));
    println!("{}", categorize_number(150));
    println!("{}", categorize_number(-200));
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_zero', 'code', '#[test]\nfn test_zero() {\n    assert_eq!(categorize_number(0), "zero");\n}', 'visible', true),
            jsonb_build_object('name', 'test_positive_small', 'code', '#[test]\nfn test_positive_small() {\n    assert_eq!(categorize_number(5), "positive (small)");\n    assert_eq!(categorize_number(10), "positive (small)");\n}', 'visible', true),
            jsonb_build_object('name', 'test_positive_medium', 'code', '#[test]\nfn test_positive_medium() {\n    assert_eq!(categorize_number(50), "positive (medium)");\n    assert_eq!(categorize_number(100), "positive (medium)");\n}', 'visible', true),
            jsonb_build_object('name', 'test_positive_large', 'code', '#[test]\nfn test_positive_large() {\n    assert_eq!(categorize_number(150), "positive (large)");\n}', 'visible', true),
            jsonb_build_object('name', 'test_negative_small', 'code', '#[test]\nfn test_negative_small() {\n    assert_eq!(categorize_number(-5), "negative (small)");\n}', 'visible', true),
            jsonb_build_object('name', 'test_negative_medium', 'code', '#[test]\nfn test_negative_medium() {\n    assert_eq!(categorize_number(-75), "negative (medium)");\n}', 'visible', true),
            jsonb_build_object('name', 'test_negative_large', 'code', '#[test]\nfn test_negative_large() {\n    assert_eq!(categorize_number(-200), "negative (large)");\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Match Structure', 'content', 'Use `match` with guards: `match n { 0 => ..., n if n > 0 => ..., _ => ... }`. Match arms can use guards with `if` conditions.'),
            jsonb_build_object('order', 2, 'title', 'Calculating Magnitude', 'content', 'Use absolute value: `let abs_value = n.abs()`. Then check: if abs_value <= 10 {"small"} else if abs_value <= 100 {"medium"} else {"large"}'),
            jsonb_build_object('order', 3, 'title', 'Building Strings', 'content', 'Use `format!` macro: `format!("positive ({})", magnitude)` or `.to_string()` for simple conversions.')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Use of match expression', 'points', 15, 'description', 'Uses match for sign checking, handles all cases, match is used as an expression'),
                jsonb_build_object('name', 'Control flow correctness', 'points', 10, 'description', 'Magnitude categorization logic is correct with proper boundary values'),
                jsonb_build_object('name', 'Code organization', 'points', 10, 'description', 'Clean readable structure, no unnecessary code duplication'),
                jsonb_build_object('name', 'Rust idioms', 'points', 5, 'description', 'Expression-based returns, appropriate use of format!')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this control flow exercise. Check: 1) Proper use of match expression (not just nested if/else), 2) Correctness of magnitude logic and boundary cases, 3) Code organization. If student uses only if/else, suggest match is more idiomatic.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Using only if/else, no match', 'guidance', 'While it works, match is more idiomatic for this case. Suggest using match with guards.'),
            jsonb_build_object('pattern', 'Wrong boundaries', 'guidance', 'Treating 10 as medium (should be small) or 100 as large (should be medium). Clarify the boundary conditions.'),
            jsonb_build_object('pattern', 'Forgetting .abs() for negatives', 'guidance', 'Leads to all negatives being "large". Explain the need for absolute value.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Use match expressions for pattern matching',
            'Understand match guards (n if n > 0)',
            'Combine match and if/else effectively',
            'Handle boundary cases correctly',
            'Format strings with format! macro'
        )
    )
);

-- ============================================
-- NODE 2: Ownership & Borrowing (3 exercises)
-- ============================================

-- Exercise 2.1: Ownership and Moves
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000004',
    '20000000-0000-0000-0000-000000000002',
    'Ownership and Moves',
    'ownership-moves',
    'Understand Rust''s ownership model where every value has one owner, recognize when ownership is transferred (moved), and distinguish between Copy types and Move types.',
    'medium',
    35,
    1,
    jsonb_build_object(
        'starter_code', 'fn main() {
    // Part 1: Basic move semantics
    println!("=== Part 1: String Moves ===");
    let s1 = String::from("hello");
    let s2 = s1;  // s1''s ownership moves to s2

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
}

fn add_exclamation(mut s: String) -> String {
    s.push_str("!");
    s
}',
        'reference_solution', 'fn main() {
    // Part 1: Basic move semantics
    println!("=== Part 1: String Moves ===");
    let s1 = String::from("hello");
    let s2 = s1.clone();  // FIX: Clone instead of move

    println!("s1: {}", s1);
    println!("s2: {}", s2);

    // Part 2: Function call moves
    println!("\n=== Part 2: Function Moves ===");
    let message = String::from("Hello, Rust!");
    print_message(message.clone());  // FIX: Clone before passing
    println!("Message again: {}", message);

    // Part 3: Return values and ownership (already works!)
    println!("\n=== Part 3: Returning Ownership ===");
    let greeting = String::from("Hi");
    let greeting = add_exclamation(greeting);
    println!("{}", greeting);

    // Part 4: Multiple uses with clone
    println!("\n=== Part 4: Cloning ===");
    let original = String::from("original");
    let copy = original.clone();  // FIX: Clone to create independent copy
    println!("Original: {}", original);
    println!("Copy: {}", copy);

    // Part 5: Copy types vs Move types (already works!)
    println!("\n=== Part 5: Copy Types ===");
    let x = 42;
    let y = x;
    println!("x: {}, y: {}", x, y);

    println!("\n✅ All parts working!");
}

fn print_message(msg: String) {
    println!("Message: {}", msg);
}

fn add_exclamation(mut s: String) -> String {
    s.push_str("!");
    s
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_string_cloning', 'code', '#[test]\nfn test_string_cloning() {\n    let s1 = String::from("hello");\n    let s2 = s1.clone();\n    assert_eq!(s1, "hello");\n    assert_eq!(s2, "hello");\n}', 'visible', true),
            jsonb_build_object('name', 'test_copy_types', 'code', '#[test]\nfn test_copy_types() {\n    let x = 42;\n    let y = x;\n    assert_eq!(x, 42);\n    assert_eq!(y, 42);\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Understanding Errors', 'content', 'When you see "borrow of moved value", you have options: 1) Don''t use the moved value, 2) Clone before moving, 3) Borrow instead of move (covered in next exercise).'),
            jsonb_build_object('order', 2, 'title', 'Cloning', 'content', 'Use `.clone()` to create a deep copy: `let copy = original.clone();`. This creates a new String with the same contents.')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Understanding of move semantics', 'points', 15, 'description', 'Recognizes when moves happen and understands that moved values become invalid'),
                jsonb_build_object('name', 'Appropriate use of .clone()', 'points', 10, 'description', 'Uses clone when truly needing multiple owners'),
                jsonb_build_object('name', 'Code organization', 'points', 10, 'description', 'Clear readable fixes with good comments'),
                jsonb_build_object('name', 'Conceptual understanding', 'points', 5, 'description', 'Comments show grasp of ownership rules')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this ownership exercise. Check: 1) Understanding of when moves happen, 2) Appropriate use of .clone(), 3) Recognition of Copy vs Move types. If student clones everything, acknowledge it works but mention we''ll learn borrowing as a more efficient approach.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Cloning everything', 'guidance', 'Works but inefficient. Prepare student for borrowing as a better alternative.'),
            jsonb_build_object('pattern', 'Not understanding Copy types', 'guidance', 'Explain stack vs heap allocation and why integers copy automatically.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Understand ownership rules (one owner, transfer on move)',
            'Recognize when ownership moves happen',
            'Know the difference between Copy and Move types',
            'Use .clone() to create independent copies'
        )
    )
);

-- Exercise 2.2: References and Borrowing
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000005',
    '20000000-0000-0000-0000-000000000002',
    'References and Borrowing',
    'references-borrowing',
    'Understand references (&T) as borrowing without taking ownership, use immutable and mutable references, and recognize that borrowing is more efficient than cloning.',
    'medium',
    40,
    2,
    jsonb_build_object(
        'starter_code', '// TODO: Implement this function to take a REFERENCE to String
// and return its length without taking ownership
fn calculate_length(/* add parameter */) -> usize {
    todo!()
}

// TODO: Implement this function to take a reference to String
// and return the length of the first word
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
fn capitalize_first_char(/* add parameter */) {
    todo!()
}

fn main() {
    let s1 = String::from("hello world");
    let len = calculate_length(&s1);
    println!("The length of ''{}'' is {}.", s1, len);

    let s2 = String::from("hello world");
    let len = first_word_length(&s2);
    println!("First word length of ''{}'' is {}", s2, len);

    let mut s3 = String::from("hello");
    append_exclamation(&mut s3);
    println!("After append: {}", s3);

    let mut s4 = String::from("rust");
    capitalize_first_char(&mut s4);
    println!("After capitalize: {}", s4);

    println!("\n✅ All tests passed!");
}',
        'reference_solution', 'fn calculate_length(s: &String) -> usize {
    s.len()
}

fn first_word_length(s: &String) -> usize {
    for (i, ch) in s.chars().enumerate() {
        if ch == '' '' {
            return i;
        }
    }
    s.len()
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
    let s1 = String::from("hello world");
    let len = calculate_length(&s1);
    println!("The length of ''{}'' is {}.", s1, len);

    let s2 = String::from("hello world");
    let len = first_word_length(&s2);
    println!("First word length of ''{}'' is {}", s2, len);

    let mut s3 = String::from("hello");
    append_exclamation(&mut s3);
    println!("After append: {}", s3);

    let mut s4 = String::from("rust");
    capitalize_first_char(&mut s4);
    println!("After capitalize: {}", s4);

    println!("\n✅ All tests passed!");
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_calculate_length', 'code', '#[test]\nfn test_calculate_length() {\n    let s = String::from("hello");\n    assert_eq!(calculate_length(&s), 5);\n    assert_eq!(s, "hello"); // s still valid\n}', 'visible', true),
            jsonb_build_object('name', 'test_first_word_length', 'code', '#[test]\nfn test_first_word_length() {\n    assert_eq!(first_word_length(&String::from("hello world")), 5);\n    assert_eq!(first_word_length(&String::from("hello")), 5);\n}', 'visible', true),
            jsonb_build_object('name', 'test_append_exclamation', 'code', '#[test]\nfn test_append_exclamation() {\n    let mut s = String::from("hello");\n    append_exclamation(&mut s);\n    assert_eq!(s, "hello!");\n}', 'visible', true),
            jsonb_build_object('name', 'test_capitalize_first_char', 'code', '#[test]\nfn test_capitalize_first_char() {\n    let mut s = String::from("hello");\n    capitalize_first_char(&mut s);\n    assert_eq!(s, "Hello");\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Reference Syntax', 'content', 'Immutable reference: `fn function_name(parameter: &String)`. Mutable reference: `fn function_name(parameter: &mut String)`.'),
            jsonb_build_object('order', 2, 'title', 'Benefits of Borrowing', 'content', 'Borrowing is more efficient than cloning - no copy, just a reference. The caller keeps ownership!')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Correct use of references', 'points', 20, 'description', 'Uses &String and &mut String appropriately'),
                jsonb_build_object('name', 'Implementation correctness', 'points', 10, 'description', 'Functions work correctly for all edge cases'),
                jsonb_build_object('name', 'Code quality', 'points', 10, 'description', 'Clean readable code following Rust conventions')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this borrowing exercise. Check: 1) Correct use of & and &mut, 2) Understanding of when to use immutable vs mutable references, 3) Implementation correctness. If student takes ownership instead of borrowing, explain the benefit of references.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Taking ownership instead of borrowing', 'guidance', 'Should use &String, not String. Explain that borrowing lets the caller keep using the value.'),
            jsonb_build_object('pattern', 'Using &mut when & would work', 'guidance', 'Explain principle of least privilege - use immutable unless you need to modify.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Use &T references to borrow data without taking ownership',
            'Use &mut T to modify data in place',
            'Understand that borrowing is more efficient than cloning',
            'Create multiple immutable references to the same data'
        )
    )
);

-- Exercise 2.3: The Borrowing Rules
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000006',
    '20000000-0000-0000-0000-000000000002',
    'The Borrowing Rules',
    'borrowing-rules',
    'Understand Rust''s borrowing rules (one mutable reference OR any number of immutable references), recognize when references conflict, and restructure code to satisfy the borrow checker.',
    'hard',
    50,
    3,
    jsonb_build_object(
        'starter_code', 'fn main() {
    // Part 1: Can''t have mutable and immutable references simultaneously
    println!("=== Part 1: Mixed Borrows ===");
    let mut s = String::from("hello");

    let r1 = &s;       // Immutable borrow
    let r2 = &s;       // Another immutable borrow - OK!
    let r3 = &mut s;   // ERROR: Can''t borrow as mutable while immutable borrows exist

    println!("r1: {}, r2: {}, r3: {}", r1, r2, r3);

    // Part 2: Can''t have multiple mutable references simultaneously
    println!("\n=== Part 2: Multiple Mutable Borrows ===");
    let mut s = String::from("hello");

    let r1 = &mut s;
    let r2 = &mut s;  // ERROR: Can''t borrow as mutable more than once

    println!("r1: {}, r2: {}", r1, r2);

    // Part 3: Scope-based borrowing (should work!)
    println!("\n=== Part 3: Reference Scopes ===");
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);

    let r3 = &mut s;  // OK - immutable refs done
    r3.push_str(" world");
    println!("r3: {}", r3);

    // Part 4: Modifying while borrowed
    println!("\n=== Part 4: Modification During Borrow ===");
    let mut numbers = vec![1, 2, 3, 4, 5];

    let first = &numbers[0];  // Immutable borrow
    numbers.push(6);          // ERROR: Can''t mutate while borrowed
    println!("First number: {}", first);

    println!("\n✅ All parts working!");
}

fn get_length(s: &String) -> usize {
    s.len()
}

fn append_exclamation(s: &mut String) {
    s.push(''!'');
}',
        'reference_solution', 'fn main() {
    // Part 1: Fixed by using references sequentially
    println!("=== Part 1: Mixed Borrows ===");
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);
    // r1 and r2 no longer used

    let r3 = &mut s;  // OK now
    println!("r3: {}", r3);

    // Part 2: Use mutable references sequentially
    println!("\n=== Part 2: Multiple Mutable Borrows ===");
    let mut s = String::from("hello");

    let r1 = &mut s;
    println!("r1: {}", r1);
    // r1 done

    let r2 = &mut s;
    println!("r2: {}", r2);

    // Part 3: Already works! (demonstrates proper scoping)
    println!("\n=== Part 3: Reference Scopes ===");
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("r1: {}, r2: {}", r1, r2);

    let r3 = &mut s;
    r3.push_str(" world");
    println!("r3: {}", r3);

    // Part 4: Copy the value instead of holding reference
    println!("\n=== Part 4: Modification During Borrow ===");
    let mut numbers = vec![1, 2, 3, 4, 5];

    let first = numbers[0];  // Copy (i32 is Copy)
    numbers.push(6);         // OK now
    println!("First number: {}", first);

    println!("\n✅ All parts working!");
}

fn get_length(s: &String) -> usize {
    s.len()
}

fn append_exclamation(s: &mut String) {
    s.push(''!'');
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_sequential_borrows', 'code', '#[test]\nfn test_sequential_borrows() {\n    let mut s = String::from("hello");\n    let r1 = &s;\n    let r2 = &s;\n    assert_eq!(r1, "hello");\n    assert_eq!(r2, "hello");\n    // References done\n    let r3 = &mut s;\n    r3.push_str(" world");\n    assert_eq!(r3, "hello world");\n}', 'visible', true),
            jsonb_build_object('name', 'test_copy_vs_borrow', 'code', '#[test]\nfn test_copy_vs_borrow() {\n    let mut numbers = vec![1, 2, 3];\n    let first = numbers[0];  // Copy\n    numbers.push(4);         // OK\n    assert_eq!(first, 1);\n    assert_eq!(numbers, vec![1, 2, 3, 4]);\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Understanding Errors', 'content', 'When you see "cannot borrow as mutable because it is also borrowed as immutable", you have both types of borrows active. Solutions: 1) Restructure to use immutable refs first, then mutable, 2) Use explicit scopes, 3) Remove conflicting borrows.'),
            jsonb_build_object('order', 2, 'title', 'Reference Scopes', 'content', 'Reference scopes end at their last use, not the closing brace (non-lexical lifetimes). Use this to your advantage!')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Understanding of borrowing rules', 'points', 20, 'description', 'Demonstrates understanding of when borrows conflict and reference scope boundaries'),
                jsonb_build_object('name', 'Solution approach', 'points', 10, 'description', 'Uses appropriate fix for each situation, maintains intended functionality'),
                jsonb_build_object('name', 'Code organization', 'points', 10, 'description', 'Clean restructuring with explanatory comments')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this borrowing rules exercise. This is HARD - the borrowing rules are Rust''s most challenging concept. Check: 1) Understanding of when mutable and immutable borrows conflict, 2) Appropriate solutions (restructuring, scoping, copying), 3) Maintains functionality. Acknowledge the difficulty and celebrate success!',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Deleting functionality', 'guidance', 'Instead of fixing borrows, removing prints. Suggest the goal is to keep functionality while satisfying borrow checker.'),
            jsonb_build_object('pattern', 'Not understanding NLL', 'guidance', 'References'' scopes end at their last use, not the closing brace. Explain non-lexical lifetimes.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Understand the "mutable XOR immutable" borrowing rule',
            'Recognize when references conflict',
            'Know that reference scopes end at their last use',
            'Choose between copying (for Copy types) and restructuring',
            'Appreciate WHY these rules exist (preventing data races)'
        )
    )
);

-- ============================================
-- NODE 3: Structs & Enums (3 exercises)
-- ============================================

-- Exercise 3.1: Defining and Using Structs
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000007',
    '20000000-0000-0000-0000-000000000003',
    'Defining and Using Structs',
    'structs-basics',
    'Define custom structs with named fields, instantiate structs, implement associated functions (constructors), and understand that structs own their data.',
    'easy',
    30,
    1,
    jsonb_build_object(
        'starter_code', '// TODO: Define a Book struct with:
// - title: String
// - author: String
// - pages: u32
// - available: bool

impl Book {
    // TODO: Implement new(title, author, pages) -> Book
    // Set available to true by default
}

// TODO: Define a Library struct with:
// - name: String
// - books: Vec<Book>

impl Library {
    // TODO: Implement new(name) -> Library
    // Initialize with empty books vector

    // TODO: Implement count_available_books(&self) -> usize
    fn count_available_books(&self) -> usize {
        todo!()
    }
}

fn main() {
    let book1 = Book::new(
        String::from("The Rust Programming Language"),
        String::from("Steve Klabnik and Carol Nichols"),
        552,
    );

    assert_eq!(book1.title, "The Rust Programming Language");
    assert_eq!(book1.available, true);
    println!("✓ Book created successfully");

    let mut library = Library::new(String::from("City Library"));
    assert_eq!(library.books.len(), 0);
    println!("✓ Library created successfully");

    library.books.push(book1);
    let mut book2 = Book::new(
        String::from("Rust for Rustaceans"),
        String::from("Jon Gjengset"),
        282,
    );
    book2.available = false;
    library.books.push(book2);

    assert_eq!(library.count_available_books(), 1);
    println!("✓ Available books counted correctly");

    println!("\n✅ All tests passed!");
}',
        'reference_solution', 'struct Book {
    title: String,
    author: String,
    pages: u32,
    available: bool,
}

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

struct Library {
    name: String,
    books: Vec<Book>,
}

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
    let book1 = Book::new(
        String::from("The Rust Programming Language"),
        String::from("Steve Klabnik and Carol Nichols"),
        552,
    );

    assert_eq!(book1.title, "The Rust Programming Language");
    assert_eq!(book1.available, true);
    println!("✓ Book created successfully");

    let mut library = Library::new(String::from("City Library"));
    assert_eq!(library.books.len(), 0);
    println!("✓ Library created successfully");

    library.books.push(book1);
    let mut book2 = Book::new(
        String::from("Rust for Rustaceans"),
        String::from("Jon Gjengset"),
        282,
    );
    book2.available = false;
    library.books.push(book2);

    assert_eq!(library.count_available_books(), 1);
    println!("✓ Available books counted correctly");

    println!("\n✅ All tests passed!");
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_book_creation', 'code', '#[test]\nfn test_book_creation() {\n    let book = Book::new(String::from("Test"), String::from("Author"), 100);\n    assert_eq!(book.title, "Test");\n    assert_eq!(book.available, true);\n}', 'visible', true),
            jsonb_build_object('name', 'test_library_count', 'code', '#[test]\nfn test_library_count() {\n    let mut library = Library::new(String::from("Lib"));\n    library.books.push(Book::new(String::from("B1"), String::from("A1"), 100));\n    let mut b2 = Book::new(String::from("B2"), String::from("A2"), 200);\n    b2.available = false;\n    library.books.push(b2);\n    assert_eq!(library.count_available_books(), 1);\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Struct Definition', 'content', 'Syntax: `struct StructName { field1: Type1, field2: Type2 }`. Use String (not &str) for owned data.'),
            jsonb_build_object('order', 2, 'title', 'Associated Functions', 'content', 'Constructors don''t take self: `fn new(params) -> StructName`. Use field shorthand when names match: `StructName { field1, field2 }`.')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Correct struct definitions', 'points', 15, 'description', 'Fields have correct names and types'),
                jsonb_build_object('name', 'Proper associated functions', 'points', 10, 'description', 'Constructors follow new convention and don''t take self'),
                jsonb_build_object('name', 'Implementation correctness', 'points', 10, 'description', 'count_available_books correctly counts and handles edge cases'),
                jsonb_build_object('name', 'Code quality', 'points', 5, 'description', 'Idiomatic Rust, proper use of &self')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this struct exercise. Check: 1) Struct definitions with correct field types, 2) Associated functions vs instance methods, 3) Correct implementation of count. If using verbose loop, suggest iterator style as more idiomatic.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Using &str instead of String', 'guidance', 'Structs should own their data - use String, not &str.'),
            jsonb_build_object('pattern', 'Constructor takes self', 'guidance', 'Constructors don''t have a self yet! They''re associated functions.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Define custom structs with named fields',
            'Understand that structs own their data',
            'Create associated functions (constructors)',
            'Distinguish between associated functions and instance methods',
            'Work with structs inside collections'
        )
    )
);

-- Exercise 3.2: Methods and Self
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000008',
    '20000000-0000-0000-0000-000000000003',
    'Methods and Self',
    'methods-self',
    'Implement instance methods that operate on struct data, understand &self, &mut self, and self in method signatures, and use method chaining patterns.',
    'medium',
    35,
    2,
    jsonb_build_object(
        'starter_code', '// TODO: Define Rectangle struct with width and height (u32)

impl Rectangle {
    // TODO: Implement new(width, height) -> Rectangle

    // TODO: Implement area(&self) -> u32
    // Returns width * height

    // TODO: Implement perimeter(&self) -> u32
    // Returns 2 * (width + height)

    // TODO: Implement is_square(&self) -> bool
    // Returns true if width == height

    // TODO: Implement can_hold(&self, other: &Rectangle) -> bool
    // Returns true if self can contain other

    // TODO: Implement scale(&mut self, factor: u32)
    // Multiplies both dimensions by factor

    // TODO: Implement with_width(mut self, width: u32) -> Self
    // Builder pattern for method chaining

    // TODO: Implement with_height(mut self, height: u32) -> Self
}

fn main() {
    let rect1 = Rectangle::new(30, 50);
    assert_eq!(rect1.area(), 1500);
    assert_eq!(rect1.perimeter(), 160);
    println!("✓ Basic methods work");

    let mut rect2 = Rectangle::new(10, 20);
    rect2.scale(3);
    assert_eq!(rect2.area(), 1800);
    println!("✓ Scale works");

    let rect3 = Rectangle::new(10, 10)
        .with_width(50)
        .with_height(100);
    assert_eq!(rect3.area(), 5000);
    println!("✓ Builder pattern works");

    println!("\n✅ All tests passed!");
}',
        'reference_solution', 'struct Rectangle {
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
    let rect1 = Rectangle::new(30, 50);
    assert_eq!(rect1.area(), 1500);
    assert_eq!(rect1.perimeter(), 160);
    println!("✓ Basic methods work");

    let mut rect2 = Rectangle::new(10, 20);
    rect2.scale(3);
    assert_eq!(rect2.area(), 1800);
    println!("✓ Scale works");

    let rect3 = Rectangle::new(10, 10)
        .with_width(50)
        .with_height(100);
    assert_eq!(rect3.area(), 5000);
    println!("✓ Builder pattern works");

    println!("\n✅ All tests passed!");
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_area', 'code', '#[test]\nfn test_area() {\n    let rect = Rectangle::new(30, 50);\n    assert_eq!(rect.area(), 1500);\n}', 'visible', true),
            jsonb_build_object('name', 'test_perimeter', 'code', '#[test]\nfn test_perimeter() {\n    let rect = Rectangle::new(30, 50);\n    assert_eq!(rect.perimeter(), 160);\n}', 'visible', true),
            jsonb_build_object('name', 'test_is_square', 'code', '#[test]\nfn test_is_square() {\n    let square = Rectangle::new(25, 25);\n    assert_eq!(square.is_square(), true);\n    let rect = Rectangle::new(30, 50);\n    assert_eq!(rect.is_square(), false);\n}', 'visible', true),
            jsonb_build_object('name', 'test_scale', 'code', '#[test]\nfn test_scale() {\n    let mut rect = Rectangle::new(10, 20);\n    rect.scale(3);\n    assert_eq!(rect.width, 30);\n    assert_eq!(rect.height, 60);\n}', 'visible', true),
            jsonb_build_object('name', 'test_builder', 'code', '#[test]\nfn test_builder() {\n    let rect = Rectangle::new(10, 10).with_width(50).with_height(100);\n    assert_eq!(rect.width, 50);\n    assert_eq!(rect.height, 100);\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Self Types', 'content', '&self for reading, &mut self for modifying, self for consuming/transforming. 99% of methods use &self or &mut self.'),
            jsonb_build_object('order', 2, 'title', 'Builder Pattern', 'content', 'Builder methods take `mut self`, modify fields, and return `self` to enable chaining.')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Correct self parameters', 'points', 20, 'description', 'Uses &self for reading, &mut self for mutation, self for builders'),
                jsonb_build_object('name', 'Implementation correctness', 'points', 10, 'description', 'All calculations and logic correct'),
                jsonb_build_object('name', 'Code quality', 'points', 10, 'description', 'Clean, idiomatic implementations')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this methods exercise. Check: 1) Correct self parameter types, 2) Implementation correctness, 3) Understanding of when to use each self type. If scale uses &self, explain it should be &mut self.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Wrong self type for scale', 'guidance', 'scale modifies the struct, needs &mut self not &self.'),
            jsonb_build_object('pattern', 'Builder pattern errors', 'guidance', 'Builders must take self, modify, and return self for chaining.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Implement methods with appropriate self types',
            'Understand when to use &self, &mut self, and self',
            'Use the builder pattern for method chaining',
            'Apply methods to solve practical problems'
        )
    )
);

-- Exercise 3.3: Enums and Pattern Matching
INSERT INTO exercises (
    id,
    node_id,
    title,
    slug,
    description,
    difficulty,
    estimated_minutes,
    order_index,
    content
) VALUES (
    '30000000-0000-0000-0000-000000000009',
    '20000000-0000-0000-0000-000000000003',
    'Enums and Pattern Matching',
    'enums-pattern-matching',
    'Define enums with multiple variants, create enums with data (tuple and struct variants), use match expressions to handle enum variants, and work with Option<T>.',
    'medium',
    45,
    3,
    jsonb_build_object(
        'starter_code', '// TODO: Define Message enum with variants:
// - Quit (no data)
// - Move { x: i32, y: i32 }
// - Write(String)
// - ChangeColor(u8, u8, u8)

impl Message {
    // TODO: Return description based on variant
    fn process_message(&self) -> String {
        todo!()
    }

    // TODO: Return true if Move variant
    fn is_move_message(&self) -> bool {
        todo!()
    }

    // TODO: Return Some(text) if Write, None otherwise
    fn extract_text(&self) -> Option<String> {
        todo!()
    }
}

// TODO: Define Shape enum with:
// - Circle { radius: f64 }
// - Rectangle { width: f64, height: f64 }
// - Triangle { base: f64, height: f64 }

fn calculate_area(shape: &Shape) -> f64 {
    todo!()
}

fn main() {
    let msg1 = Message::Quit;
    assert_eq!(msg1.process_message(), "Quit message");
    println!("✓ Quit works");

    let msg2 = Message::Move { x: 10, y: 20 };
    assert_eq!(msg2.process_message(), "Move to (10, 20)");
    assert_eq!(msg2.is_move_message(), true);
    println!("✓ Move works");

    let msg3 = Message::Write(String::from("Hello"));
    assert_eq!(msg3.extract_text(), Some(String::from("Hello")));
    println!("✓ Write works");

    let circle = Shape::Circle { radius: 5.0 };
    assert_eq!(calculate_area(&circle), 78.5);
    println!("✓ Shape works");

    println!("\n✅ All tests passed!");
}',
        'reference_solution', 'enum Message {
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
            Message::ChangeColor(r, g, b) => format!("Change color to RGB({}, {}, {})", r, g, b),
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
    let msg1 = Message::Quit;
    assert_eq!(msg1.process_message(), "Quit message");
    println!("✓ Quit works");

    let msg2 = Message::Move { x: 10, y: 20 };
    assert_eq!(msg2.process_message(), "Move to (10, 20)");
    assert_eq!(msg2.is_move_message(), true);
    println!("✓ Move works");

    let msg3 = Message::Write(String::from("Hello"));
    assert_eq!(msg3.extract_text(), Some(String::from("Hello")));
    println!("✓ Write works");

    let circle = Shape::Circle { radius: 5.0 };
    assert_eq!(calculate_area(&circle), 78.5);
    println!("✓ Shape works");

    println!("\n✅ All tests passed!");
}',
        'test_cases', jsonb_build_array(
            jsonb_build_object('name', 'test_quit', 'code', '#[test]\nfn test_quit() {\n    let msg = Message::Quit;\n    assert_eq!(msg.process_message(), "Quit message");\n    assert_eq!(msg.is_move_message(), false);\n}', 'visible', true),
            jsonb_build_object('name', 'test_move', 'code', '#[test]\nfn test_move() {\n    let msg = Message::Move { x: 10, y: 20 };\n    assert_eq!(msg.process_message(), "Move to (10, 20)");\n    assert_eq!(msg.is_move_message(), true);\n}', 'visible', true),
            jsonb_build_object('name', 'test_write', 'code', '#[test]\nfn test_write() {\n    let msg = Message::Write(String::from("test"));\n    assert_eq!(msg.extract_text(), Some(String::from("test")));\n}', 'visible', true),
            jsonb_build_object('name', 'test_shapes', 'code', '#[test]\nfn test_shapes() {\n    let c = Shape::Circle { radius: 5.0 };\n    assert_eq!(calculate_area(&c), 78.5);\n    let r = Shape::Rectangle { width: 10.0, height: 20.0 };\n    assert_eq!(calculate_area(&r), 200.0);\n}', 'visible', true)
        ),
        'hints', jsonb_build_array(
            jsonb_build_object('order', 1, 'title', 'Enum Syntax', 'content', 'Unit variant: `Quit`, struct variant: `Move { x: i32, y: i32 }`, tuple variant: `Write(String)`. Each can hold different data.'),
            jsonb_build_object('order', 2, 'title', 'Pattern Matching', 'content', 'Use match to destructure: `match self { Message::Move { x, y } => ..., _ => ... }`. Match must be exhaustive.')
        ),
        'grading_rubric', jsonb_build_object(
            'test_weight', 60,
            'ai_weight', 40,
            'ai_criteria', jsonb_build_array(
                jsonb_build_object('name', 'Correct enum definitions', 'points', 15, 'description', 'All variants defined with correct data types'),
                jsonb_build_object('name', 'Pattern matching', 'points', 15, 'description', 'Exhaustive matching with correct destructuring'),
                jsonb_build_object('name', 'Option handling', 'points', 5, 'description', 'Returns Some/None appropriately'),
                jsonb_build_object('name', 'Code quality', 'points', 5, 'description', 'Clean match expressions, idiomatic patterns')
            )
        ),
        'ai_evaluation_prompt', 'Evaluate this enums exercise. Check: 1) Enum definitions with correct variants, 2) Pattern matching (exhaustive, correct destructuring), 3) Option<T> usage. If using matches! macro, praise idiomatic approach.',
        'common_mistakes', jsonb_build_array(
            jsonb_build_object('pattern', 'Non-exhaustive match', 'guidance', 'Rust requires handling all cases. Explain the compiler enforces this to prevent bugs.'),
            jsonb_build_object('pattern', 'Wrong destructuring syntax', 'guidance', 'Struct variants use {}, tuple variants use (). Show correct syntax.')
        ),
        'expected_outcomes', jsonb_build_array(
            'Define enums with various variant types',
            'Use pattern matching to handle enum variants',
            'Destructure data from enum variants',
            'Work with Option<T> for nullable values',
            'Understand exhaustive matching requirements'
        )
    )
);

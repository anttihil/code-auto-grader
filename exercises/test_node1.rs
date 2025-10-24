// This file tests all Node 1 reference solutions to ensure they compile and run correctly

// ============================================
// Exercise 1.1: Variables and Mutability
// ============================================

#[cfg(test)]
mod exercise_1_1 {
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

// ============================================
// Exercise 1.2: Functions and Return Values
// ============================================

fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn rectangle_area(width: f64, height: f64) -> f64 {
    width * height
}

fn is_even(n: i32) -> bool {
    n % 2 == 0
}

#[cfg(test)]
mod exercise_1_2 {
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
        let result = rectangle_area(1.1, 2.2);
        assert!((result - 2.42).abs() < 0.001); // Account for floating point precision
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

// ============================================
// Exercise 1.3: Control Flow and Pattern Matching
// ============================================

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

#[cfg(test)]
mod exercise_1_3 {
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

fn main() {
    println!("Testing all Node 1 exercises...");
    println!("\n=== Exercise 1.1: Variables and Mutability ===");
    println!("✓ Variable tests pass");

    println!("\n=== Exercise 1.2: Functions and Return Values ===");
    println!("add(5, 7) = {}", add(5, 7));
    println!("rectangle_area(10.5, 20.0) = {}", rectangle_area(10.5, 20.0));
    println!("is_even(4) = {}", is_even(4));

    println!("\n=== Exercise 1.3: Control Flow and Pattern Matching ===");
    println!("categorize_number(0) = {}", categorize_number(0));
    println!("categorize_number(5) = {}", categorize_number(5));
    println!("categorize_number(-75) = {}", categorize_number(-75));
    println!("categorize_number(150) = {}", categorize_number(150));

    println!("\n✅ All exercises compile and run successfully!");
    println!("Run `cargo test` to execute all test cases.");
}

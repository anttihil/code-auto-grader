# Content Directory

This directory contains all exercise content for the Code Auto-Grader platform.

## Structure

```
content/
├── users.json           # User data for seeding
├── tech-trees.json      # Tech tree metadata
├── nodes.json           # Node metadata and relationships
└── {language}/          # Language-specific content
    └── {topic}/         # Topic-specific exercises
        └── *.md         # Individual exercises
```

## File Formats

### JSON Metadata Files

- **users.json**: Hard-coded users for development/dogfooding
- **tech-trees.json**: Learning paths (e.g., "Rust Fundamentals")
- **nodes.json**: Topics within tech trees with prerequisites

### Exercise Markdown Files

Each exercise is a markdown file with JSON frontmatter:

```markdown
{
  "id": "30000000-0000-0000-0000-000000000001",
  "node_id": "20000000-0000-0000-0000-000000000001",
  "title": "Variables and Mutability",
  "slug": "variables-mutability",
  "difficulty": "easy",
  "estimated_minutes": 15,
  "order_index": 1
}
---

# Exercise Title

## Problem Description
Your exercise description...

## Starting Code
\`\`\`rust
// Starter code
\`\`\`

## Reference Solution
\`\`\`rust
// Solution code
\`\`\`

## Test Cases
\`\`\`rust
#[test]
fn test_something() { }
\`\`\`

## Hints
1. Hint text...

## Common Mistakes
1. **Pattern**: Explanation...

## Expected Outcomes
- Outcome 1
- Outcome 2
```

## UUID Namespacing

To avoid conflicts, UUIDs follow this pattern:

- `00000000-*`: Users
- `10000000-*`: Tech Trees
- `20000000-*`: Nodes
- `30000000-*`: Exercises

## Editing Content

### Development Workflow

1. Edit markdown files or JSON metadata
2. Run `make db-seed` to populate database
3. Database is rebuilt from these files
4. Commit changes to Git

### Adding New Exercises

1. Create new markdown file in appropriate topic directory
2. Follow the frontmatter format above
3. Use next sequential UUID in the `30000000-*` range
4. Set `order_index` appropriately
5. Run `make db-seed` to test

## Important Notes

- **Git is the source of truth** for content (Phase 1-2)
- Database is ephemeral and rebuilt from these files
- All content changes should be committed to Git
- Do NOT edit database directly during development
- Content will migrate to database-first in later phases (when we have admin UI)

## Seeding Commands

```bash
make db-seed              # Seed all content
make db-seed-users        # Seed only users
make db-seed-exercises    # Seed only exercises
make db-seed-force        # Force overwrite existing data
```

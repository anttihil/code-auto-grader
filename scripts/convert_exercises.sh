#!/bin/bash
# Script to convert old exercise format to new frontmatter format

set -e

SOURCE_DIR="/home/antti/projects/code-auto-grader/exercises"
TARGET_DIR="/home/antti/projects/code-auto-grader/content/rust"

# Mapping: node number -> directory
declare -A NODE_DIRS=(
  [1]="fundamentals"
  [2]="ownership"
  [3]="structs-enums"
)

# Mapping: node number -> node_id UUID
declare -A NODE_IDS=(
  [1]="20000000-0000-0000-0000-000000000001"
  [2]="20000000-0000-0000-0000-000000000002"
  [3]="20000000-0000-0000-0000-000000000003"
)

# Mapping: file -> exercise UUID
declare -A EXERCISE_UUIDS=(
  ["node1_exercise1_variables"]="30000000-0000-0000-0000-000000000001"
  ["node1_exercise2_functions"]="30000000-0000-0000-0000-000000000002"
  ["node1_exercise3_control_flow"]="30000000-0000-0000-0000-000000000003"
  ["node2_exercise1_ownership"]="30000000-0000-0000-0000-000000000004"
  ["node2_exercise2_borrowing"]="30000000-0000-0000-0000-000000000005"
  ["node2_exercise3_borrowing_rules"]="30000000-0000-0000-0000-000000000006"
  ["node3_exercise1_structs"]="30000000-0000-0000-0000-000000000007"
  ["node3_exercise2_methods"]="30000000-0000-0000-0000-000000000008"
  ["node3_exercise3_enums"]="30000000-0000-0000-0000-000000000009"
)

# Mapping: file -> slug
declare -A SLUGS=(
  ["node1_exercise1_variables"]="variables-mutability"
  ["node1_exercise2_functions"]="functions-return-values"
  ["node1_exercise3_control_flow"]="control-flow"
  ["node2_exercise1_ownership"]="ownership-basics"
  ["node2_exercise2_borrowing"]="borrowing-references"
  ["node2_exercise3_borrowing_rules"]="borrowing-rules"
  ["node3_exercise1_structs"]="defining-structs"
  ["node3_exercise2_methods"]="implementing-methods"
  ["node3_exercise3_enums"]="enums-pattern-matching"
)

# Convert a single exercise
convert_exercise() {
  local source_file=$1
  local basename=$(basename "$source_file" .md)

  # Extract node number
  local node_num=$(echo "$basename" | grep -oP 'node\K\d+')
  local exercise_num=$(echo "$basename" | grep -oP 'exercise\K\d+')

  # Get metadata
  local uuid="${EXERCISE_UUIDS[$basename]}"
  local node_id="${NODE_IDS[$node_num]}"
  local slug="${SLUGS[$basename]}"
  local node_dir="${NODE_DIRS[$node_num]}"

  if [ -z "$uuid" ]; then
    echo "Warning: No UUID mapping for $basename, skipping"
    return
  fi

  # Extract title from first line
  local title=$(head -1 "$source_file" | sed 's/^# //')

  # Extract difficulty and time estimate
  local difficulty=$(grep -oP '(?<=Difficulty\*\*: )\w+' "$source_file" | tr '[:upper:]' '[:lower:]' || echo "easy")
  local estimated_minutes=$(grep -oP '(?<=Estimated Time\*\*: )\d+' "$source_file" || echo "20")

  # Create target directory
  mkdir -p "$TARGET_DIR/$node_dir"

  # Determine output filename
  local output_file="$TARGET_DIR/$node_dir/$(printf '%02d' $exercise_num)-$slug.md"

  echo "Converting: $basename -> $output_file"

  # Create frontmatter
  cat > "$output_file" << EOF
{
  "id": "$uuid",
  "node_id": "$node_id",
  "title": "$title",
  "slug": "$slug",
  "difficulty": "$difficulty",
  "estimated_minutes": $estimated_minutes,
  "order_index": $exercise_num
}
---

EOF

  # Append content (skip old metadata section)
  awk '
    BEGIN { skip=0; in_meta=0 }
    /^## Metadata/ { in_meta=1; next }
    /^##/ && in_meta { in_meta=0 }
    !in_meta { print }
  ' "$source_file" >> "$output_file"

  echo "✓ Created: $output_file"
}

# Convert all exercises
for file in "$SOURCE_DIR"/node*.md; do
  if [ -f "$file" ]; then
    convert_exercise "$file"
  fi
done

echo ""
echo "✅ Conversion complete!"
echo "Converted exercises are in: $TARGET_DIR"

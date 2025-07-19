# Block IDs and References Test

## Basic Block ID
- This is a block with an ID
  id:: 550e8400-e29b-41d4-a716-446655440000
- This block references the above block: ((550e8400-e29b-41d4-a716-446655440000))

## Nested Block with ID
- Parent block
  - Child block with ID
    id:: 550e8400-e29b-41d4-a716-446655440001
    - Nested content here
  - Another child referencing sibling: ((550e8400-e29b-41d4-a716-446655440001))

## Block Embeds
- Original block to be embedded
  id:: 550e8400-e29b-41d4-a716-446655440002
  - This has nested content
  - Multiple lines
    - Deep nesting

- Here's an embed of the above block:
- {{embed ((550e8400-e29b-41d4-a716-446655440002))}}

## Multiple References
- Important fact about [[Page A]]
  id:: 550e8400-e29b-41d4-a716-446655440003
- First reference: ((550e8400-e29b-41d4-a716-446655440003))
- Second reference: ((550e8400-e29b-41d4-a716-446655440003))
- Third embed: {{embed ((550e8400-e29b-41d4-a716-446655440003))}}

## Edge Cases
- Block with ID at end id:: 550e8400-e29b-41d4-a716-446655440004
- Block with properties and ID
  type:: note
  id:: 550e8400-e29b-41d4-a716-446655440005
  priority:: high
- Reference to non-existent block: ((550e8400-e29b-41d4-a716-446655440099))
- Empty embed: {{embed }}
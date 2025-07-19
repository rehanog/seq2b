# Page B
alias:: Page Beta, Secondary Page
type:: project
status:: active

- TODO [#A] Critical: Fix data loss bug in sync by [[Jan 15th, 2025]]
  id:: 660e8400-e29b-41d4-a716-426614174002
- TODO [#A] Important: Review [[Page A]] implementation (due [[Jan 14th, 2025]])
  assigned:: [[Jane Smith]]
- Block referencing Page A's task: ((550e8400-e29b-41d4-a716-446655440001))
- TODO [#B] Medium: Add tests for #parser
  - TODO [#B] Unit tests (target [[Jan 17th, 2025]])
  - TODO [#C] Integration tests (target [[Jan 19th, 2025]])
- TODO [#C] Low: Update documentation before [[Jan 31st, 2025]]
- DONE [#A] Completed high priority task on [[Jan 11th, 2025]]
- Welcome to Page B with **bold**, *italic*, and ^^highlighted^^ text
- This links back to [[Page A]] and has #important #documentation tags
- TODO Implement better linking to [[Page C]] (research started [[Jan 9th, 2025]])
  - DONE Research best practices on [[Jan 10th, 2025]]
  - TODO Apply to codebase by [[Jan 16th, 2025]]
- Here's a reference to [[Page C]] as well
  - Nested reference to [[Page D]]
  - Meeting notes from [[Jan 8th, 2025]]
    collapsed:: true
    - Discussion about architecture
    - Decision to use Go
- Solo block with no links but with ~~deprecated info~~
- Weekly sync scheduled for every Monday (next: [[Jan 20th, 2025]])
  recurring:: true
  time:: 10:00 AM

## Query Examples
- All TODO items: {{query (todo TODO)}}
- High priority tasks: {{query (and (todo TODO) (priority A))}}
- Tasks assigned to me: {{query (and (todo TODO) (property assigned [[John Doe]]))}}

## Block Embeds
- Embedding critical task from Page A:
- {{embed ((550e8400-e29b-41d4-a716-446655440001))}}
- Embedding local task:
- {{embed ((660e8400-e29b-41d4-a716-426614174002))}}

## Tables (Not Yet Supported)
| Feature | Status | Priority |
|---------|--------|----------|
| Parser  | DONE   | High     |
| Tags    | NOW    | Medium   |
| PDF     | LATER  | High     |
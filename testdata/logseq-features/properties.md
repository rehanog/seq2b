# Properties Test
alias:: Props Test, Properties Example
tags:: test, properties, demo
type:: [[documentation]]
created:: [[2025-01-19]]
modified:: [[2025-01-19]]
status:: draft

## Page Properties
- Page properties must be at the top of the file
- They apply to the entire page
- Common page properties include alias, tags, type

## Block Properties
- This is a regular block
- This block has a type property
  type:: note
- This block has multiple properties
  type:: task
  priority:: high
  assigned:: [[John Doe]]
  due:: [[2025-01-25]]
- Properties with multiple values
  tags:: logseq, productivity, notes
  related:: [[Page A]], [[Page B]], [[Page C]]

## Special Properties
- TODO Task with scheduled date
  SCHEDULED: <2025-01-20 Mon>
- Task with deadline
  DEADLINE: <2025-01-25 Sat>
- Collapsed block (UI state)
  collapsed:: true
  - Hidden child 1
  - Hidden child 2
- Block with custom ID
  id:: my-custom-block-id
- Numbered list type
  logseq.order-list-type:: number
  - First item
  - Second item

## Property Values
- Text value
  simple:: plain text value
- Page reference value
  related:: [[Page A]]
- Multiple page references
  references:: [[Page A]], [[Page B]]
- Tag value
  category:: #important
- Date value
  date:: 2025-01-19
- Number value
  priority:: 1
- Boolean value
  public:: true
  archived:: false

## Edge Cases
- Property with spaces in name (invalid)
  my property:: value
- Empty property
  empty::
- Property with special characters
  special-chars:: @#$%^&*()
- Very long property value
  description:: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.

## Query Examples
- Find all blocks with type:: task
- Find all blocks with priority:: high
- Find all pages with tags containing "test"
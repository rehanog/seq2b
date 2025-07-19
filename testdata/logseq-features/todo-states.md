# TODO States and Task Management

## Basic TODO States
- TODO Write documentation
- DONE Complete parser implementation
- TODO Add more tests

## Extended Logseq States
- NOW Working on this right now
- LATER Will do this eventually
- DOING Currently in progress
- WAIT Waiting for feedback
- CANCELLED This task was cancelled

## TODO with Properties
- TODO Call client
  SCHEDULED: <2025-01-20 Mon 14:00>
  DEADLINE: <2025-01-20 Mon 17:00>
  
- DONE Submit report
  SCHEDULED: <2025-01-19 Sun>
  completed:: 2025-01-19T15:30:00
  
- TODO Review PRs #work #code-review
  priority:: high
  assigned:: [[John Doe]]
  estimated:: 2h

## Nested TODOs
- TODO Project Alpha
  - DONE Set up repository
  - DOING Implement core features
    - DONE User authentication
    - TODO API endpoints
    - TODO Frontend integration
  - WAIT Deploy to production
    waiting-for:: [[DevOps Team]]

## Repeating Tasks
- TODO Daily standup
  SCHEDULED: <2025-01-20 Mon +1d>
  
- TODO Weekly review
  SCHEDULED: <2025-01-21 Tue +1w>
  
- TODO Monthly report
  SCHEDULED: <2025-01-31 Fri +1m>

## Task Priorities with Markers
- TODO [#A] Critical bug fix
- TODO [#B] Feature implementation
- TODO [#C] Code refactoring
- TODO [#D] Nice to have enhancement

## Time Tracking
- DONE Meeting with team
  :LOGBOOK:
  CLOCK: [2025-01-19 Mon 10:00]--[2025-01-19 Mon 11:30] => 1:30
  :END:
  
- DOING Writing documentation
  :LOGBOOK:
  CLOCK: [2025-01-19 Mon 14:00]
  :END:

## Checkbox Lists
- [ ] Shopping list
  - [x] Milk
  - [x] Bread
  - [ ] Eggs
  - [-] Vegetables (partially done)
    - [x] Tomatoes
    - [x] Lettuce
    - [ ] Carrots
    - [ ] Potatoes

## TODO Queries Use Cases
- Show all TODO items: {{query (todo TODO)}}
- Show NOW and DOING: {{query (or (todo NOW) (todo DOING))}}
- Show overdue tasks: {{query (and (todo TODO) (property deadline) (before today))}}
- Show tasks for this week: {{query (and (task) (between today +7d))}}

## Task Dependencies
- TODO Deploy to production
  depends-on:: [[Complete testing]], [[Security review]]
  blocked-by:: [[Infrastructure setup]]
  
- DONE Complete testing
  blocks:: [[Deploy to production]]

## Task Context
- TODO Call dentist @phone #personal
- TODO Buy groceries @errands #personal
- TODO Code review @computer #work
- TODO Read research paper @anywhere #learning

## Edge Cases
- TODO Task with very long description that spans multiple lines and contains various [[page links]] and #tags and even some **bold** and *italic* text to test how well the parser handles complex content in TODO items
- TO DO (with space - not a valid TODO)
- todo (lowercase - not a valid TODO)
- TODO: With colon (parser dependent)
- - TODO Double dash (parser dependent)
- [ ] TODO Checkbox and TODO marker together
# Queries Test

## Basic Query
- {{query}}
- This should show all blocks in the graph

## TODO Queries
- All TODO items:
- {{query (todo TODO)}}

- All DONE items:
- {{query (todo DONE)}}

- All tasks (any state):
- {{query (task)}}

## Tag Queries
- All blocks with #important tag:
- {{query [[#important]]}}

- Multiple tags (AND):
- {{query (and [[#work]] [[#urgent]])}}

- Multiple tags (OR):
- {{query (or [[#personal]] [[#work]])}}

## Page Reference Queries
- All blocks mentioning [[Page A]]:
- {{query [[Page A]]}}

- Blocks mentioning multiple pages:
- {{query (and [[Page A]] [[Page B]])}}

## Property Queries
- All blocks with high priority:
- {{query (property priority high)}}

- All blocks with type property:
- {{query (property type)}}

- Specific property value:
- {{query (property status draft)}}

## Date Queries
- Today's journal:
- {{query (today)}}

- Past 7 days:
- {{query (between -7d today)}}

- Specific date range:
- {{query (between [[2025-01-01]] [[2025-01-31]])}}

- All blocks with scheduled dates:
- {{query (property scheduled)}}

## Complex Queries
- TODO items tagged #work from last week:
- {{query (and (todo TODO) [[#work]] (between -7d today))}}

- High priority tasks not done:
- {{query (and (todo TODO) (property priority high))}}

- Pages with specific properties:
- {{query (and (page) (property type documentation))}}

## NOT Queries
- All blocks NOT tagged #done:
- {{query (not [[#done]])}}

- TODOs without #work tag:
- {{query (and (todo TODO) (not [[#work]]))}}

## Query Sorting and Limits
- Recent TODOs (sorted by date):
- {{query (todo TODO) :sort-by created :limit 10}}

- High priority first:
- {{query (property priority) :sort-by priority :desc true}}

## Query Result Customization
- Collapsed results:
- {{query (todo TODO) :collapsed? true}}

- Show only block content:
- {{query [[#important]] :breadcrumb? false}}

## Advanced Query Features
- Group by page:
- {{query (todo TODO) :group-by page}}

- Custom title:
- {{query (todo TODO) :title "My Open Tasks"}}

## Query Variables
- Dynamic current page:
- {{query (and (todo TODO) (link <%current-page%>))}}

- With namespace:
- {{query (namespace [[project]])}}

## Performance Queries
- These might be slow on large graphs:
- Full text search: {{query "specific text"}}
- Regex matching: {{query (re-find "pattern")}}

## Edge Cases
- Empty query: {{query}}
- Malformed query: {{query (and}}
- Non-existent property: {{query (property nonexistent value)}}
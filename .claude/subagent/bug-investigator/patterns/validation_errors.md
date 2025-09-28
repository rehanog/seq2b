# Validation Error Patterns

## Common Validation Issues in Seq2b

### 1. Markdown Parsing Validation
**Symptoms**:
- Malformed wiki links crash parser
- Unclosed brackets cause infinite loops
- Special characters break rendering

**Common Patterns**:
```go
// BAD: No validation
func ParseWikiLink(content string) string {
    start := strings.Index(content, "[[")
    end := strings.Index(content, "]]")
    return content[start+2:end]  // Panic if indices are -1
}

// GOOD: Proper validation
func ParseWikiLink(content string) (string, error) {
    start := strings.Index(content, "[[")
    if start == -1 {
        return "", fmt.Errorf("no opening brackets found")
    }
    
    end := strings.Index(content[start:], "]]")
    if end == -1 {
        return "", fmt.Errorf("unclosed wiki link at position %d", start)
    }
    
    return content[start+2:start+end], nil
}
```

### 2. File Path Validation
**Symptoms**:
- Path traversal vulnerabilities
- Invalid characters in filenames
- Platform-specific path issues

**Test Cases**:
```go
func TestPathValidation(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"valid path", "pages/test.md", false},
        {"path traversal", "../../../etc/passwd", true},
        {"absolute path", "/etc/passwd", true},
        {"null bytes", "test\x00.md", true},
        {"windows reserved", "CON.md", true},
        {"special chars", "test:file.md", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePath(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidatePath(%q) error = %v, wantErr %v", 
                    tt.path, err, tt.wantErr)
            }
        })
    }
}
```

### 3. Input Sanitization
**Symptoms**:
- XSS vulnerabilities in rendered content
- SQL injection in search queries
- Command injection in file operations

**Prevention**:
```go
// Sanitize HTML content
func SanitizeHTML(input string) string {
    // Escape dangerous characters
    input = html.EscapeString(input)
    
    // Additional sanitization for Markdown
    input = strings.ReplaceAll(input, "<script", "&lt;script")
    input = strings.ReplaceAll(input, "javascript:", "")
    
    return input
}

// Validate block content
func ValidateBlock(block Block) error {
    if block.UUID == "" {
        return errors.New("block UUID cannot be empty")
    }
    
    if len(block.Content) > MaxBlockSize {
        return fmt.Errorf("block content exceeds maximum size of %d", MaxBlockSize)
    }
    
    if containsNullBytes(block.Content) {
        return errors.New("block content contains null bytes")
    }
    
    return nil
}
```

### 4. Property Validation
**Symptoms**:
- Invalid property syntax causes parser errors
- Circular references in properties
- Type mismatches in property values

**Test Pattern**:
```go
func TestPropertyValidation(t *testing.T) {
    tests := []struct {
        name     string
        property Property
        wantErr  string
    }{
        {
            name: "valid property",
            property: Property{
                Key:   "tags",
                Value: []string{"todo", "important"},
            },
            wantErr: "",
        },
        {
            name: "empty key",
            property: Property{
                Key:   "",
                Value: "value",
            },
            wantErr: "property key cannot be empty",
        },
        {
            name: "invalid characters in key",
            property: Property{
                Key:   "my:key",
                Value: "value",
            },
            wantErr: "property key contains invalid characters",
        },
        {
            name: "circular reference",
            property: Property{
                Key:   "parent",
                Value: "[[Current Page]]",
            },
            wantErr: "circular reference detected",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateProperty(tt.property)
            if tt.wantErr == "" {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            } else {
                if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
                    t.Errorf("expected error containing %q, got %v", 
                        tt.wantErr, err)
                }
            }
        })
    }
}
```

### 5. API Request Validation
**Symptoms**:
- Missing required fields cause panics
- Invalid data types cause runtime errors
- Unbounded input sizes cause OOM

**Validation Middleware**:
```go
func ValidateRequest(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Size limits
        r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
        
        // Content type validation
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, "Invalid content type", http.StatusBadRequest)
            return
        }
        
        // Parse and validate JSON
        var req Request
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        if err := req.Validate(); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        next(w, r)
    }
}
```

## Common Validation Mistakes

### 1. Trusting User Input
```go
// BAD
pageTitle := r.FormValue("title")
content := fmt.Sprintf("<h1>%s</h1>", pageTitle)  // XSS vulnerability

// GOOD
pageTitle := html.EscapeString(r.FormValue("title"))
content := fmt.Sprintf("<h1>%s</h1>", pageTitle)
```

### 2. Incomplete Error Checking
```go
// BAD
data, _ := ioutil.ReadFile(path)  // Ignoring error
json.Unmarshal(data, &config)     // May panic on nil data

// GOOD
data, err := ioutil.ReadFile(path)
if err != nil {
    return fmt.Errorf("reading file: %w", err)
}
if err := json.Unmarshal(data, &config); err != nil {
    return fmt.Errorf("parsing JSON: %w", err)
}
```

### 3. Off-by-One Errors
```go
// BAD
if index <= len(slice) {  // Should be <
    value := slice[index]  // Panic on boundary
}

// GOOD
if index < len(slice) {
    value := slice[index]
}
```

## Validation Testing Checklist

- [ ] Test empty/nil inputs
- [ ] Test boundary values (0, -1, MAX_INT)
- [ ] Test invalid UTF-8 sequences
- [ ] Test extremely long inputs
- [ ] Test special characters (<, >, &, ", ')
- [ ] Test path traversal attempts
- [ ] Test null bytes in strings
- [ ] Test circular references
- [ ] Test type mismatches
- [ ] Test concurrent validation

## Fuzzing for Validation

```go
func FuzzMarkdownParser(f *testing.F) {
    // Seed corpus
    f.Add("[[Normal Link]]")
    f.Add("[[Link with | pipe]]")
    f.Add("[["))
    f.Add("]][[]][[")
    
    f.Fuzz(func(t *testing.T, input string) {
        // Should not panic
        result, err := ParseMarkdown(input)
        
        // If no error, result should be valid
        if err == nil {
            if err := ValidateMarkdown(result); err != nil {
                t.Errorf("invalid result from valid input: %v", err)
            }
        }
    })
}
```

## References
- [OWASP Input Validation Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html)
- [Go Security Best Practices](https://github.com/OWASP/Go-SCP)
- [Common Web Vulnerabilities](https://owasp.org/www-project-top-ten/)
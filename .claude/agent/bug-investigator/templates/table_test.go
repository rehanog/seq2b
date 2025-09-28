package templates

import (
	"testing"
)

// TableTestTemplate demonstrates table-driven testing patterns
// Use this template when testing:
// - Multiple input/output scenarios
// - Edge cases and boundary conditions
// - Input validation
// - Parser variations

// Basic table-driven test
func TestOperation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   "hello",
			want:    "HELLO",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "special characters",
			input:   "hello@world!",
			want:    "HELLO@WORLD!",
			wantErr: false,
		},
		{
			name:    "unicode",
			input:   "hello 世界",
			want:    "HELLO 世界",
			wantErr: false,
		},
		{
			name:    "nil input causes error",
			input:   "\x00",
			want:    "",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := performOperation(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("performOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if got != tt.want {
				t.Errorf("performOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Table test with setup and cleanup
func TestWithSetup(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Resource
		input   int
		want    int
		cleanup func(*Resource)
	}{
		{
			name: "with cache",
			setup: func() *Resource {
				r := NewResource()
				r.EnableCache()
				return r
			},
			input: 5,
			want:  25,
			cleanup: func(r *Resource) {
				r.ClearCache()
			},
		},
		{
			name: "without cache",
			setup: func() *Resource {
				return NewResource()
			},
			input:   5,
			want:    25,
			cleanup: func(r *Resource) {},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := tt.setup()
			defer tt.cleanup(resource)
			
			got := resource.Calculate(tt.input)
			if got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Complex table test with multiple assertions
func TestComplexScenarios(t *testing.T) {
	tests := []struct {
		name      string
		input     Input
		want      Output
		wantErr   error
		wantPanic bool
		validate  func(t *testing.T, got Output)
	}{
		{
			name: "successful processing",
			input: Input{
				Data:    "test data",
				Options: Options{Validate: true},
			},
			want: Output{
				Result: "processed",
				Status: "success",
			},
			wantErr: nil,
			validate: func(t *testing.T, got Output) {
				if got.Timestamp.IsZero() {
					t.Error("timestamp should be set")
				}
				if len(got.Metadata) == 0 {
					t.Error("metadata should not be empty")
				}
			},
		},
		{
			name: "validation failure",
			input: Input{
				Data:    "",
				Options: Options{Validate: true},
			},
			want:    Output{},
			wantErr: ErrValidation,
		},
		{
			name: "panic on nil input",
			input: Input{
				Data: nil,
			},
			wantPanic: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic but didn't get one")
					}
				}()
			}
			
			got, err := Process(tt.input)
			
			if err != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err == nil {
				if got.Result != tt.want.Result {
					t.Errorf("Result = %v, want %v", got.Result, tt.want.Result)
				}
				if got.Status != tt.want.Status {
					t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
				}
				
				if tt.validate != nil {
					tt.validate(t, got)
				}
			}
		})
	}
}

// Parameterized test with subtests
func TestParameterized(t *testing.T) {
	type args struct {
		x int
		y int
	}
	
	tests := []struct {
		name string
		args args
		want int
	}{
		{"positive numbers", args{2, 3}, 5},
		{"negative numbers", args{-2, -3}, -5},
		{"mixed numbers", args{-2, 3}, 1},
		{"with zero", args{0, 5}, 5},
		{"both zero", args{0, 0}, 0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", 
					tt.args.x, tt.args.y, got, tt.want)
			}
		})
	}
}

// Golden file test pattern
func TestGoldenFiles(t *testing.T) {
	tests := []struct {
		name       string
		inputFile  string
		goldenFile string
		update     bool // set to true to update golden files
	}{
		{
			name:       "basic markdown",
			inputFile:  "testdata/input/basic.md",
			goldenFile: "testdata/golden/basic.json",
		},
		{
			name:       "complex markdown",
			inputFile:  "testdata/input/complex.md",
			goldenFile: "testdata/golden/complex.json",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := readFile(t, tt.inputFile)
			got := ParseMarkdown(input)
			
			if tt.update {
				writeFile(t, tt.goldenFile, got)
				return
			}
			
			want := readFile(t, tt.goldenFile)
			if got != want {
				t.Errorf("ParseMarkdown() output doesn't match golden file")
				t.Errorf("got:\n%s", got)
				t.Errorf("want:\n%s", want)
			}
		})
	}
}

// Placeholder types and functions - replace with actual implementation
type Input struct {
	Data    interface{}
	Options Options
}

type Options struct {
	Validate bool
}

type Output struct {
	Result    string
	Status    string
	Timestamp time.Time
	Metadata  map[string]string
}

type Resource struct{}

func NewResource() *Resource                 { return &Resource{} }
func (r *Resource) EnableCache()            {}
func (r *Resource) ClearCache()             {}
func (r *Resource) Calculate(int) int       { return 0 }
func performOperation(string) (string, error) { return "", nil }
func Process(Input) (Output, error)         { return Output{}, nil }
func Add(x, y int) int                      { return x + y }
func ParseMarkdown(string) string           { return "" }
func readFile(t *testing.T, path string) string { return "" }
func writeFile(t *testing.T, path string, content string) {}

var ErrValidation = fmt.Errorf("validation error")
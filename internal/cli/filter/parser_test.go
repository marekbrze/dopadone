package filter

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple equality",
			input: "name=foo",
			expected: []Token{
				{Type: TokenIdent, Value: "name"},
				{Type: TokenEQ, Value: "="},
				{Type: TokenIdent, Value: "foo"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "comparison operators",
			input: "count>=10 AND count<100",
			expected: []Token{
				{Type: TokenIdent, Value: "count"},
				{Type: TokenGE, Value: ">="},
				{Type: TokenNumber, Value: "10"},
				{Type: TokenAND, Value: "AND"},
				{Type: TokenIdent, Value: "count"},
				{Type: TokenLT, Value: "<"},
				{Type: TokenNumber, Value: "100"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "string value with quotes",
			input: `name='foo bar'`,
			expected: []Token{
				{Type: TokenIdent, Value: "name"},
				{Type: TokenEQ, Value: "="},
				{Type: TokenString, Value: "foo bar"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "parentheses and OR",
			input: "(status=active OR status=pending)",
			expected: []Token{
				{Type: TokenLParen, Value: "("},
				{Type: TokenIdent, Value: "status"},
				{Type: TokenEQ, Value: "="},
				{Type: TokenIdent, Value: "active"},
				{Type: TokenOR, Value: "OR"},
				{Type: TokenIdent, Value: "status"},
				{Type: TokenEQ, Value: "="},
				{Type: TokenIdent, Value: "pending"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "not equal operator",
			input: "status!=archived",
			expected: []Token{
				{Type: TokenIdent, Value: "status"},
				{Type: TokenNE, Value: "!="},
				{Type: TokenIdent, Value: "archived"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "less than or equal",
			input: "priority<=3",
			expected: []Token{
				{Type: TokenIdent, Value: "priority"},
				{Type: TokenLE, Value: "<="},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(tokens))
				return
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i].Type, tok.Type)
				}
				if tok.Value != tt.expected[i].Value {
					t.Errorf("token %d: expected value %q, got %q", i, tt.expected[i].Value, tok.Value)
				}
			}
		})
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple equality", "name=foo", false},
		{"comparison", "count>=10", false},
		{"AND expression", "status=active AND priority=high", false},
		{"OR expression", "status=active OR status=pending", false},
		{"parentheses", "(status=active OR status=pending) AND priority=high", false},
		{"complex", "status=active AND (priority=high OR priority=urgent) AND progress>=50", false},
		{"string value", `name='foo bar'`, false},
		{"double quotes", `name="test"`, false},
		{"numeric value", "progress>=75", false},
		{"float value", "score>=3.5", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.input != "" && expr == nil {
				t.Error("Parse() returned nil expression for non-empty input")
			}
		})
	}
}

func TestParser_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing operator", "name foo"},
		{"missing value", "name="},
		{"unclosed parenthesis", "(status=active"},
		{"empty parentheses", "()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse() expected error for input %q", tt.input)
			}
		})
	}
}

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		data     map[string]interface{}
		expected bool
	}{
		{
			name:     "simple equality match",
			filter:   "status=active",
			data:     map[string]interface{}{"status": "active"},
			expected: true,
		},
		{
			name:     "simple equality no match",
			filter:   "status=active",
			data:     map[string]interface{}{"status": "completed"},
			expected: false,
		},
		{
			name:     "case insensitive equality",
			filter:   "status=ACTIVE",
			data:     map[string]interface{}{"status": "active"},
			expected: true,
		},
		{
			name:     "not equal",
			filter:   "status!=archived",
			data:     map[string]interface{}{"status": "active"},
			expected: true,
		},
		{
			name:     "greater than",
			filter:   "progress>=50",
			data:     map[string]interface{}{"progress": 75},
			expected: true,
		},
		{
			name:     "less than",
			filter:   "progress<50",
			data:     map[string]interface{}{"progress": 25},
			expected: true,
		},
		{
			name:     "AND expression both true",
			filter:   "status=active AND priority=high",
			data:     map[string]interface{}{"status": "active", "priority": "high"},
			expected: true,
		},
		{
			name:     "AND expression one false",
			filter:   "status=active AND priority=high",
			data:     map[string]interface{}{"status": "active", "priority": "low"},
			expected: false,
		},
		{
			name:     "OR expression both false",
			filter:   "status=active OR priority=high",
			data:     map[string]interface{}{"status": "completed", "priority": "low"},
			expected: false,
		},
		{
			name:     "OR expression one true",
			filter:   "status=active OR priority=high",
			data:     map[string]interface{}{"status": "active", "priority": "low"},
			expected: true,
		},
		{
			name:     "complex expression",
			filter:   "status=active AND (priority=high OR priority=urgent)",
			data:     map[string]interface{}{"status": "active", "priority": "urgent"},
			expected: true,
		},
		{
			name:     "numeric comparison int",
			filter:   "progress>50",
			data:     map[string]interface{}{"progress": 75},
			expected: true,
		},
		{
			name:     "numeric comparison equal",
			filter:   "progress=100",
			data:     map[string]interface{}{"progress": 100},
			expected: true,
		},
		{
			name:     "missing field",
			filter:   "status=active",
			data:     map[string]interface{}{"name": "test"},
			expected: false,
		},
		{
			name:     "empty filter matches all",
			filter:   "",
			data:     map[string]interface{}{"status": "active"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EvaluateFilter(tt.filter, []map[string]interface{}{tt.data})
			if err != nil {
				t.Errorf("EvaluateFilter() error = %v", err)
				return
			}

			matched := len(result) > 0
			if matched != tt.expected {
				t.Errorf("EvaluateFilter() = %v, want %v", matched, tt.expected)
			}
		})
	}
}

func TestEvaluateFilter(t *testing.T) {
	data := []map[string]interface{}{
		{"id": "1", "status": "active", "priority": "high", "progress": 75},
		{"id": "2", "status": "completed", "priority": "medium", "progress": 100},
		{"id": "3", "status": "active", "priority": "low", "progress": 25},
		{"id": "4", "status": "on_hold", "priority": "high", "progress": 50},
	}

	tests := []struct {
		name          string
		filter        string
		expectedCount int
		expectedIDs   []string
	}{
		{
			name:          "filter by status",
			filter:        "status=active",
			expectedCount: 2,
			expectedIDs:   []string{"1", "3"},
		},
		{
			name:          "filter by priority",
			filter:        "priority=high",
			expectedCount: 2,
			expectedIDs:   []string{"1", "4"},
		},
		{
			name:          "filter by progress",
			filter:        "progress>=50",
			expectedCount: 3,
			expectedIDs:   []string{"1", "2", "4"},
		},
		{
			name:          "complex filter",
			filter:        "status=active AND progress>=50",
			expectedCount: 1,
			expectedIDs:   []string{"1"},
		},
		{
			name:          "no match",
			filter:        "status=archived",
			expectedCount: 0,
			expectedIDs:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EvaluateFilter(tt.filter, data)
			if err != nil {
				t.Errorf("EvaluateFilter() error = %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("EvaluateFilter() returned %d results, want %d", len(result), tt.expectedCount)
				return
			}

			if tt.expectedCount > 0 {
				resultIDs := make([]string, len(result))
				for i, r := range result {
					resultIDs[i] = r["id"].(string)
				}
				for _, expectedID := range tt.expectedIDs {
					found := false
					for _, resultID := range resultIDs {
						if resultID == expectedID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected ID %s not found in results", expectedID)
					}
				}
			}
		})
	}
}

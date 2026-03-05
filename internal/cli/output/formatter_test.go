package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"table format", "table", false},
		{"json format", "json", false},
		{"yaml format", "yaml", false},
		{"empty defaults to table", "", false},
		{"invalid format", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := NewFormatter(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && formatter == nil {
				t.Error("NewFormatter() returned nil formatter")
			}
		})
	}
}

func TestTableFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewTableFormatterWithWriter(buf)

	headers := []string{"ID", "Name", "Status"}
	formatter.PrintHeader(headers)

	rows := [][]string{
		{"1", "Project A", "active"},
		{"2", "Project B", "completed"},
	}
	for _, row := range rows {
		formatter.PrintRow(row)
	}

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("TableFormatter produced no output")
	}
}

func TestJSONFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewJSONFormatterWithWriter(buf)

	obj := map[string]string{
		"id":     "1",
		"name":   "Test Project",
		"status": "active",
	}

	if err := formatter.PrintObject(obj); err != nil {
		t.Errorf("PrintObject() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("JSONFormatter produced no output")
	}

	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	if result["id"] != "1" {
		t.Errorf("Expected id=1, got %s", result["id"])
	}
}

func TestJSONFormatter_MultipleObjects(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewJSONFormatterWithWriter(buf)

	objs := []interface{}{
		map[string]string{"id": "1", "name": "A"},
		map[string]string{"id": "2", "name": "B"},
	}

	for _, obj := range objs {
		formatter.AddObject(obj)
	}

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("JSONFormatter produced no output")
	}

	var results []map[string]string
	if err := json.Unmarshal([]byte(output), &results); err != nil {
		t.Errorf("Output is not valid JSON array: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(results))
	}
}

func TestJSONFormatter_SingleObject(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewJSONFormatterWithWriter(buf)

	formatter.AddObject(map[string]string{"id": "1", "name": "Single"})

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()

	var obj map[string]string
	if err := json.Unmarshal([]byte(output), &obj); err != nil {
		t.Errorf("Single object should be JSON object, not array: %v", err)
	}

	if obj["name"] != "Single" {
		t.Errorf("Expected name=Single, got %s", obj["name"])
	}
}

func TestJSONFormatter_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewJSONFormatterWithWriter(buf)

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("Empty formatter should produce no output, got: %q", output)
	}
}

func TestYAMLFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewYAMLFormatterWithWriter(buf)

	obj := map[string]interface{}{
		"id":     "1",
		"name":   "Test Project",
		"status": "active",
		"count":  42,
	}

	if err := formatter.PrintObject(obj); err != nil {
		t.Errorf("PrintObject() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("YAMLFormatter produced no output")
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid YAML: %v", err)
	}

	if result["id"] != "1" {
		t.Errorf("Expected id=1, got %v", result["id"])
	}
}

func TestYAMLFormatter_MultipleObjects(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewYAMLFormatterWithWriter(buf)

	objs := []interface{}{
		map[string]string{"id": "1", "name": "A"},
		map[string]string{"id": "2", "name": "B"},
	}

	for _, obj := range objs {
		formatter.AddObject(obj)
	}

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("YAMLFormatter produced no output")
	}

	var results []map[string]interface{}
	if err := yaml.Unmarshal([]byte(output), &results); err != nil {
		t.Errorf("Output is not valid YAML array: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(results))
	}
}

func TestYAMLFormatter_SingleObject(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewYAMLFormatterWithWriter(buf)

	formatter.AddObject(map[string]string{"id": "1", "name": "Single"})

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()

	var obj map[string]interface{}
	if err := yaml.Unmarshal([]byte(output), &obj); err != nil {
		t.Errorf("Single object should be YAML object: %v", err)
	}

	if obj["name"] != "Single" {
		t.Errorf("Expected name=Single, got %v", obj["name"])
	}
}

func TestYAMLFormatter_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewYAMLFormatterWithWriter(buf)

	if err := formatter.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("Empty formatter should produce no output, got: %q", output)
	}
}

func TestYAMLFormatter_ComplexObject(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewYAMLFormatterWithWriter(buf)

	obj := map[string]interface{}{
		"id":          "project-123",
		"name":        "Complex Project",
		"status":      "active",
		"priority":    "high",
		"progress":    75,
		"description": "A complex project with many fields",
		"tags":        []string{"backend", "api", "urgent"},
		"metadata": map[string]interface{}{
			"created_by": "user-1",
			"version":    2,
		},
	}

	if err := formatter.PrintObject(obj); err != nil {
		t.Errorf("PrintObject() error = %v", err)
	}

	output := buf.String()

	var result map[string]interface{}
	if err := yaml.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid YAML: %v", err)
	}

	if result["name"] != "Complex Project" {
		t.Errorf("Expected name=Complex Project, got %v", result["name"])
	}

	tags, ok := result["tags"].([]interface{})
	if !ok {
		t.Error("Expected tags to be an array")
		return
	}
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
}

func TestAllFormatters_Interface(t *testing.T) {
	formatters := []struct {
		name      string
		formatter Formatter
	}{
		{"table", NewTableFormatterWithWriter(&bytes.Buffer{})},
		{"json", NewJSONFormatterWithWriter(&bytes.Buffer{})},
		{"yaml", NewYAMLFormatterWithWriter(&bytes.Buffer{})},
	}

	for _, tt := range formatters {
		t.Run(tt.name, func(t *testing.T) {
			tt.formatter.PrintHeader([]string{"ID", "Name"})
			tt.formatter.PrintRow([]string{"1", "Test"})
			if err := tt.formatter.Flush(); err != nil {
				t.Errorf("Flush() error = %v", err)
			}
		})
	}
}

func TestFormatConstants(t *testing.T) {
	if FormatTable != "table" {
		t.Errorf("FormatTable = %q, want 'table'", FormatTable)
	}
	if FormatJSON != "json" {
		t.Errorf("FormatJSON = %q, want 'json'", FormatJSON)
	}
	if FormatYAML != "yaml" {
		t.Errorf("FormatYAML = %q, want 'yaml'", FormatYAML)
	}
}

func TestNewFormatter_ErrorMessage(t *testing.T) {
	_, err := NewFormatter("invalid")
	if err == nil {
		t.Error("NewFormatter() should return error for invalid format")
	}
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

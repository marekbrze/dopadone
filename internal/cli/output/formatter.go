package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Padding(0, 1)

	cellStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

type Formatter interface {
	PrintHeader(headers []string)
	PrintRow(columns []string)
	Flush() error
}

type TableFormatter struct {
	writer io.Writer
}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{
		writer: os.Stdout,
	}
}

func NewTableFormatterWithWriter(w io.Writer) *TableFormatter {
	return &TableFormatter{
		writer: w,
	}
}

func (f *TableFormatter) PrintHeader(headers []string) {
	styledHeaders := make([]string, len(headers))
	for i, h := range headers {
		styledHeaders[i] = headerStyle.Render(h)
	}
	fmt.Fprintln(f.writer, lipgloss.JoinHorizontal(lipgloss.Left, styledHeaders...))
}

func (f *TableFormatter) PrintRow(columns []string) {
	styledCells := make([]string, len(columns))
	for i, c := range columns {
		styledCells[i] = cellStyle.Render(c)
	}
	fmt.Fprintln(f.writer, lipgloss.JoinHorizontal(lipgloss.Left, styledCells...))
}

func (f *TableFormatter) Flush() error {
	return nil
}

type JSONFormatter struct {
	writer  io.Writer
	objects []interface{}
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		writer:  os.Stdout,
		objects: []interface{}{},
	}
}

func NewJSONFormatterWithWriter(w io.Writer) *JSONFormatter {
	return &JSONFormatter{
		writer:  w,
		objects: []interface{}{},
	}
}

func (f *JSONFormatter) PrintHeader(headers []string) {
}

func (f *JSONFormatter) PrintRow(columns []string) {
}

func (f *JSONFormatter) AddObject(obj interface{}) {
	f.objects = append(f.objects, obj)
}

func (f *JSONFormatter) PrintObject(obj interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(obj)
}

func (f *JSONFormatter) Flush() error {
	if len(f.objects) == 0 {
		return nil
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")

	if len(f.objects) == 1 {
		return encoder.Encode(f.objects[0])
	}

	return encoder.Encode(f.objects)
}

func NewFormatter(format string) (Formatter, error) {
	switch Format(format) {
	case FormatTable, "":
		return NewTableFormatter(), nil
	case FormatJSON:
		return NewJSONFormatter(), nil
	case FormatYAML:
		return NewYAMLFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s (supported: table, json, yaml)", format)
	}
}

func PrintSuccess(message string) {
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)
	fmt.Println(successStyle.Render("✓ " + message))
}

func PrintError(message string) {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)
	fmt.Fprintln(os.Stderr, errorStyle.Render("✗ "+message))
}

func PrintInfo(message string) {
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))
	fmt.Println(infoStyle.Render(message))
}

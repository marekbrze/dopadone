package output

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLFormatter struct {
	writer  io.Writer
	objects []interface{}
}

func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{
		writer:  os.Stdout,
		objects: []interface{}{},
	}
}

func NewYAMLFormatterWithWriter(w io.Writer) *YAMLFormatter {
	return &YAMLFormatter{
		writer:  w,
		objects: []interface{}{},
	}
}

func (f *YAMLFormatter) PrintHeader(headers []string) {
}

func (f *YAMLFormatter) PrintRow(columns []string) {
}

func (f *YAMLFormatter) AddObject(obj interface{}) {
	f.objects = append(f.objects, obj)
}

func (f *YAMLFormatter) PrintObject(obj interface{}) error {
	encoder := yaml.NewEncoder(f.writer)
	encoder.SetIndent(2)
	return encoder.Encode(obj)
}

func (f *YAMLFormatter) Flush() error {
	if len(f.objects) == 0 {
		return nil
	}

	encoder := yaml.NewEncoder(f.writer)
	encoder.SetIndent(2)

	if len(f.objects) == 1 {
		return encoder.Encode(f.objects[0])
	}

	return encoder.Encode(f.objects)
}

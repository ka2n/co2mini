package co2mini

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// OutputWriter is interface for output
type OutputWriter interface {
	Write(v Value) error
}

// JSONOutputWriter is output Value as JSON
type JSONOutputWriter struct {
	W io.Writer
}

func (j JSONOutputWriter) Write(v Value) error {
	var writer = j.W
	if writer == nil {
		writer = os.Stdout
	}

	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(v); err != nil {
		return err
	}
	return nil
}

// PlainOutputWriter is output Value as plain text
type PlainOutputWriter struct {
	W io.Writer
}

func (j PlainOutputWriter) Write(v Value) error {
	var writer = j.W
	if writer == nil {
		writer = os.Stdout
	}

	if v.CO2 != nil {
		fmt.Fprintf(writer, "co2: %d ppm\n", *v.CO2)
	}
	if v.Temp != nil {
		fmt.Fprintf(writer, "temp: %d\n", int(*v.Temp))
	}

	return nil
}

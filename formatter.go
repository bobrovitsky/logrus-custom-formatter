package lcf

import (
	"bytes"
	"fmt"

	"github.com/Sirupsen/logrus"
)

const (
	// Basic formatting just logs the level name, function name, message and fields.
	Basic = `%[levelName]s:%[name]s:%[message]s%[fields]s\n`

	// Message formatting just logs the message and fields.
	Message = `%[message]s%[fields]s\n`
)

// CustomFormatter is the main formatter for the library.
type CustomFormatter struct {
	// Post-processed formatting template (e.g. `%[1]s:%[2]s:%[3]s\n`).
	Template string

	// Handler functions whose indexes match up with Template Sprintf explicit argument indexes.
	Handlers []Handler

	// Attribute names (e.g. "levelName") used in pre-processed Template.
	Attributes Attributes

	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	isTerminal bool
}

// Format is called by logrus and returns the formatted string.
func (f CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Call handlers.
	values := make([]interface{}, len(f.Handlers))
	for i, handler := range f.Handlers {
		value, err := handler(entry, &f)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	// Parse template and return.
	parsed := fmt.Sprintf(f.Template, values...)
	return bytes.NewBufferString(parsed).Bytes(), nil
}

// NewFormatter creates a new CustomFormatter, sets the Template string, and returns its pointer.
// This function is usually called just once during a running program's lifetime.
// :param template: Pre-processed formatting template (e.g. `%[message]s\n`).
// :param custom: User-defined formatters evaluated before built-in formatters. Keys are attributes to look for in the
// 	formatting string (e.g. `%[myFormatter]s`) and values are formatting functions.
func NewFormatter(template string, custom CustomHandlers) *CustomFormatter {
	formatter := CustomFormatter{isTerminal: logrus.IsTerminal()}
	formatter.Template, formatter.Handlers, formatter.Attributes = ParseTemplate(template, custom)
	return &formatter
}
package reports

import "fmt"

type ErrUnsupportedFormat struct {
	format string
}

func NewErrUnsupportedFormat(format string) ErrUnsupportedFormat {
	return ErrUnsupportedFormat{format}
}

func (e ErrUnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported report format: %s", e.format)
}

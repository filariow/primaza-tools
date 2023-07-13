package flags

import (
	"fmt"
	"strings"

	"github.com/primaza/tools/pkg/console"
)

type OutputFlag string

func (f *OutputFlag) String() string {
	return string(*f)
}

func (f *OutputFlag) Set(v string) error {
	l := strings.ToLower(v)

	aa := console.AllowedFormats()
	for _, a := range aa {
		if a == l {
			*f = OutputFlag(v)
			return nil
		}
	}
	return fmt.Errorf(
		"Invalid value for output: %s. Must be one of '%s'",
		v, strings.Join(aa, ", "))
}

func (f *OutputFlag) Type() string {
	return "OutputFlag"
}

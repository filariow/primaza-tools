package console

import (
	"encoding/json"
)

type JsonFormatter struct{}

func (f *JsonFormatter) Format(v any) ([]byte, error) {
	return json.Marshal(v)
}

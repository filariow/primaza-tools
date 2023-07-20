package console

type Format string

const (
	FormatJson    Format = "json"
	FormatMermaid Format = "mermaid"
)

func AllowedFormats() []string {
	return []string{
		string(FormatJson),
		string(FormatMermaid),
	}
}

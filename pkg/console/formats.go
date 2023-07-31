package console

type Format string

const (
	FormatHTML    Format = "html"
	FormatJson    Format = "json"
	FormatMermaid Format = "mermaid"
)

func AllowedFormats() []string {
	return []string{
		string(FormatHTML),
		string(FormatJson),
		string(FormatMermaid),
	}
}

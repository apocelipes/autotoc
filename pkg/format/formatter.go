package format

// 返回格式化后的html
type FormatterFunc func(data, indent string) string

var (
	formatters = make(map[string]FormatterFunc)
)

func SetFormatter(name string, formatter FormatterFunc) {
	if name == "" || formatter == nil {
		panic("name or formatter should not be empty")
	}

	formatters[name] = formatter
}

func NewFormatter(name string) FormatterFunc {
	formatter, ok := formatters[name]
	if !ok {
		return nil
	}

	return formatter
}

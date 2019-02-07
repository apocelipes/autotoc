package format

// 返回格式化后的html
type FormatFunc func(string) string

var (
	formatters = make(map[string]FormatFunc)
)

func SetFormatter(name string, formatter FormatFunc) {
	if name == "" || formatter == nil {
		panic("name or formatter should not be empty")
	}

	formatters[name] = formatter
}

func NewFormatter(name string) FormatFunc {
	formatter, ok := formatters[name]
	if !ok {
		return nil
	}

	return formatter
}

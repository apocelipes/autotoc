package parser

// 同时使用多个解析器解析title，返回最先解析出的结果
type multiParser struct {
	parsers []TitleParser
}

func init() {
	SetParser("multi", newMultiParser)
}

func newMultiParser(topTag string) TitleParser {
	p := &multiParser{
		parsers: make([]TitleParser, 0),
	}

	// 创建所有注册过的解析器
	for name := range parserCreators {
		if name == "multi" {
			continue
		}
		p.parsers = append(p.parsers, GetParser(name, topTag))
	}

	return p
}

func (m *multiParser) Parse(content string) *TitleNode {
	for _, p := range m.parsers {
		node := p.Parse(content)
		if node != nil {
			return node
		}
	}

	return nil
}

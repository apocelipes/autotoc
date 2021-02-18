package parser

// 同时使用多个解析器解析title，返回最先解析出的结果
type multiParser struct {
	parsers []TitleParser
}

const MultiParserName = "multi"

func init() {
	SetParser(MultiParserName, newMultiParser)
}

func newMultiParser(topTag string) TitleParser {
	p := &multiParser{
		parsers: make([]TitleParser, 0),
	}

	// 创建所有注册过的解析器
	for name := range parserCreators {
		if name == MultiParserName {
			continue
		}
		p.parsers = append(p.parsers, GetParser(name, topTag))
	}

	return p
}

func (m *multiParser) Parse(content string) *TitleNode {
	for _, p := range m.parsers {
		if node := p.Parse(content); node != nil {
			return node
		}
	}

	return nil
}

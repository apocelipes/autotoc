package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// TitleFilter 根据规则过滤标题
type TitleFilter interface {
	FilterTitleContent(content string) bool
}

type DefaultFilter struct {
	// 排除某些标题
	ExcludeTitles []string
	// 排除符合regexp匹配的title
	ExcludeRegExp *regexp.Regexp
}

// SetExcludeTitles 将由逗号分隔的一串标题分割成slice并设置
func (filter *DefaultFilter) SetExcludeTitles(titles string) {
	t := strings.Split(titles, ",")
	filter.ExcludeTitles = make([]string, 0, len(t))
	for i := range t {
		filter.ExcludeTitles = append(filter.ExcludeTitles, t[i])
	}
}

func (filter *DefaultFilter) SetExcludeRegExp(reg string) (err error) {
	if reg == "" {
		return
	}

	filter.ExcludeRegExp, err = regexp.Compile(reg)
	if err != nil {
		return fmt.Errorf("parse exclude-filter error: %v", err)
	}

	return
}

// FilterTitleContent 根据title的内容过滤标题
// 返回true表示title需要被过滤
func (filter *DefaultFilter) FilterTitleContent(content string) bool {
	for i := range filter.ExcludeTitles {
		if content == filter.ExcludeTitles[i] {
			return true
		}
	}

	if filter.ExcludeRegExp == nil {
		return false
	}

	return filter.ExcludeRegExp.MatchString(content)
}

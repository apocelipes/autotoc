package utils

import "flag"

// StringFlagWithShortName 返回一个绑定了长参数名和短参数名的flag处理器
func StringFlagWithShortName(longName, shortName, defaultValue, usage string) *string {
	p := flag.String(longName, defaultValue, usage)
	flag.StringVar(p, shortName, defaultValue, usage)

	return p
}

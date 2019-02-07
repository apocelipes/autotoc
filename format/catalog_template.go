package format

import "fmt"

// 默认的目录模板
const catalogTemplate = `<blockquote id="%s">
<h4>%s</h4>
<ul>
%s
</ul>
</blockquote>
`

// 将生成的目录嵌入catalog模板
// TODO 自定义模板
func RenderCatalog(catalogId, catalogTitle, catalogContent string) string {
	return fmt.Sprintf(catalogTemplate, catalogId, catalogTitle, catalogContent)
}

# autotoc
autoc是一个帮助你为markdown文档生成目录的自动化工具。

autotoc遵守gfm格式规范，可用于为支持gfm格式markdown的站点的内容作者为自己的文档生成易于读者阅读和跳转的文档内容目录，将你从繁重的重复劳动中解放出来。

## Installation
```bash
# 保证你的$GOPATH/bin在$PATH中
export PATH=$PATH:$GOPATH/bin
go get github.com/apocelipes/autotoc
```

## 快速入门
autotoc支持从标准输入或文件中读取信息，然后将根据文档内容生成的目录信息输出至标准输出或是写入文件中的`[TOC]`标志的地方。

文档的标题可以是html格式或是markdown格式，autotoc将会自动识别。

如果文件中存在`[TOC]`标志，autotoc将从标志的下一行开始读取文档信息。

你可以通过`--help`或`-h`选项获取帮助信息：

```bash
$ autotoc --help

Usage: ./autotoc [option]... <file>

读入file，根据其内容生成目录结构。
未提供file参数时默认读取stdin。

可选参数：
-t string, --top-tag string
        设置作为目录顶层项的tag，将从指定tag开始解析标题 (default: "h2")
-f string, --formatter string
        选择格式化html代码的方式，目前只支持default和prettyprint(output为markdown时不支持) (default: "default")
--catalog-id string
        目录的html id(output为markdown时不支持) (default: "bookmark")
--title string
        目录的标题 (default: "本文索引")
-o string, --output string
        输出的目录格式，可以为html或md(markdown) (default: "html")
-l string, --title-language string
        扫描文件的标题语法类型，可以为html，md或multi（multi同时支持所有type） (default: "multi")
-i string, --indent string
        目录的缩进，默认为2空格，输入\t以替代tab (default: 2空格)
-m string, --toc-mark string
        指定文件中写入目录的位置 (default: "[TOC]")
--exclude-title=[title1,title2]
        过滤掉内容等于参数指定值的标题 (default: "参考")
--exclude-filter=[pattern]
        过滤掉内容和参数指定的表达式匹配的标题 (default: "")
-w      是否将目录写入文件指定位置
--no-exclude    不过滤任何标题
-h, --help      显示本帮助信息并终止程序
```

假设我们有一个名为`example.md`的文件，它的内容如下：
```markdown
[TOC]
# 主标题

主标题的描述

## 一级次标题1

一级次标题1的描述

### 二级次标题1

二级次标题1的描述

### 二级次标题2

二级次标题2的描述

## 一级次标题2

一级次标题2的描述
```

输出HTML形式的并格式化的目录：
```bash
$ autotoc -t h1 example.md

<blockquote id="bookmark">
  <h4>本文索引</h4>
  <ul>
    <li>
      <a href="#主标题">主标题</a>
      <ul>
        <li>
          <a href="#一级次标题1">一级次标题1</a>
          <ul>
            <li><a href="#二级次标题1">二级次标题1</a></li>
            <li><a href="#二级次标题2">二级次标题2</a></li>
          </ul>
        </li>
        <li><a href="#一级次标题2">一级次标题2</a></li>
      </ul>
    </li>
  </ul>
</blockquote>
```

输出markdown格式的目录：
```bash
$ autotoc -t h1 -o md example.md

#### 本文索引:
- [主标题](#主标题)
  - [一级次标题1](#一级次标题1)
    - [二级次标题1](#二级次标题1)
    - [二级次标题2](#二级次标题2)
  - [一级次标题2](#一级次标题2)
```

autotoc能够从标准输入读取信息，因此可以在管道中组合使用它：
```bash
$ cat example.md | autotoc -t h1 -o md

#### 本文索引:
- [主标题](#主标题)
  - [一级次标题1](#一级次标题1)
    - [二级次标题1](#二级次标题1)
    - [二级次标题2](#二级次标题2)
  - [一级次标题2](#一级次标题2)
```

使用`-w`参数可以让目录信息写入到文件中写有`[TOC]`标志的地方：
```bash
$ autotoc -t h1 -o md -w example.md
$ head -n 10 example.md

#### 本文索引:
- [主标题](#主标题)
  - [一级次标题1](#一级次标题1)
    - [二级次标题1](#二级次标题1)
    - [二级次标题2](#二级次标题2)
  - [一级次标题2](#一级次标题2)

# 主标题

主标题的描述
```

注意，autotoc在生成目录到章节对应的链接时会将除了unicode字母和数字之外的其他字符去除。这一做法在不支持gfm格式或是对gfm兼容性较弱的地方可能导致目录无法跳转或其他问题。

## TODO
- 支持完整的gfm语法
- 补充parser的单元测试
- 自定义包裹目录信息的模板
- 添加verbose模式

欢迎提交issue或PR！

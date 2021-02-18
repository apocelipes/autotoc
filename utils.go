package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// StringFlagWithShortName 返回一个绑定了长参数名和短参数名的flag处理器
func StringFlagWithShortName(longName, shortName, defaultValue, usage string) *string {
	p := flag.String(longName, defaultValue, usage)
	flag.StringVar(p, shortName, defaultValue, usage)

	return p
}

// WriteCatalog 控制目录和文件信息的写入方向
func WriteCatalog(source *os.File, catalog, tocMark string, outputStdout bool) error {
	if outputStdout {
		return WriteStdout(catalog, tocMark, source)
	}

	return WriteBackFile(catalog, tocMark, source)
}

func hasTocMark(file *os.File, tocMark string) bool {
	scanner := bufio.NewScanner(file)
	defer file.Seek(0, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == tocMark {
			return true
		}
	}

	return false
}

// 将catalog拼接到文件内容的tocMark处，
// 如果文件内容中不存在tocMark，就拼接在内容的开头
func concatCatalog(hasToc bool, catalog, tocMark, fileData string) string {
	if hasToc {
		// 先删除catalog多余的换行符，因为tocMark所在位置已经存在一个换行符
		catalog = strings.TrimRight(catalog, "\n")
		return strings.Replace(fileData, tocMark, catalog, 1)
	}

	return catalog + fileData
}

func combine2File(file *os.File, catalog, tocMark string) (string, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	hasToc := hasTocMark(file, tocMark)
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return concatCatalog(hasToc, catalog, tocMark, string(data)), nil
}

// WriteStdout 将目录和文件内容写入标准输出
func WriteStdout(catalog, tocMark string, source *os.File) error {
	data, err := combine2File(source, catalog, tocMark)
	if err != nil {
		return err
	}

	_, err = os.Stdout.WriteString(data)
	return err
}

// WriteBackFile 将目录写回文件指定位置
func WriteBackFile(catalog, tocMark string, file *os.File) error {
	filePath, err := filepath.Abs(file.Name())
	if err != nil {
		return err
	}
	backupName := filepath.Join(filepath.Dir(filePath),
		".backup_"+filepath.Base(filePath))
	backup, err := os.OpenFile(backupName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = io.Copy(backup, file)
	if err != nil {
		return err
	}
	err = backup.Close()
	if err != nil {
		return err
	}

	fullData, err := combine2File(file, catalog, tocMark)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.WriteString(fullData)
	if err != nil {
		return err
	}

	// 成功写回，删除备份
	err = os.Remove(backupName)
	if err != nil {
		return err
	}

	return nil
}

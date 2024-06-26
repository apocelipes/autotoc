package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
)

// WriteCatalog 控制目录和文件信息的写入方向
func WriteCatalog(data []byte, catalog, tocMark string, outputStdout bool, fileName string) error {
	if outputStdout {
		return WriteStdout([]byte(catalog), []byte(tocMark), data)
	}

	return WriteBackFile([]byte(catalog), []byte(tocMark), data, fileName)
}

func insertCatalogToFile(fileData, catalog, tocMark []byte) []byte {
	hasToc := len(tocMark) != 0 && bytes.Contains(fileData, tocMark)
	if hasToc {
		// 先删除catalog多余的换行符，因为tocMark所在位置已经存在一个换行符
		catalog = bytes.TrimRight(catalog, "\n")
		return bytes.Replace(fileData, tocMark, catalog, 1)
	}

	return slices.Concat(catalog, fileData)
}

// WriteStdout 将目录和文件内容写入标准输出
func WriteStdout(catalog, tocMark, source []byte) error {
	data := insertCatalogToFile(source, catalog, tocMark)

	_, err := os.Stdout.Write([]byte(data))
	return err
}

// WriteBackFile 将目录写回文件指定位置
func WriteBackFile(catalog, tocMark, fileData []byte, fileName string) error {
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	backupName := filepath.Join(filepath.Dir(filePath), ".backup_"+filepath.Base(filePath))
	if err := os.WriteFile(backupName, fileData, 0644); err != nil {
		return err
	}

	fullData := insertCatalogToFile(fileData, catalog, tocMark)
	if err := os.WriteFile(fileName, fullData, 0644); err != nil {
		return err
	}

	// 成功写回，删除备份
	if err := os.Remove(backupName); err != nil {
		return err
	}

	return nil
}

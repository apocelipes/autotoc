package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

// 返回一个绑定了长参数名和短参数名的flag处理器
func StringFlagWithShortName(longName, shortName, defaultValue, usage string) *string {
	p := flag.String(longName, defaultValue, usage)
	flag.StringVar(p, shortName, defaultValue, usage)

	return p
}

// 将目录写回文件指定位置
func WriteBackFile(catalog, tocMark string, file *os.File) error {
	backupName := ".backup_" + file.Name()
	backup, err := os.OpenFile(backupName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer backup.Close()
	file.Seek(0, 0)
	_, err = io.Copy(backup, file)
	if err != nil {
		return err
	}

	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	buffer := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == tocMark {
			buffer.WriteString(catalog + "\n")
			continue
		}

		buffer.WriteString(line + "\n")
	}

	err = file.Truncate(0)
	file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.WriteString(buffer.String())
	if err != nil {
		return err
	}
	// 成功写回，删除备份
	os.Remove(backupName)

	return nil
}

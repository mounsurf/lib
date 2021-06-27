package file

import (
	"bufio"
	"embed"
	"strings"
)

func ReadLinesByFs(fs embed.FS, path string) ([]string, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimRight(scanner.Text(), "\r\n"))
	}
	_ = f.Close()
	return lines, scanner.Err()
}

func ReadBytesByFs(fs embed.FS, path string) ([]byte, error) {
	return fs.ReadFile(path)
}

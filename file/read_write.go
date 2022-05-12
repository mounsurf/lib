package file

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func WriteFile(filename string, content []byte) (int, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(content)
}

func WriteFileWithString(filename string, content string) (int, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.WriteString(content)
}

func ReadLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimRight(scanner.Text(), "\r\n"))
	}
	f.Close()
	return lines, scanner.Err()
}

func ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(fileName))
}

func Append2File(file, text string) (int, error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.WriteString(text)
}

func GetCurDir() string {
	curDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return curDir + "/"
}

func EmptyFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return f.Close()
}

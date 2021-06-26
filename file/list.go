package file

import (
	"io/ioutil"
	"strings"
)

func GetDirFiles(dir string, containsSubPath bool) ([]string, error) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, file := range files {
		if file.IsDir() && containsSubPath {
			curDirFiles, err := GetDirFiles(dir+file.Name(), containsSubPath)
			if err == nil {
				result = append(result, curDirFiles...)
			}
		} else {
			result = append(result, dir+file.Name())
		}
	}
	return result, nil
}

func GetFileName(fileNameWithPath string) string {
	index := strings.LastIndex(fileNameWithPath, "/")
	if index < 0 {
		return fileNameWithPath
	} else if index == len(fileNameWithPath)-1 {
		return ""
	} else {
		return fileNameWithPath[index+1:]
	}
}

func GetFileExtension(fileName string) string {
	index := strings.LastIndex(fileName, ".")
	if index < 0 || index == len(fileName)-1 {
		return ""
	} else {
		return fileName[index+1:]
	}
}

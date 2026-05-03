package res

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Res struct{}

const defaultDir = ".2fa/"

var defaultPrefix string

func exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
func isDir(path string) (bool, error) {
	info, err := os.Stat(path)

	if err == nil {
		return info.IsDir(), nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
func Exists(path string) (bool, error) {
	return exists(defaultPrefix + path)
}

func create(filename string) error {
	if ok, err := exists(filename); !ok && err == nil {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
func mkdir(dirname string) error {
	if ok, err := isDir(dirname); !ok && err == nil {
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
func CreateFile(filename string) error {
	return create(defaultPrefix + filename)
}

func write(filename string, content string) error {
	if ok, err := exists(filename); ok && err == nil {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(content)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
func WriteFile(filename string, content string) error {
	return write(defaultPrefix+filename, content)
}
func read(filename string) (*string, error) {
	var content string
	if ok, err := exists(filename); ok && err == nil {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		content = string(bytes)
	} else {
		return nil, err
	}
	return &content, nil
}
func ReadFile(filename string) (*string, error) {
	return read(defaultPrefix + filename)
}

func extend(filename string, content string) error {
	if ok, err := exists(filename); ok && err == nil {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(content)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func init() {
	home, _ := os.UserHomeDir()
	defaultPrefix = home + "/" + defaultDir
	var exists, err = isDir(defaultPrefix)
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Println("dir exists")
		return
	}
	fmt.Println("dir doesn't exists")
	fmt.Println("creating new one", defaultPrefix)
	err = mkdir(defaultPrefix)
	if err != nil {
		panic(err)
	}
}

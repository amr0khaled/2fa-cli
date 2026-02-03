package res

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Res struct{}

var defaultDir = ".2fa/"

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

func create(filename string) error {
	if ok, err := exists(filename); ok && err == nil {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
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
	dir := home + "/" + defaultDir
	var exists, _ = exists(dir)
	if exists {
		fmt.Println("dir exists")
	} else {
		fmt.Println("dir doesn't exists")
		fmt.Println("creating new one")
		create(dir)
	}
}

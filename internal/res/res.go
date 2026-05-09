package res

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic("Encountered an issue: " + err.Error())
	}
}

type Res struct {
	dir string
}

var res *Res = NewRes()

const defaultDir = ".2fa/"

func NewRes() *Res {
	home, _ := os.UserHomeDir()
	defaultPrefix, _ := filepath.Abs(home + "/" + defaultDir)
	fmt.Printf("%s\n", defaultPrefix)
	return &Res{
		dir: defaultPrefix,
	}
}

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
	return exists(Join(res.dir, path))
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
func mkdir(dirnames ...string) error {
	if ok, err := isDir(Join(dirnames...)); !ok && err == nil {
		err := os.MkdirAll(Join(dirnames...), 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
func CreateFile(filename string) error {
	return create(Join(res.dir, filename))
}

func write(filename string, content string) error {
	if _, err := exists(filename); err == nil {
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
	return write(Join(res.dir, filename), content)
}
func read(filename string) (string, error) {
	var content string
	if ok, err := exists(filename); ok && err == nil {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}
		content = string(bytes)
	} else {
		return "", err
	}
	return content, nil
}
func Join(paths ...string) string {
	path := []string{}
	return strings.Join(append(path, paths...), "/")
}
func ReadFile(filenames ...string) (string, error) {
	return read(Join(res.dir, Join(filenames...)))
}
func del(filename string) error {
	if ok, err := exists(filename); ok && err == nil {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
func DeleteFile(filename string) error {
	return del(Join(res.dir, filename))
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

type KeyFile struct {
	Content string
}

func readKeys() map[string]KeyFile {
	var keys map[string]KeyFile = map[string]KeyFile{}
	filepath.WalkDir(res.dir, func(path string, d fs.DirEntry, err error) error {
		check(err)
		if !d.IsDir() {
			filename := filepath.Base(path)
			if filename != ".key" {
				return nil
			}
			file, ok := keys[path]
			if !ok {
				abs, _ := filepath.Abs(path)
				content, err := read(abs)
				check(err)
				file = KeyFile{
					Content: content,
				}
			}
			filerel, _ := strings.CutPrefix(path, res.dir+"/")
			entry, _ := strings.CutSuffix(filerel, "/.key")
			keys[entry] = file
		}
		return nil
	})
	return keys
}

func init() {
	var exists, err = isDir(res.dir)
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Println("dir exists")
		return
	}
	fmt.Println("dir doesn't exists")
	fmt.Println("creating new one", res.dir)
	err = mkdir(res.dir)
	if err != nil {
		panic(err)
	}
}

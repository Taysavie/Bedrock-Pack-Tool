package rar

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nwaples/rardecode"
)

var blacklist = []string{".DS_Store", "desktop.ini", "Thumbs.db"}

func Unrar(in []byte) (out map[string][]byte, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("unexpected error processing RAR")
		}
	}()
	out = make(map[string][]byte)
	var single = true
	var last = ""
	archive, err := rardecode.NewReader(bytes.NewReader(in), "")
	if err != nil {
		fmt.Println(err)
		return map[string][]byte{}, err
	}
	for {
		header, err := archive.Next()
		if archive == nil || header == nil {
			break
		}
		if header.IsDir {
			continue
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return map[string][]byte{}, err
			}
		}
		filename := strings.ReplaceAll(header.Name, "\\", "/")
		baseName := filepath.Base(filename)
		if slices.Contains(blacklist, baseName) {
			continue
		}
		if strings.HasPrefix(baseName, "._") {
			continue
		}
		if strings.Contains(filename, "__MACOSX") {
			continue
		}
		name := strings.TrimPrefix(filename, "/")
		base := strings.Split(name, "/")[0]
		if single && last != "" && base != last {
			single = false
		}
		last = base
		data, err := io.ReadAll(archive)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return map[string][]byte{}, err
			}
		}
		out[name] = data
	}
	if single {
		oout := out
		out = make(map[string][]byte)
		for name, data := range oout {
			newBase := strings.Join(strings.Split(name, "/")[1:], "/")
			out[newBase] = data
		}
	}
	return out, nil
}

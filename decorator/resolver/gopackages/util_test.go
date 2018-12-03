package gopackages_test

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func tempDir(m map[string]string) (dir string, err error) {
	if dir, err = ioutil.TempDir("", ""); err != nil {
		return
	}
	for fpathrel, src := range m {
		if strings.HasSuffix(fpathrel, "/") {
			// just a dir
			if err = os.MkdirAll(filepath.Join(dir, fpathrel), 0777); err != nil {
				return
			}
		} else {
			fpath := filepath.Join(dir, fpathrel)
			fdir, _ := filepath.Split(fpath)
			if err = os.MkdirAll(fdir, 0777); err != nil {
				return
			}

			var formatted []byte
			if strings.HasSuffix(fpath, ".go") {
				formatted, err = format.Source([]byte(src))
				if err != nil {
					err = fmt.Errorf("formatting %s: %v", fpathrel, err)
					return
				}
			} else {
				formatted = []byte(src)
			}

			if err = ioutil.WriteFile(fpath, formatted, 0666); err != nil {
				return
			}
		}
	}
	return
}

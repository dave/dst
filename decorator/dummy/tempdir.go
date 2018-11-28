package dummy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func TempDir(m map[string]string) (dir string, err error) {
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
			if err = ioutil.WriteFile(fpath, []byte(src), 0666); err != nil {
				return
			}
		}
	}
	return
}

package dummy

import (
	"io/ioutil"

	"gopkg.in/src-d/go-billy.v4/osfs"
)

func TempDir(root Dir) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	fs := osfs.New(dir)
	for name, item := range root {
		item.create("/", name, fs)
	}
	return dir
}

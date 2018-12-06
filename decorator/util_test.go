package decorator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
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

func compareDir(t *testing.T, dir string, expect map[string]string) {
	t.Helper()
	found := map[string]string{}
	walk := func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relfpath, err := filepath.Rel(dir, fpath)
		if err != nil {
			t.Fatal(err)
		}
		b, err := ioutil.ReadFile(fpath)
		if err != nil {
			t.Fatal(err)
		}
		found[relfpath] = string(b)
		return nil
	}
	if err := filepath.Walk(dir, walk); err != nil {
		t.Fatal(err)
	}
	var keysFound []string
	var keysExpect []string
	for k := range found {
		keysFound = append(keysFound, k)
	}
	for k := range expect {
		keysExpect = append(keysExpect, k)
	}
	sort.Strings(keysFound)
	sort.Strings(keysExpect)
	keysFoundJoined := strings.Join(keysFound, " ")
	keysExpectJoined := strings.Join(keysExpect, " ")
	t.Run("files", func(t *testing.T) {
		compare(t, keysExpectJoined, keysFoundJoined)
	})
	done := map[string]bool{}
	for k, v := range found {
		if done[k] {
			continue
		}
		t.Run(k, func(t *testing.T) {
			if strings.HasSuffix(k, ".go") {
				compareSrc(t, expect[k], v)
			} else {
				compare(t, strings.TrimSpace(expect[k]), strings.TrimSpace(v))
			}
		})
	}
}

func compare(t *testing.T, expect, found string) {
	t.Helper()
	if expect != found {
		t.Errorf("\nexpect: %q\nfound : %q", expect, found)
	}
}

func compareSrc(t *testing.T, expect, found string) {
	t.Helper()
	bFound, err := format.Source([]byte(found))
	if err != nil {
		t.Fatal(err)
	}
	bExpect, err := format.Source([]byte(expect))
	if err != nil {
		t.Fatal(err)
	}
	expect = string(bExpect)
	found = string(bFound)
	if expect != found {
		t.Errorf("\nexpect: %q\nfound : %q", expect, found)
	}
}

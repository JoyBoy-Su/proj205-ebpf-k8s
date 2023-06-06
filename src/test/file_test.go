package test

import (
	"os"
	"strings"
	"testing"
)

func TestReadLink(t *testing.T) {
	dst_path := "/home/ubuntu/jiadisu/data"
	link_info, err := os.Lstat(dst_path)
	if err != nil {
		panic(err)
	}
	if link_info.Mode()&os.ModeSymlink != 0 {
		src_file, err := os.Readlink(link_info.Name())
		if err != nil {
			panic(err)
		}
		if strings.Compare(src_file, "bpf") != 0 {
			t.Fatal("src file error, expected bpf, got ", src_file)
		}
	}
}

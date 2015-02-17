package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var test_files []os.FileInfo
var count_of_testfiles int

func init() {
	path = "."
	test_files, _ = ioutil.ReadDir(path)
	count_of_testfiles = len(test_files)
}

func Test_filter_nil(t *testing.T) {
	c1 := filter_nil(test_files)
	if len(c1) != count_of_testfiles {
		t.Errorf("%d != %d, but should be the same", len(c1), count_of_testfiles)
	}
}

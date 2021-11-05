package ghfs

import (
	"testing"
	"testing/fstest"
	"testing/iotest"
)

func TestFS(t *testing.T) {
	fsys, err := New("golang", "time")
	if err != nil {
		t.Fatal(err)
	}
	if err := fstest.TestFS(fsys, "README.md", "LICENSE", "rate/rate.go"); err != nil {
		t.Fatal(err)
	}
}

func TestIO(t *testing.T) {
	fsys, err := New("golang", "time")
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Open("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if err := iotest.TestReader(f, []byte("module golang.org/x/time\n")); err != nil {
		t.Fatal(err)
	}
}

package ghfs

import (
	"context"
	"testing"
	"testing/fstest"
	"testing/iotest"

	"github.com/k1LoW/go-github-client/v39/factory"
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

func TestOptionClient(t *testing.T) {
	client, err := factory.NewGithubClient()
	if err != nil {
		t.Fatal(err)
	}
	fsys, err := New("golang", "time", Client(client))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fsys.Open("README.md"); err != nil {
		t.Fatal(err)
	}
}

func TestOptionContext(t *testing.T) {
	fsys, err := New("golang", "time", Context(context.TODO()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fsys.Open("README.md"); err != nil {
		t.Fatal(err)
	}
}

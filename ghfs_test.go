package ghfs

import (
	"context"
	"io"
	"testing"
	"testing/fstest"
	"testing/iotest"

	"github.com/k1LoW/go-github-client/v48/factory"
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

func TestOptionBranch(t *testing.T) {
	fsys, err := New("golang", "go", Branch("release-branch.go1"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Open("VERSION")
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if want := "go1.0.3"; got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestOptionTag(t *testing.T) {
	fsys, err := New("golang", "go", Tag("go1"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Open("VERSION")
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if want := "go1"; got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

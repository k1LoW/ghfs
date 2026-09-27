package ghfs

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"testing/fstest"
	"testing/iotest"

	"github.com/google/go-github/v67/github"
	"github.com/k1LoW/go-github-client/v67/factory"
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
	fsys, err := New("golang", "time", Tag("v0.14.0"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Open("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if err := iotest.TestReader(f, []byte("module golang.org/x/time\n\ngo 1.24.0\n")); err != nil {
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

func TestOptionContextCancelBlobRead(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"default_branch":"main"}`))
	})
	mux.HandleFunc("GET /repos/o/r/branches/main", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"name":"main","commit":{"sha":"c0"}}`))
	})
	mux.HandleFunc("GET /repos/o/r/git/trees/c0", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"sha":"c0","tree":[{"path":"README.md","mode":"100644","type":"blob","sha":"b0","size":5}]}`))
	})
	mux.HandleFunc("GET /repos/o/r/git/blobs/b0", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"sha":"b0","encoding":"base64","content":"aGVsbG8=","size":5}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	client := github.NewClient(nil)
	u, err := url.Parse(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	client.BaseURL = u

	ctx, cancel := context.WithCancel(t.Context())
	fsys, err := New("o", "r", Client(client), Context(ctx))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fsys.ReadFile("README.md"); err != nil {
		t.Fatal(err)
	}

	cancel()
	if _, err := fsys.ReadFile("README.md"); !errors.Is(err, context.Canceled) {
		t.Errorf("got %v\nwant %v", err, context.Canceled)
	}
}

package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	if err := initStore(dir, "dsp2b", "https://oss.example.com"); err != nil {
		t.Fatalf("initStore: %v", err)
	}
	return Default()
}

func TestPutAndGetObject(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	key := "images/blueprint/abc.png"
	want := []byte("hello blueprint")

	if err := s.PutObject(ctx, key, bytes.NewReader(want), int64(len(want))); err != nil {
		t.Fatalf("PutObject: %v", err)
	}

	obj, err := s.GetObject(ctx, key)
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	defer obj.Close()

	got, err := io.ReadAll(obj)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("content = %q, want %q", got, want)
	}

	info, err := obj.Stat()
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size != int64(len(want)) {
		t.Fatalf("size = %d, want %d", info.Size, len(want))
	}
}

func TestGetObjectMissingReturnsNotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetObject(ctx, "images/blueprint/missing.png")
	if err == nil {
		t.Fatal("expected error for missing object")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	var respErr RespError
	if !errors.As(err, &respErr) {
		t.Fatalf("err does not implement RespError: %v", err)
	}
	if respErr.StatusCode() != 404 {
		t.Fatalf("status = %d, want 404", respErr.StatusCode())
	}
}

func TestExists(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	key := "blueprint/collection/1/a.zip"

	if ok, err := s.Exists(ctx, key); err != nil || ok {
		t.Fatalf("Exists before put = (%v, %v), want (false, nil)", ok, err)
	}
	if err := s.PutObject(ctx, key, bytes.NewReader([]byte("zip")), 3); err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	if ok, err := s.Exists(ctx, key); err != nil || !ok {
		t.Fatalf("Exists after put = (%v, %v), want (true, nil)", ok, err)
	}
}

func TestPutObjectRejectsTraversal(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	cases := []string{
		"../outside.png",
		"images/../../outside.png",
		"/../outside.png",
	}
	for _, key := range cases {
		if err := s.PutObject(ctx, key, bytes.NewReader([]byte("x")), 1); err == nil {
			t.Fatalf("PutObject(%q) should fail", key)
		}
	}

	// 确认没有在桶目录之外落盘
	if _, err := os.Stat(filepath.Join(s.root, "outside.png")); !os.IsNotExist(err) {
		t.Fatalf("file escaped bucket dir: %v", err)
	}
}

func TestPutObjectRejectsEmptyKey(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	if err := s.PutObject(ctx, "", bytes.NewReader([]byte("x")), 1); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFileURL(t *testing.T) {
	s := newTestStore(t)
	got, err := s.FileURL(context.Background(), "/images/blueprint/a.png")
	if err != nil {
		t.Fatalf("FileURL: %v", err)
	}
	want := "https://oss.example.com/dsp2b/images/blueprint/a.png"
	if got != want {
		t.Fatalf("FileURL = %q, want %q", got, want)
	}
}

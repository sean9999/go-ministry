package graph

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"

	"github.com/sean9999/harebrain"
)

type collection = string
type key = string

const ROOT = "testdata"
const NODES_FOLDER = "nodes"
const EDGES_FOLDER = "edges"

type Store[T harebrain.EncodeHasher] interface {
	Get(string) (T, error)
	Save(T) error
	Has(string) bool
	Delete(string) error
	All() []T
	AllRecords() (map[string][]byte, error)
}

func storeBase() *harebrain.Database {
	db := harebrain.NewDatabase()
	db.Open(ROOT)
	return db
}

type GraphStore struct {
	Nodes Store[*Node]
	Edges Store[*Edge]
}

// func (gs GraphStore) Restore(filename string)

func (gs GraphStore) Zip(filename string) error {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	defer w.Close()
	subfolders := []string{"nodes", "edges"}
	for _, dir := range subfolders {
		entries, err := os.ReadDir(filepath.Join(ROOT, dir))
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.Type().IsRegular() {
				f, err := w.Create(filepath.Join(dir, entry.Name()))
				if err != nil {
					return err
				}
				contents, err := os.ReadFile(filepath.Join(ROOT, dir, entry.Name()))
				if err != nil {
					return err
				}
				_, err = f.Write(contents)
				if err != nil {
					return err
				}
			}
		}
	}
	w.Close()
	os.WriteFile(filename, buf.Bytes(), 0644)
	return nil
}

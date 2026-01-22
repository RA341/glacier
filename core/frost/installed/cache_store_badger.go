package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

type CacheStoreBadger struct {
	db *badger.DB
}

func NewCacheStoreBadger(dbPath string) (*CacheStoreBadger, error) {
	opts := badger.DefaultOptions(dbPath).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}
	return &CacheStoreBadger{db: db}, nil
}

func (c *CacheStoreBadger) Close() error {
	return c.db.Close()
}

// GetFileList returns all unique filenames stored
func (c *CacheStoreBadger) GetFileList() ([]string, error) {
	var files []string
	err := c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("f:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := string(it.Item().Key())
			files = append(files, strings.TrimPrefix(key, "f:"))
		}
		return nil
	})
	return files, err
}

// GetChunkLen returns how many chunks are registered for a file
func (c *CacheStoreBadger) GetChunkLen(file string) (int, error) {
	count := 0
	err := c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(fmt.Sprintf("c:%s:", file))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

// Add initializes the file entry and all its chunks
func (c *CacheStoreBadger) Add(file string, chunks []Chunk) error {
	return c.db.Update(func(txn *badger.Txn) error {
		// mark the file as existing
		fileKey := []byte(fmt.Sprintf("f:%s", file))
		if err := txn.Set(fileKey, []byte{1}); err != nil {
			return err
		}

		// add all chunks
		for i, chunk := range chunks {
			// Using %010d ensures indexes are sorted numerically (0000000001, 0000000002)
			key := []byte(fmt.Sprintf("c:%s:%010d", file, i))
			data, _ := json.Marshal(chunk)
			if err := txn.Set(key, data); err != nil {
				return err
			}
		}
		return nil
	})
}

// Get retrieves all chunks for a specific file in order
func (c *CacheStoreBadger) Get(file string) (chunks []Chunk, found bool, err error) {
	err = c.db.View(func(txn *badger.Txn) error {
		// check if file exists
		fileKey := []byte(fmt.Sprintf("f:%s", file))
		_, err := txn.Get(fileKey)
		if errors.Is(err, badger.ErrKeyNotFound) {
			found = false
			return nil
		}
		found = true

		// Iterate through chunks
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(fmt.Sprintf("c:%s:", file))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var chunk Chunk
			err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &chunk)
			})
			if err != nil {
				return err
			}
			chunks = append(chunks, chunk)
		}
		return nil
	})
	return chunks, found, err
}

// Update updates a specific chunk by its index
func (c *CacheStoreBadger) Update(file string, index int, chunk *Chunk) error {
	return c.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("c:%s:%010d", file, index))

		if _, err := txn.Get(key); errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("chunk index %d not found for file %s", index, file)
		}

		data, _ := json.Marshal(chunk)
		return txn.Set(key, data)
	})
}

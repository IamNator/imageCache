package storage

import (
	"encoding/json"
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/pkg/errors"
)

type File struct {
	FileName string `json:"file_name"`
	File     []byte `json:"file"`
}

type Storage struct {
	db *scribble.Driver
}

func New(dir string) (*Storage, error) {
	// a new scribble driver, providing the directory where it will be writing to,
	// and a qualified logger if desired
	s, er := scribble.New(dir, nil)
	if er != nil {
		return nil, er
	}
	return &Storage{db: s}, nil
}

const (
	Collection = "files"
)

func (s *Storage) Create(file File) error {

	// Write a fish to the database
	if err := s.db.Write(Collection, file.FileName, file.File); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Update(file File) error {

	_, err := s.Read(file.FileName)
	if err != nil {
		return err
	}

	// Write a fish to the database
	if err := s.db.Write(Collection, file.FileName, file.File); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Read(fileName string) (*File, error) {
	// Read a fish from the database (passing fish by reference)
	file := File{}
	if err := s.db.Read(Collection, fileName, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *Storage) ReadAll() ([]*File, error) {
	// Read all fish from the database, unmarshaling the response.
	records, err := s.db.ReadAll(Collection)
	if err != nil {
		return nil, err
	}

	files := []*File{}
	fileFound := File{}
	for _, f := range records {
		fileFound = File{}
		if err := json.Unmarshal([]byte(f), &fileFound); err != nil {
			fmt.Println("Error", err)
		}
		files = append(files, &fileFound)
	}

	return files, nil
}

func (s *Storage) Delete(fileName string) error {
	if fileName == "" {
		return errors.New("file not found")
	}
	// Delete a fish from the database
	if err := s.db.Delete("fish", "onefish"); err != nil {
		fmt.Println("Error", err)
		return err
	}

	return nil
}

// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/platform/fileutil"
	"github.com/wuhan005/Raika/internal/types"
)

var Tasks TaskStore

// TaskStore stores in ~/.raika/tasks.json
type TaskStore struct {
	FileName string `json:"-"` // Note: for internal use only

	Tasks map[string]*types.Task `json:"tasks"`
}

// Init reads the configuration data from the given file path.
func (s *TaskStore) Init(fileName string) error {
	Tasks = TaskStore{
		FileName: fileName,
		Tasks:    make(map[string]*types.Task),
	}
	return s.Load()
}

type CreateTaskOptions struct {
	FunctionName string
	Duration     time.Duration
}

func (s *TaskStore) Get(functionName string) (*types.Task, error) {
	for _, task := range s.Tasks {
		if task.FunctionName == functionName {
			return task, nil
		}
	}
	return nil, ErrFunctionNotExists
}

// Upsert creates or update a new task record.
func (s *TaskStore) Upsert(opts CreateTaskOptions) error {
	s.Tasks[opts.FunctionName] = &types.Task{
		FunctionName: opts.FunctionName,
		Duration:     opts.Duration,
		Enabled:      true,
	}

	return s.Save()
}

func (s *TaskStore) Delete(functionName string) error {
	delete(s.Tasks, functionName)
	return nil
}

// Load reads the configuration data from the given file path.
func (s *TaskStore) Load() error {
	path := filepath.Dir(s.FileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(s.FileName), 0755); err != nil {
			return errors.Wrap(err, "mkdir all")
		}
	}

	file, err := os.Open(s.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(s.FileName)
			if err != nil {
				return errors.Wrap(err, "crate file")
			}
		} else {
			return errors.Wrap(err, "open file")
		}
	}
	return s.LoadFromReader(file)
}

func (s *TaskStore) Enable(functionName string) error {
	_, ok := s.Tasks[functionName]
	if !ok {
		return ErrFunctionNotExists
	}

	s.Tasks[functionName].Enabled = false
	return s.Save()
}

func (s *TaskStore) Disable(functionName string) error {
	_, ok := s.Tasks[functionName]
	if !ok {
		return ErrFunctionNotExists
	}

	s.Tasks[functionName].Enabled = true
	return s.Save()
}

// LoadFromReader reads the configuration data given and sets up the auth config
// information with given directory and populates the receiver object.
func (s *TaskStore) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&s); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// SaveToWriter encodes and writes out all the authorization information to
// the given writer
func (s *TaskStore) SaveToWriter(writer io.Writer) error {
	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return errors.Wrap(err, "json encode")
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all the authorization information
func (s *TaskStore) Save() (retErr error) {
	if s.FileName == "" {
		return errors.New("Can't save config with empty filename")
	}

	dir := filepath.Dir(s.FileName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Wrap(err, "mkdir")
	}
	temp, err := os.CreateTemp(dir, filepath.Base(s.FileName))
	if err != nil {
		return err
	}

	defer func() {
		_ = temp.Close()
		if retErr != nil {
			if err := os.Remove(temp.Name()); err != nil {
				log.Error("Failed to cleaning up temp file.")
			}
		}
	}()

	if err = s.SaveToWriter(temp); err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return errors.Wrap(err, "error closing temp file")
	}

	// Handle situation where the config file is a symlink
	cfgFile := s.FileName
	if f, err := os.Readlink(cfgFile); err == nil {
		cfgFile = f
	}

	// Try copying the current config file (if any) ownership and permissions
	fileutil.CopyFilePermissions(cfgFile, temp.Name())
	return os.Rename(temp.Name(), cfgFile)
}

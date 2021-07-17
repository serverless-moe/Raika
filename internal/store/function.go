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

	"github.com/wuhan005/Raika/internal/platform"
	"github.com/wuhan005/Raika/internal/platform/fileutil"
	"github.com/wuhan005/Raika/internal/types"
)

var Functions FunctionStore

// FunctionStore stores in ~/.raika/functions.json
type FunctionStore struct {
	FileName string `json:"-"` // Note: for internal use only

	Functions map[string][]types.Function `json:"functions"`
}

// Init reads the configuration data from the given file path.
func (s *FunctionStore) Init(fileName string) error {
	Functions = FunctionStore{
		FileName:  fileName,
		Functions: make(map[string][]types.Function),
	}
	return s.Load()
}

// Set creates a new function record.
func (s *FunctionStore) Set(functionName string, platformID string, triggerURL string, opts platform.CreateFunctionOptions) error {
	if s.Functions[functionName] == nil {
		s.Functions[functionName] = make([]types.Function, 0)
	}

	f := types.Function{
		PlatformID:            platformID,
		URL:                   triggerURL,
		CreatedAt:             time.Now(),
		Name:                  opts.Name,
		Description:           opts.Description,
		MemorySize:            opts.MemorySize,
		Environment:           opts.EnvironmentVariables,
		InitializationTimeout: opts.InitializationTimeout,
		RuntimeTimeout:        opts.RuntimeTimeout,
		HTTPPort:              opts.HTTPPort,
		File:                  opts.File,
	}

	for k, function := range s.Functions[functionName] {
		if function.PlatformID == platformID {
			s.Functions[functionName][k] = f
			return s.Save()
		}
	}

	s.Functions[functionName] = append(s.Functions[functionName], f)
	return s.Save()
}

func (s *FunctionStore) Get(functionName string) ([]types.Function, error) {
	function, ok := s.Functions[functionName]
	if !ok {
		return nil, ErrFunctionNotExists
	}
	return function, nil
}

// Load reads the configuration data from the given file path.
func (s *FunctionStore) Load() error {
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

// LoadFromReader reads the configuration data given and sets up the auth config
// information with given directory and populates the receiver object.
func (s *FunctionStore) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&s); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// SaveToWriter encodes and writes out all the authorization information to
// the given writer
func (s *FunctionStore) SaveToWriter(writer io.Writer) error {
	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return errors.Wrap(err, "json encode")
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all the authorization information
func (s *FunctionStore) Save() (retErr error) {
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

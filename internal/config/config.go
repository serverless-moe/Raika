// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/types"
)

var HomePath, _ = os.UserHomeDir()
var DefaultConfigPath = filepath.Join(HomePath, "./.raika/config.json")

// File stores in ~/.raika/config.json
type File struct {
	FileName string `json:"-"` // Note: for internal use only

	AuthConfigs map[string]types.AuthConfig `json:"auths"`
}

// New initializes an empty configuration file for the given filename 'fileName'.
func New(fileName string) *File {
	return &File{
		FileName:    fileName,
		AuthConfigs: make(map[string]types.AuthConfig),
	}
}

// Load reads the configuration data from the given file path.
func (f *File) Load() error {
	path := filepath.Dir(f.FileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(f.FileName), 0755); err != nil {
			return errors.Wrap(err, "mkdir all")
		}
	}

	file, err := os.Open(f.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(f.FileName)
			if err != nil {
				return errors.Wrap(err, "crate file")
			}
		} else {
			return errors.Wrap(err, "open file")
		}
	}
	return f.LoadFromReader(file)
}

// LoadFromReader reads the configuration data given and sets up the auth config
// information with given directory and populates the receiver object.
func (f *File) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&f); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// SaveToWriter encodes and writes out all the authorization information to
// the given writer
func (f *File) SaveToWriter(writer io.Writer) error {
	data, err := json.MarshalIndent(f, "", "\t")
	if err != nil {
		return errors.Wrap(err, "json encode")
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all the authorization information
func (f *File) Save() (retErr error) {
	if f.FileName == "" {
		return errors.New("Can't save config with empty filename")
	}

	dir := filepath.Dir(f.FileName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Wrap(err, "mkdir")
	}
	temp, err := os.CreateTemp(dir, filepath.Base(f.FileName))
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

	if err = f.SaveToWriter(temp); err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return errors.Wrap(err, "error closing temp file")
	}

	// Handle situation where the confie file is a symlink
	cfgFile := f.FileName
	if f, err := os.Readlink(cfgFile); err == nil {
		cfgFile = f
	}

	// Try copying the current config file (if any) ownership and permissions
	copyFilePermissions(cfgFile, temp.Name())
	return os.Rename(temp.Name(), cfgFile)
}

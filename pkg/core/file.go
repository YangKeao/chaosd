// Copyright 2022 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"encoding/json"

	"github.com/pingcap/errors"
)

type FileCommand struct {
	CommonAttackConfig

	//	FileName is the name of the file to be created, modified, deleted, renamed, or appended.
	FileName string `json:"file-name,omitempty"`
	// DirName is the directory name to create or delete.
	DirName string `json:"dir-name,omitempty"`
	// Privilege is the file privilege to be set.
	Privilege uint32 `json:"privilege,omitempty"`
	// SourceFile is the name need to be renamed.
	SourceFile string `json:"source-file,omitempty"`
	// DestFile is the name to be renamed.
	DestFile string `json:"dest-file,omitempty"`
	// Data is the data for append.
	Data string `json:"data,omitempty"`
	// Count is the number of times to append the data.
	Count int `json:"count,omitempty"`
	// OriginPrivilege used to save the file's origin privilege.
	OriginPrivilege int `json:"origin-privilege,omitempty"`

	// OriginStr is the origin string of the file.
	OriginStr string `json:"origin-string,omitempty"`
	// DestStr is the destination string of the file.
	DestStr string `json:"dest-string,omitempty"`
	// Line is the line number of the file to be replaced.
	Line int `json:"line,omitempty"`
}

var _ AttackConfig = &FileCommand{}

const (
	FileCreateAction          = "create"
	FileModifyPrivilegeAction = "modify"
	FileDeleteAction          = "delete"
	FileRenameAction          = "rename"
	FileAppendAction          = "append"
	FileReplaceAction         = "replace"
)

func (n *FileCommand) Validate() error {
	if err := n.CommonAttackConfig.Validate(); err != nil {
		return err
	}
	switch n.Action {
	case FileCreateAction:
		return n.validFileCreate()
	case FileModifyPrivilegeAction:
		return n.validFileModify()
	case FileDeleteAction:
		return n.validFileDelete()
	case FileRenameAction:
		return n.validFileRename()
	case FileAppendAction:
		return n.validFileAppend()
	case FileReplaceAction:
		return n.validFileReplace()
	default:
		return errors.Errorf("file action %s not supported", n.Action)
	}
}

func (n *FileCommand) validFileCreate() error {
	if len(n.FileName) == 0 && len(n.DirName) == 0 {
		return errors.New("one of file-name and dir-name is required")
	}

	return nil
}

func (n *FileCommand) validFileModify() error {
	if len(n.FileName) == 0 {
		return errors.New("file name is required")
	}

	if n.Privilege == 0 {
		return errors.New("file privilege is required")
	}

	return nil
}

func (n *FileCommand) validFileDelete() error {
	if len(n.FileName) == 0 && len(n.DirName) == 0 {
		return errors.New("one of file-name and dir-name is required")
	}

	return nil
}

func (n *FileCommand) validFileRename() error {
	if len(n.SourceFile) == 0 || len(n.DestFile) == 0 {
		return errors.New("both source file and destination file are required")
	}

	return nil
}

func (n *FileCommand) validFileAppend() error {
	if len(n.FileName) == 0 {
		return errors.New("file-name is required")
	}

	if len(n.Data) == 0 {
		return errors.New("append data is required")
	}

	return nil
}

func (n *FileCommand) validFileReplace() error {
	if len(n.FileName) == 0 {
		return errors.New("file-name is required")
	}

	if len(n.OriginStr) == 0 || len(n.DestStr) == 0 {
		return errors.New("both origin and destination string are required")
	}

	return nil
}

func (n *FileCommand) CompleteDefaults() {
	switch n.Action {
	case FileAppendAction:
		n.setDefaultForFileAppend()
	}
}

func (n *FileCommand) setDefaultForFileAppend() {
	if n.Count == 0 {
		n.Count = 1
	}
}

func (n FileCommand) RecoverData() string {
	data, _ := json.Marshal(n)
	return string(data)
}

func NewFileCommand() *FileCommand {
	return &FileCommand{
		CommonAttackConfig: CommonAttackConfig{
			Kind: FileAttack,
		},
	}
}

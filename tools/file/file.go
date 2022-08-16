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

package main

import (
	"fmt"
	"os"

	_ "github.com/alecthomas/template"
	"github.com/spf13/cobra"
	_ "github.com/swaggo/swag"
)

// CommandFlags are flags that used in all Commands
var rootCmd = &cobra.Command{
	Use:   "file tool",
	Short: "A command line client to execute operations related to files",
}

func init() {
	rootCmd.AddCommand(
		NewFileOrDirCreateCommand(),
		NewFileOrDirRenameCommand(),
		NewFileModifyPrivilegeCommand(),
		NewFileAppendCommand(),
		NewFileCopyCommand(),
		NewFileDeleteCommand(),
		NewFileReplaceCommand(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

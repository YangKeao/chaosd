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

package attack

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/chaos-mesh/chaosd/cmd/server"
	"github.com/chaos-mesh/chaosd/pkg/core"
	"github.com/chaos-mesh/chaosd/pkg/server/chaosd"
	"github.com/chaos-mesh/chaosd/pkg/utils"
)

func NewFileAttackCommand(uid *string) *cobra.Command {
	options := core.NewFileCommand()
	dep := fx.Options(
		server.Module,
		fx.Provide(func() *core.FileCommand {
			options.UID = *uid
			return options
		}),
	)

	cmd := &cobra.Command{
		Use:   "file <subcommand>",
		Short: "File attack related commands",
	}

	cmd.AddCommand(
		NewFileCreateCommand(dep, options),
		NewFileModifyPrivilegeCommand(dep, options),
		NewFileDeleteCommand(dep, options),
		NewFileRenameCommand(dep, options),
		NewFileAppendCommand(dep, options),
		NewFileReplaceCommand(dep, options),
	)

	return cmd
}

func NewFileCreateCommand(dep fx.Option, options *core.FileCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create file",

		Run: func(*cobra.Command, []string) {
			options.Action = core.FileCreateAction
			options.CompleteDefaults()
			utils.FxNewAppWithoutLog(dep, fx.Invoke(commonFileAttackFunc)).Run()
		},
	}

	cmd.Flags().StringVarP(&options.FileName, "file-name", "f", "", "the name of file to be created")
	cmd.Flags().StringVarP(&options.DirName, "dir-name", "d", "", "the name of directory to be created")

	return cmd
}

func NewFileModifyPrivilegeCommand(dep fx.Option, options *core.FileCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify",
		Short: "modify file privilege",
		Run: func(cmd *cobra.Command, args []string) {
			options.Action = core.FileModifyPrivilegeAction
			utils.FxNewAppWithoutLog(dep, fx.Invoke(commonFileAttackFunc)).Run()
		},
	}

	cmd.Flags().StringVarP(&options.FileName, "file-name", "f", "", "file to be change privilege")
	cmd.Flags().Uint32VarP(&options.Privilege, "privilege", "p", 0, "privilege to be update")

	return cmd
}

func NewFileDeleteCommand(dep fx.Option, options *core.FileCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete file",

		Run: func(cmd *cobra.Command, args []string) {
			options.Action = core.FileDeleteAction
			options.CompleteDefaults()
			utils.FxNewAppWithoutLog(dep, fx.Invoke(commonFileAttackFunc)).Run()
		},
	}

	cmd.Flags().StringVarP(&options.FileName, "file-name", "f", "", "the file to be deleted")
	cmd.Flags().StringVarP(&options.DirName, "dir-name", "d", "", "the directory to be deleted")

	return cmd
}

func NewFileRenameCommand(dep fx.Option, options *core.FileCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename",
		Short: "rename file",

		Run: func(cmd *cobra.Command, args []string) {
			options.Action = core.FileRenameAction
			utils.FxNewAppWithoutLog(dep, fx.Invoke(commonFileAttackFunc)).Run()
		},
	}

	cmd.Flags().StringVarP(&options.SourceFile, "source-file", "s", "", "the source file/dir of rename")
	cmd.Flags().StringVarP(&options.DestFile, "dest-file", "d", "", "the destination file/dir of rename")

	return cmd
}

func NewFileAppendCommand(dep fx.Option, options *core.FileCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append",
		Short: "append file",

		Run: func(cmd *cobra.Command, args []string) {
			options.Action = core.FileAppendAction
			utils.FxNewAppWithoutLog(dep, fx.Invoke(commonFileAttackFunc)).Run()
		},
	}

	cmd.Flags().StringVarP(&options.FileName, "file-name", "f", "", "append data to the file")
	cmd.Flags().StringVarP(&options.Data, "data", "d", "", "append data")
	cmd.Flags().IntVarP(&options.Count, "count", "c", 1, "append count with default value is 1")

	return cmd
}

func NewFileReplaceCommand(dep fx.Option, options *core.FileCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace",
		Short: "replace data in file",

		Run: func(cmd *cobra.Command, args []string) {
			options.Action = core.FileReplaceAction
			utils.FxNewAppWithoutLog(dep, fx.Invoke(commonFileAttackFunc)).Run()
		},
	}

	cmd.Flags().StringVarP(&options.FileName, "file-name", "f", "", "replace data in the file")
	cmd.Flags().StringVarP(&options.OriginStr, "origin-string", "o", "", "the origin string to be replaced")
	cmd.Flags().StringVarP(&options.DestStr, "dest-string", "d", "", "the destination string to replace the origin string")
	cmd.Flags().IntVarP(&options.Line, "line", "l", 0, "the line number to replace, default is 0, means replace all lines")

	return cmd
}

func commonFileAttackFunc(options *core.FileCommand, chaos *chaosd.Server) {
	if err := options.Validate(); err != nil {
		utils.ExitWithError(utils.ExitBadArgs, err)
	}

	uid, err := chaos.ExecuteAttack(chaosd.FileAttack, options, core.CommandMode)
	if err != nil {
		utils.ExitWithError(utils.ExitError, err)
	}

	utils.NormalExit(fmt.Sprintf("Attack file successfully, uid: %s", uid))
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	vfs "github.com/twpayne/go-vfs"
)

var unmanagedCmd = &cobra.Command{
	Use:     "unmanaged",
	Args:    cobra.NoArgs,
	Short:   "List the unmanaged files in the destination directory",
	Long:    mustGetLongHelp("unmanaged"),
	Example: getExample("unmanaged"),
	PreRunE: config.ensureNoError,
	RunE:    config.runUnmanagedCmd,
}

func init() {
	rootCmd.AddCommand(unmanagedCmd)
}

func (c *Config) runUnmanagedCmd(cmd *cobra.Command, args []string) error {
	ts, err := c.getTargetState(nil)
	if err != nil {
		return err
	}
	return vfs.Walk(c.fs, c.DestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == c.DestDir {
			return nil
		}
		entry, _ := ts.Get(c.fs, path)
		managed := entry != nil
		ignored := ts.TargetIgnore.Match(strings.TrimPrefix(path, c.DestDir+"/"))
		if !managed && !ignored {
			fmt.Println(path)
		}
		if info.IsDir() && (!managed || ignored) {
			return filepath.SkipDir
		}
		return nil
	})
}
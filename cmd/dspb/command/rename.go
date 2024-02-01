package command

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dsp2b/dsp2b-go/pkg/blueprint"
	"github.com/spf13/cobra"
)

func renameCmd(cmd *cobra.Command, args []string) error {
	repo, err := ReadRepository()
	if err != nil {
		return err
	}

	scan := newScan(repo)

	err = scan.Scan(".", func(path string, entry os.DirEntry) error {
		filename := filepath.Join(path, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		newData, err := blueprint.Rename(string(data), filepath.Base(strings.TrimSuffix(filename, ".txt")))
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, []byte(newData), 0644)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

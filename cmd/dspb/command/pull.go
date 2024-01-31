package command

import "github.com/spf13/cobra"

func pull(cmd *cobra.Command, args []string) error {
	// 读取dspb.json
	repo, err := ReadRepository()
	if err != nil {
		return err
	}

	if err := cloneRepo(repo.ID); err != nil {
		return err
	}

	return nil
}

package main

import (
	"github.com/codfrm/cago/pkg/logger"
	"github.com/dsp2b/dsp2b-go/cmd/dspb/command"
	"github.com/dsp2b/dsp2b-go/configs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
)

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger.SetLogger(l)

	rootCmd := &cobra.Command{
		Use:     "dspb",
		Short:   "dspb controls the dspb service.",
		Version: configs.Version,
	}
	command.AddCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute err: %v", err)
	}
}

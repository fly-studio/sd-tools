package main

import (
	"github.com/google/gops/agent"
	"github.com/spf13/cobra"
	"gopkg.in/go-mixed/go-common.v1/utils/io"
	"log"
	"path/filepath"
	"sd-downloader/internal/downloader"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	currentDir := ioUtils.GetCurrentDir()

	rootCmd := &cobra.Command{
		Use:   "sd",
		Short: "Stable diffusion tools",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringSliceP("config", "c", []string{filepath.Join(currentDir, "conf/sd.yml")}, "stable diffusion config files")

	rootCmd.AddCommand(downloader.Cmd())
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

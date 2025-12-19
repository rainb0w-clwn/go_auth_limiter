package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Start cli application",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		var err error
		configR, err := os.Open(configFile)
		if err != nil {
			log.Printf("%s", "Error opening config file: "+err.Error())
			return nil
		}

		cfg, err = config.New(configR)
		if err != nil {
			log.Printf("%s", "Error parsing config file: "+err.Error())
			return nil
		}
		if err != nil {
			return fmt.Errorf("process load config: %w", err)
		}
		return nil
	},
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("Starting cli...")
		if err := runApp(); err != nil {
			log.Fatalf("Fail start: %v", err)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"/configs/config.yml",
		"Path to configuration file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Fail start: %v", err)
	}
}

func runApp() error {
	return nil
}

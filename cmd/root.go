package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/cvhariharan/flowctl/internal/config"
	"github.com/spf13/cobra"
)

var appConfig config.Config

// rootCmd represents the base command when called without any subcommands
// If --new-config flag is added, a new config file (config.toml) will be created in the current directory
var rootCmd = &cobra.Command{
	Use:   "flowctl",
	Short: "Self-service workflow execution engine",
	Run: func(cmd *cobra.Command, args []string) {
		if ok, _ := cmd.Flags().GetBool("new-config"); ok {
			cfg, err := template.ParseFS(StaticFiles, "configs/*.toml")
			if err != nil {
				log.Fatal(err)
			}
			cfgFile, err := os.Create("config.toml")
			if err != nil {
				log.Fatal(err)
			}
			if err := cfg.ExecuteTemplate(cfgFile, "config.sample.toml", config.GetDefaultConfig()); err != nil {
				log.Fatal(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// LoadConfig loads the config file into the global appConfig variable
// If the config path is empty, the current working directory will be used
func LoadConfig(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	appConfig = cfg
	return nil
}

func init() {
	rootCmd.Flags().Bool("new-config", false, "Generate a new default config.toml file")
	rootCmd.PersistentFlags().String("config", "", "Path to config file")
}

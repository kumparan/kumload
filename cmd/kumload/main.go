package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/kumparan/kumload"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kumload",
	Short:   "kumload CLI",
	Long:    "CLI for database migration",
	Version: "1.0.0",
	Run:     run,
}

func run(cmd *cobra.Command, args []string) {
	workDir, err := os.Getwd()
	if err != nil {
		log.Error(err)
		return
	}

	configPath := cmd.Flag("config").Value.String()
	cfg, err := kumload.ParseConfig(path.Join(workDir, configPath))
	if err != nil {
		log.Error(err)
		return
	}

	scriptPath := cmd.Flag("script").Value.String()
	if scriptPath != "" {
		scriptByte, err := ioutil.ReadFile(path.Join(workDir, scriptPath))
		if err != nil {
			log.Error(err)
			return
		}
		cfg.Script = string(scriptByte)
	}

	if scriptPath == "" && cfg.Mappings == nil {
		log.Error("mappings or script must provided")
		return
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Target.Name))
	if err != nil {
		log.Error(err)
		return
	}

	if cfg.LogLevel == kumload.LogLevelDebug {
		db.LogMode(true)
	}

	migration := kumload.NewMigration(db, cfg)
	err = migration.Run()
	if err != nil {
		log.Error(err)
		return
	}
}

func main() {
	rootCmd.Flags().String("config", "config.yml", "config file path")
	rootCmd.Flags().String("script", "", "script file path, this is optional")

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		return
	}
}

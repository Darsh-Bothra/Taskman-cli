// Package config loads and exposes application configuration via Viper.
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	KeyFile = "file"
)

// Init sets up Viper defaults and attempts to read a config file
// (~/.config/taskman/taskman.yaml).  Missing config files are silently
// ignored so the CLI works out-of-the-box.
func Init() {
	viper.SetConfigName("taskman")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "taskman"))

	// Sensible default — tasks live in the user's home directory.
	viper.SetDefault(KeyFile, filepath.Join(os.Getenv("HOME"), ".taskman.json"))

	// Env vars are prefixed TASKMAN_ (e.g. TASKMAN_FILE=/tmp/tasks.json).
	viper.SetEnvPrefix("TASKMAN")
	viper.AutomaticEnv()

	// Ignore "config not found" — it simply means the user hasn't created
	// one yet, which is perfectly fine.
	_ = viper.ReadInConfig()
}

// TaskFile returns the resolved path to the task storage file.
// CLI flag value takes precedence over config / env / default.
func TaskFile(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return viper.GetString(KeyFile)
}
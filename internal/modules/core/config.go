package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	EnvDevelopment = "development"

	runtimeMaxProcsConfigKey = "runtime.max_procs"
)

func setupProjectDir() {
	projectDir = os.Getenv("PROJECT_DIR")
	if projectDir != "" {
		return
	}
	if env == "development" {
		projectDir, _ = filepath.Abs("")
		return
	}
	panic("please set PROJECT_DIR in env")
}

var (
	env        string
	projectDir string
	hostname   string

	ServiceRole string
	ConfigStore *viper.Viper
)

func GetHostname() string {
	return hostname
}

func GetProjectDir() string {
	return projectDir
}

func GetConfigDir() string {
	return path.Join(projectDir, "configs")
}

func init() {
	time.Local = time.UTC
	env = os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	setupProjectDir()
	configDir := GetConfigDir()
	hostname = os.Getenv("HOSTNAME")
	if hostname == "" {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			panic(err)
		}
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName("default")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.SetConfigName(env)
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		} else {
			fmt.Printf("WARNING: No %s config file found, skipping\n", env)
		}
	}

	viper.SetConfigName("override")
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		} else {
			fmt.Println("WARNING: No override config file found, skipping")
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	ConfigStore = viper.GetViper()

	maxProc := ConfigStore.GetInt(runtimeMaxProcsConfigKey)
	if maxProc > 0 {
		runtime.GOMAXPROCS(maxProc)
		fmt.Printf("Set GOMAXPROCS to %d\n", maxProc)
	}
}

func Env() string {
	return env
}

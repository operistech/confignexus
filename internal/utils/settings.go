package utils

import (
	"github.com/spf13/viper"
)

type Settings struct {
	HTTPPort      string
	HTTPSPort     string
	ListenAddress string
	HTTPEnabled   string
	HTTPRedirect  string
	CertPath      string
	KeyPath       string
	HTTPAddr      string // Combined address for HTTP
	HTTPSAddr     string // Combined address for HTTPS
	DebugLog      string
	RepoAddress   string
	RepoBranch    string
}

func LoadSettings() (*Settings, error) {
	viper.SetDefault("HTTPPort", "9000")
	viper.SetDefault("HTTPSPort", "9443")
	viper.SetDefault("ListenAddress", "localhost")
	viper.SetDefault("HTTPEnabled", "true")
	viper.SetDefault("HTTPRedirect", "true")
	viper.SetDefault("CertPath", "./cert.pem")
	viper.SetDefault("KeyPath", "./key.pem")
	viper.SetDefault("DebugLog", "false")
	viper.SetDefault("RepoBranch", "main")

	viper.SetEnvPrefix("CN")
	viper.AutomaticEnv()

	// Look for config file in current working directory
	viper.AddConfigPath(".")
	// Look for config file in /etc
	viper.AddConfigPath("/etc/configNexus")

	// Name of the configuration file without extension
	viper.SetConfigName("config")

	// Try to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	httpaddr := viper.GetString("ListenAddress") + ":" + viper.GetString("HTTPPort")
	httpsaddr := viper.GetString("ListenAddress") + ":" + viper.GetString("HTTPSPort")

	return &Settings{
		HTTPPort:      viper.GetString("HTTPPort"),
		HTTPSPort:     viper.GetString("HTTPSPort"),
		ListenAddress: viper.GetString("ListenAddress"),
		HTTPEnabled:   viper.GetString("HTTPEnabled"),
		HTTPRedirect:  viper.GetString("HTTPRedirect"),
		CertPath:      viper.GetString("CertPath"),
		KeyPath:       viper.GetString("KeyPath"),
		HTTPAddr:      httpaddr,
		HTTPSAddr:     httpsaddr,
		DebugLog:      viper.GetString("DebugLog"),
		RepoAddress:   viper.GetString("RepoAddress"),
		RepoBranch:    viper.GetString("RepoBranch"),
	}, nil
}

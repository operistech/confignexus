/*
   This file is part of configNexus.

   configNexus is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   configNexus is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with configNexus.  If not, see <https://www.gnu.org/licenses/>.

   Copyright (C) 2023 Operistech Inc.
*/

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

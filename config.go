/*
 *
 *    Copyright 2021 Boris Barnier <bozzo@users.noreply.github.com>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type configuration struct {
	Version int              `yaml:"version"`
	Ejp     ejpConfiguration `yaml:"ejp"`
	Knx     knxConfiguration `yaml:"knx"`
}

type ejpConfiguration struct {
	Url       string  `yaml:"url"`
	DateParam string  `yaml:"dateParam"`
	UserAgent string  `yaml:"userAgent"`
	Zone      ejpZone `yaml:"zone"` // EjpNord, EjpOuest, EjpPaca, EjpSud
}

type knxConfiguration struct {
	Ip           string `yaml:"ip"`
	Port         string `yaml:"port"`
	PreavisGroup string `yaml:"preavisGroup"`
	AsservGroup  string `yaml:"asservGroup"`
}

type ejpZone string

const (
	EjpNord  ejpZone = "EjpNord"
	EjpOuest ejpZone = "EjpOuest"
	EjpPaca  ejpZone = "EjpPaca"
	EjpSud   ejpZone = "EjpSud"
)

func (config *configuration) loadConfigurationFromFile() error {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		file = "config.yml"
	}
	ymlFile, err := os.Open(filepath.Clean(file))
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(ymlFile)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}

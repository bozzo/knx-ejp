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
	"github.com/bozzo/knx-go/knx"
	"github.com/bozzo/knx-go/knx/cemi"
	"github.com/bozzo/knx-go/knx/dpt"
	"github.com/bozzo/knx-go/knx/util"
	"github.com/sirupsen/logrus"
	"os"
)

var config = configuration{}

func init() {
	format := os.Getenv("LOG_FORMAT")
	// LOG_LEVEL not set, let's default to debug
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func sendKnxCommand(client knx.GroupRouter, group string, status bool) error {
	knxGroup, err := cemi.NewGroupAddrString(group)
	if err != nil {
		return err
	}
	err = client.Send(knx.GroupEvent{
		Command:     knx.GroupWrite,
		Destination: knxGroup,
		Data:        dpt.DPT_1001(status).Pack(),
	})
	return err
}

func main() {
	// Setup logger for auxiliary logging. This enables us to see log messages from internal
	// routines.
	util.Logger = logrus.StandardLogger()

	if err := config.loadConfigurationFromFile(); err != nil {
		logrus.Fatalf("unable to load configuration file: %s", err)
	}

	// Connect to the gateway.
	client, err := knx.NewGroupRouter(config.Knx.IP+":"+config.Knx.Port, knx.DefaultRouterConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	// Close upon exiting. Even if the gateway closes the connection, we still have to clean up.
	defer client.Close()

	ejp := ejpClient{
		baseURL:         config.Ejp.URL,
		userAgent:       config.Ejp.UserAgent,
		asservBeginHour: 6,
		asservEndHour:   0,
		zone:            config.Ejp.Zone,
	}

	resp, err := ejp.getEjpStatus()
	if err != nil {
		logrus.Fatal(err)
	}

	// Send "preavis" to KNX Bus
	err = sendKnxCommand(client, config.Knx.PreavisGroup, resp.preavis)
	if err != nil {
		logrus.Fatal(err)
	}

	// Send "asservissement" to KNX Bus
	err = sendKnxCommand(client, config.Knx.AsservGroup, resp.asserv)
	if err != nil {
		logrus.Fatal(err)
	}
}

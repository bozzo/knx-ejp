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
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ejpResponse struct {
	JourJ  map[ejpZone]ejpType `json:"JourJ"`
	JourJ1 map[ejpZone]ejpType `json:"JourJ1"`
}

type ejpResults struct {
	preavis bool
	asserv  bool
}

type ejpType string

// EJP Types returned by the API
const (
	EJP ejpType = "EST_EJP"
	//NOT_EJP ejpType = "NON_EJP"
	//ND      ejpType = "ND"
)

type ejpClient struct {
	baseURL         string
	userAgent       string
	asservBeginHour int
	asservEndHour   int
	zone            ejpZone
}

func (ejp *ejpClient) _buildURL() (*url.URL, error) {
	u, err := url.Parse(ejp.baseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("Date_a_remonter", time.Now().UTC().Format("2006-01-02"))
	q.Set("TypeAlerte", "EJP")
	q.Set("_", strconv.FormatInt(time.Now().Unix(), 10))
	u.RawQuery = q.Encode()
	return u, nil
}

func (ejp *ejpClient) getEjpStatus() (*ejpResults, error) {
	u, err := ejp._buildURL()
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-agent", ejp.userAgent)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result ejpResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logrus.Fatal(err)
	}

	return ejp._parseResults(result), nil
}

func (ejp *ejpClient) _parseResults(result ejpResponse) *ejpResults {
	now := time.Now().UTC()

	var results ejpResults

	if now.Before(time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.UTC)) {
		results = ejpResults{
			asserv:  false,
			preavis: result.JourJ[ejp.zone] == EJP,
		}
	} else {
		results = ejpResults{
			asserv:  result.JourJ[ejp.zone] == EJP,
			preavis: result.JourJ1[ejp.zone] == EJP,
		}
	}
	return &results
}

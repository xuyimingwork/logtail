/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/vogo/logtail"
	"github.com/vogo/logtail/repeater"
)

func main() {
	config := &logtail.Config{
		DefaultFormat: &logtail.Format{Prefix: "!!!!-!!-!!"},
		Transfers: map[string]*logtail.TransferConfig{
			"console": {
				Name: "console",
				Type: "console",
			},
			"ding": {
				Name: "ding",
				Type: "ding",
				URL:  "http://localhost:55321",
			},
		},
		Routers: map[string]*logtail.RouterConfig{
			"error-console": {
				Name: "error-console",
				Matchers: []*logtail.MatcherConfig{
					{
						Contains: []string{"ERROR"},
					},
				},
				Transfers: []string{"console", "ding"},
			},
		},
		DefaultRouters: []string{"error-console"},
		Servers: map[string]*logtail.ServerConfig{
			"server-test": {
				Name: "server-test",
			},
		},
	}

	runner, err := logtail.NewRunner(config)
	if err != nil {
		panic(err)
	}

	if err = runner.Start(); err != nil {
		panic(err)
	}

	serverID := "server-test"
	server := runner.Servers[config.Servers[serverID].Name]

	c := make(chan []byte)
	go repeater.Repeat("/Users/gelnyang/temp/logtail/test.log", c)

	for {
		b := <-c
		_ = server.Fire(b)
	}
}

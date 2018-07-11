// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package add_path_metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

func TestExportPath(t *testing.T) {
	testConfig, err := common.NewConfigFrom(map[string]interface{}{
		"format": "/data/log1/$app/$pod.name/$containerid/$ignore",
	})
	if err != nil {
		t.Fatal(err)
	}

	input := common.MapStr{}
	input["source"] = "/data/log1/riven/test-pod/container4324234/access.log"
	actual := getActualValue(t, testConfig, input)

	expected := common.MapStr{
		"app":         "riven",
		"containerid": "container4324234",
		"source":      "/data/log1/riven/test-pod/container4324234/access.log",
		"pod": map[string]string{
			"name": "test-pod",
		},
	}

	assert.Equal(t, expected.String(), actual.String())
}

func getActualValue(t *testing.T, config *common.Config, input common.MapStr) common.MapStr {
	logp.TestingSetup()

	p, err := newAddPathMetadata(config)
	if err != nil {
		logp.Err("Error initializing add_locale")
		t.Fatal(err)
	}

	actual, err := p.Run(&beat.Event{Fields: input})
	return actual.Fields
}

func TestParse(t *testing.T) {
	format := "/data/log1/$app/$pod.name/$containerid/$ignore"
	source := "/data/log1/riven/pod-123132/container4324234/access.log"
	re, err := parse(format, source)
	if err != nil {
		t.Fatal(err)
	}
	expected := make(map[string]string)
	expected["app"] = "riven"
	expected["pod.name"] = "pod-123132"
	expected["containerid"] = "container4324234"
	assert.Equal(t, expected, re, "should be equal")
}

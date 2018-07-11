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
	"fmt"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

type addPathMetadata struct {
	format string
}

// TimezoneFormat type
type TimezoneFormat int

// Timezone formats
const ()

func init() {
	processors.RegisterPlugin("add_path_metadata", newAddPathMetadata)
}

func newAddPathMetadata(c *common.Config) (processors.Processor, error) {
	config := struct {
		Format string `config:"format"`
	}{
		Format: "offset",
	}

	err := c.Unpack(&config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unpack the add_path_metadata configuration")
	}

	pm := addPathMetadata{
		format: config.Format,
	}

	return pm, nil
}

func (l addPathMetadata) Run(event *beat.Event) (*beat.Event, error) {
	if value, ok := event.Fields["source"]; ok {
		source := value.(string)
		re, err := parse(l.format, source)
		if err != nil {
			return nil, fmt.Errorf("failed to parse source when add_path_metadata, format:%s, souce:%s", l.format, source)
		}
		for k, v := range re {
			event.PutValue(k, v)
		}
	}
	return event, nil
}

func (l addPathMetadata) String() string {
	return "add_path_metadata=[format:" + l.format + "]"
}

//parse the source with format args to a map result
//format: /data/log1/$app/$podname/$containerid/$ignore
//source: /data/log1/riven/pod-123132/container4324234/access.log
func parse(format, source string) (map[string]string, error) {
	i, j := 0, 0
	re := make(map[string]string)
	reading := false
	startI, startJ := 0, 0
	for i < len(source) && j < len(format) {
		if !reading {
			if source[i] == format[j] {
				i++
				j++
				continue
			}
			if format[j] != '$' {
				return nil, errors.New("format not match source")
			}
			if format[j] == '$' {
				j++
				startI, startJ = i, j
				reading = true
			}
		}

		for i+1 < len(source) && source[i+1] != '/' {
			i++
		}
		for j+1 < len(format) && format[j+1] != '/' {
			j++
		}

		//stop read
		k := format[startJ : j+1]
		v := source[startI : i+1]
		if k != "ignore" {
			re[k] = v
		}
		startI, startJ = 0, 0
		reading = false
		i++
		j++
	}
	return re, nil
}

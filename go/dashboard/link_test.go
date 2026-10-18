// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dashboard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func dashboardWithLinks(links string) string {
	return `{
	  "duration": "1h",
	  "panels": {},
	  "layouts": [],
	  "links": ` + links + `
	}`
}

func TestUnmarshalURLLink(t *testing.T) {
	var spec Spec
	err := json.Unmarshal([]byte(dashboardWithLinks(`[{"name":"Docs","url":"https://perses.dev","renderVariables":true}]`)), &spec)
	require.NoError(t, err)
	require.Len(t, spec.Links, 1)
	assert.Equal(t, "Docs", spec.Links[0].Name)
	assert.Equal(t, "https://perses.dev", spec.Links[0].URL)
	assert.True(t, spec.Links[0].RenderVariables)
	assert.Empty(t, spec.Links[0].Tags)
}

func TestUnmarshalTagLinkTrimsTags(t *testing.T) {
	var spec Spec
	raw := dashboardWithLinks(`[{"name":" Service detail ","tags":[" traffic ","billing"],"keepTime":true,"includeVariables":true}]`)
	require.NoError(t, json.Unmarshal([]byte(raw), &spec))
	require.Len(t, spec.Links, 1)
	assert.Equal(t, "Service detail", spec.Links[0].Name)
	assert.Equal(t, []string{"traffic", "billing"}, spec.Links[0].Tags)
	assert.True(t, spec.Links[0].KeepTime)
	assert.True(t, spec.Links[0].IncludeVariables)
	assert.Empty(t, spec.Links[0].URL)

	data, err := json.Marshal(spec.Links[0])
	require.NoError(t, err)
	assert.NotContains(t, string(data), `"url"`)
}

func TestUnmarshalTagLinkYAML(t *testing.T) {
	const raw = `
duration: 1h
panels: {}
layouts: []
links:
  - name: Service detail
    tags: [traffic]
    keepTime: true
`
	var spec Spec
	require.NoError(t, yaml.Unmarshal([]byte(raw), &spec))
	require.Len(t, spec.Links, 1)
	assert.Equal(t, []string{"traffic"}, spec.Links[0].Tags)
	assert.True(t, spec.Links[0].KeepTime)
}

func TestUnmarshalRejectsLinkWithNeitherTarget(t *testing.T) {
	var spec Spec
	err := json.Unmarshal([]byte(dashboardWithLinks(`[{"name":"Docs"}]`)), &spec)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "links[0]: link must set a url or at least one tag")
}

func TestUnmarshalRejectsLinkWithBothTargets(t *testing.T) {
	var spec Spec
	err := json.Unmarshal([]byte(dashboardWithLinks(`[{"name":"Docs","url":"https://perses.dev","tags":["traffic"]}]`)), &spec)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "links[0]: link must set either a url or tags, not both")
}

func TestUnmarshalRejectsEmptyTag(t *testing.T) {
	var spec Spec
	err := json.Unmarshal([]byte(dashboardWithLinks(`[{"name":"Docs","tags":["traffic"," "]}]`)), &spec)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "links[0]: link tag must not be empty")
}

func TestUnmarshalRejectsTagLinkWithoutName(t *testing.T) {
	var spec Spec
	err := json.Unmarshal([]byte(dashboardWithLinks(`[{"tags":["traffic"]}]`)), &spec)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "links[0]: a tag link must set a name")
}

func TestUnmarshalRejectsPanelTagLinkWithoutName(t *testing.T) {
	const raw = `{
	  "duration": "1h",
	  "layouts": [],
	  "panels": {
	    "cpu": {
	      "kind": "Panel",
	      "spec": {
	        "plugin": {"kind": "TimeSeriesChart", "spec": {}},
	        "links": [{"tags": ["traffic"]}]
	      }
	    }
	  }
	}`
	var spec Spec
	err := json.Unmarshal([]byte(raw), &spec)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `panels[cpu].spec.links[0]: a tag link must set a name`)
}

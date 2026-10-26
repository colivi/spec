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
	"errors"
	"strings"
)

// Link is a URL or a same-project tag menu.
// An empty URL with no tags is still valid: dashboards stored before tag links used that placeholder.
// When Tags is set, URL must be empty, every tag must be non-empty, and Name is required.
// KeepTime and IncludeVariables apply to both. RenderVariables only interpolates $var inside a URL.
type Link struct {
	Name             string   `json:"name,omitempty" yaml:"name,omitempty"`
	URL              string   `json:"url,omitempty" yaml:"url,omitempty"`
	Tooltip          string   `json:"tooltip,omitempty" yaml:"tooltip,omitempty"`
	RenderVariables  bool     `json:"renderVariables,omitempty" yaml:"renderVariables,omitempty"`
	TargetBlank      bool     `json:"targetBlank,omitempty" yaml:"targetBlank,omitempty"`
	Tags             []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	KeepTime         bool     `json:"keepTime,omitempty" yaml:"keepTime,omitempty"`
	IncludeVariables bool     `json:"includeVariables,omitempty" yaml:"includeVariables,omitempty"`
}

func (l *Link) UnmarshalJSON(data []byte) error {
	var tmp Link
	type plain Link
	if err := json.Unmarshal(data, (*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).validate(); err != nil {
		return err
	}
	*l = tmp
	return nil
}

func (l *Link) UnmarshalYAML(unmarshal func(any) error) error {
	var tmp Link
	type plain Link
	if err := unmarshal((*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).validate(); err != nil {
		return err
	}
	*l = tmp
	return nil
}

// MarshalJSON keeps "url" on a URL link even when it is empty.
// omitempty would drop it and the document would no longer say which kind of link it is.
func (l Link) MarshalJSON() ([]byte, error) {
	type plain Link
	if l.URL != "" || len(l.Tags) > 0 {
		return json.Marshal(plain(l))
	}
	payload := struct {
		plain
		URL string `json:"url"`
	}{plain: plain(l), URL: ""}
	return json.Marshal(payload)
}

func (l *Link) validate() error {
	l.URL = strings.TrimSpace(l.URL)
	l.Name = strings.TrimSpace(l.Name)

	var tags []string
	for _, tag := range l.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return errors.New("link tag must not be empty")
		}
		tags = append(tags, tag)
	}
	l.Tags = tags

	if l.URL != "" && len(l.Tags) > 0 {
		return errors.New("link must set either a url or tags, not both")
	}
	if len(l.Tags) > 0 && l.Name == "" {
		return errors.New("a tag link must set a name")
	}
	return nil
}

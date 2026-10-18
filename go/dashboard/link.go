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
	"fmt"
	"strings"
)

func (d *Spec) validateLinks() error {
	for i := range d.Links {
		if err := normalizeLink(&d.Links[i], fmt.Sprintf("links[%d]", i)); err != nil {
			return err
		}
	}
	for key, panel := range d.Panels {
		if panel == nil {
			continue
		}
		for i := range panel.Spec.Links {
			if err := normalizeLink(&panel.Spec.Links[i], fmt.Sprintf("panels[%s].spec.links[%d]", key, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// normalizeLink trims the URL and tags, then requires exactly one of a URL or
// at least one tag. A tag link also needs a name, which is the menu title.
// Validation lives here rather than in UnmarshalJSON so cue get go keeps a
// structured #Link definition.
func normalizeLink(link *Link, where string) error {
	link.URL = strings.TrimSpace(link.URL)
	link.Name = strings.TrimSpace(link.Name)

	trimmed := make([]string, 0, len(link.Tags))
	for _, tag := range link.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return fmt.Errorf("%s: link tag must not be empty", where)
		}
		trimmed = append(trimmed, tag)
	}
	if len(trimmed) == 0 {
		link.Tags = nil
	} else {
		link.Tags = trimmed
	}

	hasURL := link.URL != ""
	hasTags := len(link.Tags) > 0
	if hasURL == hasTags {
		if hasURL {
			return fmt.Errorf("%s: link must set either a url or tags, not both", where)
		}
		return fmt.Errorf("%s: link must set a url or at least one tag", where)
	}
	if hasTags && link.Name == "" {
		return fmt.Errorf("%s: a tag link must set a name", where)
	}
	return nil
}

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

export interface Link {
  name?: string;
  /** Set for a URL link. Mutually exclusive with tags. */
  url?: string;
  tooltip?: string;
  /** Interpolate $var inside url. Independent of includeVariables. */
  renderVariables?: boolean;
  targetBlank?: boolean;
  /**
   * When set, the link is a menu of dashboards in the current project that
   * carry every tag. Mutually exclusive with url. A tag link requires name.
   */
  tags?: string[];
  /** Forward the current time range when opening a dashboard from a tag menu. */
  keepTime?: boolean;
  /** Forward variables whose names exist on the destination dashboard. */
  includeVariables?: boolean;
}

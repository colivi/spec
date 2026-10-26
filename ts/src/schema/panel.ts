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

import { z } from 'zod';

import type { Link, PanelDefinition, PanelDisplay, PanelSpec, QueryDefinition } from '../dashboard';
import { annotationSpecSchema } from './annotation';
import { pluginSchema } from './plugin';
import type { PluginSchema } from './plugin';

export const panelDisplaySpec: z.ZodType<PanelDisplay, PanelDisplay> = z.object({
  name: z.string().optional(),
  description: z.string().optional(),
});

export const querySpecSchema: z.ZodType<QueryDefinition, QueryDefinition> = z.object({
  kind: z.string().min(1),
  spec: z.object({
    name: z.string().optional(),
    plugin: pluginSchema,
  }),
});

const linkInputSchema = z.object({
  name: z.string().optional(),
  url: z.string().optional(),
  tooltip: z.string().optional(),
  renderVariables: z.boolean().optional(),
  targetBlank: z.boolean().optional(),
  tags: z.array(z.string()).optional(),
  keepTime: z.boolean().optional(),
  includeVariables: z.boolean().optional(),
});

export const linkSchema: z.ZodType<Link, z.input<typeof linkInputSchema>> = linkInputSchema
  .superRefine((link, ctx) => {
    const url = link.url?.trim() ?? '';
    const name = link.name?.trim() ?? '';
    const tags = (link.tags ?? []).map((tag) => tag.trim()).filter((tag) => tag.length > 0);
    const rawTagCount = link.tags?.length ?? 0;
    if (rawTagCount > 0 && tags.length !== rawTagCount) {
      ctx.addIssue({ code: 'custom', message: 'link tag must not be empty' });
    }
    if (url.length > 0 && tags.length > 0) {
      ctx.addIssue({ code: 'custom', message: 'link must set either a url or tags, not both' });
    }
    if (tags.length > 0 && name.length === 0) {
      ctx.addIssue({ code: 'custom', message: 'a tag link must set a name' });
    }
  })
  .transform((link): Link => {
    const url = link.url?.trim() ?? '';
    const name = link.name?.trim() || undefined;
    const tags = (link.tags ?? []).map((tag) => tag.trim()).filter((tag) => tag.length > 0);
    const shared = {
      tooltip: link.tooltip,
      keepTime: link.keepTime,
      includeVariables: link.includeVariables,
    };
    if (tags.length > 0) {
      return { ...shared, name: name ?? '', tags };
    }
    return {
      ...shared,
      name,
      url,
      renderVariables: link.renderVariables,
      targetBlank: link.targetBlank,
    };
  });

// Input stays the loose JSON object. Output is the URLLink | TagLink union.
export const panelSpecSchema: z.ZodType<PanelSpec> = z.object({
  display: panelDisplaySpec.optional(),
  plugin: pluginSchema,
  queries: z.array(querySpecSchema).optional(),
  links: z.array(linkSchema).optional(),
  annotations: z.array(annotationSpecSchema).optional(),
});

export function buildPanelSpecSchema(customPluginSchema: PluginSchema): z.ZodType<PanelSpec> {
  return z.object({
    display: panelDisplaySpec.optional(),
    plugin: customPluginSchema,
    queries: z.array(querySpecSchema).optional(),
    links: z.array(linkSchema).optional(),
    annotations: z.array(annotationSpecSchema).optional(),
  });
}

export const panelDefinitionSchema: z.ZodType<PanelDefinition> = z.object({
  kind: z.literal('Panel'),
  spec: panelSpecSchema,
});

export function buildPanelDefinitionSchema(customPluginSchema: PluginSchema): z.ZodType<PanelDefinition> {
  return z.object({
    kind: z.literal('Panel'),
    spec: buildPanelSpecSchema(customPluginSchema),
  });
}

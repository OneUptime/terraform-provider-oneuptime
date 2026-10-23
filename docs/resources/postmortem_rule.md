---
page_title: "oneuptime_postmortem_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Choose which resolved incidents OneUptime AI drafts a postmortem for, and the templates a draft may follow. With no rule, every resolved incident gets a draft; with rules, the first rule that matches decides.
---

# oneuptime_postmortem_rule (Resource)

Choose which resolved incidents OneUptime AI drafts a postmortem for, and the templates a draft may follow. With no rule, every resolved incident gets a draft; with rules, the first rule that matches decides.

## Example Usage

```terraform
resource "oneuptime_postmortem_rule" "example" {
  name        = "Example postmortem rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this postmortem rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this postmortem rule.
- `description_pattern` (String) Case-insensitive regex matched against the incident's description. Leave empty to match any description.
- `incident_severities` (Set of String) Match only incidents with these severities. Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `is_enabled` (Boolean) Whether this rule is enabled. A disabled rule is skipped, as if it were not there. Defaults to `true`.
- `labels` (Set of String) Match only incidents that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `llm_provider_id` (String) ID of the LLM provider OneUptime AI uses for everything this rule drafts. Empty: the project's default provider, else the global one. The ID of a `oneuptime_llm_provider`.
- `monitor_labels` (Set of String) Match only when one of the incident's monitors carries at least one of these labels - the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Match only incidents from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `postmortem_templates` (Set of String) The postmortem templates a draft for the incidents this rule matches may follow. Leave empty for OneUptime AI's standard outline. With more than one, OneUptime AI picks the one that fits the incident best. At most 10. IDs of `oneuptime_incident_postmortem_template` resources.
- `title_pattern` (String) Case-insensitive regex matched against the incident's title. Leave empty to match any title.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing postmortem rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_postmortem_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_postmortem_rule.example <id>
```

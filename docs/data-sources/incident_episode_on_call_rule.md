---
page_title: "oneuptime_incident_episode_on_call_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically executing on-call duty policies when matching incident episodes are created
---

# oneuptime_incident_episode_on_call_rule (Data Source)

Configure rules for automatically executing on-call duty policies when matching incident episodes are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode on call rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode_on_call_rule" "example" {
  name = "Example incident episode on call rule"
}

# Or by id:
data "oneuptime_incident_episode_on_call_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this incident episode on-call rule.
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description.
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this incident episode on-call rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `episode_labels` (Set of String) Only trigger for episodes that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only trigger for episodes with these severities. Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `on_call_duty_policies` (Set of String) On-call duty policies to execute when an incident episode matches this rule. IDs of `oneuptime_on_call_policy` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

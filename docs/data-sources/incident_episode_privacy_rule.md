---
page_title: "oneuptime_incident_episode_privacy_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically marking matching incident episodes as private
---

# oneuptime_incident_episode_privacy_rule (Data Source)

Configure rules for automatically marking matching incident episodes as private

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode privacy rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode_privacy_rule" "example" {
  name = "Example incident episode privacy rule"
}

# Or by id:
data "oneuptime_incident_episode_privacy_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this incident episode privacy rule.
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description. Leave empty to match any description.
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title. Leave empty to match any title.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this incident episode privacy rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `episode_labels` (Set of String) Only trigger for episodes that have at least one of these labels. Leave empty to match regardless of episode labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only trigger for episodes with these severities. Leave empty to match episodes of any severity. IDs of `oneuptime_incident_severity` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

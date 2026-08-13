---
page_title: "oneuptime_incident_episode_privacy_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically marking matching incident episodes as private
---

# oneuptime_incident_episode_privacy_rule (Data Source)

Configure rules for automatically marking matching incident episodes as private Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_episode_privacy_rule" "by_name" {
  name = "example-incident_episode_privacy_rule"
}

data "oneuptime_incident_episode_privacy_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this incident episode privacy rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `incident_severities` (Set) Only trigger for episodes with these severities. Leave empty to match episodes of any severity... Computed.
- `episode_labels` (Set) Only trigger for episodes that have at least one of these labels. Leave empty to match regardless of episode labels... Computed.
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title. Leave empty to match any title... Computed.
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description. Leave empty to match any description... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

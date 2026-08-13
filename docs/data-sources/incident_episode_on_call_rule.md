---
page_title: "oneuptime_incident_episode_on_call_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically executing on-call duty policies when matching incident episodes are created
---

# oneuptime_incident_episode_on_call_rule (Data Source)

Configure rules for automatically executing on-call duty policies when matching incident episodes are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_episode_on_call_rule" "by_name" {
  name = "example-incident_episode_on_call_rule"
}

data "oneuptime_incident_episode_on_call_rule" "by_id" {
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
- `description` (String) Description of this incident episode on-call rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `incident_severities` (Set) Only trigger for episodes with these severities. Leave empty to match any severity... Computed.
- `episode_labels` (Set) Only trigger for episodes that have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title... Computed.
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description... Computed.
- `on_call_duty_policies` (Set) On-call duty policies to execute when an incident episode matches this rule... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

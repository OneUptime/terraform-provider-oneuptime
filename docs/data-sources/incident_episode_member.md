---
page_title: "oneuptime_incident_episode_member Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Link between incidents and episodes
---

# oneuptime_incident_episode_member (Data Source)

Link between incidents and episodes Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_episode_member" "by_name" {
  name = "example-incident_episode_member"
}

data "oneuptime_incident_episode_member" "by_id" {
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
- `incident_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `added_at` (String) A date time object.. Computed.
- `added_by` (String) How this incident was added to the episode (rule, manual, or api).. Computed.
- `added_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `matched_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified_of_incident_added` (Bool) Has the owner been notified that this incident was added to the episode?.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

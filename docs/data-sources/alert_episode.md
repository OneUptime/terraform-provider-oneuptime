---
page_title: "oneuptime_alert_episode Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert episodes (groups of related alerts) for your project
---

# oneuptime_alert_episode (Data Source)

Manage alert episodes (groups of related alerts) for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_episode" "by_name" {
  name = "example-alert_episode"
}

data "oneuptime_alert_episode" "by_id" {
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
- `title` (String) Title of this alert episode.. Computed.
- `description` (String) Description of this alert episode. This is in markdown format... Computed.
- `episode_number` (Number) Auto-incrementing episode number per project.. Computed.
- `episode_number_with_prefix` (String) Episode number with prefix (e.g., 'AE-42' or '#42').. Computed.
- `current_alert_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `root_cause` (String) User-documented root cause of this episode.. Computed.
- `last_alert_added_at` (String) A date time object.. Computed.
- `resolved_at` (String) A date time object.. Computed.
- `all_alerts_resolved_at` (String) A date time object.. Computed.
- `assigned_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `assigned_to_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_grouping_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policies` (Set) List of on-call duty policies to execute for this episode... Computed.
- `is_on_call_policy_executed` (Bool) Whether the on-call policy has been executed for this episode.. Computed.
- `alert_count` (Number) Denormalized count of alerts in this episode.. Computed.
- `title_template` (String) Template used to generate the episode title. Stored for dynamic variable updates... Computed.
- `description_template` (String) Template used to generate the episode description. Stored for dynamic variable updates... Computed.
- `is_manually_created` (Bool) Whether this episode was manually created vs auto-created by a rule.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified_of_episode_creation` (Bool) Are owners notified when this episode is created?.. Computed.
- `grouping_key` (String) Key used for grouping alerts into this episode. Generated from groupByFields of the matching rule... Computed.
- `remediation_notes` (String) User-documented remediation steps and notes for this episode.. Computed.
- `post_updates_to_workspace_channels` (String) Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams).. Computed.
- `is_private` (Bool) If true, this alert episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners... Computed.

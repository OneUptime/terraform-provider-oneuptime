---
page_title: "oneuptime_alert_episode Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert episodes (groups of related alerts) for your project
---

# oneuptime_alert_episode (Resource)

Manage alert episodes (groups of related alerts) for your project

## Example Usage

```terraform
resource "oneuptime_alert_episode" "example" {
  title = "This is an example of longer text content that might be stored in this field."
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `title` (String) Title of this alert episode..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this alert episode. This is in markdown format...
- `episode_number` (Number) Auto-incrementing episode number per project..
- `current_alert_state_id` (String) A unique identifier for an object, represented as a UUID..
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `root_cause` (String) User-documented root cause of this episode..
- `last_alert_added_at` (String) A date time object..
- `resolved_at` (String) A date time object..
- `assigned_to_user_id` (String) A unique identifier for an object, represented as a UUID..
- `assigned_to_team_id` (String) A unique identifier for an object, represented as a UUID..
- `alert_grouping_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policies` (Set) List of on-call duty policies to execute for this episode...
- `title_template` (String) Template used to generate the episode title. Stored for dynamic variable updates...
- `description_template` (String) Template used to generate the episode description. Stored for dynamic variable updates...
- `is_manually_created` (Bool) Whether this episode was manually created vs auto-created by a rule..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `grouping_key` (String) Key used for grouping alerts into this episode. Generated from groupByFields of the matching rule...
- `remediation_notes` (String) User-documented remediation steps and notes for this episode..
- `post_updates_to_workspace_channels` (String) Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams)..
- `is_private` (Bool) If true, this alert episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `episode_number_with_prefix` (String) Episode number with prefix (e.g., 'AE-42' or '#42')..
- `all_alerts_resolved_at` (String) A date time object..
- `is_on_call_policy_executed` (Bool) Whether the on-call policy has been executed for this episode..
- `alert_count` (Number) Denormalized count of alerts in this episode..
- `is_owner_notified_of_episode_creation` (Bool) Are owners notified when this episode is created?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_episode.example <id>
```

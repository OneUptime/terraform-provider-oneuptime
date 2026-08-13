---
page_title: "oneuptime_incident_episode Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident episodes (groups of related incidents) for your project
---

# oneuptime_incident_episode (Resource)

Manage incident episodes (groups of related incidents) for your project

## Example Usage

```terraform
resource "oneuptime_incident_episode" "example" {
  title = "This is an example of longer text content that might be stored in this field."
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `title` (String) Title of this incident episode..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this incident episode. This is in markdown format...
- `episode_number` (Number) Auto-incrementing episode number per project..
- `current_incident_state_id` (String) A unique identifier for an object, represented as a UUID..
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `root_cause` (String) User-documented root cause of this episode..
- `last_incident_added_at` (String) A date time object..
- `resolved_at` (String) A date time object..
- `assigned_to_user_id` (String) A unique identifier for an object, represented as a UUID..
- `assigned_to_team_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policies` (Set) List of on-call duty policies to execute for this episode...
- `title_template` (String) Template used to generate the episode title. Stored for dynamic variable updates...
- `description_template` (String) Template used to generate the episode description. Stored for dynamic variable updates...
- `is_manually_created` (Bool) Whether this episode was manually created vs auto-created by a rule..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `grouping_key` (String) Key used for grouping incidents into this episode. Generated from groupByFields of the matching rule...
- `incident_grouping_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `remediation_notes` (String) User-documented remediation steps and notes for this episode..
- `postmortem_note` (String) User-documented postmortem summary for this episode..
- `post_updates_to_workspace_channels` (String) Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams)..
- `is_visible_on_status_page` (Bool) Should this episode be visible on the status page?..
- `declared_at` (String) A date time object..
- `should_status_page_subscribers_be_notified_on_episode_created` (Bool) Should status page subscribers be notified when this episode is created?..
- `is_private` (Bool) If true, this incident episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `episode_number_with_prefix` (String) Episode number with prefix (e.g., 'IE-42' or '#42')..
- `all_incidents_resolved_at` (String) A date time object..
- `is_on_call_policy_executed` (Bool) Whether the on-call policy has been executed for this episode..
- `incident_count` (Number) Denormalized count of incidents in this episode..
- `is_owner_notified_of_episode_creation` (Bool) Are owners notified when this episode is created?..
- `subscriber_notification_status_on_episode_created` (String) Status of notification sent to subscribers when this episode was created..
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_episode.example <id>
```

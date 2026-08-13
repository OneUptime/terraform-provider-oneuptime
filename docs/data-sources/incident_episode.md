---
page_title: "oneuptime_incident_episode Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident episodes (groups of related incidents) for your project
---

# oneuptime_incident_episode (Data Source)

Manage incident episodes (groups of related incidents) for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_episode" "by_name" {
  name = "example-incident_episode"
}

data "oneuptime_incident_episode" "by_id" {
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
- `title` (String) Title of this incident episode.. Computed.
- `description` (String) Description of this incident episode. This is in markdown format... Computed.
- `episode_number` (Number) Auto-incrementing episode number per project.. Computed.
- `episode_number_with_prefix` (String) Episode number with prefix (e.g., 'IE-42' or '#42').. Computed.
- `current_incident_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `root_cause` (String) User-documented root cause of this episode.. Computed.
- `last_incident_added_at` (String) A date time object.. Computed.
- `resolved_at` (String) A date time object.. Computed.
- `all_incidents_resolved_at` (String) A date time object.. Computed.
- `assigned_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `assigned_to_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policies` (Set) List of on-call duty policies to execute for this episode... Computed.
- `is_on_call_policy_executed` (Bool) Whether the on-call policy has been executed for this episode.. Computed.
- `incident_count` (Number) Denormalized count of incidents in this episode.. Computed.
- `title_template` (String) Template used to generate the episode title. Stored for dynamic variable updates... Computed.
- `description_template` (String) Template used to generate the episode description. Stored for dynamic variable updates... Computed.
- `is_manually_created` (Bool) Whether this episode was manually created vs auto-created by a rule.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified_of_episode_creation` (Bool) Are owners notified when this episode is created?.. Computed.
- `grouping_key` (String) Key used for grouping incidents into this episode. Generated from groupByFields of the matching rule... Computed.
- `incident_grouping_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `remediation_notes` (String) User-documented remediation steps and notes for this episode.. Computed.
- `postmortem_note` (String) User-documented postmortem summary for this episode.. Computed.
- `post_updates_to_workspace_channels` (String) Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams).. Computed.
- `is_visible_on_status_page` (Bool) Should this episode be visible on the status page?.. Computed.
- `declared_at` (String) A date time object.. Computed.
- `should_status_page_subscribers_be_notified_on_episode_created` (Bool) Should status page subscribers be notified when this episode is created?.. Computed.
- `subscriber_notification_status_on_episode_created` (String) Status of notification sent to subscribers when this episode was created.. Computed.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.. Computed.
- `is_private` (Bool) If true, this incident episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners... Computed.

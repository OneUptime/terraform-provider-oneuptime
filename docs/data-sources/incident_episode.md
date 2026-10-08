---
page_title: "oneuptime_incident_episode Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident episodes (groups of related incidents) for your project
---

# oneuptime_incident_episode (Data Source)

Manage incident episodes (groups of related incidents) for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode" "example" {
  title = "example-title"
}

# Or by id:
data "oneuptime_incident_episode" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `assigned_to_team_id` (String) Team ID that is assigned to this episode. The ID of a `oneuptime_team`.
- `assigned_to_user_id` (String) User ID who is assigned to this episode. The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_incident_state_id` (String) Current Incident State ID. The ID of a `oneuptime_incident_state`.
- `description` (String) Description of this incident episode. This is in markdown format.
- `description_template` (String) Template used to generate the episode description. Stored for dynamic variable updates.
- `episode_number` (Number) Auto-incrementing episode number per project.
- `episode_number_with_prefix` (String) Episode number with prefix (e.g., 'IE-42' or '#42').
- `grouping_key` (String) Key used for grouping incidents into this episode. Generated from groupByFields of the matching rule. When a private incident opened the episode, its title is in the key only as a keyed hash.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_count` (Number) Denormalized count of incidents in this episode.
- `incident_grouping_rule_id` (String) ID of the Incident Grouping Rule that created this episode (if applicable). The ID of a `oneuptime_incident_grouping_rule`.
- `incident_severity_id` (String) Incident Severity ID. The ID of a `oneuptime_incident_severity`.
- `is_manually_created` (Boolean) Whether this episode was manually created vs auto-created by a rule.
- `is_on_call_policy_executed` (Boolean) Whether the on-call policy has been executed for this episode.
- `is_owner_notified_of_episode_creation` (Boolean) Are owners notified when this episode is created?
- `is_private` (Boolean) If true, this incident episode is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners.
- `is_visible_on_status_page` (Boolean) Should this episode be visible on the status page?
- `postmortem_note` (String) User-documented postmortem summary for this episode.
- `remediation_notes` (String) User-documented remediation steps and notes for this episode.
- `root_cause` (String) User-documented root cause of this episode.
- `should_status_page_subscribers_be_notified_on_episode_created` (Boolean) Should status page subscribers be notified when this episode is created?
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_episode_created` (String) Status of notification sent to subscribers when this episode was created.
- `title` (String) Title of this incident episode.
- `title_template` (String) Template used to generate the episode title. Stored for dynamic variable updates.

### Read-Only

- `all_incidents_resolved_at` (String) When all incidents in this episode were first detected as resolved. Used for resolve delay calculation.
- `created_at` (String) Date and Time when the object was created.
- `declared_at` (String) When this episode was declared.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_incident_added_at` (String) When the last incident was added to this episode.
- `on_call_duty_policies` (Set of String) List of on-call duty policies to execute for this episode. IDs of `oneuptime_on_call_policy` resources.
- `post_updates_to_workspace_channels` (String) Workspace channels to post episode updates to (e.g., Slack, Microsoft Teams). A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resolved_at` (String) When this episode was resolved.
- `updated_at` (String) Date and Time when the object was updated.

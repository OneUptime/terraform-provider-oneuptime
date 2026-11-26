---
page_title: "oneuptime_team_compliance_setting Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Compliance settings for your OneUptime team
---

# oneuptime_team_compliance_setting (Data Source)

Compliance settings for your OneUptime team

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team compliance setting may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_compliance_setting" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_team_compliance_setting" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `enabled` (Boolean) Whether this compliance rule is enabled.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `notification_channel` (String) Deprecated: use notificationChannels. The first of the rule's notification channels, or empty when it accepts any channel. Sending this field without notificationChannels sets the rule to that one channel (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams, Mattermost, Webhook or CustomChannel).
- `rule_type` (String) Type of compliance rule.
- `team_id` (String) ID of Team this compliance setting belongs to. The ID of a `oneuptime_team`.

### Read-Only

- `alert_severities` (Set of String) Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity. IDs of `oneuptime_alert_severity` resources.
- `created_at` (String) Date and Time when the object was created.
- `incident_severities` (Set of String) Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity. IDs of `oneuptime_incident_severity` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `notification_channels` (String) On-call rules only: the channels members must be notified on, as a list - each member needs a rule on every one of them (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams, Mattermost, Webhook or CustomChannel, which any custom channel that is on meets). Leave empty to accept any channel. A JSON value: read it with `jsondecode()`.
- `options` (String) Additional options for this compliance rule. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

---
page_title: "oneuptime_incident_reminder_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure reminder rules to periodically notify incident owners while an incident is still open
---

# oneuptime_incident_reminder_rule (Resource)

Configure reminder rules to periodically notify incident owners while an incident is still open

## Example Usage

```terraform
resource "oneuptime_incident_reminder_rule" "example" {
  name                         = "Example incident reminder rule"
  reminder_interval_in_minutes = 42
  description                  = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this incident reminder rule's place among the incident reminder rules of its project, lowest first. Setting it places this incident reminder rule and moves no other, so a configuration that numbers every incident reminder rule of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.incident_reminder_rules`:

```terraform
order = index(local.incident_reminder_rules, each.key) + 1
```

- Give each incident reminder rule a number of its own: two incident reminder rules at one number are shown the older first.
- An incident reminder rule created without `order` goes after the last one of its project.
- Deleting an incident reminder rule leaves a gap in the numbers, which changes nothing: incident reminder rules go by their numbers, not by counting.
- Moving an incident reminder rule in the OneUptime Dashboard numbers the project's incident reminder rules 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `name` (String) Name of this reminder rule.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind incident owners while the incident is still open. For example, set to 30 to remind owners every 30 minutes.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this reminder rule.
- `incident_severities` (Set of String) Only apply this reminder rule to incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `is_enabled` (Boolean) Whether this reminder rule is enabled. Defaults to `true`.
- `labels` (Set of String) Only apply this reminder rule to incidents with these labels. Leave empty to match incidents with any labels. IDs of `oneuptime_label` resources.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `stop_reminders_on_state` (String) Stop sending reminders once the incident reaches this state. Select Acknowledged to stop reminders when the incident is acknowledged, or Resolved to keep reminding until the incident is resolved. Allowed values: `Acknowledged`, `Resolved`. Defaults to `Resolved`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident reminder rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_reminder_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_reminder_rule.example <id>
```

---
page_title: "oneuptime_incident_on_call_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically executing on-call policies when matching incidents are created
---

# oneuptime_incident_on_call_rule (Data Source)

Configure rules for automatically executing on-call policies when matching incidents are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident on call rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_on_call_rule" "example" {
  name = "Example incident on call rule"
}

# Or by id:
data "oneuptime_incident_on_call_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_hold_time_in_minutes` (Number) When When to Page is After OneUptime AI: how many minutes OneUptime AI may hold the page while it investigates and fixes (1 to 30). Empty: 10. If OneUptime AI has not resolved the record by then, the policies page as usual.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this incident on-call rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description.
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.
- `name` (String) Name of this incident on-call rule.
- `page_as` (String) What the people this rule's on-call policies page are asked to do: Responders take the record on - acknowledging their page acknowledges it and stops every other page about it; Observers are kept informed while somebody else responds - acknowledging their page only says they have seen it, and the record is not acknowledged.
- `when_to_page` (String) When the on-call policies of this rule page: Right Away, as soon as a record matches it; Working Hours, only during the working hours of this rule - a record that matches it outside them is paged about when they begin; or After OneUptime AI, for alert and incident rules - when OneUptime AI may fix what the record is about, it responds first, and the policies page only if it has not resolved the record within the hold time.
- `working_hours_timezone` (String) The IANA time zone the working hours of this rule are in: Europe/London.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `incident_labels` (Set of String) Only trigger for incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only trigger for incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `monitor_labels` (Set of String) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) On-call policies to execute when an incident matches this rule. IDs of `oneuptime_on_call_policy` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `working_hours` (String) When this rule may page, when When to Page is Working Hours: the same daily or weekly windows as an on-call schedule layer's restriction times, read in the time zone of this rule. A JSON value: read it with `jsondecode()`.

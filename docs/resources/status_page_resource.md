---
page_title: "oneuptime_status_page_resource Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Add resources like monitors to your status page
---

# oneuptime_status_page_resource (Resource)

Add resources like monitors to your status page

## Example Usage

```terraform
resource "oneuptime_status_page_resource" "example" {
  status_page_id = oneuptime_status_page.example.id
  display_name   = "API Server"
}
```

<a id="order"></a>
## Order

`order` is this resource's place among the resources of its group (or, in no group, of its status page). Setting it places this resource and moves no other, so a configuration that numbers every resource of a group is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.resources`:

```terraform
order = index(local.resources, each.key) + 1
```

- Give each resource a number of its own: two resources at one number are shown in the order of their IDs.
- A resource created without `order` goes after the last resource of its group.
- Deleting a resource leaves a gap in the numbers, which changes nothing: resources go by their numbers, not by counting.
- Moving a resource in the OneUptime Dashboard numbers the group's resources 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `display_name` (String) Display name of the monitor on the Status Page.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `column_axis_value` (String) Column this resource belongs to when its status page group is rendered as a grid. Should match one of the column axis values defined on the group.
- `display_description` (String) Display description of the monitor on the Status Page. This is in markdown format.
- `display_tooltip` (String) Tooltip of the monitor on the Status Page.
- `monitor_group_id` (String) ID of the monitor group this resource shows, as one row. A resource shows one monitor or one monitor group: send monitorGroupId or monitorId, never both and never neither. The ID of a `oneuptime_monitor_group`.
- `monitor_id` (String) ID of the monitor this resource shows. A resource shows one monitor or one monitor group: send monitorId or monitorGroupId, never both and never neither. The ID of a `oneuptime_monitor`.
- `order` (Number) This resource's place in its group, or among the page's resources in no group, lowest first. Resources that share a place are shown in the order of their IDs. Setting it places this resource and moves no other, so a configuration that numbers every resource of a group ends up exactly as written. Left out when you create a resource, it goes after the last resource of its group. See [Order](#order).
- `row_axis_value` (String) Row this resource belongs to when its status page group is rendered as a grid. Should match one of the row axis values defined on the group.
- `show_current_status` (Boolean) Show current status like offline, operational or degraded. Defaults to `true`.
- `show_status_history_chart` (Boolean) Show a 90 day uptime history of this monitor. Defaults to `true`.
- `show_uptime_percent` (Boolean) Show uptime percent of this monitor for the last 90 days. Defaults to `false`.
- `status_page_group_id` (String) Does this monitor belong to a status page group? The ID of a `oneuptime_status_page_group`.
- `uptime_percent_precision` (String) Precision of uptime percent of this monitor for the last 90 days. Allowed values: `99% (No Decimal)`, `99.9% (One Decimal)`, `99.99% (Two Decimal)`, `99.999% (Three Decimal)`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_page_monitor_rule_id` (String) ID of the rule that added this resource, if it was added by a rule instead of by hand. The ID of a `oneuptime_status_page_monitor_rule`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page resource by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_resource.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_resource.example <id>
```

---
page_title: "oneuptime_status_page_group Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage groups on your status page and categorize resources like monitors into these groups.
---

# oneuptime_status_page_group (Resource)

Manage groups on your status page and categorize resources like monitors into these groups.

## Example Usage

```terraform
resource "oneuptime_status_page_group" "example" {
  status_page_id = oneuptime_status_page.example.id
  name           = "Example status page group"
  description    = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this group's place among the groups of its status page. Setting it places this group and moves no other, so a configuration that numbers every group of a status page is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.groups`:

```terraform
order = index(local.groups, each.key) + 1
```

- Give each group a number of its own: two groups at one number are shown in the order of their IDs.
- A group created without `order` goes after the page's last group.
- Deleting a group leaves a gap in the numbers, which changes nothing: groups go by their numbers, not by counting.
- Moving a group in the OneUptime Dashboard numbers the status page's groups 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `name` (String) Name of the Group.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `column_axis_label` (String) Label shown above the column axis when the group is rendered as a grid (e.g. 'Region', 'Environment'). Free-form so you can use any dimension you like.
- `column_axis_values` (String) Comma-separated list of column labels for the grid (e.g. 'US-East, EU-West, Asia'). Determines column order in the grid layout.
- `description` (String) Description for this group. This is visible on Status Page. This can be in markdown format.
- `is_expanded_by_default` (Boolean) Is this group expanded by default on Status Page? Defaults to `true`.
- `order` (Number) This group's place on its status page, lowest first. Groups that share a place are shown in the order of their IDs. Setting it places this group and moves no other, so a configuration that numbers every group ends up exactly as written. Left out when you create a group, it goes after the page's last group. See [Order](#order).
- `parent_status_page_group_id` (String) ID of the Status Page Group this group is nested under. Empty for top level groups. The ID of a `oneuptime_status_page_group`.
- `row_axis_label` (String) Label shown above the row axis when the group is rendered as a grid (e.g. 'Service', 'Tenant'). Free-form so you can use any dimension you like.
- `row_axis_values` (String) Comma-separated list of row labels for the grid (e.g. 'Auth, API, Database'). Determines row order in the grid layout.
- `show_current_status` (Boolean) Show current status like offline, operational or degraded. Defaults to `true`.
- `show_uptime_percent` (Boolean) Show uptime percent of this group for the last 90 days. Defaults to `false`.
- `uptime_percent_precision` (String) Precision of uptime percent of this group for the last 90 days. Allowed values: `99% (No Decimal)`, `99.9% (One Decimal)`, `99.99% (Two Decimal)`, `99.999% (Three Decimal)`.
- `view_mode` (String) Layout of this group on the status page. 'List' renders resources stacked vertically (default). 'Grid' renders resources as a matrix using row and column axes. Allowed values: `List`, `Grid`. Defaults to `List`.

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
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page group by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_group.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_group.example <id>
```

---
page_title: "oneuptime_status_page_history_chart_bar_color Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Modify the colors of the history chart bars on Status Page
---

# oneuptime_status_page_history_chart_bar_color (Resource)

Modify the colors of the history chart bars on Status Page

## Example Usage

```terraform
resource "oneuptime_status_page_history_chart_bar_color" "example" {
  status_page_id                          = oneuptime_status_page.example.id
  uptime_percent_greater_than_or_equal_to = 42
  bar_color                               = "#ff0000"
}
```

<a id="order"></a>
## Order

`order` is this status page history chart bar color's place among the status page history chart bar colors of its status page, lowest first. Setting it places this status page history chart bar color and moves no other, so a configuration that numbers every status page history chart bar color of a status page is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.status_page_history_chart_bar_colors`:

```terraform
order = index(local.status_page_history_chart_bar_colors, each.key) + 1
```

- Give each status page history chart bar color a number of its own: two status page history chart bar colors at one number are shown the older first.
- A status page history chart bar color created without `order` goes after the last one of its status page.
- Deleting a status page history chart bar color leaves a gap in the numbers, which changes nothing: status page history chart bar colors go by their numbers, not by counting.
- Moving a status page history chart bar color in the OneUptime Dashboard numbers the status page's status page history chart bar colors 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `bar_color` (String) Color of the bar chart when this rule matches (#32a852 for example).
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `uptime_percent_greater_than_or_equal_to` (Number) Uptime percent, from 0 to 100, that a day must reach for this rule to color its bar (99.65 for example).

### Optional

- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

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
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page history chart bar color by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_history_chart_bar_color.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_history_chart_bar_color.example <id>
```

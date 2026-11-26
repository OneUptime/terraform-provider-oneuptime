---
page_title: "oneuptime_alert_severity Resource - oneuptime"
subcategory: "Alerts"
description: |-
  How much an alert matters. A new project starts with High and Low: rename them, recolor them or add your own.
---

# oneuptime_alert_severity (Resource)

How much an alert matters. A new project starts with High and Low: rename them, recolor them or add your own.

## Example Usage

```terraform
resource "oneuptime_alert_severity" "example" {
  name        = "Example alert severity"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this alert severity's place among the alert severities of its project, lowest first. Setting it places this alert severity and moves no other, so a configuration that numbers every alert severity of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.alert_severities`:

```terraform
order = index(local.alert_severities, each.key) + 1
```

- Give each alert severity a number of its own: two alert severities at one number are shown the older first.
- An alert severity created without `order` goes after the last one of its project.
- Deleting an alert severity leaves a gap in the numbers, which changes nothing: alert severities go by their numbers, not by counting.
- Moving an alert severity in the OneUptime Dashboard numbers the project's alert severities 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Archived severities stay on the alerts and episodes that have them, but can't be picked again. A project keeps at least one severity that is not archived. Defaults to `false`.
- `order` (Number) Where this severity ranks among the project's alert severities: 1 is the most severe. A new severity without a number goes to the end of the list. Setting a number places the severity there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `archived_at` (String) When this alert severity was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
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

Import an existing alert severity by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_severity.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_severity.example <id>
```

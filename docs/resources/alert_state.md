---
page_title: "oneuptime_alert_state Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Where an alert is in your response. A new project starts with Identified, Acknowledged and Resolved: rename them or add your own between them.
---

# oneuptime_alert_state (Resource)

Where an alert is in your response. A new project starts with Identified, Acknowledged and Resolved: rename them or add your own between them.

## Example Usage

```terraform
resource "oneuptime_alert_state" "example" {
  name        = "Example alert state"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this alert state's place among the alert states of its project, lowest first. Setting it places this alert state and moves no other, so a configuration that numbers every alert state of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.alert_states`:

```terraform
order = index(local.alert_states, each.key) + 1
```

- Give each alert state a number of its own: two alert states at one number are shown the older first.
- An alert state created without `order` goes just above the resolved state: it takes the resolved state's number, and the resolved state and the alert states right below it step down one. Give it `order` when the configuration numbers the others, or the next plan shows theirs changed.
- The created state, acknowledged state and resolved state stay in that order, from the top: a write that would put two of them at one number, or out of order, is refused. An apply that moves them past each other's numbers can meet that refusal on one write; apply again and it finishes.
- Deleting an alert state leaves a gap in the numbers, which changes nothing: alert states go by their numbers, not by counting.
- Moving an alert state in the OneUptime Dashboard numbers the project's alert states 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_acknowledged_state` (Boolean) Is it the acknowledged state of the alert?
- `is_archived` (Boolean) Archived states stay on the alerts and episodes that are in them, but can't be picked again. The created, acknowledged and resolved states can't be archived. Defaults to `false`.
- `is_created_state` (Boolean) Is it the created state of the alert?
- `is_resolved_state` (Boolean) Is it the resolved state of the alert?
- `order` (Number) Where this state sits in the project's list of alert states: 1 is the top. Alerts only ever move down the list, and an alert in a state at or below the acknowledged (or resolved) state counts as acknowledged (or resolved), so the created, acknowledged and resolved states have to stay in that order. A new state without a number goes just above the resolved state. Setting a number places the state there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `archived_at` (String) When this alert state was archived. Empty while it is not archived.
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
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert state by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_state.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_state.example <id>
```

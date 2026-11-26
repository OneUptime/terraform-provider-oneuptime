---
page_title: "oneuptime_dashboard_playlist Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Dashboards shown one after another on a screen, such as a NOC wall, each for the same time.
---

# oneuptime_dashboard_playlist (Data Source)

Dashboards shown one after another on a screen, such as a NOC wall, each for the same time.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard playlist may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard_playlist" "example" {
  name = "Example dashboard playlist"
}

# Or by id:
data "oneuptime_dashboard_playlist" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `interval_in_seconds` (Number) How long each dashboard is shown before the next one, in seconds: 10 to 86400.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `dashboards` (String) The dashboards shown, in order: a list like [{"dashboardId": "<id>"}]. One may appear more than once. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

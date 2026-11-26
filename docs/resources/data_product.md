---
page_title: "oneuptime_data_product Resource - oneuptime"
subcategory: "Other"
description: |-
  A named group of watched tables with an owner, its consumers and a freshness SLA, on a status page and tracked as an SLO.
---

# oneuptime_data_product (Resource)

A named group of watched tables with an owner, its consumers and a freshness SLA, on a status page and tracked as an SLO.

## Example Usage

```terraform
resource "oneuptime_data_product" "example" {
  name        = "Example data product"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) What the data product is called, as its consumers know it.

### Optional

- `consumers` (String) Who uses the data: teams, dashboards, reports, AI agents.
- `data_assets` (Set of String) The watched tables the data product is made of. IDs of `oneuptime_data_asset` resources.
- `description` (String) What the data is, and what it is for.
- `fresh_by_time_of_day` (String) The time of day, in UTC (HH:MM), by which each day's data lands. Empty for no daily deadline.
- `max_data_age_in_minutes` (Number) How old the data may be, in minutes, before the product is late. Empty for no limit.
- `owner_team_id` (String) ID of the team that owns the data product. The ID of a `oneuptime_team`.
- `owner_user_email` (String) The email address of the person in `owner_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `owner_user_id` (String) ID of the person who owns the data product. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `owner_user_email`.
- `sla_target_percentage` (Number) The share of time the data must be on time, which its SLO tracks: 99 means late for at most 1% of the time.

### Read-Only

- `badge_token` (String, Sensitive) The secret in the badge's URL, so a badge can be embedded where people cannot sign in to OneUptime.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_as_of` (String) The oldest last change among the tables: what the data is as of.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `late_data_asset_ids` (String) The IDs of the tables that are late. A JSON value: read it with `jsondecode()`.
- `late_table_count` (Number) How many of the tables are late.
- `monitor_id` (String) ID of the Metrics monitor that is Degraded, with an incident, while a table is late. The ID of a `oneuptime_monitor`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `service_level_objective_id` (String) ID of the SLO that tracks the data product's SLA, with burn-rate alerts. The ID of a `oneuptime_service_level_objective`.
- `status` (String) Whether the data is fresh by the SLA: On Time, Late, or Unknown while no table has been seen changing.
- `status_changed_at` (String) When the status last changed.
- `status_checked_at` (String) When the tables were last checked against the SLA.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data product by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_product.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_product.example <id>
```

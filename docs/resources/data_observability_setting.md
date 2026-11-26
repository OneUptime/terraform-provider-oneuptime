---
page_title: "oneuptime_data_observability_setting Resource - oneuptime"
subcategory: "Other"
description: |-
  How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept.
---

# oneuptime_data_observability_setting (Resource)

How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept.

## Example Usage

```terraform
resource "oneuptime_data_observability_setting" "example" {

}
```

## Schema

### Optional

- `keep_failing_rows` (Boolean) A failing check keeps up to 20 of the rows that broke it, with sensitive columns masked. They are stored in OneUptime. Turning this off deletes them. Defaults to `false`.
- `max_gigabytes_scanned_per_check` (Number) BigQuery and Athena bill by bytes read. A check that would read more is stopped, and says so. Defaults to `10`.
- `wait_when_database_is_busy` (Boolean) Checks wait for the next collection when a database is close to its connection limit or has many statements running. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project these settings belong to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data observability setting by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_observability_setting.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_observability_setting.example <id>
```

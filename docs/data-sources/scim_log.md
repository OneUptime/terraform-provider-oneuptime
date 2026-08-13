---
page_title: "oneuptime_scim_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all SCIM provisioning operations for this project.
---

# oneuptime_scim_log (Data Source)

Logs of all SCIM provisioning operations for this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scim_log" "by_name" {
  name = "example-scim_log"
}

data "oneuptime_scim_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_scim_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `operation_type` (String) Type of SCIM operation (e.g., CreateUser, UpdateUser, DeleteUser, ListUsers, GetUser, CreateGroup, UpdateGroup, DeleteGroup, ListGroups, GetGroup, BulkOperation).. Computed.
- `status` (String) Status of the SCIM operation.. Computed.
- `status_message` (String) Short error or status description.. Computed.
- `log_body` (String) Detailed JSON with request/response data.. Computed.
- `http_method` (String) HTTP method used (GET, POST, PUT, PATCH, DELETE).. Computed.
- `request_path` (String) The SCIM endpoint path.. Computed.
- `http_status_code` (Number) Response HTTP status code.. Computed.
- `affected_user_email` (String) Email object. Computed.
- `affected_group_name` (String) Name of the group/team affected by this operation.. Computed.

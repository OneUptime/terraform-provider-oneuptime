---
page_title: "oneuptime_status_page_scim_log Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Logs of all SCIM provisioning operations for status pages.
---

# oneuptime_status_page_scim_log (Data Source)

Logs of all SCIM provisioning operations for status pages.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page scim log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_scim_log" "example" {
  status_page_id = oneuptime_status_page.example.id
}

# Or by id:
data "oneuptime_status_page_scim_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `http_method` (String) HTTP method used (GET, POST, PUT, PATCH, DELETE).
- `http_status_code` (Number) Response HTTP status code.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `log_body` (String) Detailed JSON with request/response data.
- `operation_type` (String) Type of SCIM operation (e.g., CreateUser, UpdateUser, DeleteUser, ListUsers, GetUser, BulkOperation).
- `request_path` (String) The SCIM endpoint path.
- `status` (String) Status of the SCIM operation.
- `status_message` (String) Short error or status description.
- `status_page_id` (String) ID of the Status Page. The ID of a `oneuptime_status_page`.
- `status_page_scim_id` (String) ID of your Status Page SCIM configuration. The ID of a `oneuptime_status_page_scim`.

### Read-Only

- `affected_user_email` (String) Email of the user affected by this operation.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

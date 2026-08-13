---
page_title: "oneuptime_audit_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Audit Log
---

# oneuptime_audit_log (Data Source)

API endpoints for Audit Log Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_audit_log" "by_name" {
  name = "example-audit_log"
}

data "oneuptime_audit_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `resource_type` (String) Resource Type. Computed.
- `resource_id` (String) Resource ID. Computed.
- `resource_name` (String) Resource Name. Computed.
- `action` (String) Action. Computed.
- `user_id` (String) User ID. Computed.
- `user_name` (String) User Name. Computed.
- `user_email` (String) User Email. Computed.
- `user_type` (String) User Type. Computed.
- `api_key_id` (String) API Key ID. Computed.
- `api_key_name` (String) API Key Name. Computed.
- `changes` (Set) Changes. Computed.

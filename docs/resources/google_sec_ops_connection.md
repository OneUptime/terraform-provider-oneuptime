---
page_title: "oneuptime_google_sec_ops_connection Resource - oneuptime"
subcategory: "Other"
description: |-
  Connections to Google SecOps (Chronicle) tenants. Detection alerts are polled on an interval and ingested as security events.
---

# oneuptime_google_sec_ops_connection (Resource)

Connections to Google SecOps (Chronicle) tenants. Detection alerts are polled on an interval and ingested as security events.

## Example Usage

```terraform
resource "oneuptime_google_sec_ops_connection" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  region = "Example short text"
  instance_resource_name = "This is an example of longer text content that might be stored in this field."
  service_account_json = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `name` (String) Name object.
- `region` (String) Google SecOps regional endpoint prefix, e.g. 'us' or 'europe'. Used to build the API base URL...
- `instance_resource_name` (String) The Chronicle instance resource name: projects/{project}/locations/{location}/instances/{instance}...
- `service_account_json` (String) Google Cloud service-account key (JSON) with Chronicle API read access. Encrypted at rest and never returned by the API...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `is_enabled` (Bool) Whether this connection is polled...
- `poll_interval_in_minutes` (Number) How often detection alerts are polled, in minutes...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_polled_at` (String) A date time object..
- `cursor` (String) Poll cursor: the newest detection timestamp already ingested, as an ISO string...
- `last_error` (String) The most recent poll error, if any. Cleared on the next successful poll...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_google_sec_ops_connection.example <id>
```

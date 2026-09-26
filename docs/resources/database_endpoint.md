---
page_title: "oneuptime_database_endpoint Resource - oneuptime"
subcategory: "Other"
description: |-
  Endpoints (host:port) a database is reached at. Telemetry that names one of these endpoints is shown on that database. Each endpoint belongs to at most one database in a project.
---

# oneuptime_database_endpoint (Resource)

Endpoints (host:port) a database is reached at. Telemetry that names one of these endpoints is shown on that database. Each endpoint belongs to at most one database in a project.

## Example Usage

```terraform
resource "oneuptime_database_endpoint" "example" {
  database_server_id = "123e4567-e89b-12d3-a456-426614174000"
  endpoint = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `database_server_id` (String) A unique identifier for an object, represented as a UUID..
- `endpoint` (String) Canonical endpoint of the database: host:port, with an @cluster qualifier for Kubernetes-internal names and private IPs (e.g. orders-db.data.svc.cluster.local:5432@prod-cluster). What you type is canonicalized; the default port of the engine is filled in...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `is_primary` (Bool) Is this the endpoint the database was created from? The primary endpoint cannot be removed...
- `source` (String) Who added this endpoint: auto (found in telemetry), workload (a Service name of the Kubernetes workload the database runs as) or user (added as an alias by a person)...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_matched_at` (String) A date time object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_database_endpoint.example <id>
```

---
page_title: "oneuptime_database_endpoint Data Source - oneuptime"
subcategory: "Other"
description: |-
  Endpoints (host:port) a database is reached at. Telemetry that names one of these endpoints is shown on that database. Each endpoint belongs to at most one database in a project.
---

# oneuptime_database_endpoint (Data Source)

Endpoints (host:port) a database is reached at. Telemetry that names one of these endpoints is shown on that database. Each endpoint belongs to at most one database in a project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_database_endpoint" "by_name" {
  name = "example-database_endpoint"
}

data "oneuptime_database_endpoint" "by_id" {
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
- `database_server_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `endpoint` (String) Canonical endpoint of the database: host:port, with an @cluster qualifier for Kubernetes-internal names and private IPs (e.g. orders-db.data.svc.cluster.local:5432@prod-cluster). What you type is canonicalized; the default port of the engine is filled in... Computed.
- `is_primary` (Bool) Is this the endpoint the database was created from? The primary endpoint cannot be removed... Computed.
- `source` (String) Who added this endpoint: auto (found in telemetry), workload (a Service name of the Kubernetes workload the database runs as) or user (added as an alias by a person)... Computed.
- `last_matched_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

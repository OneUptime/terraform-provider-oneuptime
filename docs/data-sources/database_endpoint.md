---
page_title: "oneuptime_database_endpoint Data Source - oneuptime"
subcategory: "Other"
description: |-
  Endpoints (host:port) a database is reached at. Telemetry that names one of these endpoints is shown on that database. Each endpoint belongs to at most one database in a project.
---

# oneuptime_database_endpoint (Data Source)

Endpoints (host:port) a database is reached at. Telemetry that names one of these endpoints is shown on that database. Each endpoint belongs to at most one database in a project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one database endpoint may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_database_endpoint" "example" {
  database_server_id = oneuptime_database.example.id
}

# Or by id:
data "oneuptime_database_endpoint" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `database_server_id` (String) ID of the database this endpoint belongs to. The ID of a `oneuptime_database`.
- `endpoint` (String) Canonical endpoint of the database: host:port, with an @cluster qualifier for Kubernetes-internal names and private IPs (e.g. orders-db.data.svc.cluster.local:5432@prod-cluster). What you type is canonicalized; the default port of the engine is filled in.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_primary` (Boolean) Is this the endpoint the database was created from? The primary endpoint cannot be removed.
- `source` (String) Who added this endpoint: auto (found in telemetry), workload (a Service name of the Kubernetes workload the database runs as) or user (added as an alias by a person).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_matched_at` (String) When telemetry or discovery last matched this endpoint to its database, refreshed at most once an hour. For a workload endpoint, when the workload last produced it.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

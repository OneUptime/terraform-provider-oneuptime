---
page_title: "oneuptime_database Resource - oneuptime"
subcategory: "Other"
description: |-
  Database servers this project runs or talks to. Each database is discovered from OpenTelemetry Collector database receivers, from database calls in application traces, from database workloads on monitored Kubernetes clusters and Docker / Podman hosts, or added manually.
---

# oneuptime_database (Resource)

Database servers this project runs or talks to. Each database is discovered from OpenTelemetry Collector database receivers, from database calls in application traces, from database workloads on monitored Kubernetes clusters and Docker / Podman hosts, or added manually.

## Example Usage

```terraform
resource "oneuptime_database" "example" {
  name = "Example short text"
  database_identifier = "This is an example of longer text content that might be stored in this field."
  db_system = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this database. Discovered databases are named after their engine and endpoint, e.g. PostgreSQL orders-db.internal:5432. Not unique - rename freely...
- `database_identifier` (String) Stable identity of this database within the project: the engine family and its canonical endpoint joined with '|' (e.g. postgresql|orders-db.internal:5432), or the workload form for databases discovered on Kubernetes / Docker / Podman (e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db). A fork keys like the engine it forks (a MariaDB database is mysql|..., a Valkey one redis|...), so learning which of the two a database is never changes its identity. Computed by the server for manually added databases...
- `db_system` (String) Database engine family as an OpenTelemetry db.system.name value, e.g. postgresql, mysql, redis, mongodb, microsoft.sql_server...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this database..
- `server_address` (String) Host name or IP address of the primary endpoint of this database (the OpenTelemetry server.address attribute)...
- `server_port` (Number) Port of the primary endpoint of this database (the OpenTelemetry server.port attribute). Defaults to the engine's standard port when not given...
- `discovery_source` (String) How this database was first discovered: collector, client-spans, kubernetes, docker, podman or manual...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this database archived? Archived databases are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry collected from this database by a collector / agent. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for telemetry collected from this database by a collector / agent (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the database default, then the project's retention settings...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `workload_identifier` (String) Identity of the Kubernetes / Docker / Podman workload this database runs as, e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db. The prefix is the engine family, not the engine (a MariaDB workload is mysql|..., a Valkey one redis|...). Empty for databases that are only known by their endpoint...
- `db_version` (String) Engine version last reported for this database, from the collector receiver or the container image tag...
- `kubernetes_cluster_id` (String) A unique identifier for an object, represented as a UUID..
- `kubernetes_namespace` (String) Kubernetes namespace of the workload this database runs as...
- `workload_kind` (String) Kind of workload this database runs as, e.g. StatefulSet, Deployment, Cluster (an operator-managed cluster) or Container...
- `workload_name` (String) Name of the workload (StatefulSet, Deployment, operator cluster or container group) this database runs as...
- `docker_host_id` (String) A unique identifier for an object, represented as a UUID..
- `podman_host_id` (String) A unique identifier for an object, represented as a UUID..
- `member_entity_keys` (String) Telemetry entity keys of the pods / containers this database runs as, each with the time it was last seen. Maintained by discovery...
- `instance_count` (Number) Number of running instances (pods / containers) last seen for this database's workload...
- `otel_collector_status` (String) Whether engine metrics are currently being received from a collector / agent for this database (connected) or a collector reported and has stopped (disconnected). Empty when no collector has ever reported...
- `collector_last_seen_at` (String) A date time object..
- `agent_version` (String) Version of the OneUptime agent reporting this database, as self-reported via the oneuptime.agent.version resource attribute..
- `last_seen_at` (String) A date time object..
- `auto_archived_at` (String) A date time object..
- `manually_restored_at` (String) A date time object..
- `db_system_source` (String) What the database engine was last determined from: manual, collector, container or client-spans. Stronger evidence may correct the engine; weaker evidence never changes it...
- `workload_last_seen_at` (String) A date time object..
- `automatic_assignments` (String) Label and owner ids that label rules, owner rules or telemetry attached automatically. Maintained by OneUptime...
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_database.example <id>
```

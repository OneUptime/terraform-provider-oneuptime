---
page_title: "oneuptime_database Data Source - oneuptime"
subcategory: "Other"
description: |-
  Database servers this project runs or talks to. Each database is discovered from OpenTelemetry Collector database receivers, from database calls in application traces, from database workloads on monitored Kubernetes clusters and Docker / Podman hosts, or added manually.
---

# oneuptime_database (Data Source)

Database servers this project runs or talks to. Each database is discovered from OpenTelemetry Collector database receivers, from database calls in application traces, from database workloads on monitored Kubernetes clusters and Docker / Podman hosts, or added manually. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_database" "by_name" {
  name = "example-database"
}

data "oneuptime_database" "by_id" {
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
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description for this database.. Computed.
- `database_identifier` (String) Stable identity of this database within the project: the engine family and its canonical endpoint joined with '|' (e.g. postgresql|orders-db.internal:5432), or the workload form for databases discovered on Kubernetes / Docker / Podman (e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db). A fork keys like the engine it forks (a MariaDB database is mysql|..., a Valkey one redis|...), so learning which of the two a database is never changes its identity. Computed by the server for manually added databases... Computed.
- `workload_identifier` (String) Identity of the Kubernetes / Docker / Podman workload this database runs as, e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db. The prefix is the engine family, not the engine (a MariaDB workload is mysql|..., a Valkey one redis|...). Empty for databases that are only known by their endpoint... Computed.
- `db_system` (String) Database engine family as an OpenTelemetry db.system.name value, e.g. postgresql, mysql, redis, mongodb, microsoft.sql_server... Computed.
- `server_address` (String) Host name or IP address of the primary endpoint of this database (the OpenTelemetry server.address attribute)... Computed.
- `server_port` (Number) Port of the primary endpoint of this database (the OpenTelemetry server.port attribute). Defaults to the engine's standard port when not given... Computed.
- `db_version` (String) Engine version last reported for this database, from the collector receiver or the container image tag... Computed.
- `discovery_source` (String) How this database was first discovered: collector, client-spans, kubernetes, docker, podman or manual... Computed.
- `kubernetes_cluster_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `kubernetes_namespace` (String) Kubernetes namespace of the workload this database runs as... Computed.
- `workload_kind` (String) Kind of workload this database runs as, e.g. StatefulSet, Deployment, Cluster (an operator-managed cluster) or Container... Computed.
- `workload_name` (String) Name of the workload (StatefulSet, Deployment, operator cluster or container group) this database runs as... Computed.
- `docker_host_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `podman_host_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `member_entity_keys` (String) Telemetry entity keys of the pods / containers this database runs as, each with the time it was last seen. Maintained by discovery... Computed.
- `instance_count` (Number) Number of running instances (pods / containers) last seen for this database's workload... Computed.
- `otel_collector_status` (String) Whether engine metrics are currently being received from a collector / agent for this database (connected) or a collector reported and has stopped (disconnected). Empty when no collector has ever reported... Computed.
- `collector_last_seen_at` (String) A date time object.. Computed.
- `agent_version` (String) Version of the OneUptime agent reporting this database, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `auto_archived_at` (String) A date time object.. Computed.
- `manually_restored_at` (String) A date time object.. Computed.
- `db_system_source` (String) What the database engine was last determined from: manual, collector, container or client-spans. Stronger evidence may correct the engine; weaker evidence never changes it... Computed.
- `workload_last_seen_at` (String) A date time object.. Computed.
- `automatic_assignments` (String) Label and owner ids that label rules, owner rules or telemetry attached automatically. Maintained by OneUptime... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this database archived? Archived databases are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry collected from this database by a collector / agent. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for telemetry collected from this database by a collector / agent (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the database default, then the project's retention settings... Computed.

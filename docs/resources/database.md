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
  name                = "Example database"
  database_identifier = "This is an example of longer text content that might be stored in this field."
  db_system           = "Example short text"
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `database_identifier` (String) Stable identity of this database within the project: the engine family and its canonical endpoint joined with '|' (e.g. postgresql|orders-db.internal:5432), or the workload form for databases discovered on Kubernetes / Docker / Podman (e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db). A fork keys like the engine it forks (a MariaDB database is mysql|..., a Valkey one redis|...), so learning which of the two a database is never changes its identity. Computed by the server for manually added databases.
- `db_system` (String) Database engine family as an OpenTelemetry db.system.name value, e.g. postgresql, mysql, redis, mongodb, microsoft.sql_server.
- `name` (String) Name of this database. Discovered databases are named after their engine and endpoint, e.g. PostgreSQL orders-db.internal:5432. Not unique - rename freely.

### Optional

- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this database server without approval even though they are riskier changes. Each pattern is one command line for this database server's agent (db, a fixed catalog of diagnostics, never free-form SQL) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this database server is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the database server can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this database server. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the database server's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the database server can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule. Defaults to `Disabled`.
- `description` (String) Friendly description for this database.
- `discovery_source` (String) How this database was first discovered: collector, client-spans, kubernetes, docker, podman or manual.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (db diagnostics from a fixed catalog such as ping, version and sessions; never free-form SQL) on this database server, through its Database AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the database server can turn it on or off. Defaults to `true`.
- `is_archived` (Boolean) Is this database archived? Archived databases are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry collected from this database by a collector / agent. Leave blank to use the project-wide default.
- `server_address` (String) Host name or IP address of the primary endpoint of this database (the OpenTelemetry server.address attribute).
- `server_port` (Number) Port of the primary endpoint of this database (the OpenTelemetry server.port attribute). Defaults to the engine's standard port when not given.
- `telemetry_retention_config` (String) Per-pillar retention overrides for telemetry collected from this database by a collector / agent (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the database default, then the project's retention settings. A JSON value: write it with `jsonencode()`.

### Read-Only

- `agent_version` (String) Version of the OneUptime agent reporting this database, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_access_configured_at` (String) When OneUptime AI access to this database server was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Database AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this database server, kept until the next successful command. Set by the server.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this database server through its Database AI agent. Set by the server.
- `archived_at` (String) When was this database archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `auto_archived_at` (String) When this database was archived automatically because no discovery source had seen it for a while. Empty when it was archived by a person or is not archived.
- `automatic_assignments` (String) Label and owner ids that label rules, owner rules or telemetry attached automatically. Maintained by OneUptime. A JSON value: write it with `jsonencode()`.
- `collector_last_seen_at` (String) When engine telemetry from a collector / agent was last received for this database.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `db_system_source` (String) What the database engine was last determined from: manual, collector, container or client-spans. Stronger evidence may correct the engine; weaker evidence never changes it.
- `db_version` (String) Engine version last reported for this database, from the collector receiver or the container image tag.
- `docker_host_id` (String) ID of the Docker host this database runs on, when it was discovered as a container on a monitored host. The ID of a `oneuptime_docker_host`.
- `id` (String) Unique identifier for the resource.
- `instance_count` (Number) Number of running instances (pods / containers) last seen for this database's workload.
- `kubernetes_cluster_id` (String) ID of the Kubernetes cluster this database runs on, when it was discovered as a workload on a monitored cluster. The ID of a `oneuptime_kubernetes_cluster`.
- `kubernetes_namespace` (String) Kubernetes namespace of the workload this database runs as.
- `last_seen_at` (String) When this database was last seen by any discovery source - engine telemetry, application database calls or the workload inventory.
- `manually_restored_at` (String) When a person last restored this database from the archive. Automatic archiving leaves it alone until it is seen again or a grace period passes.
- `member_entity_keys` (String) Telemetry entity keys of the pods / containers this database runs as, each with the time it was last seen. Maintained by discovery. A JSON value: write it with `jsonencode()`.
- `otel_collector_status` (String) Whether engine metrics are currently being received from a collector / agent for this database (connected) or a collector reported and has stopped (disconnected). Empty when no collector has ever reported.
- `podman_host_id` (String) ID of the Podman host this database runs on, when it was discovered as a container on a monitored host. The ID of a `oneuptime_podman_host`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.
- `workload_identifier` (String) Identity of the Kubernetes / Docker / Podman workload this database runs as, e.g. postgresql|kubernetes:prod-cluster/data/statefulset/orders-db. The prefix is the engine family, not the engine (a MariaDB workload is mysql|..., a Valkey one redis|...). Empty for databases that are only known by their endpoint.
- `workload_kind` (String) Kind of workload this database runs as, e.g. StatefulSet, Deployment, Cluster (an operator-managed cluster) or Container.
- `workload_last_seen_at` (String) When Kubernetes, Docker or Podman discovery last saw the workload this database runs as. Empty for a database that is not a discovered workload.
- `workload_name` (String) Name of the workload (StatefulSet, Deployment, operator cluster or container group) this database runs as.

## Import

Import an existing database by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_database.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_database.example <id>
```

---
page_title: "oneuptime_docker_swarm_cluster Data Source - oneuptime"
subcategory: "Other"
description: |-
  Docker Swarm clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime DockerSwarm Agent sends metrics, or can be manually registered.
---

# oneuptime_docker_swarm_cluster (Data Source)

Docker Swarm clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime DockerSwarm Agent sends metrics, or can be manually registered. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_docker_swarm_cluster" "by_name" {
  name = "example-docker_swarm_cluster"
}

data "oneuptime_docker_swarm_cluster" "by_id" {
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
- `description` (String) Friendly description for this DockerSwarm cluster.. Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).. Computed.
- `agent_version` (String) Version of the OneUptime DockerSwarm agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `docker_version` (String) Docker Engine version reported by the swarm manager this agent talks to... Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `node_count` (Number) Cached count of nodes in this cluster.. Computed.
- `ready_node_count` (Number) Cached count of nodes whose status is 'ready' in this cluster. Rendered as 'Nodes X/Y ready' next to nodeCount... Computed.
- `manager_node_count` (Number) Cached count of nodes with the manager role in this cluster... Computed.
- `service_count` (Number) Cached count of swarm services in this cluster.. Computed.
- `task_count` (Number) Cached count of swarm tasks (service instances) in this cluster.. Computed.
- `running_task_count` (Number) Cached count of tasks in the running state. Rendered as 'Tasks X/Y running' next to taskCount... Computed.
- `stack_count` (Number) Cached count of deployed compose stacks in this cluster.. Computed.
- `network_count` (Number) Cached count of swarm-scoped (overlay) networks in this cluster.. Computed.
- `swarm_id` (String) The Docker Swarm cluster ID (docker info -> Swarm.Cluster.ID) reported by the manager. Stable for the lifetime of the swarm; informational only — the join key is the cluster name... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this Docker Swarm cluster archived? Archived Docker Swarm clusters are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this DockerSwarm cluster. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this DockerSwarm cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the DockerSwarm cluster default, then the project's retention settings... Computed.
- `is_ai_investigation_enabled` (Bool) When on, OneUptime AI runs read-only commands (docker service ls, service ps, service logs, node ls) on this Docker Swarm cluster, through its Docker Swarm AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. Off by default. Anyone who may edit the Docker Swarm cluster can turn it on or off... Computed.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Docker Swarm cluster. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Docker Swarm cluster's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Docker Swarm cluster can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule... Computed.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Docker Swarm cluster without approval even though they are riskier changes. Each pattern is one command line for this Docker Swarm cluster's agent (docker) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Docker Swarm cluster is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Docker Swarm cluster can remove patterns or clear the list... Computed.
- `ai_access_last_verified_at` (String) A date time object.. Computed.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Docker Swarm cluster, kept until the next successful command. Set by the server... Computed.
- `ai_access_configured_at` (String) A date time object.. Computed.

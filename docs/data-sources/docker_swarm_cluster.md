---
page_title: "oneuptime_docker_swarm_cluster Data Source - oneuptime"
subcategory: "Other"
description: |-
  Docker Swarm clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime DockerSwarm Agent sends metrics, or can be manually registered.
---

# oneuptime_docker_swarm_cluster (Data Source)

Docker Swarm clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime DockerSwarm Agent sends metrics, or can be manually registered.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one docker swarm cluster may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_docker_swarm_cluster" "example" {
  name = "Example docker swarm cluster"
}

# Or by id:
data "oneuptime_docker_swarm_cluster" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the OneUptime DockerSwarm agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Docker Swarm cluster, kept until the next successful command. Set by the server.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Docker Swarm cluster. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Docker Swarm cluster's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Docker Swarm cluster can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description for this DockerSwarm cluster.
- `docker_version` (String) Docker Engine version reported by the swarm manager this agent talks to.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (docker service ls, service ps, service logs, node ls) on this Docker Swarm cluster, through its Docker Swarm AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Docker Swarm cluster can turn it on or off.
- `is_archived` (Boolean) Is this Docker Swarm cluster archived? Archived Docker Swarm clusters are hidden from lists but keep collecting telemetry.
- `manager_node_count` (Number) Cached count of nodes with the manager role in this cluster.
- `name` (String) Name of this DockerSwarm cluster. This is the join key — it must match the docker.swarm.cluster.name OTel resource attribute stamped by the OneUptime DockerSwarm Agent.
- `network_count` (Number) Cached count of swarm-scoped (overlay) networks in this cluster.
- `node_count` (Number) Cached count of nodes in this cluster.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `ready_node_count` (Number) Cached count of nodes whose status is 'ready' in this cluster. Rendered as 'Nodes X/Y ready' next to nodeCount.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this DockerSwarm cluster. Leave blank to use the project-wide default.
- `running_task_count` (Number) Cached count of tasks in the running state. Rendered as 'Tasks X/Y running' next to taskCount.
- `service_count` (Number) Cached count of swarm services in this cluster.
- `slug` (String) Friendly globally unique name for your object.
- `stack_count` (Number) Cached count of deployed compose stacks in this cluster.
- `swarm_id` (String) The Docker Swarm cluster ID (docker info -> Swarm.Cluster.ID) reported by the manager. Stable for the lifetime of the swarm; informational only — the join key is the cluster name.
- `task_count` (Number) Cached count of swarm tasks (service instances) in this cluster.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this Docker Swarm cluster was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Docker Swarm AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this Docker Swarm cluster through its Docker Swarm AI agent. Set by the server.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Docker Swarm cluster without approval even though they are riskier changes. Each pattern is one command line for this Docker Swarm cluster's agent (docker) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Docker Swarm cluster is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Docker Swarm cluster can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `archived_at` (String) When was this Docker Swarm cluster archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this cluster.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this DockerSwarm cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the DockerSwarm cluster default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

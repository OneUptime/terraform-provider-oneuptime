---
page_title: "oneuptime_docker_swarm_cluster Resource - oneuptime"
subcategory: "Other"
description: |-
  Docker Swarm clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime DockerSwarm Agent sends metrics, or can be manually registered.
---

# oneuptime_docker_swarm_cluster (Resource)

Docker Swarm clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime DockerSwarm Agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_docker_swarm_cluster" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this DockerSwarm cluster. This is the join key — it must match the docker.swarm.cluster.name OTel resource attribute stamped by the OneUptime DockerSwarm Agent...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this DockerSwarm cluster..
- `swarm_id` (String) The Docker Swarm cluster ID (docker info -> Swarm.Cluster.ID) reported by the manager. Stable for the lifetime of the swarm; informational only — the join key is the cluster name...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this Docker Swarm cluster archived? Archived Docker Swarm clusters are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this DockerSwarm cluster. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this DockerSwarm cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the DockerSwarm cluster default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime DockerSwarm agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `docker_version` (String) Docker Engine version reported by the swarm manager this agent talks to...
- `last_seen_at` (String) A date time object..
- `node_count` (Number) Cached count of nodes in this cluster..
- `ready_node_count` (Number) Cached count of nodes whose status is 'ready' in this cluster. Rendered as 'Nodes X/Y ready' next to nodeCount...
- `manager_node_count` (Number) Cached count of nodes with the manager role in this cluster...
- `service_count` (Number) Cached count of swarm services in this cluster..
- `task_count` (Number) Cached count of swarm tasks (service instances) in this cluster..
- `running_task_count` (Number) Cached count of tasks in the running state. Rendered as 'Tasks X/Y running' next to taskCount...
- `stack_count` (Number) Cached count of deployed compose stacks in this cluster..
- `network_count` (Number) Cached count of swarm-scoped (overlay) networks in this cluster..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_docker_swarm_cluster.example <id>
```

---
page_title: "oneuptime_proxmox_cluster Resource - oneuptime"
subcategory: "Other"
description: |-
  Proxmox VE clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Proxmox Agent sends metrics, or can be manually registered.
---

# oneuptime_proxmox_cluster (Resource)

Proxmox VE clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Proxmox Agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_proxmox_cluster" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Proxmox cluster. This is the join key — it must match the proxmox.cluster.name OTel resource attribute stamped by the OneUptime Proxmox Agent...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this Proxmox cluster..
- `ceph_cluster_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this Proxmox cluster archived? Archived Proxmox clusters are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Proxmox cluster. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Proxmox cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Proxmox cluster default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime Proxmox agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `pve_version` (String) Proxmox VE version reported by this cluster..
- `last_seen_at` (String) A date time object..
- `node_count` (Number) Cached count of nodes in this cluster..
- `online_node_count` (Number) Cached count of nodes currently online (pve_up == 1) in this cluster. Rendered as 'Nodes X/Y online' next to nodeCount...
- `guest_count` (Number) Cached count of guests (VMs and containers) in this cluster..
- `storage_count` (Number) Cached count of storage pools in this cluster..
- `guests_without_backup_count` (Number) Cached count of guests not covered by ANY backup job (pve_not_backed_up_total). NULL until the exporter's cluster-level backup-info collector reports. Coverage by a job is NOT the same as recent/successful backups — freshness needs the PVE task log or PBS API...
- `is_ai_investigation_enabled` (Bool) When on, OneUptime AI runs read-only commands (pvesh get on cluster, node, VM and container status) on this Proxmox cluster, through its Proxmox AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. Off by default. Anyone who may edit the Proxmox cluster can turn it on or off...
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Proxmox cluster. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Proxmox cluster's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Proxmox cluster can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule...
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Proxmox cluster without approval even though they are riskier changes. Each pattern is one command line for this Proxmox cluster's agent (pvesh) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Proxmox cluster is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Proxmox cluster can remove patterns or clear the list...

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
- `ai_access_last_verified_at` (String) A date time object..
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Proxmox cluster, kept until the next successful command. Set by the server...
- `ai_access_configured_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_proxmox_cluster.example <id>
```

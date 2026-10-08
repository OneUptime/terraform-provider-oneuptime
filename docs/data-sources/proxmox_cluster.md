---
page_title: "oneuptime_proxmox_cluster Data Source - oneuptime"
subcategory: "Other"
description: |-
  Proxmox VE clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Proxmox Agent sends metrics, or can be manually registered.
---

# oneuptime_proxmox_cluster (Data Source)

Proxmox VE clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Proxmox Agent sends metrics, or can be manually registered.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one proxmox cluster may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_proxmox_cluster" "example" {
  name = "Example proxmox cluster"
}

# Or by id:
data "oneuptime_proxmox_cluster" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the OneUptime Proxmox agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Proxmox cluster, kept until the next successful command. Set by the server.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Proxmox cluster. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Proxmox cluster's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Proxmox cluster can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `ceph_cluster_id` (String) Optional FK to the CephCluster providing storage for this hyperconverged Proxmox cluster (manually linked). The ID of a `oneuptime_ceph_cluster`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description for this Proxmox cluster.
- `guest_count` (Number) Cached count of guests (VMs and containers) in this cluster.
- `guests_without_backup_count` (Number) Cached count of guests not covered by ANY backup job (pve_not_backed_up_total). NULL until the exporter's cluster-level backup-info collector reports. Coverage by a job is NOT the same as recent/successful backups — freshness needs the PVE task log or PBS API.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (pvesh get on cluster, node, VM and container status) on this Proxmox cluster, through its Proxmox AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Proxmox cluster can turn it on or off.
- `is_archived` (Boolean) Is this Proxmox cluster archived? Archived Proxmox clusters are hidden from lists but keep collecting telemetry.
- `name` (String) Name of this Proxmox cluster. This is the join key — it must match the proxmox.cluster.name OTel resource attribute stamped by the OneUptime Proxmox Agent.
- `node_count` (Number) Cached count of nodes in this cluster.
- `online_node_count` (Number) Cached count of nodes currently online (pve_up == 1) in this cluster. Rendered as 'Nodes X/Y online' next to nodeCount.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `pve_version` (String) Proxmox VE version reported by this cluster.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Proxmox cluster. Leave blank to use the project-wide default.
- `slug` (String) Friendly globally unique name for your object.
- `storage_count` (Number) Cached count of storage pools in this cluster.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this Proxmox cluster was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Proxmox AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this Proxmox cluster through its Proxmox AI agent. Set by the server.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Proxmox cluster without approval even though they are riskier changes. Each pattern is one command line for this Proxmox cluster's agent (pvesh) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Proxmox cluster is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Proxmox cluster can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `archived_at` (String) When was this Proxmox cluster archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this cluster.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Proxmox cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Proxmox cluster default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

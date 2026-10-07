---
page_title: "oneuptime_ceph_cluster Resource - oneuptime"
subcategory: "Other"
description: |-
  Ceph clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Ceph Agent sends metrics, or can be manually registered.
---

# oneuptime_ceph_cluster (Resource)

Ceph clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Ceph Agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_ceph_cluster" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Ceph cluster. This is the join key — it must match the ceph.cluster.name OTel resource attribute stamped by the OneUptime Ceph Agent...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this Ceph cluster..
- `is_archived` (Bool) Is this Ceph cluster archived? Archived Ceph clusters are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Ceph cluster. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Ceph cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Ceph cluster default, then the project's retention settings...
- `fsid` (String) Ceph cluster fsid, sourced from the ceph.cluster.fsid OTel resource attribute when known..
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime Ceph agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `ceph_version` (String) Ceph version reported by this cluster..
- `last_seen_at` (String) A date time object..
- `mon_count` (Number) Cached count of Ceph monitors (mons) in this cluster..
- `osd_count` (Number) Cached count of OSDs in this cluster..
- `osd_up_count` (Number) Cached count of OSDs that are up (ceph_osd_up == 1) in this cluster. Rendered as 'X up / Y in / Z total' next to osdCount...
- `osd_in_count` (Number) Cached count of OSDs that are in the cluster (ceph_osd_in == 1). Rendered as 'X up / Y in / Z total' next to osdCount...
- `pool_count` (Number) Cached count of pools in this cluster..
- `health_status` (Number) Cached latest ceph_health_status value: 0 = HEALTH_OK, 1 = HEALTH_WARN, 2 = HEALTH_ERR. Rendered as the OK/WARN/ERR health pill. Null until the first metric batch arrives...
- `capacity_used_percent` (Number) Cached cluster capacity usage percent (ceph_cluster_total_used_bytes / ceph_cluster_total_bytes * 100). Stored as decimal so sub-percent precision survives the round trip. Null until both series appear in one metric batch...
- `is_ai_investigation_enabled` (Bool) When on, OneUptime AI runs read-only commands (ceph health, status, osd tree, df) on this Ceph cluster, through its Ceph AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Ceph cluster can turn it on or off...
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Ceph cluster. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Ceph cluster's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Ceph cluster can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule...
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Ceph cluster without approval even though they are riskier changes. Each pattern is one command line for this Ceph cluster's agent (ceph) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Ceph cluster is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Ceph cluster can remove patterns or clear the list...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `ai_access_last_verified_at` (String) A date time object..
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Ceph cluster, kept until the next successful command. Set by the server...
- `ai_access_configured_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_ceph_cluster.example <id>
```

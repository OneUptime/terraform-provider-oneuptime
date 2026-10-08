---
page_title: "oneuptime_vcenter Resource - oneuptime"
subcategory: "Other"
description: |-
  vSphere endpoints (a vCenter Server, or a standalone ESXi host) that are being monitored in this project. Each vCenter is auto-discovered when the OneUptime VMware Agent sends metrics, or can be manually registered.
---

# oneuptime_vcenter (Resource)

vSphere endpoints (a vCenter Server, or a standalone ESXi host) that are being monitored in this project. Each vCenter is auto-discovered when the OneUptime VMware Agent sends metrics, or can be manually registered.

~> **Renamed:** this resource was called `oneuptime_v_center` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing vcenter:

```terraform
moved {
  from = oneuptime_v_center.example
  to   = oneuptime_vcenter.example
}
```

## Example Usage

```terraform
resource "oneuptime_vcenter" "example" {
  name        = "Example vcenter"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this vCenter (a vCenter Server, or a standalone ESXi host). This is the join key — it must match the vmware.vcenter.name OTel resource attribute stamped by the OneUptime VMware Agent (VMWARE_VCENTER_NAME).

### Optional

- `agent_version` (String) Version of the OneUptime VMware Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this vCenter without approval even though they are riskier changes. Each pattern is one command line for this vCenter's agent (govc) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this vCenter is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the vCenter can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this vCenter. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the vCenter's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the vCenter can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule. Defaults to `Disabled`.
- `cluster_count` (Number) Cached count of vSphere clusters reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries cluster metrics.
- `datacenter_count` (Number) Cached count of vSphere datacenters reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries datacenter metrics.
- `datastore_capacity_bytes` (Number) Cached total capacity in bytes summed over every datastore reported through this vCenter (used + available vcenter.datastore.disk.usage). The denominator for the storage-used bar on the overview page. Stored as bigint.
- `datastore_count` (Number) Cached count of datastores reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries datastore metrics.
- `datastore_used_bytes` (Number) Cached used space in bytes summed over every datastore reported through this vCenter (vcenter.datastore.disk.usage{disk_state=used}). Stored as bigint.
- `description` (String) Friendly description for this vCenter.
- `host_count` (Number) Cached count of ESXi hosts reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries host metrics.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (govc about, ls, vm.info, host.info, events) on this vCenter, through its VMware AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the vCenter can turn it on or off. Defaults to `true`.
- `is_archived` (Boolean) Is this vCenter archived? Archived vCenters are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this vCenter.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `powered_on_vm_count` (Number) Cached count of virtual machines inferred to be powered on (the vcenter receiver emits CPU metrics only for powered-on VMs). Rendered as 'VMs X/Y powered on' next to vmCount.
- `resource_pool_count` (Number) Cached count of resource pools reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries resource pool metrics.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this vCenter. Leave blank to use the project-wide default.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this vCenter (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the vCenter default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `vm_count` (Number) Cached count of virtual machines (excluding VM templates) reported through this vCenter. Written by the ingest snapshot scan; only updated when a batch carries virtual machine metrics.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this vCenter was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a VMware AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this vCenter, kept until the next successful command. Set by the server.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this vCenter through its VMware AI agent. Set by the server.
- `archived_at` (String) When was this vCenter archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing vcenter by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_vcenter.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_vcenter.example <id>
```

---
page_title: "oneuptime_podman_host Resource - oneuptime"
subcategory: "Other"
description: |-
  Podman Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Podman Agent sends metrics, or can be manually registered.
---

# oneuptime_podman_host (Resource)

Podman Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Podman Agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_podman_host" "example" {
  name            = "Example podman host"
  host_identifier = "Example short text"
  description     = "Managed by Terraform"
}
```

## Schema

### Required

- `host_identifier` (String) Unique identifier for this Podman host, sourced from the host.name OTel resource attribute.
- `name` (String) Friendly name for this Podman host.

### Optional

- `agent_version` (String) Version of the OneUptime Podman agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Podman host without approval even though they are riskier changes. Each pattern is one command line for this Podman host's agent (docker, against the Podman API) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Podman host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Podman host can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Podman host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Podman host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Podman host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule. Defaults to `Disabled`.
- `containers_paused` (Number) Cached count of paused containers on this host.
- `containers_running` (Number) Cached count of running containers on this host.
- `containers_stopped` (Number) Cached count of stopped containers on this host.
- `description` (String) Friendly description for this Podman host.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (docker ps, inspect, logs, stats, events, against the Podman API) on this Podman host, through its Podman AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Podman host can turn it on or off. Defaults to `true`.
- `is_archived` (Boolean) Is this Podman host archived? Archived Podman hosts are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this host.
- `os_type` (String) Operating system type of the Podman host.
- `os_version` (String) Operating system version of the Podman host.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Podman host. Leave blank to use the project-wide default.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Podman host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Podman host default, then the project's retention settings. A JSON value: write it with `jsonencode()`.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this Podman host was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Podman AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Podman host, kept until the next successful command. Set by the server.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this Podman host through its Podman AI agent. Set by the server.
- `archived_at` (String) When was this Podman host archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing podman host by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_podman_host.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_podman_host.example <id>
```

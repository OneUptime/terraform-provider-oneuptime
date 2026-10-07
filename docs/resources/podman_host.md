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
  name = "Example short text"
  host_identifier = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this Podman host..
- `host_identifier` (String) Unique identifier for this Podman host, sourced from the host.name OTel resource attribute..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this Podman host..
- `is_archived` (Bool) Is this Podman host archived? Archived Podman hosts are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Podman host. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Podman host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Podman host default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime Podman agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `last_seen_at` (String) A date time object..
- `containers_running` (Number) Cached count of running containers on this host..
- `containers_stopped` (Number) Cached count of stopped containers on this host..
- `containers_paused` (Number) Cached count of paused containers on this host..
- `os_type` (String) Operating system type of the Podman host..
- `os_version` (String) Operating system version of the Podman host..
- `is_ai_investigation_enabled` (Bool) When on, OneUptime AI runs read-only commands (docker ps, inspect, logs, stats, events, against the Podman API) on this Podman host, through its Podman AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Podman host can turn it on or off...
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Podman host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Podman host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Podman host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule...
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Podman host without approval even though they are riskier changes. Each pattern is one command line for this Podman host's agent (docker, against the Podman API) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Podman host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Podman host can remove patterns or clear the list...

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
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Podman host, kept until the next successful command. Set by the server...
- `ai_access_configured_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_podman_host.example <id>
```

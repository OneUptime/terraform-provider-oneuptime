---
page_title: "oneuptime_docker_host Data Source - oneuptime"
subcategory: "Other"
description: |-
  Docker Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Docker Agent sends metrics, or can be manually registered.
---

# oneuptime_docker_host (Data Source)

Docker Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Docker Agent sends metrics, or can be manually registered. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_docker_host" "by_name" {
  name = "example-docker_host"
}

data "oneuptime_docker_host" "by_id" {
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
- `description` (String) Friendly description for this Docker host.. Computed.
- `host_identifier` (String) Unique identifier for this Docker host, sourced from the host.name OTel resource attribute.. Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).. Computed.
- `agent_version` (String) Version of the OneUptime Docker agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `containers_running` (Number) Cached count of running containers on this host.. Computed.
- `containers_stopped` (Number) Cached count of stopped containers on this host.. Computed.
- `containers_paused` (Number) Cached count of paused containers on this host.. Computed.
- `os_type` (String) Operating system type of the Docker host.. Computed.
- `os_version` (String) Operating system version of the Docker host.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this Docker host archived? Archived Docker hosts are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Docker host. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Docker host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Docker host default, then the project's retention settings... Computed.
- `is_ai_investigation_enabled` (Bool) When on, OneUptime AI runs read-only commands (docker ps, inspect, logs, stats, events) on this Docker host, through its Docker AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. Off by default. Anyone who may edit the Docker host can turn it on or off... Computed.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Docker host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Docker host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Docker host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule... Computed.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Docker host without approval even though they are riskier changes. Each pattern is one command line for this Docker host's agent (docker) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Docker host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Docker host can remove patterns or clear the list... Computed.
- `ai_access_last_verified_at` (String) A date time object.. Computed.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Docker host, kept until the next successful command. Set by the server... Computed.
- `ai_access_configured_at` (String) A date time object.. Computed.

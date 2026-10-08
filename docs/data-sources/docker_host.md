---
page_title: "oneuptime_docker_host Data Source - oneuptime"
subcategory: "Other"
description: |-
  Docker Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Docker Agent sends metrics, or can be manually registered.
---

# oneuptime_docker_host (Data Source)

Docker Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Docker Agent sends metrics, or can be manually registered.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one docker host may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_docker_host" "example" {
  name = "Example docker host"
}

# Or by id:
data "oneuptime_docker_host" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the OneUptime Docker agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit while running a command on this Docker host, kept until the next successful command. Set by the server.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this Docker host. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the Docker host's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). Anyone who may edit the Docker host can lower the mode; raising it needs Project Owner, Project Admin or Edit Auto Remediation Rule.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `containers_paused` (Number) Cached count of paused containers on this host.
- `containers_running` (Number) Cached count of running containers on this host.
- `containers_stopped` (Number) Cached count of stopped containers on this host.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description for this Docker host.
- `host_identifier` (String) Unique identifier for this Docker host, sourced from the host.name OTel resource attribute.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI runs read-only commands (docker ps, inspect, logs, stats, events) on this Docker host, through its Docker AI agent, while investigating incidents and alerts linked to it, and uses their output, with secret values redacted, as evidence. Nothing is ever changed by an investigation. On by default. Anyone who may edit the Docker host can turn it on or off.
- `is_archived` (Boolean) Is this Docker host archived? Archived Docker hosts are hidden from lists but keep collecting telemetry.
- `name` (String) Friendly name for this Docker host.
- `os_type` (String) Operating system type of the Docker host.
- `os_version` (String) Operating system version of the Docker host.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Docker host. Leave blank to use the project-wide default.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this Docker host was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Docker AI agent that registers later never overwrites a setting an operator chose.
- `ai_access_last_verified_at` (String) When a command from OneUptime AI last succeeded on this Docker host through its Docker AI agent. Set by the server.
- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run on this Docker host without approval even though they are riskier changes. Each pattern is one command line for this Docker host's agent (docker) and is compared with the command word by word: * stands for exactly one word (a name, an id), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid write command for this Docker host is refused. Destructive commands (Denied tier) never run regardless, and a command that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit Auto Remediation Rule; anyone who may edit the Docker host can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `archived_at` (String) When was this Docker host archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this host.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Docker host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Docker host default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

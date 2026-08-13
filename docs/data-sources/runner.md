---
page_title: "oneuptime_runner Data Source - oneuptime"
subcategory: "Other"
description: |-
  A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them.
---

# oneuptime_runner (Data Source)

A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_runner" "by_name" {
  name = "example-runner"
}

data "oneuptime_runner" "by_id" {
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
- `description` (String) Optional description for this agent.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `key` (String) Secret key the agent presents on every request. Anyone who can read this key can claim work as this Runner and receive its secrets in plaintext. Never share it; reset it to revoke the agent... Computed.
- `agent_version` (String) Version object. Computed.
- `last_alive` (String) A date time object.. Computed.
- `connection_status` (String) Connected if the agent has heartbeated recently... Computed.
- `can_run_runbooks` (Bool) Whether this Runner executes runbook steps. On by default — this is why most Runners are installed... Computed.
- `can_run_code_fix_tasks` (Bool) Whether this Runner works in your code repository to open AI fix pull requests. Off by default; it requires a connected code repository... Computed.
- `can_run_ai_commands` (Bool) Whether AI auto-remediation may execute commands on this Runner. Off by default. Commands are policy-checked and either match an operator allowlist or require one-click human approval... Computed.
- `host_info` (String) Self-reported host info (hostname, OS, arch). Updated on each heartbeat... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.

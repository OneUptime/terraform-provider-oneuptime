---
page_title: "oneuptime_runner Resource - oneuptime"
subcategory: "Other"
description: |-
  A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them.
---

# oneuptime_runner (Resource)

A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them.

## Example Usage

```terraform
resource "oneuptime_runner" "example" {
  name = "Example short text"
  key = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this agent..
- `key` (String) Secret key the agent presents on every request. Anyone who can read this key can claim work as this Runner and receive its secrets in plaintext. Never share it; reset it to revoke the agent...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Optional description for this agent..
- `agent_version` (String) Version object.
- `connection_status` (String) Connected if the agent has heartbeated recently...
- `can_run_runbooks` (Bool) Whether this Runner executes runbook steps. On by default — this is why most Runners are installed...
- `can_run_code_fix_tasks` (Bool) Whether this Runner works in your code repository to open AI fix pull requests. Off by default; it requires a connected code repository...
- `can_run_ai_commands` (Bool) Whether AI auto-remediation may execute commands on this Runner. Off by default. Commands are policy-checked and either match an operator allowlist or require one-click human approval...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `last_alive` (String) A date time object..
- `host_info` (String) Self-reported host info (hostname, OS, arch). Updated on each heartbeat...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runner.example <id>
```

---
page_title: "oneuptime_runner Data Source - oneuptime"
subcategory: "Other"
description: |-
  A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them.
---

# oneuptime_runner (Data Source)

A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure and, when the capability is enabled, works in your code repository to open AI fix pull requests. Runbook steps pick the Runner that should execute them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one runner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_runner" "example" {
  name = "Example runner"
}

# Or by id:
data "oneuptime_runner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `can_run_ai_commands` (Boolean) Whether OneUptime AI may run commands through this Runner: read-only kubectl while investigating a cluster it is bound to, and policy-checked remediation commands (Bash, SSH, kubectl) that either match an allowlist or wait for one-click human approval. Off by default; the in-cluster Runner installed by the Kubernetes agent chart turns it on.
- `can_run_code_fix_tasks` (Boolean) Whether this Runner works in your code repository to open AI fix pull requests. Off by default; it requires a connected code repository. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only.
- `can_run_runbooks` (Boolean) Whether this Runner executes runbook steps. On by default — this is why most Runners are installed. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only.
- `connection_status` (String) Connected if the agent has heartbeated recently.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Optional description for this agent.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `key` (String) Secret key the agent presents on every request. Anyone who can read this key can claim work as this Runner and receive its secrets in plaintext. Never share it; reset it to revoke the agent.
- `name` (String) Friendly name for this agent. Names starting with "kubernetes-agent/" are reserved for the in-cluster Runners the Kubernetes agent chart registers, and such a Runner cannot be renamed.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `agent_version` (String) Self-reported version of the Runner binary. Updated on each heartbeat.
- `created_at` (String) Date and Time when the object was created.
- `host_info` (String) Self-reported host info (hostname, OS, arch). Updated on each heartbeat. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_alive` (String) Most recent heartbeat from this agent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

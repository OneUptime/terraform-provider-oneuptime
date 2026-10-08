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
  name        = "Example runner"
  key         = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `key` (String) Secret key the agent presents on every request. Anyone who can read this key can claim work as this Runner and receive its secrets in plaintext. Never share it; reset it to revoke the agent.
- `name` (String) Friendly name for this agent. Names starting with "kubernetes-agent/" are reserved for the in-cluster Runners the Kubernetes agent chart registers, and such a Runner cannot be renamed.

### Optional

- `agent_version` (String) Self-reported version of the Runner binary. Updated on each heartbeat.
- `can_run_ai_commands` (Boolean) Whether OneUptime AI may run commands through this Runner: read-only kubectl while investigating a cluster it is bound to, and policy-checked remediation commands (Bash, SSH, kubectl) that either match an allowlist or wait for one-click human approval. Off by default; the in-cluster Runner installed by the Kubernetes agent chart turns it on. Defaults to `false`.
- `can_run_code_fix_tasks` (Boolean) Whether this Runner works in your code repository to open AI fix pull requests. Off by default; it requires a connected code repository. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only. Defaults to `false`.
- `can_run_runbooks` (Boolean) Whether this Runner executes runbook steps. On by default — this is why most Runners are installed. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only. Defaults to `true`.
- `connection_status` (String) Connected if the agent has heartbeated recently.
- `description` (String) Optional description for this agent.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `host_info` (String) Self-reported host info (hostname, OS, arch). Updated on each heartbeat. A JSON value: write it with `jsonencode()`.
- `id` (String) Unique identifier for the resource.
- `last_alive` (String) Most recent heartbeat from this agent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runner.example <id>
```

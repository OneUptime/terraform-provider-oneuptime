---
page_title: "oneuptime_runner Resource - oneuptime"
subcategory: "Runbooks & Runners"
description: |-
  A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure, and AI works on it in the code repositories that pick it. Runbook steps and code repositories each pick the Runner that runs them.
---

# oneuptime_runner (Resource)

A self-hosted OneUptime Runner: it executes runbook steps in your own infrastructure, and AI works on it in the code repositories that pick it. Runbook steps and code repositories each pick the Runner that runs them.

## Example Usage

```terraform
resource "oneuptime_runner" "example" {
  name        = "Example runner"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this agent. Names starting with "kubernetes-agent/" are reserved for the in-cluster Runners the Kubernetes agent chart registers, and such a Runner cannot be renamed.

### Optional

- `agent_version` (String) Self-reported version of the Runner binary. Updated on each heartbeat.
- `can_run_runbooks` (Boolean) Whether this Runner executes runbook steps. On by default — this is why most Runners are installed. It cannot be turned on for an in-cluster Runner the Kubernetes agent chart registered, which runs kubectl only. Defaults to `true`.
- `connection_status` (String) Connected if the agent has heartbeated recently. Allowed values: `connected`, `disconnected`.
- `description` (String) Optional description for this agent.
- `is_project_default` (Boolean) Whether this is the project's default Runner: runbook steps that pick no Runner run on it instead of the Global Runner. Defaults to `false`.
- `key` (String, Sensitive) Secret key the agent presents on every request. Anyone who can read this key can claim work as this Runner and receive its secrets in plaintext. Never share it; reset it to revoke the agent. Leave it out and OneUptime generates one.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.

### Read-Only

- `can_run_ai_commands` (Boolean) Whether OneUptime AI may run commands through this Runner: read-only kubectl while investigating a cluster it is bound to, and policy-checked remediation commands (Bash, kubectl) that either match an allowlist or wait for one-click human approval. Set by OneUptime: on once an auto remediation rule picks this Runner, and on for the in-cluster Runner the Kubernetes agent chart installs.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `host_info` (String) Self-reported host info (hostname, OS, arch). Updated on each heartbeat. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `is_global_runner` (Boolean) Whether OneUptime runs this Runner for every project. Global Runners run JavaScript steps in a sandbox that can reach the internet but not your network. Only a master admin creates them.
- `last_alive` (String) Most recent heartbeat from this agent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `takes_code_fixes` (Boolean) Whether this Runner takes OneUptime AI's code fixes, as it said on its last heartbeat. Empty for a Runner of a release that does not say.
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

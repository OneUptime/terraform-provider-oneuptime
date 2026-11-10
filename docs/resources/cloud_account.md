---
page_title: "oneuptime_cloud_account Resource - oneuptime"
subcategory: "Cloud"
description: |-
  AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them.
---

# oneuptime_cloud_account (Resource)

AWS accounts, Azure subscriptions and Google Cloud projects: each with every resource its cloud connections list, and what changed in them.

## Example Usage

```terraform
resource "oneuptime_cloud_account" "example" {
  name           = "Example cloud account"
  provider_value = "aws"
  account_id     = "123456789012"
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `account_id` (String) The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.
- `name` (String) What the account is called in OneUptime. Starts as the name its cloud gives it, and can be changed.
- `provider_value` (String) Which cloud the account is in: 'aws', 'azure' or 'gcp'. Allowed values: `aws`, `gcp`, `azure`.

### Optional

- `ai_command_allowlist` (String) Optional JSON array of command patterns that Automatic mode may run in this account without approval even though they are riskier changes. Each pattern is one cloud catalog command (cloud <operation> RESOURCE [--flag value]) and is compared with the command word by word: * stands for exactly one word (a resource id, a count), never for extra words or flags, and every flag the command uses must be written out in the pattern. At most 50 patterns of at most 500 characters each; a pattern that is not one valid change for this account's cloud is refused. Destructive operations (Denied tier) never run regardless, and an operation that always needs a human still asks. Adding a pattern needs Project Owner, Project Admin or Edit AI Rule; anyone who may edit the account can remove patterns or clear the list. A JSON value: write it with `jsonencode()`.
- `ai_remediation_mode` (String) Disabled: AI never proposes or runs a change on this account. RequireApproval: AI composes a command plan and a human approves it with one click before anything runs. Automatic: safe changes (SafeWrite) run without a human; a riskier change is proposed for approval unless the account's allowlist names its exact shape. BypassApproval: every change the policy allows — safe AND riskier — runs on its own, except what always needs a human. In EVERY mode: Denied commands never run, commands the policy marks requiresHuman always ask, and the agent itself refuses every write unless it was started with ONEUPTIME_AI_ALLOW_WRITES=true (and then only on the targets ONEUPTIME_AI_WRITE_TARGETS allows, never its protected targets). The Cloud AI agent's own configuration (ONEUPTIME_AI_FIXES, which its setup command writes) sets the mode once it connects. Anyone who may edit the account can lower the mode; raising it needs Project Owner, Project Admin or Edit AI Rule. Allowed values: `Disabled`, `RequireApproval`, `Automatic`, `BypassApproval`. Defaults to `Disabled`.
- `description` (String) What the account is for, in your team's words.
- `is_ai_investigation_enabled` (Boolean) When on, OneUptime AI reads this account live while investigating incidents and alerts linked to its resources — a resource's settings and status, its recent events, load balancer target health and quota usage — through the account's read-only cloud connection, and uses what it read as evidence; it needs no Cloud AI agent for that. Nothing is ever changed by an investigation. On by default. Anyone who may edit the account can turn it on or off. Defaults to `true`.
- `is_ai_log_read_enabled` (Boolean) When on, an investigation may also read this account's recent log lines (CloudWatch Logs, Log Analytics, Cloud Logging) through the read-only connection. Log lines can hold customer data, so this is off until someone turns it on. Anyone who may edit the account can turn it on or off. Defaults to `false`.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.

### Read-Only

- `ai_access_configured_at` (String) When OneUptime AI access to this account was first configured by anyone saving an AI access setting. Set by the server; never cleared, so a Cloud AI agent that connects later never overwrites a setting an operator chose.
- `ai_access_last_error` (String) The most recent failure OneUptime AI hit in this account — an operation that failed, or a Cloud AI agent that could not connect to it — kept until the next success. Set by the server.
- `ai_access_last_verified_at` (String) When an operation from OneUptime AI last succeeded in this account through its Cloud AI agent. Set by the server.
- `ai_agent_reset_at` (String) When this account's Cloud AI agent was last reset: every key it had stopped working, and it stays disconnected until its setup command runs again. Set by the server.
- `cloud_connection_id` (String) ID of the connection that reads this account, if one does. The ID of a `oneuptime_cloud_connection`.
- `cost_summary` (String) What the account cost this month so far, at its pace and last month, from its cloud's billing, and whether its costs could be read. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this account. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `inventory_listed_at` (String) When the account's inventory last listed every resource it has.
- `last_seen_at` (String) When a cloud connection last reached the account.
- `project_id` (String) ID of the project this account belongs to. The ID of a `oneuptime_project`.
- `provider_account_name` (String) What the cloud calls the account: an AWS account alias, an Azure subscription's or a Google Cloud project's name.
- `region_count` (Number) How many regions the account's live resources are in.
- `resource_count` (Number) How many live resources the account has in OneUptime.
- `resource_summary` (String) The account's live resources counted by type and by region (CloudAccountSummary). A JSON value: read it with `jsondecode()`.
- `savings_summary` (String) How many ways to spend less the account's cloud recommends and what they would save a month, and whether they could be read. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing cloud account by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_account.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_account.example <id>
```

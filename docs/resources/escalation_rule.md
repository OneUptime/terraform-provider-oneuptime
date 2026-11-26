---
page_title: "oneuptime_escalation_rule Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call escalation rule for the on-call policy.
---

# oneuptime_escalation_rule (Resource)

Manage on-call escalation rule for the on-call policy.

## Example Usage

```terraform
resource "oneuptime_escalation_rule" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
  name                   = "Example escalation rule"
  description            = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this rule's place among the levels of its on-call policy, the lowest paged first. Setting it places this rule and moves no other, so a configuration that numbers every rule of a policy is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.levels`:

```terraform
order = index(local.levels, each.key) + 1
```

- Give each rule a number of its own: two rules at one number are paged one after the other, in the order of their IDs.
- A rule created without `order` becomes the policy's last level.
- Deleting a rule leaves a gap in the numbers, which changes nothing: rules go by their numbers, not by counting.
- Moving a rule in the OneUptime Dashboard numbers the policy's rules 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `name` (String) Any friendly name of this object.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs. The ID of a `oneuptime_on_call_policy`.

### Optional

- `description` (String) Friendly description that will help you remember.
- `escalate_after_in_minutes` (Number) How long should we wait before we execute the next escalation rule?
- `order` (Number) This rule's level in its on-call policy: the lowest is paged first, then the next, and so on. Rules that share a level are paged in the order of their IDs. Setting it places this rule and moves no other, so a configuration that numbers every rule of a policy ends up exactly as written. Left out when you create a rule, it becomes the policy's last level. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing escalation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_escalation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_escalation_rule.example <id>
```

---
page_title: "oneuptime_incoming_call_policy_escalation_rule Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage escalation rules for incoming call policies that define who to call and in what order
---

# oneuptime_incoming_call_policy_escalation_rule (Resource)

Manage escalation rules for incoming call policies that define who to call and in what order

## Example Usage

```terraform
resource "oneuptime_incoming_call_policy_escalation_rule" "example" {
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
  name                    = "Example incoming call policy escalation rule"
  description             = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this incoming call policy escalation rule's place among the incoming call policy escalation rules of its incoming call policy, lowest first. Setting it places this incoming call policy escalation rule and moves no other, so a configuration that numbers every incoming call policy escalation rule of an incoming call policy is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.incoming_call_policy_escalation_rules`:

```terraform
order = index(local.incoming_call_policy_escalation_rules, each.key) + 1
```

- Give each incoming call policy escalation rule a number of its own: two incoming call policy escalation rules at one number are shown the older first.
- An incoming call policy escalation rule created without `order` goes after the last one of its incoming call policy.
- Deleting an incoming call policy escalation rule leaves a gap in the numbers, which changes nothing: incoming call policy escalation rules go by their numbers, not by counting.
- Moving an incoming call policy escalation rule in the OneUptime Dashboard numbers the incoming call policy's incoming call policy escalation rules 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `incoming_call_policy_id` (String) ID of the parent Incoming Call Policy. The ID of a `oneuptime_incoming_call_policy`.

### Optional

- `description` (String) Optional description of this escalation rule.
- `escalate_after_seconds` (Number) How long, in seconds, the phone rings before the call moves on to the next rule. 20 when left out; a time below 5 or above 600 rings for 5 or 600, the limits Twilio takes. Defaults to `20`.
- `name` (String) Rule name (e.g., 'Primary On-Call', 'Backup Engineer').
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule to route to (mutually exclusive with userId). The ID of a `oneuptime_on_call_policy_schedule`.
- `order` (Number) Where this rule sits in the escalation, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user to route to directly (mutually exclusive with onCallDutyPolicyScheduleId). The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

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

Import an existing incoming call policy escalation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_call_policy_escalation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_call_policy_escalation_rule.example <id>
```

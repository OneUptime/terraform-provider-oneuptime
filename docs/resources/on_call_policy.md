---
page_title: "oneuptime_on_call_policy Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Who is paged about an incident, an alert or an episode, and who is paged next if nobody acknowledges. Its levels are escalation rules.
---

# oneuptime_on_call_policy (Resource)

Who is paged about an incident, an alert or an episode, and who is paged next if nobody acknowledges. Its levels are escalation rules.

## Example Usage

```terraform
resource "oneuptime_on_call_policy" "example" {
  name        = "Example on call policy"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `acknowledgement_timeout_in_minutes` (Number) Page again when something this policy paged about is acknowledged but not resolved within this many minutes (5 to 1440). Empty: never.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Archived on-call policies are hidden from the On-Call Policies list and page no one: incidents and alerts that use them skip them. Unarchiving puts them back in service. Defaults to `false`.
- `labels` (Set of String) Labels to group and filter by. OneUptime adds to it too, on its own: the entries listed here are added and taken off, and every other entry stays. Left out, it keeps the entries OneUptime gives it. IDs of `oneuptime_label` resources.
- `repeat_policy_if_no_one_acknowledges` (Boolean) Repeat the policy if no one acknowledges the alert. Defaults to `false`.
- `repeat_policy_if_no_one_acknowledges_no_of_times` (Number) How many more times the policy starts again from its first escalation rule when nobody acknowledges: a whole number of 1 or more while Repeat Policy If No One Acknowledges is on. Defaults to `0`.

### Read-Only

- `archived_at` (String) When this on-call policy was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
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
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call policy by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_policy.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_policy.example <id>
```

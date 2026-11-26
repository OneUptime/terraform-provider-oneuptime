---
page_title: "oneuptime_incoming_request_secret Data Source - oneuptime"
subcategory: "Other"
description: |-
  A secret that requests to the project's Incoming Request monitors must send while the project requires one: in the X-OneUptime-Secret header, as a Bearer token, or as a signature of the request's body.
---

# oneuptime_incoming_request_secret (Data Source)

A secret that requests to the project's Incoming Request monitors must send while the project requires one: in the X-OneUptime-Secret header, as a Bearer token, or as a signature of the request's body.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming request secret may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_request_secret" "example" {
  name = "Example incoming request secret"
}

# Or by id:
data "oneuptime_incoming_request_secret" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `monitor_access` (String) Which Incoming Request monitors accept this secret. All Monitors: every one in this project, including monitors created later. Specific Monitors: only the monitors in Monitors. Monitors With Labels: monitors that carry at least one of the labels in Labels. Setting this empties whichever of Monitors and Labels it does not use.
- `name` (String) What the secret is for, so you can tell it apart from the project's other incoming request secrets - the senders that use it, say.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) When Monitor Access is Monitors With Labels, the Incoming Request monitors that carry at least one of these labels accept this secret. Ignored otherwise. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_used_at` (String) About when a request was last accepted with this secret, to the nearest five minutes. Empty means none has been accepted with it yet. Use it to find secrets that are safe to delete.
- `monitors` (Set of String) The Incoming Request monitors that accept this secret when Monitor Access is Specific Monitors. Ignored otherwise. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `secret_value` (String, Sensitive) The value senders send in the X-OneUptime-Secret header or as a Bearer token in the Authorization header, or sign each request's body with (X-OneUptime-Signature, or GitHub's, Slack's or Stripe's signature). Leave it empty and OneUptime generates a strong one. Set your own with 16 to 256 characters: letters, digits and keyboard symbols, no spaces. It cannot be changed after the secret is created - create a new secret instead.
- `updated_at` (String) Date and Time when the object was updated.

---
page_title: "oneuptime_monitor_secret Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor Secret is a secret variable that can be used in monitors. For example you can store auth tokens, passwords, etc. in Monitor Secret and use them in your monitors. Monitor Secret is encrypted and only accessible by the probe.
---

# oneuptime_monitor_secret (Data Source)

Monitor Secret is a secret variable that can be used in monitors. For example you can store auth tokens, passwords, etc. in Monitor Secret and use them in your monitors. Monitor Secret is encrypted and only accessible by the probe.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor secret may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_secret" "example" {
  name = "Example monitor secret"
}

# Or by id:
data "oneuptime_monitor_secret" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `monitor_access` (String) Which monitors can use this secret. All Monitors: every monitor in this project, including monitors created later. Specific Monitors: only the monitors in Monitors. Monitors With Labels: monitors that carry at least one of the labels in Labels. Setting this empties whichever of Monitors and Labels it does not use.
- `name` (String) Any friendly name of this object.
- `secret_type` (String) What the secret holds. Text: a password, API key or token, filled in as written. TOTP: an authenticator app's setup key or otpauth://totp link; a monitor gets the current one-time code, never the key. Changes only together with a new Secret Value, which is read for the new type.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) When Monitor Access is Monitors With Labels, monitors that carry at least one of these labels can use this secret. Ignored otherwise. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `monitors` (Set of String) The monitors that can use this secret when Monitor Access is Specific Monitors. Ignored otherwise. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

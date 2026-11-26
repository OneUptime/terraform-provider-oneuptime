---
page_title: "oneuptime_monitor_secret Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor Secret is a secret variable that can be used in monitors. For example you can store auth tokens, passwords, etc. in Monitor Secret and use them in your monitors. Monitor Secret is encrypted and only accessible by the probe.
---

# oneuptime_monitor_secret (Resource)

Monitor Secret is a secret variable that can be used in monitors. For example you can store auth tokens, passwords, etc. in Monitor Secret and use them in your monitors. Monitor Secret is encrypted and only accessible by the probe.

## Example Usage

```terraform
resource "oneuptime_monitor_secret" "example" {
  name        = "Example monitor secret"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `labels` (Set of String) When Monitor Access is Monitors With Labels, monitors that carry at least one of these labels can use this secret. Ignored otherwise. IDs of `oneuptime_label` resources.
- `monitor_access` (String) Which monitors can use this secret. All Monitors: every monitor in this project, including monitors created later. Specific Monitors: only the monitors in Monitors. Monitors With Labels: monitors that carry at least one of the labels in Labels. Setting this empties whichever of Monitors and Labels it does not use. Allowed values: `All Monitors`, `Specific Monitors`, `Monitors With Labels`. Defaults to `Specific Monitors`.
- `monitors` (Set of String) The monitors that can use this secret when Monitor Access is Specific Monitors. Ignored otherwise. IDs of `oneuptime_monitor` resources.
- `secret_type` (String) What the secret holds. Text: a password, API key or token, filled in as written. TOTP: an authenticator app's setup key or otpauth://totp link; a monitor gets the current one-time code, never the key. Changes only together with a new Secret Value, which is read for the new type. A plan that changes `secret_type` and keeps `secret_value` fails before anything is written: give `secret_value` its new value in the same change. Allowed values: `Text`, `TOTP`. Defaults to `Text`.
- `secret_value` (String, Sensitive) Secret value that you want to store in this object. For a TOTP secret, the setup key or the otpauth://totp link, which is stored as an otpauth://totp link. This value will be encrypted and only accessible by the probe.

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

Import an existing monitor secret by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_secret.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_secret.example <id>
```

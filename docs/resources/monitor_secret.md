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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `secret_value` (String) Secret value that you want to store in this object. This value will be encrypted and only accessible by the probe...
- `monitor_access` (String) Which monitors can use this secret. All Monitors: every monitor in this project, including monitors created later. Specific Monitors: only the monitors in Monitors. Monitors With Labels: monitors that carry at least one of the labels in Labels. Setting this empties whichever of Monitors and Labels it does not use...
- `monitors` (Set) The monitors that can use this secret when Monitor Access is Specific Monitors. Ignored otherwise...
- `labels` (Set) When Monitor Access is Monitors With Labels, monitors that carry at least one of these labels can use this secret. Ignored otherwise...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor_secret.example <id>
```

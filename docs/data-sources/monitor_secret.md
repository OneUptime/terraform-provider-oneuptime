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

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_access` (String) Which monitors can use this secret. All Monitors: every monitor in this project, including monitors created later. Specific Monitors: only the monitors in Monitors. Monitors With Labels: monitors that carry at least one of the labels in Labels. Setting this empties whichever of Monitors and Labels it does not use.
- `name` (String) Any friendly name of this object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) When Monitor Access is Monitors With Labels, monitors that carry at least one of these labels can use this secret. Ignored otherwise. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) The monitors that can use this secret when Monitor Access is Specific Monitors. Ignored otherwise. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

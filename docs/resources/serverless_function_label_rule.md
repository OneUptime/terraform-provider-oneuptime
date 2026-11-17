---
page_title: "oneuptime_serverless_function_label_rule Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Rules for automatically attaching labels to serverless functions when matching functions are created.
---

# oneuptime_serverless_function_label_rule (Resource)

Rules for automatically attaching labels to serverless functions when matching functions are created.

## Example Usage

```terraform
resource "oneuptime_serverless_function_label_rule" "example" {
  name        = "Example serverless function label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this rule.
- `description_regex_pattern` (String) Regex (case-insensitive) matched against the function description. Leave empty to match any description.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the function when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `match_labels` (Set of String) Only trigger for functions that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `name_regex_pattern` (String) Regex (case-insensitive) matched against the function name. Leave empty to match any name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing serverless function label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_serverless_function_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_serverless_function_label_rule.example <id>
```

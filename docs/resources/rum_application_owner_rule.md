---
page_title: "oneuptime_rum_application_owner_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Rules for automatically assigning owners to RUM applications when matching applications are created.
---

# oneuptime_rum_application_owner_rule (Resource)

Rules for automatically assigning owners to RUM applications when matching applications are created.

## Example Usage

```terraform
resource "oneuptime_rum_application_owner_rule" "example" {
  name        = "Example rum application owner rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this rule.
- `description_regex_pattern` (String) Regex (case-insensitive) matched against the application description. Leave empty to match any description.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `match_labels` (Set of String) Only trigger for applications that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `name_regex_pattern` (String) Regex (case-insensitive) matched against the application name. Leave empty to match any name.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `owner_teams` (Set of String) Teams to add as owners on the application when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the application when this rule matches. IDs of `oneuptime_user` records.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing rum application owner rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_rum_application_owner_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_rum_application_owner_rule.example <id>
```

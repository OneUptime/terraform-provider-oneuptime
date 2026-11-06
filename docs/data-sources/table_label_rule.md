---
page_title: "oneuptime_table_label_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Rules for automatically attaching labels to tables when matching resources are created.
---

# oneuptime_table_label_rule (Data Source)

Rules for automatically attaching labels to tables when matching resources are created.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one table label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_table_label_rule" "example" {
  name = "Example table label rule"
}

# Or by id:
data "oneuptime_table_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `database_regex_pattern` (String) Case-insensitive regex matched against the name of the database the table is in.
- `description` (String) Description of this rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this rule.
- `name_regex_pattern` (String) Case-insensitive regex matched against the table's name.
- `schema_regex_pattern` (String) Case-insensitive regex matched against the table's schema.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `labels_to_add` (Set of String) Labels to attach to the resource when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `match_labels` (Set of String) Only tables with one of these labels. Leave empty to skip the filter. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

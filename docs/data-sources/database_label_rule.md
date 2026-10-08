---
page_title: "oneuptime_database_label_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to databases when matching databases are created
---

# oneuptime_database_label_rule (Data Source)

Configure rules for automatically attaching labels to databases when matching databases are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one database label rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_database_label_rule" "example" {
  name = "Example database label rule"
}

# Or by id:
data "oneuptime_database_label_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `database_server_description_pattern` (String) Regex (case-insensitive) matched against the database description. Leave empty to match any description.
- `database_server_name_pattern` (String) Regex (case-insensitive) matched against the database name. Discovered databases are named after their engine (e.g. PostgreSQL orders-db:5432), so ^PostgreSQL matches every PostgreSQL database. Leave empty to match any name.
- `description` (String) Description of this database label rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this database label rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `database_server_labels` (Set of String) Only trigger for databases that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `labels_to_add` (Set of String) Labels to attach to the database when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

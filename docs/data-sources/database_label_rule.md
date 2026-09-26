---
page_title: "oneuptime_database_label_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to databases when matching databases are created
---

# oneuptime_database_label_rule (Data Source)

Configure rules for automatically attaching labels to databases when matching databases are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_database_label_rule" "by_name" {
  name = "example-database_label_rule"
}

data "oneuptime_database_label_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this database label rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `database_server_labels` (Set) Only trigger for databases that already have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `database_server_name_pattern` (String) Regex (case-insensitive) matched against the database name. Discovered databases are named after their engine (e.g. PostgreSQL orders-db:5432), so ^PostgreSQL matches every PostgreSQL database. Leave empty to match any name... Computed.
- `database_server_description_pattern` (String) Regex (case-insensitive) matched against the database description. Leave empty to match any description... Computed.
- `labels_to_add` (Set) Labels to attach to the database when this rule matches. Already-attached labels are not duplicated... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

---
page_title: "oneuptime_database_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to databases when matching databases are created
---

# oneuptime_database_label_rule (Resource)

Configure rules for automatically attaching labels to databases when matching databases are created

## Example Usage

```terraform
resource "oneuptime_database_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this database label rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this database label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `database_server_labels` (Set) Only trigger for databases that already have at least one of these labels. Leave empty to match regardless of labels...
- `database_server_name_pattern` (String) Regex (case-insensitive) matched against the database name. Discovered databases are named after their engine (e.g. PostgreSQL orders-db:5432), so ^PostgreSQL matches every PostgreSQL database. Leave empty to match any name...
- `database_server_description_pattern` (String) Regex (case-insensitive) matched against the database description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the database when this rule matches. Already-attached labels are not duplicated...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_database_label_rule.example <id>
```

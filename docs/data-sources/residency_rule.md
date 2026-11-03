---
page_title: "oneuptime_residency_rule Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Where data may live: the regions a category of data, or data about people in some countries, may be stored and copied to.
---

# oneuptime_residency_rule (Data Source)

Where data may live: the regions a category of data, or data about people in some countries, may be stored and copied to.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one residency rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_residency_rule" "example" {
  name = "Example residency rule"
}

# Or by id:
data "oneuptime_residency_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this residency rule. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the rule is for.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether the rule is checked.
- `name` (String) The rule's name, such as EU personal data stays in the EU.

### Read-Only

- `allowed_regions` (String) The regions the data may live in. A * matches any characters: eu-*, europe-*, westeurope. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `data_categories` (String) The categories of data the rule covers. Empty means every kind of personal data. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this residency rule belongs to. The ID of a `oneuptime_project`.
- `subject_countries` (String) Only data about people in these countries (ISO 3166-1 alpha-2). Empty means data about anyone. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.

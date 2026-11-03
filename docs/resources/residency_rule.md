---
page_title: "oneuptime_residency_rule Resource - oneuptime"
subcategory: "Data Security"
description: |-
  Where data may live: the regions a category of data, or data about people in some countries, may be stored and copied to.
---

# oneuptime_residency_rule (Resource)

Where data may live: the regions a category of data, or data about people in some countries, may be stored and copied to.

## Example Usage

```terraform
resource "oneuptime_residency_rule" "example" {
  name        = "Example residency rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) The rule's name, such as EU personal data stays in the EU.

### Optional

- `allowed_regions` (String) The regions the data may live in. A * matches any characters: eu-*, europe-*, westeurope. A JSON value: write it with `jsonencode()`.
- `data_categories` (String) The categories of data the rule covers. Empty means every kind of personal data. A JSON value: write it with `jsonencode()`.
- `description` (String) What the rule is for.
- `is_enabled` (Boolean) Whether the rule is checked. Defaults to `true`.
- `subject_countries` (String) Only data about people in these countries (ISO 3166-1 alpha-2). Empty means data about anyone. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this residency rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this residency rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing residency rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_residency_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_residency_rule.example <id>
```

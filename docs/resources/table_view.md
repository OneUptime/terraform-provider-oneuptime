---
page_title: "oneuptime_table_view Resource - oneuptime"
subcategory: "Other"
description: |-
  Table View is view settings for a table in a project. It contains columns, filters, and other settings.
---

# oneuptime_table_view (Resource)

Table View is view settings for a table in a project. It contains columns, filters, and other settings.

## Example Usage

```terraform
resource "oneuptime_table_view" "example" {
  name = "Example short text"
  table_id = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `table_id` (String) ID of the table this view is for..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `query` (String) Filters for this table view..
- `sort` (String) Sort for this table view..
- `items_on_page` (Number) Items on page..
- `facets` (String) Facet selections (owner, labels, status, etc.) for this table view..
- `columns` (String) Which columns are shown, and in what order, for this table view..

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
terraform import oneuptime_table_view.example <id>
```

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_scheduled_maintenance_template_user_owner_data Data Source - oneuptime"
subcategory: ""
description: |-
  Scheduled maintenance template user owner data data source
---

# oneuptime_scheduled_maintenance_template_user_owner_data (Data Source)

Scheduled maintenance template user owner data data source

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_template_user_owner_data" "example" {
  name = "example-scheduled_maintenance_template_user_owner_data"
}
```

## Schema

- `id` (String) Identifier to filter by. Optional.
- `name` (String) Name to filter by. Optional.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

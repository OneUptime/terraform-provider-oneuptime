---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_scheduled_maintenance_note_template_data Data Source - oneuptime"
subcategory: ""
description: |-
  Scheduled maintenance note template data data source
---

# oneuptime_scheduled_maintenance_note_template_data (Data Source)

Scheduled maintenance note template data data source

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_note_template_data" "example" {
  name = "example-scheduled_maintenance_note_template_data"
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
- `note` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Note Template], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Note Template], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance Note Template]. Computed.
- `template_name` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Note Template], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Note Template], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance Note Template]. Computed.
- `template_description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Note Template], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Note Template], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance Note Template]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

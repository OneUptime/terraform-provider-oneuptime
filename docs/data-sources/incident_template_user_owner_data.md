---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_incident_template_user_owner_data Data Source - oneuptime"
subcategory: ""
description: |-
  Incident template user owner data data source
---

# oneuptime_incident_template_user_owner_data (Data Source)

Incident template user owner data data source

## Example Usage

```terraform
data "oneuptime_incident_template_user_owner_data" "example" {
  name = "example-incident_template_user_owner_data"
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
- `incident_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Project Member, Read IncidentTemplate User Owner], Update: [No access - you don't have permission for this operation]. Computed.

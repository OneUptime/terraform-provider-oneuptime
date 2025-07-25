---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_incident_state_data Data Source - oneuptime"
subcategory: ""
description: |-
  Incident state data data source
---

# oneuptime_incident_state_data (Data Source)

Incident state data data source

## Example Usage

```terraform
data "oneuptime_incident_state_data" "example" {
  name = "example-incident_state_data"
}
```

## Schema

- `id` (String) Identifier to filter by. Optional.
- `name` (String) Name to filter by. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident State], Update: [No access - you don't have permission for this operation]. Computed.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State], Read: [Project Owner, Project Admin, Project Member, Read Incident State], Update: [Project Owner, Project Admin, Project Member, Edit Incident State]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `color` (String) Color object. Computed.
- `is_created_state` (Bool) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State], Read: [Project Owner, Project Admin, Project Member, Read Incident State], Update: [Project Owner, Project Admin, Project Member, Edit Incident State]. Computed.
- `is_acknowledged_state` (Bool) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State], Read: [Project Owner, Project Admin, Project Member, Read Incident State], Update: [Project Owner, Project Admin, Project Member, Edit Incident State]. Computed.
- `is_resolved_state` (Bool) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State], Read: [Project Owner, Project Admin, Project Member, Read Incident State], Update: [Project Owner, Project Admin, Project Member, Edit Incident State]. Computed.
- `order` (Number) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident State], Read: [Project Owner, Project Admin, Project Member, Read Incident State], Update: [Project Owner, Project Admin, Project Member, Edit Incident State]. Computed.

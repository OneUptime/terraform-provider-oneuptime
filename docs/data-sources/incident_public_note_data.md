---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_incident_public_note_data Data Source - oneuptime"
subcategory: ""
description: |-
  Incident public note data data source
---

# oneuptime_incident_public_note_data (Data Source)

Incident public note data data source

## Example Usage

```terraform
data "oneuptime_incident_public_note_data" "example" {
  name = "example-incident_public_note_data"
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
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `note` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [Project Owner, Project Admin, Project Member, Edit Incident Status Page Note]. Computed.
- `is_status_page_subscribers_notified_on_note_created` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]. Computed.
- `should_status_page_subscribers_be_notified_on_note_created` (Bool) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident Status Page Note], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]. Computed.
- `is_owner_notified` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident Status Page Note], Update: [No access - you don't have permission for this operation]. Computed.
- `posted_at` (String) A date time object.. Computed.

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_exception_data Data Source - oneuptime"
subcategory: ""
description: |-
  Exception data data source
---

# oneuptime_exception_data (Data Source)

Exception data data source

## Example Usage

```terraform
data "oneuptime_exception_data" "example" {
  name = "example-exception_data"
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
- `telemetry_service_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `message` (String) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.
- `stack_trace` (String) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.
- `exception_type` (String) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.
- `fingerprint` (String) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `marked_as_resolved_at` (String) A date time object.. Computed.
- `marked_as_archived_at` (String) A date time object.. Computed.
- `first_seen_at` (String) A date time object.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `assign_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `assign_to_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `marked_as_resolved_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `marked_as_archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_resolved` (Bool) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.
- `is_archived` (Bool) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.
- `occurance_count` (Number) Permissions - Create: [Project Owner, Project Admin, Create Telemetry Service Exception], Read: [Project Owner, Project Admin, Project Member, Read Telemetry Service Exception], Update: [Project Owner, Project Admin, Edit Telemetry Service Exception]. Computed.

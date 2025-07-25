---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_telemetry_service_data Data Source - oneuptime"
subcategory: ""
description: |-
  Telemetry service data data source
---

# oneuptime_telemetry_service_data (Data Source)

Telemetry service data data source

## Example Usage

```terraform
data "oneuptime_telemetry_service_data" "example" {
  name = "example-telemetry_service_data"
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
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Telemetry Service], Update: [No access - you don't have permission for this operation]. Computed.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Telemetry Service], Update: [Project Owner, Project Admin, Project Member, Edit Telemetry Service]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (List) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Telemetry Service], Update: [Project Owner, Project Admin, Project Member, Edit Telemetry Service]. Computed.
- `telemetry_service_token` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `retain_telemetry_data_for_days` (Number) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Telemetry Service], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Telemetry Service], Update: [Project Owner, Project Admin, Project Member, Edit Telemetry Service]. Computed.
- `service_color` (String) Color object. Computed.

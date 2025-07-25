---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_alert_severity_data Data Source - oneuptime"
subcategory: ""
description: |-
  Alert severity data data source
---

# oneuptime_alert_severity_data (Data Source)

Alert severity data data source

## Example Usage

```terraform
data "oneuptime_alert_severity_data" "example" {
  name = "example-alert_severity_data"
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
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Alert Severity], Update: [No access - you don't have permission for this operation]. Computed.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert Severity], Read: [Project Owner, Project Admin, Project Member, Read Alert Severity], Update: [Project Owner, Project Admin, Project Member, Edit Alert Severity]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `color` (String) Color object. Computed.
- `order` (Number) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Alert Severity], Read: [Project Owner, Project Admin, Project Member, Read Alert Severity], Update: [Project Owner, Project Admin, Project Member, Edit Alert Severity]. Computed.

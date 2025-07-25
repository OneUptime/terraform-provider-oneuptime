---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_probe_data Data Source - oneuptime"
subcategory: ""
description: |-
  Probe data data source
---

# oneuptime_probe_data (Data Source)

Probe data data source

## Example Usage

```terraform
data "oneuptime_probe_data" "example" {
  name = "example-probe_data"
}
```

## Schema

- `id` (String) Identifier to filter by. Optional.
- `name` (String) Name to filter by. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `key` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Probe], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Edit Probe]. Computed.
- `description` (String) Name object. Computed.
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Public], Update: [No access - you don't have permission for this operation]. Computed.
- `probe_version` (String) Version object. Computed.
- `last_alive` (String) A date time object.. Computed.
- `icon_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `should_auto_enable_probe_on_new_monitors` (Bool) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]. Computed.
- `connection_status` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]. Computed.
- `labels` (List) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [Project Owner, Project Admin, Project Member, Edit Status Page]. Computed.

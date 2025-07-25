---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_code_repository_data Data Source - oneuptime"
subcategory: ""
description: |-
  Code repository data data source
---

# oneuptime_code_repository_data (Data Source)

Code repository data data source

## Example Usage

```terraform
data "oneuptime_code_repository_data" "example" {
  name = "example-code_repository_data"
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
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [No access - you don't have permission for this operation]. Computed.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (List) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]. Computed.
- `secret_token` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `main_branch_name` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]. Computed.
- `repository_hosted_at` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]. Computed.
- `organization_name` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]. Computed.
- `repository_name` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Code Repository], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Code Repository], Update: [Project Owner, Project Admin, Project Member, Edit Code Repository]. Computed.
- `last_copilot_run_date_time` (String) A date time object.. Computed.

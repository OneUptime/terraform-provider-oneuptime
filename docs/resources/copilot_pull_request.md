---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_copilot_pull_request Resource - oneuptime"
subcategory: ""
description: |-
  Copilot pull request resource
---

# oneuptime_copilot_pull_request (Resource)

Copilot pull request resource

## Example Usage

```terraform
resource "oneuptime_copilot_pull_request" "example" {

}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `code_repository_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `service_catalog_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `service_repository_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `pull_request_id` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]. Computed.
- `copilot_pull_request_status` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]. Computed.
- `is_setup_pull_request` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member], Update: [No access - you don't have permission for this operation]. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_copilot_pull_request.example <id>
```

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_service_in_service_catalog Resource - oneuptime"
subcategory: ""
description: |-
  Service in service catalog resource
---

# oneuptime_service_in_service_catalog (Resource)

Service in service catalog resource

## Example Usage

```terraform
resource "oneuptime_service_in_service_catalog" "example" {
  name = "example-resource"
  description = "Example resource"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `name` (String) Name. Required.
- `description` (String) Description. Optional.
- `labels` (List) Labels. Optional.
- `service_color` (Map) Color object. Optional.
- `service_language` (String) Service Language. Optional.
- `tech_stack` (Map) Tech Stack. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Project Member, Read Service Catalog], Update: [No access - you don't have permission for this operation]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_service_in_service_catalog.example <id>
```

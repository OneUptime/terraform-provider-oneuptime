---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_status_page_domain Resource - oneuptime"
subcategory: ""
description: |-
  Status page domain resource
---

# oneuptime_status_page_domain (Resource)

Status page domain resource

## Example Usage

```terraform
resource "oneuptime_status_page_domain" "example" {
  domain_id = "123e4567-e89b-12d3-a456-426614174000"
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  subdomain = "example-subdomain"
  full_domain = "example-full_domain"
  cname_verification_token = "example-cname_verification_token"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `domain_id` (String) A unique identifier for an object, represented as a UUID.. Required.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Required.
- `subdomain` (String) Sumdomain. Required.
- `full_domain` (String) Full Domain. Required.
- `cname_verification_token` (String) CNAME Verification Token. Required.
- `custom_certificate` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]. Optional.
- `custom_certificate_key` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]. Optional.
- `is_custom_certificate` (Bool) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_cname_verified` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]. Computed.
- `is_ssl_ordered` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]. Computed.
- `is_ssl_provisioned` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_domain.example <id>
```

---
page_title: "oneuptime_status_page_sso Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure Status Page SSO
---

# oneuptime_status_page_sso (Data Source)

Configure Status Page SSO

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page sso may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_sso" "example" {
  name = "Example status page sso"
}

# Or by id:
data "oneuptime_status_page_sso" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `digest_method` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `is_tested` (Boolean) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [No access - you don't have permission for this operation]
- `issuer_url` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `name` (String) Any friendly name of this object.
- `public_certificate` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `sign_on_url` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `signature_method` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page SSO], Read: [Project Owner, Project Admin, Read Status Page SSO], Update: [Project Owner, Project Admin, Edit Status Page SSO]
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

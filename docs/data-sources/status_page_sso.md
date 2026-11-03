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
- `description` (String) What this sign-in is for, such as which identity provider it uses.
- `digest_method` (String) The digest algorithm your identity provider's SAML signatures use.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) When on, visitors of a private status page can sign in with this provider.
- `is_tested` (Boolean) Whether a sign-in with this provider has been tested. It does not turn the provider on.
- `issuer_url` (String) Your identity provider's issuer (entity ID), as its SAML responses name it.
- `name` (String) Any friendly name of this object.
- `public_certificate` (String) Your identity provider's signing certificate, in PEM, to check its SAML responses with.
- `sign_on_url` (String) Your identity provider's single sign-on URL, where visitors are sent to sign in.
- `signature_method` (String) The algorithm your identity provider signs SAML responses with.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

---
page_title: "oneuptime_status_page_sso Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure Status Page SSO
---

# oneuptime_status_page_sso (Resource)

Configure Status Page SSO

## Example Usage

```terraform
resource "oneuptime_status_page_sso" "example" {
  status_page_id     = oneuptime_status_page.example.id
  name               = "Example status page sso"
  description        = "Managed by Terraform"
  signature_method   = "RSA-SHA1"
  digest_method      = "SHA1"
  sign_on_url        = "https://www.example.com/path/to/resource?param=value"
  issuer_url         = "https://idp.example.com/saml/metadata"
  public_certificate = "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJAKL0UG+mRKU...\n-----END CERTIFICATE-----"
}
```

## Schema

### Required

- `description` (String) What this sign-in is for, such as which identity provider it uses.
- `digest_method` (String) The digest algorithm your identity provider's SAML signatures use. Allowed values: `SHA1`, `SHA256`, `SHA384`, `SHA512`.
- `issuer_url` (String) Your identity provider's issuer (entity ID), as its SAML responses name it.
- `name` (String) Any friendly name of this object.
- `public_certificate` (String) Your identity provider's signing certificate, in PEM, to check its SAML responses with.
- `sign_on_url` (String) Your identity provider's single sign-on URL, where visitors are sent to sign in.
- `signature_method` (String) The algorithm your identity provider signs SAML responses with. Allowed values: `RSA-SHA1`, `RSA-SHA256`, `RSA-SHA384`, `RSA-SHA512`.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `is_enabled` (Boolean) When on, visitors of a private status page can sign in with this provider. Defaults to `false`.
- `is_tested` (Boolean) Whether a sign-in with this provider has been tested. It does not turn the provider on. Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page sso by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_sso.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_sso.example <id>
```

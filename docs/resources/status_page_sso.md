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
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
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

---
page_title: "oneuptime_smtp_config Resource - oneuptime"
subcategory: "Notification Channels"
description: |-
  Manage Custom SMTP Servers for your project
---

# oneuptime_smtp_config (Resource)

Manage Custom SMTP Servers for your project

## Example Usage

```terraform
resource "oneuptime_smtp_config" "example" {
  name        = "Example smtp config"
  from_email  = "user@example.com"
  from_name   = "OneUptime Notifications"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `from_email` (String)
- `from_name` (String)
- `name` (String)

### Optional

- `amazon_ses_access_key_id` (String) The access key ID of an IAM user or role allowed to send with Amazon SES (ses:SendEmail). Required when Transport is Amazon SES.
- `amazon_ses_region` (String) The AWS Region your sender identity is verified in, such as us-east-1. Required when Transport is Amazon SES.
- `amazon_ses_secret_access_key` (String, Sensitive) The secret access key that goes with the access key ID. Required when Transport is Amazon SES. Encrypted at rest.
- `auth_type` (String) The type of authentication to use for this SMTP server. Options: Username and Password, OAuth, or None. Allowed values: `Username and Password`, `OAuth`, `None`. Defaults to `Username and Password`.
- `client_id` (String) The Client ID from your OAuth application registration. Required for OAuth authentication.
- `client_secret` (String, Sensitive) The Client Secret from your OAuth application registration. Required for OAuth authentication. For Google service accounts, this is the private key from the JSON key file.
- `description` (String) Friendly description that will help you remember.
- `hostname` (String) SMTP server hostname. Required when Transport is SMTP. Not used by HTTP-API transports like Microsoft Graph.
- `mailgun_api_key` (String, Sensitive) A Mailgun API key, or a sending API key for the domain. Required when Transport is Mailgun. Encrypted at rest.
- `mailgun_domain` (String) The Mailgun sending domain mail goes out from, such as mg.example.com. Required when Transport is Mailgun.
- `mailgun_region` (String) The region the Mailgun domain was created in: US or EU. US when not set. Allowed values: `US`, `EU`.
- `oauth_provider_type` (String) The OAuth grant type to use. 'Client Credentials' for Microsoft 365 and most providers. 'JWT Bearer' for Google Workspace service accounts. Allowed values: `Client Credentials`, `JWT Bearer`.
- `password` (String, Sensitive)
- `port` (Number) SMTP server port. Required when Transport is SMTP. Not used by HTTP-API transports.
- `postmark_server_token` (String, Sensitive) The API token of the Postmark server that sends, from its API Tokens tab. Required when Transport is Postmark. Encrypted at rest.
- `resend_api_key` (String, Sensitive) A Resend API key with sending access. Required when Transport is Resend. Encrypted at rest.
- `scope` (String) The OAuth scope(s) required for SMTP access. For Microsoft 365: https://outlook.office365.com/.default. For Google: https://mail.google.com/.
- `secure` (Boolean) Enable secure SMTP connection. Used only for SMTP transport — HTTP-API transports always use HTTPS. Defaults to `true`.
- `sendgrid_api_key` (String, Sensitive) A SendGrid API key with the Mail Send permission. Required when Transport is SendGrid. Encrypted at rest.
- `token_url` (String) The OAuth token endpoint URL. For Microsoft 365: https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token. For Google: https://oauth2.googleapis.com/token.
- `transport_type` (String) How OneUptime delivers mail for this config. 'SMTP' uses the hostname and port. 'Microsoft Graph' sends through the Microsoft Graph REST API, for a Microsoft 365 tenant with SMTP AUTH turned off. 'Amazon SES', 'SendGrid', 'Mailgun', 'Postmark' and 'Resend' send through that provider's HTTPS API with its own key, set in the columns named for it. Allowed values: `SMTP`, `Microsoft Graph`, `Amazon SES`, `SendGrid`, `Mailgun`, `Postmark`, `Resend`. Defaults to `SMTP`.
- `username` (String)

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing smtp config by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_smtp_config.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_smtp_config.example <id>
```

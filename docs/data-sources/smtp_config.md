---
page_title: "oneuptime_smtp_config Data Source - oneuptime"
subcategory: "Notification Channels"
description: |-
  Manage Custom SMTP Servers for your project
---

# oneuptime_smtp_config (Data Source)

Manage Custom SMTP Servers for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one smtp config may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_smtp_config" "example" {
  name = "Example smtp config"
}

# Or by id:
data "oneuptime_smtp_config" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `amazon_ses_access_key_id` (String) The access key ID of an IAM user or role allowed to send with Amazon SES (ses:SendEmail). Required when Transport is Amazon SES.
- `amazon_ses_authentication_type` (String) How OneUptime signs in to Amazon SES: 'Access Key' with the access key ID and secret access key, or 'IAM Role' with an IAM role in your own AWS account (Amazon SES Role ARN), which OneUptime signs in to with this project's external ID. Read only when Transport is Amazon SES.
- `amazon_ses_region` (String) The AWS Region your sender identity is verified in, such as us-east-1. Required when Transport is Amazon SES.
- `amazon_ses_role_arn` (String) The ARN of the IAM role in your AWS account that OneUptime signs in to, when Amazon SES Authentication Type is 'IAM Role'. Its trust policy must let OneUptime's AWS account in with this project's external ID, and the role must be allowed ses:SendEmail.
- `auth_type` (String) The type of authentication to use for this SMTP server. Options: Username and Password, OAuth, or None.
- `client_id` (String) The Client ID from your OAuth application registration. Required for OAuth authentication.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `from_name` (String) Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]
- `hostname` (String) SMTP server hostname. Required when Transport is SMTP. Not used by HTTP-API transports like Microsoft Graph.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `mailgun_domain` (String) The Mailgun sending domain mail goes out from, such as mg.example.com. Required when Transport is Mailgun.
- `mailgun_region` (String) The region the Mailgun domain was created in: US or EU. US when not set.
- `name` (String) Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]
- `oauth_provider_type` (String) The OAuth grant type to use. 'Client Credentials' for Microsoft 365 and most providers. 'JWT Bearer' for Google Workspace service accounts.
- `port` (Number) SMTP server port. Required when Transport is SMTP. Not used by HTTP-API transports.
- `scope` (String) The OAuth scope(s) required for SMTP access. For Microsoft 365: https://outlook.office365.com/.default. For Google: https://mail.google.com/.
- `secure` (Boolean) Enable secure SMTP connection. Used only for SMTP transport — HTTP-API transports always use HTTPS.
- `slug` (String) Friendly globally unique name for your object.
- `token_url` (String) The OAuth token endpoint URL. For Microsoft 365: https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token. For Google: https://oauth2.googleapis.com/token.
- `transport_type` (String) How OneUptime delivers mail for this config. 'SMTP' uses the hostname and port. 'Microsoft Graph' sends through the Microsoft Graph REST API, for a Microsoft 365 tenant with SMTP AUTH turned off. 'Amazon SES', 'SendGrid', 'Mailgun', 'Postmark' and 'Resend' send through that provider's HTTPS API with its own key, set in the columns named for it.
- `username` (String) Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `from_email` (String) Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

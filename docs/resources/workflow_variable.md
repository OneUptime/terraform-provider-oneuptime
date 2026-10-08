---
page_title: "oneuptime_workflow_variable Resource - oneuptime"
subcategory: "Workflows"
description: |-
  Store environment variables or secrets for your workflows.
---

# oneuptime_workflow_variable (Resource)

Store environment variables or secrets for your workflows.

## Example Usage

```terraform
resource "oneuptime_workflow_variable" "example" {
  name        = "Example workflow variable"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Variable Name.

### Optional

- `content` (String) Content of the variable. Required for Static variables. Not used by OAuth 2.0 variables, whose value is the access token OneUptime fetches.
- `description` (String) Friendly description that will help you remember.
- `is_secret` (Boolean) Is this variable a secret. If true, then it'll not be in the logs. Defaults to `false`.
- `oauth_additional_parameters` (String) OAuth 2.0 variables only. Extra form parameters sent with every token request, such as audience for Auth0 or resource for Azure AD v1. Readable by anyone who can read the variable, so do not put secrets here. A JSON value: write it with `jsonencode()`.
- `oauth_client_authentication_method` (String) OAuth 2.0 variables only. How the client ID and secret are sent: in an HTTP Basic header (client_secret_basic, the default) or in the request body (client_secret_post).
- `oauth_client_id` (String) OAuth 2.0 variables only. The client ID of the application registered with your identity provider.
- `oauth_client_secret` (String) OAuth 2.0 variables only. The client secret of the application. Required for the Client Credentials grant; optional for the Refresh Token grant (public clients have none). Encrypted, and never readable through the API.
- `oauth_grant_type` (String) OAuth 2.0 variables only. Client Credentials for machine-to-machine access, or Refresh Token to keep delegated access alive with a refresh token you obtained once.
- `oauth_refresh_token` (String) OAuth 2.0 variables using the Refresh Token grant only. OneUptime exchanges it for access tokens and stores the replacement when your identity provider rotates it. Encrypted, and never readable through the API.
- `oauth_scope` (String) OAuth 2.0 variables only. Space-separated scopes to request. Leave empty to use the scopes your identity provider grants by default.
- `oauth_token_url` (String) OAuth 2.0 variables only. The token endpoint of your identity provider.
- `variable_type` (String) Static: the content you save is used as is. OAuth 2.0: OneUptime fetches an access token from your identity provider and refreshes it automatically when a workflow uses it after it has expired. Defaults to `Static`.
- `workflow_id` (String) ID of Workflow this variable belong to. If this is null then this variable will be a global variable. The ID of a `oneuptime_workflow`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `oauth_access_token_expires_at` (String) When the cached access token expires, as reported by the identity provider (expires_in) or by the token itself (the JWT exp claim). Empty when neither says.
- `oauth_last_refresh_error` (String) Why the last attempt to fetch an access token failed. Cleared by the next successful refresh.
- `oauth_last_refresh_error_at` (String) When the last failed attempt to fetch an access token happened.
- `oauth_last_refreshed_at` (String) When OneUptime last fetched an access token for this variable. Cleared when the OAuth settings change.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing workflow variable by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_workflow_variable.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_workflow_variable.example <id>
```

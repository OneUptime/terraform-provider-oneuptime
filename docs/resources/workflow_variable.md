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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Variable Name..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `workflow_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `content` (String) Content of the variable. Required for Static variables. Not used by OAuth 2.0 variables, whose value is the access token OneUptime fetches...
- `is_secret` (Bool) Is this variable a secret. If true, then it'll not be in the logs..
- `variable_type` (String) Static: the content you save is used as is. OAuth 2.0: OneUptime fetches an access token from your identity provider and refreshes it automatically when a workflow uses it after it has expired...
- `oauth_grant_type` (String) OAuth 2.0 variables only. Client Credentials for machine-to-machine access, or Refresh Token to keep delegated access alive with a refresh token you obtained once...
- `oauth_token_url` (String) OAuth 2.0 variables only. The token endpoint of your identity provider...
- `oauth_client_id` (String) OAuth 2.0 variables only. The client ID of the application registered with your identity provider...
- `oauth_client_secret` (String) OAuth 2.0 variables only. The client secret of the application. Required for the Client Credentials grant; optional for the Refresh Token grant (public clients have none). Encrypted, and never readable through the API...
- `oauth_refresh_token` (String) OAuth 2.0 variables using the Refresh Token grant only. OneUptime exchanges it for access tokens and stores the replacement when your identity provider rotates it. Encrypted, and never readable through the API...
- `oauth_scope` (String) OAuth 2.0 variables only. Space-separated scopes to request. Leave empty to use the scopes your identity provider grants by default...
- `oauth_additional_parameters` (String) OAuth 2.0 variables only. Extra form parameters sent with every token request, such as audience for Auth0 or resource for Azure AD v1. Readable by anyone who can read the variable, so do not put secrets here...
- `oauth_client_authentication_method` (String) OAuth 2.0 variables only. How the client ID and secret are sent: in an HTTP Basic header (client_secret_basic, the default) or in the request body (client_secret_post)...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `oauth_access_token_expires_at` (String) A date time object..
- `oauth_last_refreshed_at` (String) A date time object..
- `oauth_last_refresh_error` (String) Why the last attempt to fetch an access token failed. Cleared by the next successful refresh...
- `oauth_last_refresh_error_at` (String) A date time object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_workflow_variable.example <id>
```

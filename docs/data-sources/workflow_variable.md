---
page_title: "oneuptime_workflow_variable Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Store environment variables or secrets for your workflows.
---

# oneuptime_workflow_variable (Data Source)

Store environment variables or secrets for your workflows. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workflow_variable" "by_name" {
  name = "example-workflow_variable"
}

data "oneuptime_workflow_variable" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `workflow_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `is_secret` (Bool) Is this variable a secret. If true, then it'll not be in the logs.. Computed.
- `variable_type` (String) Static: the content you save is used as is. OAuth 2.0: OneUptime fetches an access token from your identity provider and refreshes it automatically when a workflow uses it after it has expired... Computed.
- `oauth_grant_type` (String) OAuth 2.0 variables only. Client Credentials for machine-to-machine access, or Refresh Token to keep delegated access alive with a refresh token you obtained once... Computed.
- `oauth_token_url` (String) OAuth 2.0 variables only. The token endpoint of your identity provider... Computed.
- `oauth_client_id` (String) OAuth 2.0 variables only. The client ID of the application registered with your identity provider... Computed.
- `oauth_scope` (String) OAuth 2.0 variables only. Space-separated scopes to request. Leave empty to use the scopes your identity provider grants by default... Computed.
- `oauth_additional_parameters` (String) OAuth 2.0 variables only. Extra form parameters sent with every token request, such as audience for Auth0 or resource for Azure AD v1. Readable by anyone who can read the variable, so do not put secrets here... Computed.
- `oauth_client_authentication_method` (String) OAuth 2.0 variables only. How the client ID and secret are sent: in an HTTP Basic header (client_secret_basic, the default) or in the request body (client_secret_post)... Computed.
- `oauth_access_token_expires_at` (String) A date time object.. Computed.
- `oauth_last_refreshed_at` (String) A date time object.. Computed.
- `oauth_last_refresh_error` (String) Why the last attempt to fetch an access token failed. Cleared by the next successful refresh... Computed.
- `oauth_last_refresh_error_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

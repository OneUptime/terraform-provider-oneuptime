---
page_title: "oneuptime_external_mcp_server Data Source - oneuptime"
subcategory: "Other"
description: |-
  MCP servers outside OneUptime - such as AWS, Azure or GitHub - whose tools OneUptime AI may use while it investigates and fixes incidents and alerts.
---

# oneuptime_external_mcp_server (Data Source)

MCP servers outside OneUptime - such as AWS, Azure or GitHub - whose tools OneUptime AI may use while it investigates and fixes incidents and alerts.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one external mcp server may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_external_mcp_server" "example" {
  name = "Example external mcp server"
}

# Or by id:
data "oneuptime_external_mcp_server" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `auth_type` (String) How OneUptime signs in to the server: None, BearerToken (an API token), Headers (headers of the server's own, such as X-API-Key) or OAuthClientCredentials (a token from an OAuth token URL, with a client ID and secret).
- `created_by_user_id` (String) ID of the user who added this MCP server. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the server is for, such as 'Our production AWS account, read-only'. OneUptime AI reads it to choose the right server.
- `header_names` (String) The names of the headers the server signs in with, without their values. Worked out from the headers on every read.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether OneUptime AI may use this server. A server that is off is never used, whatever a rule names.
- `last_error` (String) Why the last connection failed, with secrets removed. Cleared by the next connection that works.
- `name` (String) A name for the server, such as AWS or Azure. OneUptime AI and the people reading what it did see the server by this name.
- `oauth_audience` (String) For OAuthClientCredentials: the audience to ask for, when the token URL needs one (Auth0 does). Leave empty otherwise.
- `oauth_client_id` (String) For OAuthClientCredentials: the OAuth client ID.
- `oauth_scope` (String) For OAuthClientCredentials: the scope to ask for, such as api://<application-id>/.default. Leave empty when the token URL needs none.
- `oauth_token_url` (String) For OAuthClientCredentials: the OAuth token endpoint, such as https://login.microsoftonline.com/<tenant>/oauth2/v2.0/token.
- `protocol_version` (String) The MCP protocol version the server agreed on when OneUptime last connected.
- `server_name` (String) What the server called itself when OneUptime last connected.
- `server_version` (String) The version the server gave when OneUptime last connected, if any.
- `transport` (String) How the server talks, as OneUptime found out when it last connected: StreamableHttp or Sse.
- `url` (String) The server's MCP address, such as https://mcp.example.com/mcp. Streamable HTTP and SSE servers both work: OneUptime finds out which when it connects. Never put a token in it: use Authentication.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_checked_at` (String) When OneUptime last tried to connect to the server. Empty means never.
- `last_connected_at` (String) When OneUptime last connected to the server and read its tools.
- `project_id` (String) ID of the project this MCP server belongs to. The ID of a `oneuptime_project`.
- `tools` (String) The tools the server offered when OneUptime last connected: each one's name, description and whether the server marks it read-only. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

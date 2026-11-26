---
page_title: "oneuptime_security_case_pull_request Data Source - oneuptime"
subcategory: "Security"
description: |-
  A pull request a person linked to a security case as its fix, which OneUptime follows to its deploy.
---

# oneuptime_security_case_pull_request (Data Source)

A pull request a person linked to a security case as its fix, which OneUptime follows to its deploy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case pull request may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_pull_request" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_security_case_pull_request" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `code_repository_id` (String) ID of the repository the pull request is in. Found by OneUptime from the address. The ID of a `oneuptime_code_repository`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who linked the pull request to the case. The ID of a `oneuptime_user` (see the data source).
- `deployed_commit_sha` (String) The commit that deploy ran: the merge commit or a later one. Written by OneUptime.
- `deployed_version` (String) The version that deploy shipped, when it named one. Written by OneUptime.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `merge_commit_sha` (String) The commit the pull request merged as: what a deploy must contain to ship it. Written by OneUptime.
- `pull_request_number` (Number) The pull request's number in its repository. Written by OneUptime.
- `pull_request_state` (String) Whether the pull request is open, merged or closed, as OneUptime last read it from its host.
- `pull_request_url` (String) The pull request's page, on a repository the project connected under Code: GitHub, GitLab, Forgejo or Gitea.
- `repo_name` (String) Name of the repository. Written by OneUptime.
- `repo_organization_name` (String) Organization, user or group that owns the repository. Written by OneUptime.
- `security_case_id` (String) ID of the case this pull request fixes. The ID of a `oneuptime_security_case`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `deployed_at` (String) When a recorded deploy that contains the pull request's merge commit started. Written by OneUptime.
- `followed_at` (String) When OneUptime last read the pull request's state or looked for its deploy.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `merged_at` (String) When the pull request merged. Written by OneUptime.
- `project_id` (String) ID of the project this pull request belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

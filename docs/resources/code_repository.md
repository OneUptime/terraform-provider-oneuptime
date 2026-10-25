---
page_title: "oneuptime_code_repository Resource - oneuptime"
subcategory: "Reliability Copilot"
description: |-
  Connect and manage code repositories from GitHub, GitLab, Forgejo, Gitea and any other Git server
---

# oneuptime_code_repository (Resource)

Connect and manage code repositories from GitHub, GitLab, Forgejo, Gitea and any other Git server

## Example Usage

```terraform
resource "oneuptime_code_repository" "example" {
  name                 = "Example code repository"
  repository_hosted_at = "Example short text"
  organization_name    = "Example short text"
  repository_name      = "Example short text"
  description          = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) A friendly name for this code repository.
- `organization_name` (String) The organization, group or user that owns this repository, as its host names it (a GitLab group's full path, such as my-group/backend).
- `repository_hosted_at` (String) Where this repository is hosted, and so how OneUptime reaches it: GitHub (the GitHub App), GitLab, Forgejo (Forgejo, Gitea and Codeberg) or Git (any Git server over HTTPS).
- `repository_name` (String) The name of the repository.

### Optional

- `build_command` (String) Optional. Leave it empty: the code agent learns how to build the repository from the repository itself (its AGENTS.md first) and checks its change that way. When set, a Runner of the project's own also runs it at the repository root to check that an AI-authored fix builds (e.g. 'npm run build'); a failure is fed back to the code agent for bounded repair attempts before the pull request opens.
- `description` (String) A description of this code repository.
- `git_hub_trigger_label` (String) The issue label that hands an issue to the OneUptime GitHub App. Adding this label to an issue starts the same work an '@mention implement this' would, which is how you assign work to the app from the GitHub UI. Unset means 'oneuptime'.
- `is_git_hub_commands_enabled` (Boolean) Whether the OneUptime GitHub App acts on mentions, assignments and trigger labels in this repository. Only people with write access to the repository can command it, and it never merges anything. On by default. Defaults to `true`.
- `is_production_review_enabled` (Boolean) Whether OneUptime comments on this repository's pull requests with what production says about the code they change: the exceptions thrown in it, how busy and slow the operations it runs are, and the incidents and deploys around it. Only pull requests from people who can write to the repository are reviewed on their own.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `main_branch_name` (String) The name of the main/default branch.
- `max_open_fix_pull_requests` (Number) Maximum AI-authored fix pull requests that may be open on this repository at the same time. At the cap, new AI fix runs are refused a repository token, so they cannot push branches or open pull requests. Unset means no cap; 0 blocks AI fix pull requests for this repository entirely.
- `production_review_pull_request_content` (String) What a production review posts on a pull request. Review (the default): the review, or only a link to it in OneUptime on a public repository. Link: only a link, whatever the repository. ReviewEvenIfPublic: the review, even on a public repository.
- `repository_url` (String) The repository's address: its page on GitHub, GitLab or Forgejo, or the HTTPS address a Git repository is cloned with.
- `runner_id` (String) ID of the Runner AI works in this repository on: one of the project's own Runners. Empty means a global Runner, one OneUptime runs for every project. The ID of a `oneuptime_runner`.
- `secret_token` (String) Secret token used to verify incoming webhooks.
- `setup_command` (String) Optional. Leave it empty: the code agent learns how to set up, build and test the repository from the repository itself (its AGENTS.md first) and checks its change that way. When set, a Runner of the project's own also runs it at the repository root, in the cloned workspace, to install dependencies before the build and test commands check an AI-authored fix (e.g. 'npm ci').
- `test_command` (String) Optional. Leave it empty: the code agent learns how to test the repository from the repository itself (its AGENTS.md first) and checks its change that way. When set, a Runner of the project's own also runs it at the repository root to run the test suite against an AI-authored fix (e.g. 'npm test'); a failure is fed back to the code agent for bounded repair attempts before the pull request opens.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `git_hub_app_installation_id` (String) The GitHub App installation ID used to authenticate with this repository.
- `git_lab_project_id` (String) The GitLab project ID OneUptime reads this repository by. Set when the repository is connected.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing code repository by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_code_repository.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_code_repository.example <id>
```

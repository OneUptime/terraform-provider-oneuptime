---
page_title: "oneuptime_cloud_resource_owner_rule Resource - oneuptime"
subcategory: "Cloud"
description: |-
  Rules for automatically assigning owners to cloud resources (Cloud) or cloud environments (App Platforms) when matching ones are created.
---

# oneuptime_cloud_resource_owner_rule (Resource)

Rules for automatically assigning owners to cloud resources (Cloud) or cloud environments (App Platforms) when matching ones are created.

## Example Usage

```terraform
resource "oneuptime_cloud_resource_owner_rule" "example" {
  name        = "Example cloud resource owner rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this rule.

### Optional

- `cloud_resource_kind` (String) environment: App Platforms' environments. resource: the resources of Cloud's accounts. Empty: both, as rules made before the two products split. Allowed values: `environment`, `resource`.
- `cloud_tag_pattern` (String) Regex (case-insensitive) matched against each of the resource's cloud tags, written key=value: ^team=payments$ matches every resource tagged team=payments, ^env= every resource with an env tag. Leave empty to match any tags.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this rule.
- `description_regex_pattern` (String) Regex (case-insensitive) matched against the resource description. Leave empty to match any description.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `match_labels` (Set of String) Only trigger for resources that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `name_regex_pattern` (String) Regex (case-insensitive) matched against the resource name. Leave empty to match any name.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `owner_teams` (Set of String) Teams to add as owners on the resource when this rule matches. IDs of `oneuptime_team` resources.
- `owner_user_emails` (Set of String) The email addresses of the people in `owner_users`, instead of their IDs: each one of someone in this project, or invited to it. Give one of the two.
- `owner_users` (Set of String) Users to add as owners on the resource when this rule matches. IDs of `oneuptime_user` records. Or name them by email address, in `owner_user_emails`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing cloud resource owner rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_resource_owner_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_resource_owner_rule.example <id>
```

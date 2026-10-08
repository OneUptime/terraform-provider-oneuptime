---
page_title: "oneuptime_ai_agent Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.
---

# oneuptime_ai_agent (Data Source)

Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai agent may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_agent" "example" {
  name = "Example ai agent"
}

# Or by id:
data "oneuptime_ai_agent" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `connection_status` (String) Connection Status of the AI Agent.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create AI Agent], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit AI Agent]
- `icon_file_id` (String) AI Agent Icon File ID. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Is this the default AI Agent for the project? When set, this agent will be used for automated tasks.
- `key` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create AI Agent], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit AI Agent]
- `name` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create AI Agent], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit AI Agent]
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `ai_agent_version` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create AI Agent], Read: [Public], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit AI Agent]
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_alive` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create AI Agent], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read AI Agent], Update: [No access - you don't have permission for this operation]
- `project_id` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create AI Agent], Read: [Public], Update: [No access - you don't have permission for this operation]
- `updated_at` (String) Date and Time when the object was updated.

---
page_title: "oneuptime_runbook_secret Data Source - oneuptime"
subcategory: "Other"
description: |-
  Runbook Secret is a secret variable that can be used by runbook agents. For example you can store auth tokens, passwords, etc. in Runbook Secret and use them in your runbook steps. Runbook Secret is encrypted and only accessible by the assigned agent.
---

# oneuptime_runbook_secret (Data Source)

Runbook Secret is a secret variable that can be used by runbook agents. For example you can store auth tokens, passwords, etc. in Runbook Secret and use them in your runbook steps. Runbook Secret is encrypted and only accessible by the assigned agent.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one runbook secret may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_runbook_secret" "example" {
  name = "Example runbook secret"
}

# Or by id:
data "oneuptime_runbook_secret" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runners` (Set of String) List of runbook agents that can access this secret. IDs of `oneuptime_runner` resources.
- `updated_at` (String) Date and Time when the object was updated.

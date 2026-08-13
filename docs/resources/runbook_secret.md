---
page_title: "oneuptime_runbook_secret Resource - oneuptime"
subcategory: "Other"
description: |-
  Runbook Secret is a secret variable that can be used by runbook agents. For example you can store auth tokens, passwords, etc. in Runbook Secret and use them in your runbook steps. Runbook Secret is encrypted and only accessible by the assigned agent.
---

# oneuptime_runbook_secret (Resource)

Runbook Secret is a secret variable that can be used by runbook agents. For example you can store auth tokens, passwords, etc. in Runbook Secret and use them in your runbook steps. Runbook Secret is encrypted and only accessible by the assigned agent.

## Example Usage

```terraform
resource "oneuptime_runbook_secret" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `secret_value` (String) Secret value that you want to store in this object. This value will be encrypted and only accessible by the assigned runbook agent...
- `runners` (Set) List of runbook agents that can access this secret..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runbook_secret.example <id>
```

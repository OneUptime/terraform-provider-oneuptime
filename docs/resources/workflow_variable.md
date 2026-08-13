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
  content = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Variable Name..
- `content` (String) Content of the variable..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `workflow_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `is_secret` (Bool) Is this variable a secret. If true, then it'll not be in the logs..
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
terraform import oneuptime_workflow_variable.example <id>
```

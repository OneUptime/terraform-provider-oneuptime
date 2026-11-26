---
page_title: "oneuptime_security_case_custom_field Resource - oneuptime"
subcategory: "Security"
description: |-
  Fields a project adds to its security cases, filled in on each case.
---

# oneuptime_security_case_custom_field (Resource)

Fields a project adds to its security cases, filled in on each case.

## Example Usage

```terraform
resource "oneuptime_security_case_custom_field" "example" {
  name        = "Example security case custom field"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) The field's name, as each case shows it and as its value is kept under in the case's customFields.

### Optional

- `custom_field_type` (String) What the field holds: text, a number, yes or no, a date, a choice from a list, and the rest. Allowed values: `Text`, `Number`, `Boolean`, `Dropdown`, `MultiSelectDropdown`, `Date`, `DateTime`, `LongText`, `Markdown`, `File`, `MultipleFiles`.
- `description` (String) What the field is for, shown under it when a case's fields are edited.
- `dropdown_options` (String) Options and optional colors for dropdown fields, in the order they are listed. Plain one-per-line values remain supported. Records store an option as its text: changing an option here keeps the values records already hold. To rename an option and move those values with it, send the renames with the update, in miscDataProps: {"renamedDropdownOptions": [{"from": "Old text", "to": "New text"}]}. "to" must be one of the options.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand. Allowed values: `Monitor`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the field. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this custom field belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case custom field by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_custom_field.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_custom_field.example <id>
```

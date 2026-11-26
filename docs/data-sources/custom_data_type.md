---
page_title: "oneuptime_custom_data_type Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Your own kinds of sensitive data, such as employee IDs or policy numbers: a pattern, the words and field names that confirm it, and an optional checksum.
---

# oneuptime_custom_data_type (Data Source)

Your own kinds of sensitive data, such as employee IDs or policy numbers: a pattern, the words and field names that confirm it, and an optional checksum.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one custom data type may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_custom_data_type" "example" {
  name = "Example custom data type"
}

# Or by id:
data "oneuptime_custom_data_type" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `category` (String) The category it belongs to.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this custom data type. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the data type is.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ignore_case` (Boolean) Whether the pattern ignores upper and lower case.
- `is_enabled` (Boolean) Whether scans look for this data type.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) The data type's name, such as Employee ID.
- `pattern` (String) A JavaScript regular expression for one value. Leave empty to find it from field names alone.
- `requires_keyword` (Boolean) Count a value only when a keyword is nearby or the field's name hints at the type.
- `sensitivity` (String) How sensitive it is: Restricted, Confidential, Internal or Public.
- `validator` (String) A checksum a value must pass: None, Luhn, Mod97 or Verhoeff.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `keywords` (String) Words that, near a value, confirm it is this type. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `name_hints` (String) Field, column and attribute names that hold this type, such as employee id. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this custom data type belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

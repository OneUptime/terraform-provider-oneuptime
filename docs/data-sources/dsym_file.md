---
page_title: "oneuptime_dsym_file Data Source - oneuptime"
subcategory: "Other"
description: |-
  The function symbols of one build of an iOS app's binary, from its dSYM. Used to name the frames of the app's crash and hang stacks by their functions, and to group them by function. Matched to a frame by the binary's UUID.
---

# oneuptime_dsym_file (Data Source)

The function symbols of one build of an iOS app's binary, from its dSYM. Used to name the frames of the app's crash and hang stacks by their functions, and to group them by function. Matched to a frame by the binary's UUID.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dsym file may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_d_sym` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_dsym_file" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_dsym_file" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `app_name` (String) The app the dSYM was uploaded for, as the upload named it (its service.name). Shown beside the dSYM; a frame is matched by the binary's UUID alone.
- `architecture` (String) The architecture this build of the binary is for.
- `binary_name` (String) The binary the symbols are of, as its dSYM names it: the app's executable or one of its frameworks.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `debug_id` (String) The binary's UUID (its LC_UUID), upper case with dashes: the UUID a crash's or hang's frames name the build by.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `size_in_bytes` (Number) Size of the stored symbols in bytes.
- `symbol_count` (Number) How many function symbols this build of the binary has.
- `symbols` (String) Every function symbol of this build of the binary: its offset from the binary's load address and its name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

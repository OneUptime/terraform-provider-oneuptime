---
page_title: "oneuptime_dsym_file Resource - oneuptime"
subcategory: "Other"
description: |-
  The function symbols of one build of an iOS app's binary, from its dSYM. Used to name the frames of the app's crash and hang stacks by their functions, and to group them by function. Matched to a frame by the binary's UUID.
---

# oneuptime_dsym_file (Resource)

The function symbols of one build of an iOS app's binary, from its dSYM. Used to name the frames of the app's crash and hang stacks by their functions, and to group them by function. Matched to a frame by the binary's UUID.

~> **Renamed:** this resource was called `oneuptime_d_sym` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing dsym file:

```terraform
moved {
  from = oneuptime_d_sym.example
  to   = oneuptime_dsym_file.example
}
```

## Example Usage

```terraform
resource "oneuptime_dsym_file" "example" {
  debug_id     = "70B89F27-1634-3580-A695-57CDB41D7743"
  binary_name  = "MyApp"
  architecture = "arm64"
  symbols      = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `architecture` (String) The architecture this build of the binary is for.
- `binary_name` (String) The binary the symbols are of, as its dSYM names it: the app's executable or one of its frameworks.
- `debug_id` (String) The binary's UUID (its LC_UUID), upper case with dashes: the UUID a crash's or hang's frames name the build by.
- `symbols` (String) Every function symbol of this build of the binary: its offset from the binary's load address and its name.

### Optional

- `app_name` (String) The app the dSYM was uploaded for, as the upload named it (its service.name). Shown beside the dSYM; a frame is matched by the binary's UUID alone.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `size_in_bytes` (Number) Size of the stored symbols in bytes.
- `symbol_count` (Number) How many function symbols this build of the binary has.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing dsym file by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_dsym_file.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_dsym_file.example <id>
```

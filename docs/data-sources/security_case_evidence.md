---
page_title: "oneuptime_security_case_evidence Data Source - oneuptime"
subcategory: "Other"
description: |-
  Files and links kept on a security case as evidence, each with who added it and when, and a file with the SHA-256 of what was uploaded.
---

# oneuptime_security_case_evidence (Data Source)

Files and links kept on a security case as evidence, each with who added it and when, and a file with the SHA-256 of what was uploaded.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case evidence may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_evidence" "example" {
  name = "Example security case evidence"
}

# Or by id:
data "oneuptime_security_case_evidence" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who added the evidence. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the evidence shows, and where and how it was collected.
- `file_id` (String) ID of the uploaded file, for File evidence: upload it first, then name it here. The ID of a `oneuptime_file`.
- `file_name` (String) The file's name as it was uploaded, for File evidence.
- `file_size_in_bytes` (Number) The file's size as it was uploaded, in bytes.
- `file_type` (String) The file's type as it was uploaded, like text/plain.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What the evidence is: a File uploaded to OneUptime, or a Link to evidence kept somewhere else. It can't be changed.
- `name` (String) What the evidence is, in a few words. A file's own name when none is given.
- `security_case_id` (String) ID of the case this evidence is kept on. The ID of a `oneuptime_security_case`.
- `sha256` (String) The SHA-256 of the file's bytes as they were uploaded, in hex: compare it with a download's to check it is what was kept.
- `url` (String) Where the evidence is kept, for Link evidence: an http or https address.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this evidence belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

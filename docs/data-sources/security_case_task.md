---
page_title: "oneuptime_security_case_task Data Source - oneuptime"
subcategory: "Other"
description: |-
  What has to be done on a security case, who is doing it, and whether it is done.
---

# oneuptime_security_case_task (Data Source)

What has to be done on a security case, who is doing it, and whether it is done.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case task may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_task" "example" {
  security_case_id = oneuptime_security_case.example.id
}

# Or by id:
data "oneuptime_security_case_task" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `assigned_to_user_id` (String) ID of the person doing the task. Empty while nobody is. The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) ID of the person who added the task. The ID of a `oneuptime_user` (see the data source).
- `done_by_user_id` (String) ID of the person who marked the task done. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_done` (Boolean) Whether the task is done. Setting it records who did it and when; clearing it opens the task again.
- `security_case_id` (String) ID of the case this task is on. The ID of a `oneuptime_security_case`.
- `title` (String) What has to be done, in a line.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `done_at` (String) When the task was marked done. Empty while it is open.
- `project_id` (String) ID of the project this task belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

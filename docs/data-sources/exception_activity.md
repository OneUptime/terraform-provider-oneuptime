---
page_title: "oneuptime_exception_activity Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  The activity of an exception: who resolved, reopened, archived, assigned, merged or unmerged it and when, and when it came back after it was resolved. Read one exception's activity at a time, by its telemetryExceptionId.
---

# oneuptime_exception_activity (Data Source)

The activity of an exception: who resolved, reopened, archived, assigned, merged or unmerged it and when, and when it came back after it was resolved. Read one exception's activity at a time, by its telemetryExceptionId.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one exception activity may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_exception_activity" "example" {
  telemetry_exception_id = oneuptime_exception.example.id
}

# Or by id:
data "oneuptime_exception_activity" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `activity_type` (String) What happened: Resolved, ResolvedInNextRelease, Reopened, Regressed, Archived, Unarchived, Assigned, Unassigned, Merged, MergedInto, Unmerged or UnmergedFrom.
- `api_key_name` (String) The name the API key that made the change had at the time, when an API key did.
- `assignee_team_id` (String) For Assigned, ID of the team it was assigned to; for Unassigned, of the team it was taken from. The ID of a `oneuptime_team`.
- `assignee_user_id` (String) For Assigned, ID of the person it was assigned to; for Unassigned, of the person it was taken from. The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) ID of the person whose action wrote the entry. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_assignee_notified` (Boolean) For Assigned, whether the person, or the members of the team, it was assigned to were notified. True for every other entry.
- `related_telemetry_exception_id` (String) For a merge or an unmerge, ID of the other exception. The ID of a `oneuptime_exception`.
- `release` (String) For ResolvedInNextRelease, the release it waits out; for Regressed, the release of the occurrence that brought it back, when it names one.
- `telemetry_exception_id` (String) ID of the exception this entry is about. The ID of a `oneuptime_exception`.
- `user_id` (String) ID of the person who did it, when a person did. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this entry belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

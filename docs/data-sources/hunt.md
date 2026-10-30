---
page_title: "oneuptime_hunt Data Source - oneuptime"
subcategory: "Other"
description: |-
  Searches over security events saved with a name, to open again or to run on a schedule that tells people when it finds something.
---

# oneuptime_hunt (Data Source)

Searches over security events saved with a name, to open again or to run on a schedule that tells people when it finds something.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one hunt may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_hunt" "example" {
  name = "Example hunt"
}

# Or by id:
data "oneuptime_hunt" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who saved the hunt. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Why the hunt matters and what to do when it finds something, for whoever opens it later.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_found_event_count` (Number) How many events the last run that found something found.
- `last_run_error` (String) Why the hunt's last scheduled run could not be made. Empty when it ran.
- `last_run_event_count` (Number) How many events the hunt's last scheduled run found.
- `name` (String) What the hunt looks for, in a few words. At most 100 characters.
- `schedule` (String) How often the hunt runs on its own and tells people what it found: Not Scheduled, Every Hour, Every 6 Hours, Every Day or Every Week. Each run searches the events that arrived since the run before.
- `scheduled_by_user_id` (String) ID of the person the hunt's scheduled runs run as. The ID of a `oneuptime_user` (see the data source).
- `search_query` (String) What the hunt searches for, in the Events search bar's language: words a message contains, fields such as severity:High or user:alice, and source attributes such as @threat.matched:true. At most 2,000 characters.
- `time_range` (String) How far back the hunt searches when it is opened, as the Events time picker names it: Past 1 Hour, Past 1 Day, Past 1 Week, Past 1 Month and the rest.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `filters` (String) The Events sidebar's filters kept beside the search: each field's values, or with ! before the field the values it leaves out. At most 50 values. A JSON value: write it with `jsonencode()`.
- `last_found_at` (String) The end of the events the last run that found something searched: what it found arrived before this.
- `last_found_since_at` (String) The start of the events the last run that found something searched: what it found arrived after this.
- `last_run_at` (String) When the hunt last ran on its schedule, whether it found anything or not.
- `next_run_at` (String) When the hunt next runs on its schedule. Empty for a hunt that is not scheduled.
- `notify_teams` (Set of String) The teams whose members are told when a scheduled run finds events, each as their notification settings say. Only members who can read security events are told. IDs of `oneuptime_team` resources.
- `notify_users` (Set of String) The people told when a scheduled run finds events, each as their notification settings say. Only people who can read security events are told. IDs of `oneuptime_user` records.
- `project_id` (String) ID of the project this hunt belongs to. The ID of a `oneuptime_project`.
- `searched_up_to_at` (String) How far the scheduled runs have searched: the next run reads the events from here on. Set to the moment the schedule is set or the search changes.
- `updated_at` (String) Date and Time when the object was updated.

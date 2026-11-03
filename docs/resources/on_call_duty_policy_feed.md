---
page_title: "oneuptime_on_call_duty_policy_feed Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Log of the entire onCallDutyPolicy state change. This is a log of all the on call duty policy changes.
---

# oneuptime_on_call_duty_policy_feed (Resource)

Log of the entire onCallDutyPolicy state change. This is a log of all the on call duty policy changes.

## Example Usage

```terraform
resource "oneuptime_on_call_duty_policy_feed" "example" {
  on_call_duty_policy_id              = oneuptime_on_call_policy.example.id
  feed_info_in_markdown               = "# Heading\n\nThis is **markdown** content"
  on_call_duty_policy_feed_event_type = "UserAdded"
  display_color                       = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for the onCallDutyPolicy log.
- `feed_info_in_markdown` (String) Log of the entire onCallDutyPolicy state change in Markdown.
- `on_call_duty_policy_feed_event_type` (String) On Call Duty Policy Feed Event. Allowed values: `UserAdded`, `UserRemoved`, `TeamAdded`, `TeamRemoved`, `OnCallDutyScheduleAdded`, `OnCallDutyScheduleRemoved`, `OnCallDutyPolicyCreated`, `RosterHandoff`, `OwnerUserAdded`, `OwnerUserRemoved`, `OwnerTeamAdded`, `OwnerTeamRemoved`, `UserOverrideAdded`, `UserOverrideRemoved`, `CoverageGapStarted`.
- `on_call_duty_policy_id` (String) Relation to OnCallDutyPolicy ID in which this resource belongs. The ID of a `oneuptime_on_call_policy`.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call duty policy feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_duty_policy_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_duty_policy_feed.example <id>
```

---
page_title: "oneuptime_scheduled_maintenance_external_calendar Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Calendars from outside the project - public holidays, change freezes, another project's maintenance - read from their iCalendar (.ics) link and shown on the scheduled maintenance calendar. Read-only, and never part of the project's own calendar links.
---

# oneuptime_scheduled_maintenance_external_calendar (Resource)

Calendars from outside the project - public holidays, change freezes, another project's maintenance - read from their iCalendar (.ics) link and shown on the scheduled maintenance calendar. Read-only, and never part of the project's own calendar links.

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_external_calendar" "example" {
  name         = "Example scheduled maintenance external calendar"
  calendar_url = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `calendar_url` (String, Sensitive) The calendar's iCalendar (.ics) link, starting with https:// or webcal://. It can be a secret link: it is encrypted at rest and never returned by the API.
- `name` (String) What the calendar holds, as the maintenance calendar names it: Public holidays, Change freezes.

### Optional

- `color` (String) The colour the calendar's events are drawn in on the maintenance calendar, in hex (#0d9488 for example). Picked for you when left out.
- `is_enabled` (Boolean) When off, the calendar is not read and its events are not shown, until it is turned on again. Defaults to `true`.
- `refresh_interval_in_minutes` (Number) How often OneUptime reads the calendar: 15, 30, 60, 360 or 1440 minutes. Defaults to `60`.

### Read-Only

- `calendar_host` (String) The host the calendar link points at, such as calendar.google.com - what the page shows in place of the secret link.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `event_count` (Number) How many of the calendar's events OneUptime keeps: those from 90 days back to a year ahead, up to 2,000.
- `has_more_events` (Boolean) True when the calendar holds more events in the window than OneUptime keeps: the 2,000 nearest today are kept.
- `id` (String) Unique identifier for the resource.
- `last_read_at` (String) When OneUptime last read the calendar without a problem. Empty until the first read.
- `last_read_attempt_at` (String) When OneUptime last tried to read the calendar.
- `last_read_error` (String) Why the last read failed - NotFound, AccessDenied, ServerError, Unreachable, NotAllowed, Timeout, TooLarge, TooManyRedirects, NotACalendar or TooComplex. Empty when it worked.
- `last_read_http_status` (Number) The HTTP status the calendar's server answered the last failed read with, when it answered.
- `next_read_at` (String) When OneUptime next reads the calendar.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance external calendar by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_external_calendar.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_external_calendar.example <id>
```

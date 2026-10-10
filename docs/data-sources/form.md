---
page_title: "oneuptime_form Data Source - oneuptime"
subcategory: "Other"
description: |-
  Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident or a scheduled maintenance event in this project.
---

# oneuptime_form (Data Source)

Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident or a scheduled maintenance event in this project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one form may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_form" "example" {
  name = "Example form"
}

# Or by id:
data "oneuptime_form" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown.
- `favicon_file_id` (String) ID of the file the form's public page shows as the browser tab's icon: upload the image to /api/file first. Leave it empty to show the OneUptime favicon. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ip_whitelist` (String) The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network.
- `is_enabled` (Boolean) Whether the form's link works. While it is off, the public page shows a not-available message and nothing can be submitted.
- `logo_alt_text` (String) What the logo says, read out by screen readers: usually your organization's name. Leave it empty and screen readers skip the logo.
- `logo_file_id` (String) ID of the file the form's public page shows as its logo: upload the image to /api/file first. Leave it empty to show the OneUptime logo. The ID of a `oneuptime_file`.
- `name` (String) The form's name, shown as the heading of its public page. Unique within the project.
- `share_key` (String) The key in the form's public link, /accounts/form/<shareKey>. Generated when the form is created. Resetting the link in the dashboard replaces it, and the old link stops working.
- `success_message` (String) Shown after the form is submitted, together with the number of what the submission created. Markdown.
- `target_type` (String) What each submission creates: Incident, or ScheduledMaintenance (a scheduled maintenance event).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `fields` (String) The questions the form asks, in order. Each has an id, a source (Question: one of the form's own, answered by type; TargetField: a built-in field of what the form creates, by targetField; TargetCustomField: one of its custom fields, by customFieldId; Submitter: the submitter's Name or Email), a label, optional help text, isRequired and isHidden (not shown on the public form unless the template a submission starts from asks it, and otherwise answered only from the template a submission started from; never required, and never a field the target cannot be created without). isRequired and isHidden are the form's default: each template can make a question Required, Optional or Hidden (its fieldSettings). A new form starts with a title, a description and the submitter's name and email. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `target_settings` (String) What every submission starts with besides the answers. For incidents: defaultTitle, incidentSeverityId, incidentTemplateId, monitorIds, labelIds, onCallDutyPolicyIds, ownerUserIds and ownerTeamIds. For scheduled maintenance events: defaultTitle, monitorIds, statusPageIds, labelIds, ownerUserIds, ownerTeamIds, showOnStatusPages and notifySubscribers. A JSON value: write it with `jsonencode()`.
- `templates` (String) Named sets of answers a submission can start from, in the order the form lists them. Each has an id, a name (unique within the form), isDefault (the form opens with it; at most one template), answers: an object keyed by question id, each answer as a submission sends it - text, a number, true or false, an option's value, or a list of values for a multi-select - and fieldSettings: an object keyed by question id that makes a question Required, Optional or Hidden when a submission starts from the template; a question it does not list is asked as the form asks it, and a question the form cannot create its record without is always asked and required. The public form lists the templates above its questions and fills in a template's answers when one is chosen, or when its link names one (?template=<id>). A question the submission was not asked - hidden by the form or by the template - is answered only from the template a submission started from. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

---
page_title: "oneuptime_incident_severity Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident severity for your project (Created, Acknowledged for example). Add / edit or remove severities.
---

# oneuptime_incident_severity (Data Source)

Manage incident severity for your project (Created, Acknowledged for example). Add / edit or remove severities.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident severity may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_severity" "example" {
  name = "Example incident severity"
}

# Or by id:
data "oneuptime_incident_severity" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `order` (Number) Where this severity ranks among the project's incident severities: 1 is the most severe. A new severity without a number goes to the end of the list. Setting a number moves the severity to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

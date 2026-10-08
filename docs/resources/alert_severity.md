---
page_title: "oneuptime_alert_severity Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert severity for your project (Created, Acknowledged for example). Add / edit or remove severities.
---

# oneuptime_alert_severity (Resource)

Manage alert severity for your project (Created, Acknowledged for example). Add / edit or remove severities.

## Example Usage

```terraform
resource "oneuptime_alert_severity" "example" {
  name        = "Example alert severity"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `order` (Number) Where this severity ranks among the project's alert severities: 1 is the most severe. A new severity without a number goes to the end of the list. Setting a number moves the severity to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert severity by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_severity.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_severity.example <id>
```

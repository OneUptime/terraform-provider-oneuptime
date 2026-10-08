---
page_title: "oneuptime_oid_collection_template Resource - oneuptime"
subcategory: "Other"
description: |-
  A reusable set of SNMP health OIDs. Every network device linked to a template collects its OIDs, and editing the template changes what every linked device collects on its next poll.
---

# oneuptime_oid_collection_template (Resource)

A reusable set of SNMP health OIDs. Every network device linked to a template collects its OIDs, and editing the template changes what every linked device collects on its next poll.

## Example Usage

```terraform
resource "oneuptime_oid_collection_template" "example" {
  name        = "Example oid collection template"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) The device type this template describes. Devices linked to it all collect the same OIDs.

### Optional

- `description` (String) Friendly description that will help you remember.
- `oids` (String) SNMP OIDs (CPU, memory, temperature, or any custom OID) collected by every device linked to this template. You do not need OIDs for interfaces - bits in/out, errors, utilization and up/down are walked for every port automatically. A JSON value: write it with `jsonencode()`.
- `tables` (String) SNMP tables walked by every device linked to this template - one row per IPsec tunnel, Wi-Fi radio, routing neighbour, fan or power supply. Each table lists the column OIDs to collect and, optionally, the columns that name each row. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing oid collection template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_oid_collection_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_oid_collection_template.example <id>
```

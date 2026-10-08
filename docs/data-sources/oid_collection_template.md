---
page_title: "oneuptime_oid_collection_template Data Source - oneuptime"
subcategory: "Other"
description: |-
  A reusable set of SNMP health OIDs. Every network device linked to a template collects its OIDs, and editing the template changes what every linked device collects on its next poll.
---

# oneuptime_oid_collection_template (Data Source)

A reusable set of SNMP health OIDs. Every network device linked to a template collects its OIDs, and editing the template changes what every linked device collects on its next poll.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one oid collection template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_oid_collection_template" "example" {
  name = "Example oid collection template"
}

# Or by id:
data "oneuptime_oid_collection_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) The device type this template describes. Devices linked to it all collect the same OIDs.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `oids` (String) SNMP OIDs (CPU, memory, temperature, or any custom OID) collected by every device linked to this template. You do not need OIDs for interfaces - bits in/out, errors, utilization and up/down are walked for every port automatically. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `tables` (String) SNMP tables walked by every device linked to this template - one row per IPsec tunnel, Wi-Fi radio, routing neighbour, fan or power supply. Each table lists the column OIDs to collect and, optionally, the columns that name each row. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

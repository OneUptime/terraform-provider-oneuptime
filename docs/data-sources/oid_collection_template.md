---
page_title: "oneuptime_oid_collection_template Data Source - oneuptime"
subcategory: "Network"
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

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) The device type this template describes. Devices linked to it all collect the same OIDs.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `oids` (String) SNMP OIDs (CPU, memory, temperature, or any custom OID) collected by every device linked to this template. You do not need OIDs for interfaces - bits in/out, errors, utilization and up/down are walked for every port automatically. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `tables` (String) SNMP tables walked by every device linked to this template - one row per IPsec tunnel, Wi-Fi radio, routing neighbour, fan or power supply. Each table lists the column OIDs to collect and, optionally, the columns that name each row. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.

---
page_title: "oneuptime_inventory_item_relationship Data Source - oneuptime"
subcategory: "Other"
description: |-
  Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence.
---

# oneuptime_inventory_item_relationship (Data Source)

Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one inventory item relationship may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_inventory_item_relationship" "example" {
  from_entity_key = "example-from-entity-key"
}

# Or by id:
data "oneuptime_inventory_item_relationship" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `avg_duration_ms` (Number) Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only).
- `call_count` (Number) Calls observed over this edge in the most recent computation window (depends-on edges only).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `error_count` (Number) Errored calls observed over this edge in the most recent computation window (depends-on edges only).
- `from_entity_key` (String) Stable identity key of the source entity of this edge.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `relationship_type` (String) The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of).
- `source` (String) Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies.
- `to_entity_key` (String) Stable identity key of the target entity of this edge.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When this relationship was first observed in telemetry.
- `last_seen_at` (String) Most recent time this relationship was observed (bumped, throttled). Drives staleness pruning.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

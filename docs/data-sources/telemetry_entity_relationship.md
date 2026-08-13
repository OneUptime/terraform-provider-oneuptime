---
page_title: "oneuptime_telemetry_entity_relationship Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence.
---

# oneuptime_telemetry_entity_relationship (Data Source)

Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_telemetry_entity_relationship" "by_name" {
  name = "example-telemetry_entity_relationship"
}

data "oneuptime_telemetry_entity_relationship" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `from_entity_key` (String) Stable identity key of the source entity of this edge... Computed.
- `to_entity_key` (String) Stable identity key of the target entity of this edge... Computed.
- `relationship_type` (String) The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of)... Computed.
- `source` (String) Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies... Computed.
- `first_seen_at` (String) A date time object.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `call_count` (Number) Calls observed over this edge in the most recent computation window (depends-on edges only)... Computed.
- `error_count` (Number) Errored calls observed over this edge in the most recent computation window (depends-on edges only)... Computed.
- `avg_duration_ms` (Number) Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only)... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

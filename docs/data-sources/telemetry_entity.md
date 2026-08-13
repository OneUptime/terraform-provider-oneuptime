---
page_title: "oneuptime_telemetry_entity Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Catalog of OpenTelemetry entities (service, host, k8s.pod, container, ...) discovered from telemetry resource attributes.
---

# oneuptime_telemetry_entity (Data Source)

Catalog of OpenTelemetry entities (service, host, k8s.pod, container, ...) discovered from telemetry resource attributes. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_telemetry_entity" "by_name" {
  name = "example-telemetry_entity"
}

data "oneuptime_telemetry_entity" "by_id" {
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
- `entity_type` (String) The OpenTelemetry entity type (service, host, k8s.pod, container, ...)... Computed.
- `entity_key` (String) Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns)... Computed.
- `display_name` (String) Human-readable name derived for the entity explorer UI... Computed.
- `source` (String) How this row came to exist: discovered from telemetry, mirrored from a OneUptime inventory table, or created manually by a user. Determines whether stale-entity pruning applies... Computed.
- `description` (String) Free-text description. Primarily for manually created entities, where there are no telemetry attributes to explain what the thing is... Computed.
- `identifying_attributes` (String) The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key... Computed.
- `descriptive_attributes` (String) Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity... Computed.
- `labels` (String) Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up... Computed.
- `resource_type` (String) Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster)... Computed.
- `resource_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `first_seen_at` (String) A date time object.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

---
page_title: "oneuptime_telemetry_entity Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Catalog of OpenTelemetry entities (service, host, k8s.pod, container, ...) discovered from telemetry resource attributes.
---

# oneuptime_telemetry_entity (Resource)

Catalog of OpenTelemetry entities (service, host, k8s.pod, container, ...) discovered from telemetry resource attributes.

## Example Usage

```terraform
resource "oneuptime_telemetry_entity" "example" {
  entity_type = "Example short text"
  entity_key = "Example short text"
  source = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `entity_type` (String) The OpenTelemetry entity type (service, host, k8s.pod, container, ...)...
- `entity_key` (String) Stable identity hash derived from the entity's identifying attributes (matches the keys stamped into signal entityKeys columns)...
- `source` (String) How this row came to exist: discovered from telemetry, mirrored from a OneUptime inventory table, or created manually by a user. Determines whether stale-entity pruning applies...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `display_name` (String) Human-readable name derived for the entity explorer UI...
- `description` (String) Free-text description. Primarily for manually created entities, where there are no telemetry attributes to explain what the thing is...
- `identifying_attributes` (String) The immutable identifying attribute set (the entity's identity). Descriptive attributes are deliberately excluded so they can change without changing the entity key...
- `descriptive_attributes` (String) Mutable descriptive metadata (image tag, version, IP, ...) merged last-writer-wins. Never part of the identity...
- `labels` (String) Labels observed on this entity's telemetry (e.g. promoted from oneuptime.label.* resource attributes), merged as a set union. Simple string array in v1 — a relation to the Label table is a follow-up...
- `resource_type` (String) Polymorphic pointer type to a rich typed row, if one exists (Service / Host / DockerHost / KubernetesCluster)...
- `resource_id` (String) A unique identifier for an object, represented as a UUID..
- `first_seen_at` (String) A date time object..
- `last_seen_at` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_telemetry_entity.example <id>
```

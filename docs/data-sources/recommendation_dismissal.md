---
page_title: "oneuptime_recommendation_dismissal Data Source - oneuptime"
subcategory: "Other"
description: |-
  Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.
---

# oneuptime_recommendation_dismissal (Data Source)

Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_recommendation_dismissal" "by_name" {
  name = "example-recommendation_dismissal"
}

data "oneuptime_recommendation_dismissal" "by_id" {
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
- `recommendation_type` (String) Which family of recommendation this dismissal belongs to. See the RecommendationType enum... Computed.
- `recommendation_id` (String) The catalog-wide id of the dismissed recommendation, for example Kubernetes:k8s-hpa-at-max-replicas... Computed.
- `resource_type` (String) The kind of resource this recommendation was shown on, for example Kubernetes or Docker. Empty for recommendations that are not scoped to a resource... Computed.
- `resource_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `dismissal_reason` (String) Optional note explaining why this recommendation was dismissed, shown to whoever finds it in the dismissed list later... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

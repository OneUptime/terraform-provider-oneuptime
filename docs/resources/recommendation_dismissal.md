---
page_title: "oneuptime_recommendation_dismissal Resource - oneuptime"
subcategory: "Other"
description: |-
  Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.
---

# oneuptime_recommendation_dismissal (Resource)

Recommendations your team has dismissed. Dismissing hides a recommendation for everyone on the project until it is restored; it never deletes anything that was already created from it.

## Example Usage

```terraform
resource "oneuptime_recommendation_dismissal" "example" {
  recommendation_type = "Example short text"
  recommendation_id = "Example short text"
}
```

## Schema

### Required

- `recommendation_type` (String) Which family of recommendation this dismissal belongs to. See the RecommendationType enum...
- `recommendation_id` (String) The catalog-wide id of the dismissed recommendation, for example Kubernetes:k8s-hpa-at-max-replicas...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `resource_type` (String) The kind of resource this recommendation was shown on, for example Kubernetes or Docker. Empty for recommendations that are not scoped to a resource...
- `resource_id` (String) A unique identifier for an object, represented as a UUID..
- `dismissal_reason` (String) Optional note explaining why this recommendation was dismissed, shown to whoever finds it in the dismissed list later...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_recommendation_dismissal.example <id>
```

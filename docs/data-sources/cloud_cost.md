---
page_title: "oneuptime_cloud_cost Data Source - oneuptime"
subcategory: "Cloud"
description: |-
  API endpoints for Cloud Cost
---

# oneuptime_cloud_cost (Data Source)

API endpoints for Cloud Cost

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one cloud cost may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_cloud_cost" "example" {
  cloud_account_id = "example-cloud-account-id"
}

# Or by id:
data "oneuptime_cloud_cost" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `account_id` (String) Account.
- `cloud_account_id` (String) Cloud Account ID.
- `cloud_connection_id` (String) Cloud Connection ID.
- `cost` (Number) Cost.
- `currency` (String) Currency.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `provider_value` (String) Cloud Provider.
- `region` (String) Region.
- `resource_id` (String) Resource.
- `service` (String) Service.
- `tags` (String) Tags.

### Read-Only

- `day` (String) Day.
- `project_id` (String) Project ID.
- `read_at` (String) Read At.

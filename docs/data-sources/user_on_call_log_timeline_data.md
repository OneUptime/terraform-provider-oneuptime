---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_user_on_call_log_timeline_data Data Source - oneuptime"
subcategory: ""
description: |-
  User on call log timeline data data source
---

# oneuptime_user_on_call_log_timeline_data (Data Source)

User on call log timeline data data source

## Example Usage

```terraform
data "oneuptime_user_on_call_log_timeline_data" "example" {
  name = "example-user_on_call_log_timeline_data"
}
```

## Schema

- `id` (String) Identifier to filter by. Optional.
- `name` (String) Name to filter by. Optional.

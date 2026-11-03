---
page_title: "oneuptime_user Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Someone in this project, or invited to it: look them up by email address, and use their `id` wherever a resource takes a user ID. Every resource that names a person also takes their address directly (`user_email` beside `user_id`, `owner_user_emails` beside `owner_users`), so this is needed only to read who someone is, or which teams they are in. Only the API key's own project is searched.
---

# oneuptime_user (Data Source)

Someone in this project, or invited to it: look them up by email address, and use their `id` wherever a resource takes a user ID. Every resource that names a person also takes their address directly (`user_email` beside `user_id`, `owner_user_emails` beside `owner_users`), so this is needed only to read who someone is, or which teams they are in. Only the API key's own project is searched.

Look someone up by `email`, `name` or `id`: exactly one of them. Only the people of the API key's project are searched - its members, and those invited to it - so nobody outside it is ever found, and no match is an error rather than an empty result.

## Example Usage

```terraform
data "oneuptime_user" "example" {
  email = "alice@example.com"
}

# Or by id:
data "oneuptime_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `email` (String) The person's email address, in any case. Set it to look them up by it: the usual way.
- `id` (String) The person's user ID. Set it to look them up by ID instead of by email address or name.
- `name` (String) The person's name, exactly as it is shown in OneUptime. Set it to look them up by it; two people of the same name are an error, so prefer `email`.

### Read-Only

- `has_accepted_invitation` (Boolean) Whether they have joined the project: false while every invitation they have is still waiting for them.
- `team_ids` (Set of String) IDs of the `oneuptime_team` resources of this project they are in, or invited to.

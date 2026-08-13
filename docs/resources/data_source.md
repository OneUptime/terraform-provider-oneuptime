---
page_title: "oneuptime_data_source Resource - oneuptime"
subcategory: "Other"
description: |-
  Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data.
---

# oneuptime_data_source (Resource)

Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data.

## Example Usage

```terraform
resource "oneuptime_data_source" "example" {
  name = "Example short text"
  data_source_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `data_source_type` (String) The kind of external system this data source connects to (Prometheus, PostgreSQL, MySQL, Microsoft SQL Server, ClickHouse, Loki, Elasticsearch, or REST API)...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `url` (String) Base URL for HTTP-based sources (Prometheus, Loki, Elasticsearch, REST API) — e.g. https://prometheus.example.com..
- `database_host` (String) Hostname or IP address for database sources (PostgreSQL, MySQL, SQL Server, ClickHouse)...
- `database_port` (Number) Port for database sources. Defaults per engine: PostgreSQL 5432, MySQL 3306, SQL Server 1433, ClickHouse 8123...
- `database_name` (String) Database (or ClickHouse database) to connect to...
- `username` (String) Username for database sources, or HTTP basic-auth username for HTTP sources. Use a READ-ONLY account — dashboards only ever read...
- `password` (String) Password for database sources, or HTTP basic-auth password. Encrypted at rest and never returned by the API...
- `api_token` (String) Bearer token for HTTP sources (Prometheus behind an auth proxy, Elasticsearch API key, REST APIs). Encrypted at rest and never returned by the API...
- `custom_headers` (String) Extra HTTP headers sent to HTTP-based sources (e.g. auth headers for a proxy). Values are encrypted at rest and never returned by the API...
- `additional_options` (String) Per-type options that are not secrets — e.g. { "sslEnabled": true, "elasticsearchIndex": "logs-*" }...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_data_source.example <id>
```

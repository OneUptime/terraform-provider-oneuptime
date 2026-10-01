---
page_title: "oneuptime_queue Resource - oneuptime"
subcategory: "Other"
description: |-
  Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually.
---

# oneuptime_queue (Resource)

Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually.

## Example Usage

```terraform
resource "oneuptime_queue" "example" {
  name = "Example short text"
  queue_identifier = "This is an example of longer text content that might be stored in this field."
  messaging_system = "Example short text"
  destination_name = "This is an example of longer text content that might be stored in this field."
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this queue. Discovered queues are named after their destination, e.g. orders.created. Not unique - two brokers may each have a queue of the same name. Rename freely...
- `queue_identifier` (String) Stable identity of this queue within the project: the messaging system's identity family, the broker scope (the Azure Service Bus / Event Hubs namespace, empty for every other system) and the destination, canonicalized and joined with '|' (e.g. kafka||orders.created, servicebus|orders-prod|orders). An ActiveMQ queue keys as jms, like the JMS clients that use it, so both are one queue. Computed by the server for manually added queues...
- `messaging_system` (String) The broker this queue lives on, as an OpenTelemetry messaging.system value, e.g. kafka, rabbitmq, activemq, jms, aws_sqs, aws.sns, gcp_pubsub, servicebus, eventhubs, pulsar, rocketmq, nats or bullmq. A queue first seen through JMS is refined to its broker (activemq) once that broker's metrics are seen...
- `destination_name` (String) The queue, topic or subscription name as applications and the broker name it (the OpenTelemetry messaging.destination.name value), normalized the way discovery normalizes it: an SQS queue URL or an SNS topic ARN becomes its name, a Pub/Sub resource path its id, a Pulsar short name its persistent://public/default/ topic, and each UUID {uuid}...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this queue..
- `broker_scope` (String) The Azure Service Bus or Event Hubs namespace this queue lives in, lowercased (the first part of <namespace>.servicebus.windows.net). Part of the queue's identity, because two namespaces can each hold a queue of the same name. Empty for every other messaging system...
- `discovery_source` (String) How this queue was first discovered: traces (messaging spans of instrumented applications), broker-metrics (the broker's own metrics) or manual...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this queue archived? Archived queues are hidden from lists but their telemetry is still collected...
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `broker_address` (String) The broker address (host[:port]) last reported for this queue by its telemetry, with any credentials removed. Display only - it is not part of the queue's identity...
- `last_seen_at` (String) A date time object..
- `broker_metrics_last_seen_at` (String) A date time object..
- `auto_archived_at` (String) A date time object..
- `manually_restored_at` (String) A date time object..
- `automatic_assignments` (String) Label and owner ids that label rules or owner rules attached automatically. Maintained by OneUptime...
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_queue.example <id>
```

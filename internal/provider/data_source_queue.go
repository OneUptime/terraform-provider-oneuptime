package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QueueDataSource{}

func NewQueueDataSource() datasource.DataSource {
    return &QueueDataSource{}
}

// QueueDataSource defines the data source implementation.
type QueueDataSource struct {
    client *Client
}

// QueueDataSourceModel describes the data source data model.
type QueueDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    QueueIdentifier types.String `tfsdk:"queue_identifier"`
    MessagingSystem types.String `tfsdk:"messaging_system"`
    DestinationName types.String `tfsdk:"destination_name"`
    BrokerScope types.String `tfsdk:"broker_scope"`
    BrokerAddress types.String `tfsdk:"broker_address"`
    DiscoverySource types.String `tfsdk:"discovery_source"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    BrokerMetricsLastSeenAt types.String `tfsdk:"broker_metrics_last_seen_at"`
    AutoArchivedAt types.String `tfsdk:"auto_archived_at"`
    ManuallyRestoredAt types.String `tfsdk:"manually_restored_at"`
    AutomaticAssignments types.String `tfsdk:"automatic_assignments"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
}

func (d *QueueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_queue"
}

func (d *QueueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually. Look up an existing queue by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this queue.",
                Computed: true,
            },
            "queue_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identity of this queue within the project: the messaging system's identity family, the broker scope (the Azure Service Bus / Event Hubs namespace, empty for every other system) and the destination, canonicalized and joined with '|' (e.g. kafka||orders.created, servicebus|orders-prod|orders). An ActiveMQ queue keys as jms, like the JMS clients that use it, so both are one queue. Computed by the server for manually added queues..",
                Computed: true,
            },
            "messaging_system": schema.StringAttribute{
                MarkdownDescription: "The broker this queue lives on, as an OpenTelemetry messaging.system value, e.g. kafka, rabbitmq, activemq, jms, aws_sqs, aws.sns, gcp_pubsub, servicebus, eventhubs, pulsar, rocketmq, nats or bullmq. A queue first seen through JMS is refined to its broker (activemq) once that broker's metrics are seen..",
                Computed: true,
            },
            "destination_name": schema.StringAttribute{
                MarkdownDescription: "The queue, topic or subscription name as applications and the broker name it (the OpenTelemetry messaging.destination.name value), normalized the way discovery normalizes it: an SQS queue URL or an SNS topic ARN becomes its name, a Pub/Sub resource path its id, a Pulsar short name its persistent://public/default/ topic, and each UUID {uuid}..",
                Computed: true,
            },
            "broker_scope": schema.StringAttribute{
                MarkdownDescription: "The Azure Service Bus or Event Hubs namespace this queue lives in, lowercased (the first part of <namespace>.servicebus.windows.net). Part of the queue's identity, because two namespaces can each hold a queue of the same name. Empty for every other messaging system..",
                Computed: true,
            },
            "broker_address": schema.StringAttribute{
                MarkdownDescription: "The broker address (host[:port]) last reported for this queue by its telemetry, with any credentials removed. Display only - it is not part of the queue's identity..",
                Computed: true,
            },
            "discovery_source": schema.StringAttribute{
                MarkdownDescription: "How this queue was first discovered: traces (messaging spans of instrumented applications), broker-metrics (the broker's own metrics) or manual..",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "broker_metrics_last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "auto_archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "manually_restored_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "automatic_assignments": schema.StringAttribute{
                MarkdownDescription: "Label and owner ids that label rules or owner rules attached automatically. Maintained by OneUptime..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this queue archived? Archived queues are hidden from lists but their telemetry is still collected..",
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *QueueDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *QueueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data QueueDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a queue.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "slug": true,
        "description": true,
        "queueIdentifier": true,
        "messagingSystem": true,
        "destinationName": true,
        "brokerScope": true,
        "brokerAddress": true,
        "discoverySource": true,
        "lastSeenAt": true,
        "brokerMetricsLastSeenAt": true,
        "autoArchivedAt": true,
        "manuallyRestoredAt": true,
        "automaticAssignments": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "labels": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/message-queue/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read queue, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No queue found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read queue: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/message-queue/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list queue, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list queue: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No queue found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one queue matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for queue.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["queueIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.QueueIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.QueueIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.QueueIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.QueueIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.QueueIdentifier = types.StringNull()
        }
    } else if val, ok := item["queueIdentifier"].(string); ok {
        data.QueueIdentifier = types.StringValue(val)
    } else {
        data.QueueIdentifier = types.StringNull()
    }
    if obj, ok := item["messagingSystem"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MessagingSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MessagingSystem = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MessagingSystem = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MessagingSystem = types.StringValue(string(jsonBytes))
        } else {
            data.MessagingSystem = types.StringNull()
        }
    } else if val, ok := item["messagingSystem"].(string); ok {
        data.MessagingSystem = types.StringValue(val)
    } else {
        data.MessagingSystem = types.StringNull()
    }
    if obj, ok := item["destinationName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DestinationName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DestinationName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DestinationName = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationName = types.StringNull()
        }
    } else if val, ok := item["destinationName"].(string); ok {
        data.DestinationName = types.StringValue(val)
    } else {
        data.DestinationName = types.StringNull()
    }
    if obj, ok := item["brokerScope"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BrokerScope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BrokerScope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BrokerScope = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BrokerScope = types.StringValue(string(jsonBytes))
        } else {
            data.BrokerScope = types.StringNull()
        }
    } else if val, ok := item["brokerScope"].(string); ok {
        data.BrokerScope = types.StringValue(val)
    } else {
        data.BrokerScope = types.StringNull()
    }
    if obj, ok := item["brokerAddress"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BrokerAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BrokerAddress = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BrokerAddress = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BrokerAddress = types.StringValue(string(jsonBytes))
        } else {
            data.BrokerAddress = types.StringNull()
        }
    } else if val, ok := item["brokerAddress"].(string); ok {
        data.BrokerAddress = types.StringValue(val)
    } else {
        data.BrokerAddress = types.StringNull()
    }
    if obj, ok := item["discoverySource"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoverySource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DiscoverySource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DiscoverySource = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DiscoverySource = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoverySource = types.StringNull()
        }
    } else if val, ok := item["discoverySource"].(string); ok {
        data.DiscoverySource = types.StringValue(val)
    } else {
        data.DiscoverySource = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    if obj, ok := item["brokerMetricsLastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BrokerMetricsLastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BrokerMetricsLastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BrokerMetricsLastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BrokerMetricsLastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.BrokerMetricsLastSeenAt = types.StringNull()
        }
    } else if val, ok := item["brokerMetricsLastSeenAt"].(string); ok {
        data.BrokerMetricsLastSeenAt = types.StringValue(val)
    } else {
        data.BrokerMetricsLastSeenAt = types.StringNull()
    }
    if obj, ok := item["autoArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AutoArchivedAt = types.StringNull()
        }
    } else if val, ok := item["autoArchivedAt"].(string); ok {
        data.AutoArchivedAt = types.StringValue(val)
    } else {
        data.AutoArchivedAt = types.StringNull()
    }
    if obj, ok := item["manuallyRestoredAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ManuallyRestoredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ManuallyRestoredAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ManuallyRestoredAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ManuallyRestoredAt = types.StringValue(string(jsonBytes))
        } else {
            data.ManuallyRestoredAt = types.StringNull()
        }
    } else if val, ok := item["manuallyRestoredAt"].(string); ok {
        data.ManuallyRestoredAt = types.StringValue(val)
    } else {
        data.ManuallyRestoredAt = types.StringNull()
    }
    if obj, ok := item["automaticAssignments"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutomaticAssignments = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutomaticAssignments = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutomaticAssignments = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutomaticAssignments = types.StringValue(string(jsonBytes))
        } else {
            data.AutomaticAssignments = types.StringNull()
        }
    } else if val, ok := item["automaticAssignments"].(string); ok {
        data.AutomaticAssignments = types.StringValue(val)
    } else {
        data.AutomaticAssignments = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
    }
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

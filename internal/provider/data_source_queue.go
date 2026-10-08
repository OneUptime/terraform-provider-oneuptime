package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
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
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
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
    Labels types.Set `tfsdk:"labels"`
}

func (d *QueueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_queue"
}

func (d *QueueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Message queues, topics and subscriptions this project's applications publish to and consume from. Each queue is discovered from the OpenTelemetry messaging spans of instrumented applications and from broker metrics sent by an OpenTelemetry Collector, or added manually. Look up an existing queue by `id`, or by any of its other arguments (`name`, `archived_by_user_id`, `broker_address`, ...): each one set must match, and exactly one queue may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this queue. Discovered queues are named after their destination, e.g. orders.created. Not unique - two brokers may each have a queue of the same name. Rename freely.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description for this queue.",
                Optional: true,
                Computed: true,
            },
            "queue_identifier": schema.StringAttribute{
                MarkdownDescription: "Stable identity of this queue within the project: the messaging system's identity family, the broker scope (the Azure Service Bus / Event Hubs namespace, empty for every other system) and the destination, canonicalized and joined with '|' (e.g. kafka||orders.created, servicebus|orders-prod|orders). An ActiveMQ queue keys as jms, like the JMS clients that use it, so both are one queue. Computed by the server for manually added queues.",
                Optional: true,
                Computed: true,
            },
            "messaging_system": schema.StringAttribute{
                MarkdownDescription: "The broker this queue lives on, as an OpenTelemetry messaging.system value, e.g. kafka, rabbitmq, activemq, jms, aws_sqs, aws.sns, gcp_pubsub, servicebus, eventhubs, pulsar, rocketmq, nats or bullmq. A queue first seen through JMS is refined to its broker (activemq) once that broker's metrics are seen.",
                Optional: true,
                Computed: true,
            },
            "destination_name": schema.StringAttribute{
                MarkdownDescription: "The queue, topic or subscription name as applications and the broker name it (the OpenTelemetry messaging.destination.name value), normalized the way discovery normalizes it: an SQS queue URL or an SNS topic ARN becomes its name, a Pub/Sub resource path its id, a Pulsar short name its persistent://public/default/ topic, and each UUID {uuid}.",
                Optional: true,
                Computed: true,
            },
            "broker_scope": schema.StringAttribute{
                MarkdownDescription: "The Azure Service Bus or Event Hubs namespace this queue lives in, lowercased (the first part of <namespace>.servicebus.windows.net). Part of the queue's identity, because two namespaces can each hold a queue of the same name. Empty for every other messaging system.",
                Optional: true,
                Computed: true,
            },
            "broker_address": schema.StringAttribute{
                MarkdownDescription: "The broker address (host[:port]) last reported for this queue by its telemetry, with any credentials removed. Display only - it is not part of the queue's identity.",
                Optional: true,
                Computed: true,
            },
            "discovery_source": schema.StringAttribute{
                MarkdownDescription: "How this queue was first discovered: traces (messaging spans of instrumented applications), broker-metrics (the broker's own metrics) or manual.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When this queue was last seen by any discovery source - application messaging spans or broker metrics.",
                Computed: true,
            },
            "broker_metrics_last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the broker's own metrics for this queue (backlog, lag, dead letters, ...) were last received. Empty when they never were.",
                Computed: true,
            },
            "auto_archived_at": schema.StringAttribute{
                MarkdownDescription: "When this queue was archived automatically because no discovery source had seen it for a while. Empty when it was archived by a person or is not archived.",
                Computed: true,
            },
            "manually_restored_at": schema.StringAttribute{
                MarkdownDescription: "When a person last restored this queue from the archive. Automatic archiving leaves it alone until it is seen again or a grace period passes.",
                Computed: true,
            },
            "automatic_assignments": schema.StringAttribute{
                MarkdownDescription: "Label and owner ids that label rules or owner rules attached automatically. Maintained by OneUptime. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this queue archived? Archived queues are hidden from lists but their telemetry is still collected.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this queue archived?",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
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

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.QueueIdentifier.IsNull() && !data.QueueIdentifier.IsUnknown() {
        filters["queueIdentifier"] = data.QueueIdentifier.ValueString()
        filterNames = append(filterNames, "queue_identifier = "+fmt.Sprintf("%q", data.QueueIdentifier.ValueString()))
    }
    if !data.MessagingSystem.IsNull() && !data.MessagingSystem.IsUnknown() {
        filters["messagingSystem"] = data.MessagingSystem.ValueString()
        filterNames = append(filterNames, "messaging_system = "+fmt.Sprintf("%q", data.MessagingSystem.ValueString()))
    }
    if !data.DestinationName.IsNull() && !data.DestinationName.IsUnknown() {
        filters["destinationName"] = data.DestinationName.ValueString()
        filterNames = append(filterNames, "destination_name = "+fmt.Sprintf("%q", data.DestinationName.ValueString()))
    }
    if !data.BrokerScope.IsNull() && !data.BrokerScope.IsUnknown() {
        filters["brokerScope"] = data.BrokerScope.ValueString()
        filterNames = append(filterNames, "broker_scope = "+fmt.Sprintf("%q", data.BrokerScope.ValueString()))
    }
    if !data.BrokerAddress.IsNull() && !data.BrokerAddress.IsUnknown() {
        filters["brokerAddress"] = data.BrokerAddress.ValueString()
        filterNames = append(filterNames, "broker_address = "+fmt.Sprintf("%q", data.BrokerAddress.ValueString()))
    }
    if !data.DiscoverySource.IsNull() && !data.DiscoverySource.IsUnknown() {
        filters["discoverySource"] = data.DiscoverySource.ValueString()
        filterNames = append(filterNames, "discovery_source = "+fmt.Sprintf("%q", data.DiscoverySource.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the queue up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the queue up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
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
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No queue matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one queue matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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

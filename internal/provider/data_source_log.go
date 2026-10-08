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
var _ datasource.DataSource = &LogDataSource{}

func NewLogDataSource() datasource.DataSource {
    return &LogDataSource{}
}

// LogDataSource defines the data source implementation.
type LogDataSource struct {
    client *Client
}

// LogDataSourceModel describes the data source data model.
type LogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    Time types.String `tfsdk:"time"`
    TimeUnixNano types.String `tfsdk:"time_unix_nano"`
    SeverityText types.String `tfsdk:"severity_text"`
    SeverityNumber types.Number `tfsdk:"severity_number"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
    EntityKeys types.Set `tfsdk:"entity_keys"`
    ServiceEntityKey types.String `tfsdk:"service_entity_key"`
    HostEntityKey types.String `tfsdk:"host_entity_key"`
    K8sPodEntityKey types.String `tfsdk:"k8s_pod_entity_key"`
    K8sNodeEntityKey types.String `tfsdk:"k8s_node_entity_key"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    ContainerEntityKey types.String `tfsdk:"container_entity_key"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    SessionId types.String `tfsdk:"session_id"`
    Body types.String `tfsdk:"body"`
    ObservedTimeUnixNano types.String `tfsdk:"observed_time_unix_nano"`
    DroppedAttributesCount types.Number `tfsdk:"dropped_attributes_count"`
    Flags types.Number `tfsdk:"flags"`
}

func (d *LogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log"
}

func (d *LogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Log Look up an existing log by `id`, or by any of its other arguments (`attributes`, `body`, `container_entity_key`, ...): each one set must match, and exactly one log may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID.",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Service ID.",
                Optional: true,
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Service Type.",
                Optional: true,
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time.",
                Optional: true,
                Computed: true,
            },
            "time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Time (in Unix Nano).",
                Optional: true,
                Computed: true,
            },
            "severity_text": schema.StringAttribute{
                MarkdownDescription: "Severity Text.",
                Optional: true,
                Computed: true,
            },
            "severity_number": schema.NumberAttribute{
                MarkdownDescription: "Severity Number.",
                Optional: true,
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes.",
                Optional: true,
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "entity_keys": schema.SetAttribute{
                MarkdownDescription: "Entity Keys.",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_entity_key": schema.StringAttribute{
                MarkdownDescription: "Service Entity Key.",
                Optional: true,
                Computed: true,
            },
            "host_entity_key": schema.StringAttribute{
                MarkdownDescription: "Host Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_pod_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Pod Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_node_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Node Entity Key.",
                Optional: true,
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key.",
                Optional: true,
                Computed: true,
            },
            "container_entity_key": schema.StringAttribute{
                MarkdownDescription: "Container Entity Key.",
                Optional: true,
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID.",
                Optional: true,
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID.",
                Optional: true,
                Computed: true,
            },
            "session_id": schema.StringAttribute{
                MarkdownDescription: "Session ID.",
                Optional: true,
                Computed: true,
            },
            "body": schema.StringAttribute{
                MarkdownDescription: "Log Body.",
                Optional: true,
                Computed: true,
            },
            "observed_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Observed Time (in Unix Nano).",
                Optional: true,
                Computed: true,
            },
            "dropped_attributes_count": schema.NumberAttribute{
                MarkdownDescription: "Dropped Attributes Count.",
                Optional: true,
                Computed: true,
            },
            "flags": schema.NumberAttribute{
                MarkdownDescription: "Flags.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *LogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.PrimaryEntityId.IsNull() && !data.PrimaryEntityId.IsUnknown() {
        filters["primaryEntityId"] = data.PrimaryEntityId.ValueString()
        filterNames = append(filterNames, "primary_entity_id = "+fmt.Sprintf("%q", data.PrimaryEntityId.ValueString()))
    }
    if !data.PrimaryEntityType.IsNull() && !data.PrimaryEntityType.IsUnknown() {
        filters["primaryEntityType"] = data.PrimaryEntityType.ValueString()
        filterNames = append(filterNames, "primary_entity_type = "+fmt.Sprintf("%q", data.PrimaryEntityType.ValueString()))
    }
    if !data.Time.IsNull() && !data.Time.IsUnknown() {
        filters["time"] = data.Time.ValueString()
        filterNames = append(filterNames, "time = "+fmt.Sprintf("%q", data.Time.ValueString()))
    }
    if !data.TimeUnixNano.IsNull() && !data.TimeUnixNano.IsUnknown() {
        filters["timeUnixNano"] = data.TimeUnixNano.ValueString()
        filterNames = append(filterNames, "time_unix_nano = "+fmt.Sprintf("%q", data.TimeUnixNano.ValueString()))
    }
    if !data.SeverityText.IsNull() && !data.SeverityText.IsUnknown() {
        filters["severityText"] = data.SeverityText.ValueString()
        filterNames = append(filterNames, "severity_text = "+fmt.Sprintf("%q", data.SeverityText.ValueString()))
    }
    if !data.SeverityNumber.IsNull() && !data.SeverityNumber.IsUnknown() {
        filters["severityNumber"] = lookupNumber(data.SeverityNumber)
        filterNames = append(filterNames, "severity_number = "+data.SeverityNumber.ValueBigFloat().String())
    }
    if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
        filters["attributes"] = data.Attributes.ValueString()
        filterNames = append(filterNames, "attributes = "+fmt.Sprintf("%q", data.Attributes.ValueString()))
    }
    if !data.ServiceEntityKey.IsNull() && !data.ServiceEntityKey.IsUnknown() {
        filters["serviceEntityKey"] = data.ServiceEntityKey.ValueString()
        filterNames = append(filterNames, "service_entity_key = "+fmt.Sprintf("%q", data.ServiceEntityKey.ValueString()))
    }
    if !data.HostEntityKey.IsNull() && !data.HostEntityKey.IsUnknown() {
        filters["hostEntityKey"] = data.HostEntityKey.ValueString()
        filterNames = append(filterNames, "host_entity_key = "+fmt.Sprintf("%q", data.HostEntityKey.ValueString()))
    }
    if !data.K8sPodEntityKey.IsNull() && !data.K8sPodEntityKey.IsUnknown() {
        filters["k8sPodEntityKey"] = data.K8sPodEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_pod_entity_key = "+fmt.Sprintf("%q", data.K8sPodEntityKey.ValueString()))
    }
    if !data.K8sNodeEntityKey.IsNull() && !data.K8sNodeEntityKey.IsUnknown() {
        filters["k8sNodeEntityKey"] = data.K8sNodeEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_node_entity_key = "+fmt.Sprintf("%q", data.K8sNodeEntityKey.ValueString()))
    }
    if !data.K8sClusterEntityKey.IsNull() && !data.K8sClusterEntityKey.IsUnknown() {
        filters["k8sClusterEntityKey"] = data.K8sClusterEntityKey.ValueString()
        filterNames = append(filterNames, "k8s_cluster_entity_key = "+fmt.Sprintf("%q", data.K8sClusterEntityKey.ValueString()))
    }
    if !data.ContainerEntityKey.IsNull() && !data.ContainerEntityKey.IsUnknown() {
        filters["containerEntityKey"] = data.ContainerEntityKey.ValueString()
        filterNames = append(filterNames, "container_entity_key = "+fmt.Sprintf("%q", data.ContainerEntityKey.ValueString()))
    }
    if !data.TraceId.IsNull() && !data.TraceId.IsUnknown() {
        filters["traceId"] = data.TraceId.ValueString()
        filterNames = append(filterNames, "trace_id = "+fmt.Sprintf("%q", data.TraceId.ValueString()))
    }
    if !data.SpanId.IsNull() && !data.SpanId.IsUnknown() {
        filters["spanId"] = data.SpanId.ValueString()
        filterNames = append(filterNames, "span_id = "+fmt.Sprintf("%q", data.SpanId.ValueString()))
    }
    if !data.SessionId.IsNull() && !data.SessionId.IsUnknown() {
        filters["sessionId"] = data.SessionId.ValueString()
        filterNames = append(filterNames, "session_id = "+fmt.Sprintf("%q", data.SessionId.ValueString()))
    }
    if !data.Body.IsNull() && !data.Body.IsUnknown() {
        filters["body"] = data.Body.ValueString()
        filterNames = append(filterNames, "body = "+fmt.Sprintf("%q", data.Body.ValueString()))
    }
    if !data.ObservedTimeUnixNano.IsNull() && !data.ObservedTimeUnixNano.IsUnknown() {
        filters["observedTimeUnixNano"] = data.ObservedTimeUnixNano.ValueString()
        filterNames = append(filterNames, "observed_time_unix_nano = "+fmt.Sprintf("%q", data.ObservedTimeUnixNano.ValueString()))
    }
    if !data.DroppedAttributesCount.IsNull() && !data.DroppedAttributesCount.IsUnknown() {
        filters["droppedAttributesCount"] = lookupNumber(data.DroppedAttributesCount)
        filterNames = append(filterNames, "dropped_attributes_count = "+data.DroppedAttributesCount.ValueBigFloat().String())
    }
    if !data.Flags.IsNull() && !data.Flags.IsUnknown() {
        filters["flags"] = lookupNumber(data.Flags)
        filterNames = append(filterNames, "flags = "+data.Flags.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "time": true,
        "timeUnixNano": true,
        "severityText": true,
        "severityNumber": true,
        "attributes": true,
        "attributeKeys": true,
        "entityKeys": true,
        "serviceEntityKey": true,
        "hostEntityKey": true,
        "k8sPodEntityKey": true,
        "k8sNodeEntityKey": true,
        "k8sClusterEntityKey": true,
        "containerEntityKey": true,
        "traceId": true,
        "spanId": true,
        "sessionId": true,
        "body": true,
        "observedTimeUnixNano": true,
        "droppedAttributesCount": true,
        "flags": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/logs/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/logs/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for log.")
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
    if obj, ok := item["primaryEntityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := item["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := item["primaryEntityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := item["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
    }
    if obj, ok := item["time"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Time = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Time = types.StringValue(string(jsonBytes))
        } else {
            data.Time = types.StringNull()
        }
    } else if val, ok := item["time"].(string); ok {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if obj, ok := item["timeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.TimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["timeUnixNano"].(string); ok {
        data.TimeUnixNano = types.StringValue(val)
    } else {
        data.TimeUnixNano = types.StringNull()
    }
    if obj, ok := item["severityText"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SeverityText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SeverityText = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SeverityText = types.StringValue(string(jsonBytes))
        } else {
            data.SeverityText = types.StringNull()
        }
    } else if val, ok := item["severityText"].(string); ok {
        data.SeverityText = types.StringValue(val)
    } else {
        data.SeverityText = types.StringNull()
    }
    if val, ok := item["severityNumber"].(float64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["severityNumber"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SeverityNumber = types.NumberValue(big.NewFloat(val))
        } else {
            data.SeverityNumber = types.NumberNull()
        }
    } else {
        data.SeverityNumber = types.NumberNull()
    }
    if obj, ok := item["attributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := item["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := item["attributeKeys"].([]interface{}); ok {
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
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AttributeKeys = types.SetNull(types.StringType)
    }
    if val, ok := item["entityKeys"].([]interface{}); ok {
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
        data.EntityKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.EntityKeys = types.SetNull(types.StringType)
    }
    if obj, ok := item["serviceEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceEntityKey = types.StringNull()
        }
    } else if val, ok := item["serviceEntityKey"].(string); ok {
        data.ServiceEntityKey = types.StringValue(val)
    } else {
        data.ServiceEntityKey = types.StringNull()
    }
    if obj, ok := item["hostEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HostEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HostEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HostEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HostEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.HostEntityKey = types.StringNull()
        }
    } else if val, ok := item["hostEntityKey"].(string); ok {
        data.HostEntityKey = types.StringValue(val)
    } else {
        data.HostEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sPodEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sPodEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sPodEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sPodEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sPodEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sPodEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sPodEntityKey"].(string); ok {
        data.K8sPodEntityKey = types.StringValue(val)
    } else {
        data.K8sPodEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sNodeEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sNodeEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sNodeEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sNodeEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sNodeEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sNodeEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sNodeEntityKey"].(string); ok {
        data.K8sNodeEntityKey = types.StringValue(val)
    } else {
        data.K8sNodeEntityKey = types.StringNull()
    }
    if obj, ok := item["k8sClusterEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.K8sClusterEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.K8sClusterEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.K8sClusterEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.K8sClusterEntityKey = types.StringNull()
        }
    } else if val, ok := item["k8sClusterEntityKey"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    } else {
        data.K8sClusterEntityKey = types.StringNull()
    }
    if obj, ok := item["containerEntityKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContainerEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContainerEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContainerEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContainerEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ContainerEntityKey = types.StringNull()
        }
    } else if val, ok := item["containerEntityKey"].(string); ok {
        data.ContainerEntityKey = types.StringValue(val)
    } else {
        data.ContainerEntityKey = types.StringNull()
    }
    if obj, ok := item["traceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := item["traceId"].(string); ok {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
    }
    if obj, ok := item["spanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SpanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SpanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SpanId = types.StringValue(string(jsonBytes))
        } else {
            data.SpanId = types.StringNull()
        }
    } else if val, ok := item["spanId"].(string); ok {
        data.SpanId = types.StringValue(val)
    } else {
        data.SpanId = types.StringNull()
    }
    if obj, ok := item["sessionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SessionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SessionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SessionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SessionId = types.StringValue(string(jsonBytes))
        } else {
            data.SessionId = types.StringNull()
        }
    } else if val, ok := item["sessionId"].(string); ok {
        data.SessionId = types.StringValue(val)
    } else {
        data.SessionId = types.StringNull()
    }
    if obj, ok := item["body"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Body = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Body = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Body = types.StringValue(string(jsonBytes))
        } else {
            data.Body = types.StringNull()
        }
    } else if val, ok := item["body"].(string); ok {
        data.Body = types.StringValue(val)
    } else {
        data.Body = types.StringNull()
    }
    if obj, ok := item["observedTimeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ObservedTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ObservedTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ObservedTimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ObservedTimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.ObservedTimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["observedTimeUnixNano"].(string); ok {
        data.ObservedTimeUnixNano = types.StringValue(val)
    } else {
        data.ObservedTimeUnixNano = types.StringNull()
    }
    if val, ok := item["droppedAttributesCount"].(float64); ok {
        data.DroppedAttributesCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["droppedAttributesCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DroppedAttributesCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.DroppedAttributesCount = types.NumberNull()
        }
    } else {
        data.DroppedAttributesCount = types.NumberNull()
    }
    if val, ok := item["flags"].(float64); ok {
        data.Flags = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["flags"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Flags = types.NumberValue(big.NewFloat(val))
        } else {
            data.Flags = types.NumberNull()
        }
    } else {
        data.Flags = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

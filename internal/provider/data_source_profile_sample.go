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
var _ datasource.DataSource = &ProfileSampleDataSource{}

func NewProfileSampleDataSource() datasource.DataSource {
    return &ProfileSampleDataSource{}
}

// ProfileSampleDataSource defines the data source implementation.
type ProfileSampleDataSource struct {
    client *Client
}

// ProfileSampleDataSourceModel describes the data source data model.
type ProfileSampleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    ProfileId types.String `tfsdk:"profile_id"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    Time types.String `tfsdk:"time"`
    TimeUnixNano types.String `tfsdk:"time_unix_nano"`
    Stacktrace types.Set `tfsdk:"stacktrace"`
    StacktraceHash types.String `tfsdk:"stacktrace_hash"`
    FrameTypes types.Set `tfsdk:"frame_types"`
    Value types.String `tfsdk:"value"`
    ProfileType types.String `tfsdk:"profile_type"`
    Labels types.String `tfsdk:"labels"`
    EntityKeys types.Set `tfsdk:"entity_keys"`
    ServiceEntityKey types.String `tfsdk:"service_entity_key"`
    HostEntityKey types.String `tfsdk:"host_entity_key"`
    K8sPodEntityKey types.String `tfsdk:"k8s_pod_entity_key"`
    K8sNodeEntityKey types.String `tfsdk:"k8s_node_entity_key"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    ContainerEntityKey types.String `tfsdk:"container_entity_key"`
}

func (d *ProfileSampleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_profile_sample"
}

func (d *ProfileSampleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for ProfileSample Look up an existing profile sample by `id`, or by any of its other arguments (`container_entity_key`, `host_entity_key`, `k8s_cluster_entity_key`, ...): each one set must match, and exactly one profile sample may match them all.",

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
            "profile_id": schema.StringAttribute{
                MarkdownDescription: "Profile ID.",
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
            "stacktrace": schema.SetAttribute{
                MarkdownDescription: "Stacktrace.",
                Computed: true,
                ElementType: types.StringType,
            },
            "stacktrace_hash": schema.StringAttribute{
                MarkdownDescription: "Stacktrace Hash.",
                Optional: true,
                Computed: true,
            },
            "frame_types": schema.SetAttribute{
                MarkdownDescription: "Frame Types.",
                Computed: true,
                ElementType: types.StringType,
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "Value.",
                Optional: true,
                Computed: true,
            },
            "profile_type": schema.StringAttribute{
                MarkdownDescription: "Profile Type.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.StringAttribute{
                MarkdownDescription: "Labels.",
                Optional: true,
                Computed: true,
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
        },
    }
}

func (d *ProfileSampleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProfileSampleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProfileSampleDataSourceModel

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
    if !data.ProfileId.IsNull() && !data.ProfileId.IsUnknown() {
        filters["profileId"] = data.ProfileId.ValueString()
        filterNames = append(filterNames, "profile_id = "+fmt.Sprintf("%q", data.ProfileId.ValueString()))
    }
    if !data.TraceId.IsNull() && !data.TraceId.IsUnknown() {
        filters["traceId"] = data.TraceId.ValueString()
        filterNames = append(filterNames, "trace_id = "+fmt.Sprintf("%q", data.TraceId.ValueString()))
    }
    if !data.SpanId.IsNull() && !data.SpanId.IsUnknown() {
        filters["spanId"] = data.SpanId.ValueString()
        filterNames = append(filterNames, "span_id = "+fmt.Sprintf("%q", data.SpanId.ValueString()))
    }
    if !data.Time.IsNull() && !data.Time.IsUnknown() {
        filters["time"] = data.Time.ValueString()
        filterNames = append(filterNames, "time = "+fmt.Sprintf("%q", data.Time.ValueString()))
    }
    if !data.TimeUnixNano.IsNull() && !data.TimeUnixNano.IsUnknown() {
        filters["timeUnixNano"] = data.TimeUnixNano.ValueString()
        filterNames = append(filterNames, "time_unix_nano = "+fmt.Sprintf("%q", data.TimeUnixNano.ValueString()))
    }
    if !data.StacktraceHash.IsNull() && !data.StacktraceHash.IsUnknown() {
        filters["stacktraceHash"] = data.StacktraceHash.ValueString()
        filterNames = append(filterNames, "stacktrace_hash = "+fmt.Sprintf("%q", data.StacktraceHash.ValueString()))
    }
    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        filters["value"] = data.Value.ValueString()
        filterNames = append(filterNames, "value = "+fmt.Sprintf("%q", data.Value.ValueString()))
    }
    if !data.ProfileType.IsNull() && !data.ProfileType.IsUnknown() {
        filters["profileType"] = data.ProfileType.ValueString()
        filterNames = append(filterNames, "profile_type = "+fmt.Sprintf("%q", data.ProfileType.ValueString()))
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        filters["labels"] = data.Labels.ValueString()
        filterNames = append(filterNames, "labels = "+fmt.Sprintf("%q", data.Labels.ValueString()))
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

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the profile sample up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the profile sample up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "profileId": true,
        "traceId": true,
        "spanId": true,
        "time": true,
        "timeUnixNano": true,
        "stacktrace": true,
        "stacktraceHash": true,
        "frameTypes": true,
        "value": true,
        "profileType": true,
        "labels": true,
        "entityKeys": true,
        "serviceEntityKey": true,
        "hostEntityKey": true,
        "k8sPodEntityKey": true,
        "k8sNodeEntityKey": true,
        "k8sClusterEntityKey": true,
        "containerEntityKey": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/profile-sample/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile_sample, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No profile sample found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read profile_sample: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/profile-sample/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list profile_sample, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list profile_sample: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No profile sample matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one profile sample matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for profile_sample.")
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
    if obj, ok := item["profileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProfileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProfileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProfileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProfileId = types.StringValue(string(jsonBytes))
        } else {
            data.ProfileId = types.StringNull()
        }
    } else if val, ok := item["profileId"].(string); ok {
        data.ProfileId = types.StringValue(val)
    } else {
        data.ProfileId = types.StringNull()
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
    if val, ok := item["stacktrace"].([]interface{}); ok {
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
        data.Stacktrace = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Stacktrace = types.SetNull(types.StringType)
    }
    if obj, ok := item["stacktraceHash"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StacktraceHash = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StacktraceHash = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StacktraceHash = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StacktraceHash = types.StringValue(string(jsonBytes))
        } else {
            data.StacktraceHash = types.StringNull()
        }
    } else if val, ok := item["stacktraceHash"].(string); ok {
        data.StacktraceHash = types.StringValue(val)
    } else {
        data.StacktraceHash = types.StringNull()
    }
    if val, ok := item["frameTypes"].([]interface{}); ok {
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
        data.FrameTypes = types.SetValueMust(types.StringType, setItems)
    } else {
        data.FrameTypes = types.SetNull(types.StringType)
    }
    if obj, ok := item["value"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := item["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := item["profileType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProfileType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProfileType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProfileType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProfileType = types.StringValue(string(jsonBytes))
        } else {
            data.ProfileType = types.StringNull()
        }
    } else if val, ok := item["profileType"].(string); ok {
        data.ProfileType = types.StringValue(val)
    } else {
        data.ProfileType = types.StringNull()
    }
    if obj, ok := item["labels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Labels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Labels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Labels = types.StringValue(string(jsonBytes))
        } else {
            data.Labels = types.StringNull()
        }
    } else if val, ok := item["labels"].(string); ok {
        data.Labels = types.StringValue(val)
    } else {
        data.Labels = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

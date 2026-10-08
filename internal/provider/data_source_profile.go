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
var _ datasource.DataSource = &ProfileDataSource{}

func NewProfileDataSource() datasource.DataSource {
    return &ProfileDataSource{}
}

// ProfileDataSource defines the data source implementation.
type ProfileDataSource struct {
    client *Client
}

// ProfileDataSourceModel describes the data source data model.
type ProfileDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    ProfileId types.String `tfsdk:"profile_id"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    StartTime types.String `tfsdk:"start_time"`
    EndTime types.String `tfsdk:"end_time"`
    StartTimeUnixNano types.String `tfsdk:"start_time_unix_nano"`
    EndTimeUnixNano types.String `tfsdk:"end_time_unix_nano"`
    DurationNano types.String `tfsdk:"duration_nano"`
    ProfileType types.String `tfsdk:"profile_type"`
    Unit types.String `tfsdk:"unit"`
    PeriodType types.String `tfsdk:"period_type"`
    Period types.String `tfsdk:"period"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
    EntityKeys types.Set `tfsdk:"entity_keys"`
    ServiceEntityKey types.String `tfsdk:"service_entity_key"`
    HostEntityKey types.String `tfsdk:"host_entity_key"`
    K8sPodEntityKey types.String `tfsdk:"k8s_pod_entity_key"`
    K8sNodeEntityKey types.String `tfsdk:"k8s_node_entity_key"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    ContainerEntityKey types.String `tfsdk:"container_entity_key"`
    SampleCount types.Number `tfsdk:"sample_count"`
    OriginalPayloadFormat types.String `tfsdk:"original_payload_format"`
}

func (d *ProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_profile"
}

func (d *ProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Profile Look up an existing profile by `id`, or by any of its other arguments (`attributes`, `container_entity_key`, `duration_nano`, ...): each one set must match, and exactly one profile may match them all.",

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
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time.",
                Optional: true,
                Computed: true,
            },
            "end_time": schema.StringAttribute{
                MarkdownDescription: "End Time.",
                Optional: true,
                Computed: true,
            },
            "start_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Start Time in Unix Nano.",
                Optional: true,
                Computed: true,
            },
            "end_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "End Time in Unix Nano.",
                Optional: true,
                Computed: true,
            },
            "duration_nano": schema.StringAttribute{
                MarkdownDescription: "Duration in Nanoseconds.",
                Optional: true,
                Computed: true,
            },
            "profile_type": schema.StringAttribute{
                MarkdownDescription: "Profile Type.",
                Optional: true,
                Computed: true,
            },
            "unit": schema.StringAttribute{
                MarkdownDescription: "Unit.",
                Optional: true,
                Computed: true,
            },
            "period_type": schema.StringAttribute{
                MarkdownDescription: "Period Type.",
                Optional: true,
                Computed: true,
            },
            "period": schema.StringAttribute{
                MarkdownDescription: "Period.",
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
            "sample_count": schema.NumberAttribute{
                MarkdownDescription: "Sample Count.",
                Optional: true,
                Computed: true,
            },
            "original_payload_format": schema.StringAttribute{
                MarkdownDescription: "Original Payload Format.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *ProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProfileDataSourceModel

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
    if !data.StartTime.IsNull() && !data.StartTime.IsUnknown() {
        filters["startTime"] = data.StartTime.ValueString()
        filterNames = append(filterNames, "start_time = "+fmt.Sprintf("%q", data.StartTime.ValueString()))
    }
    if !data.EndTime.IsNull() && !data.EndTime.IsUnknown() {
        filters["endTime"] = data.EndTime.ValueString()
        filterNames = append(filterNames, "end_time = "+fmt.Sprintf("%q", data.EndTime.ValueString()))
    }
    if !data.StartTimeUnixNano.IsNull() && !data.StartTimeUnixNano.IsUnknown() {
        filters["startTimeUnixNano"] = data.StartTimeUnixNano.ValueString()
        filterNames = append(filterNames, "start_time_unix_nano = "+fmt.Sprintf("%q", data.StartTimeUnixNano.ValueString()))
    }
    if !data.EndTimeUnixNano.IsNull() && !data.EndTimeUnixNano.IsUnknown() {
        filters["endTimeUnixNano"] = data.EndTimeUnixNano.ValueString()
        filterNames = append(filterNames, "end_time_unix_nano = "+fmt.Sprintf("%q", data.EndTimeUnixNano.ValueString()))
    }
    if !data.DurationNano.IsNull() && !data.DurationNano.IsUnknown() {
        filters["durationNano"] = data.DurationNano.ValueString()
        filterNames = append(filterNames, "duration_nano = "+fmt.Sprintf("%q", data.DurationNano.ValueString()))
    }
    if !data.ProfileType.IsNull() && !data.ProfileType.IsUnknown() {
        filters["profileType"] = data.ProfileType.ValueString()
        filterNames = append(filterNames, "profile_type = "+fmt.Sprintf("%q", data.ProfileType.ValueString()))
    }
    if !data.Unit.IsNull() && !data.Unit.IsUnknown() {
        filters["unit"] = data.Unit.ValueString()
        filterNames = append(filterNames, "unit = "+fmt.Sprintf("%q", data.Unit.ValueString()))
    }
    if !data.PeriodType.IsNull() && !data.PeriodType.IsUnknown() {
        filters["periodType"] = data.PeriodType.ValueString()
        filterNames = append(filterNames, "period_type = "+fmt.Sprintf("%q", data.PeriodType.ValueString()))
    }
    if !data.Period.IsNull() && !data.Period.IsUnknown() {
        filters["period"] = data.Period.ValueString()
        filterNames = append(filterNames, "period = "+fmt.Sprintf("%q", data.Period.ValueString()))
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
    if !data.SampleCount.IsNull() && !data.SampleCount.IsUnknown() {
        filters["sampleCount"] = lookupNumber(data.SampleCount)
        filterNames = append(filterNames, "sample_count = "+data.SampleCount.ValueBigFloat().String())
    }
    if !data.OriginalPayloadFormat.IsNull() && !data.OriginalPayloadFormat.IsUnknown() {
        filters["originalPayloadFormat"] = data.OriginalPayloadFormat.ValueString()
        filterNames = append(filterNames, "original_payload_format = "+fmt.Sprintf("%q", data.OriginalPayloadFormat.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the profile up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the profile up by.",
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
        "startTime": true,
        "endTime": true,
        "startTimeUnixNano": true,
        "endTimeUnixNano": true,
        "durationNano": true,
        "profileType": true,
        "unit": true,
        "periodType": true,
        "period": true,
        "attributes": true,
        "attributeKeys": true,
        "entityKeys": true,
        "serviceEntityKey": true,
        "hostEntityKey": true,
        "k8sPodEntityKey": true,
        "k8sNodeEntityKey": true,
        "k8sClusterEntityKey": true,
        "containerEntityKey": true,
        "sampleCount": true,
        "originalPayloadFormat": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/profile/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No profile found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read profile: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/profile/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list profile, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list profile: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No profile matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one profile matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for profile.")
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
    if obj, ok := item["startTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartTime = types.StringValue(string(jsonBytes))
        } else {
            data.StartTime = types.StringNull()
        }
    } else if val, ok := item["startTime"].(string); ok {
        data.StartTime = types.StringValue(val)
    } else {
        data.StartTime = types.StringNull()
    }
    if obj, ok := item["endTime"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndTime = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndTime = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndTime = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndTime = types.StringValue(string(jsonBytes))
        } else {
            data.EndTime = types.StringNull()
        }
    } else if val, ok := item["endTime"].(string); ok {
        data.EndTime = types.StringValue(val)
    } else {
        data.EndTime = types.StringNull()
    }
    if obj, ok := item["startTimeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartTimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartTimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.StartTimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["startTimeUnixNano"].(string); ok {
        data.StartTimeUnixNano = types.StringValue(val)
    } else {
        data.StartTimeUnixNano = types.StringNull()
    }
    if obj, ok := item["endTimeUnixNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndTimeUnixNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndTimeUnixNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndTimeUnixNano = types.StringValue(string(jsonBytes))
        } else {
            data.EndTimeUnixNano = types.StringNull()
        }
    } else if val, ok := item["endTimeUnixNano"].(string); ok {
        data.EndTimeUnixNano = types.StringValue(val)
    } else {
        data.EndTimeUnixNano = types.StringNull()
    }
    if obj, ok := item["durationNano"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DurationNano = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DurationNano = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DurationNano = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DurationNano = types.StringValue(string(jsonBytes))
        } else {
            data.DurationNano = types.StringNull()
        }
    } else if val, ok := item["durationNano"].(string); ok {
        data.DurationNano = types.StringValue(val)
    } else {
        data.DurationNano = types.StringNull()
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
    if obj, ok := item["unit"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Unit = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Unit = types.StringValue(string(jsonBytes))
        } else {
            data.Unit = types.StringNull()
        }
    } else if val, ok := item["unit"].(string); ok {
        data.Unit = types.StringValue(val)
    } else {
        data.Unit = types.StringNull()
    }
    if obj, ok := item["periodType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PeriodType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PeriodType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PeriodType = types.StringValue(string(jsonBytes))
        } else {
            data.PeriodType = types.StringNull()
        }
    } else if val, ok := item["periodType"].(string); ok {
        data.PeriodType = types.StringValue(val)
    } else {
        data.PeriodType = types.StringNull()
    }
    if obj, ok := item["period"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Period = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Period = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Period = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Period = types.StringValue(string(jsonBytes))
        } else {
            data.Period = types.StringNull()
        }
    } else if val, ok := item["period"].(string); ok {
        data.Period = types.StringValue(val)
    } else {
        data.Period = types.StringNull()
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
    if val, ok := item["sampleCount"].(float64); ok {
        data.SampleCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sampleCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SampleCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.SampleCount = types.NumberNull()
        }
    } else {
        data.SampleCount = types.NumberNull()
    }
    if obj, ok := item["originalPayloadFormat"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OriginalPayloadFormat = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OriginalPayloadFormat = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OriginalPayloadFormat = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OriginalPayloadFormat = types.StringValue(string(jsonBytes))
        } else {
            data.OriginalPayloadFormat = types.StringNull()
        }
    } else if val, ok := item["originalPayloadFormat"].(string); ok {
        data.OriginalPayloadFormat = types.StringValue(val)
    } else {
        data.OriginalPayloadFormat = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataFlowDataSource{}

func NewDataFlowDataSource() datasource.DataSource {
    return &DataFlowDataSource{}
}

// DataFlowDataSource defines the data source implementation.
type DataFlowDataSource struct {
    client *Client
}

// DataFlowDataSourceModel describes the data source data model.
type DataFlowDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    FlowKey types.String `tfsdk:"flow_key"`
    Kind types.String `tfsdk:"kind"`
    SourceDataStoreId types.String `tfsdk:"source_data_store_id"`
    ServiceId types.String `tfsdk:"service_id"`
    DestinationDataStoreId types.String `tfsdk:"destination_data_store_id"`
    Destination types.String `tfsdk:"destination"`
    DataTypeIds types.String `tfsdk:"data_type_ids"`
    Sensitivity types.String `tfsdk:"sensitivity"`
    IsConfirmed types.Bool `tfsdk:"is_confirmed"`
    OverlapPercent types.Number `tfsdk:"overlap_percent"`
    CallCount types.Number `tfsdk:"call_count"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
}

func (d *DataFlowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_flow"
}

func (d *DataFlowDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Where sensitive data goes: the services that read each store, where those services send data, and copies of a store's data found in another store. Look up an existing data flow by `id`, or by any of its other arguments (`call_count`, `destination`, `destination_data_store_id`, ...): each one set must match, and exactly one data flow may match them all.",

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
                MarkdownDescription: "ID of the project this data flow belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "flow_key": schema.StringAttribute{
                MarkdownDescription: "The flow's kind, source and destination, which identify it within the project.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What the flow is: Read, Write, AI Provider, External Host or Copy.",
                Optional: true,
                Computed: true,
            },
            "source_data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store the data comes from. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "The service the data passes through. The ID of a `oneuptime_service`.",
                Optional: true,
                Computed: true,
            },
            "destination_data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store the data reaches, for a write or a copy. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
            },
            "destination": schema.StringAttribute{
                MarkdownDescription: "Where the data goes, when it is not a store: a host or an AI provider.",
                Optional: true,
                Computed: true,
            },
            "data_type_ids": schema.StringAttribute{
                MarkdownDescription: "The sensitive data types that can travel along the flow. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "sensitivity": schema.StringAttribute{
                MarkdownDescription: "The sensitivity of the data that can travel along it.",
                Optional: true,
                Computed: true,
            },
            "is_confirmed": schema.BoolAttribute{
                MarkdownDescription: "A scan found the data at the destination, not only the path to it.",
                Optional: true,
                Computed: true,
            },
            "overlap_percent": schema.NumberAttribute{
                MarkdownDescription: "For a copy: how much of the source's values were found in the destination, from 0 to 100.",
                Optional: true,
                Computed: true,
            },
            "call_count": schema.NumberAttribute{
                MarkdownDescription: "How many calls the service map counted along the flow in its last window.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the flow was first seen.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the flow was last seen.",
                Computed: true,
            },
        },
    }
}

func (d *DataFlowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataFlowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataFlowDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.FlowKey.IsNull() && !data.FlowKey.IsUnknown() {
        filters["flowKey"] = data.FlowKey.ValueString()
        filterNames = append(filterNames, "flow_key = "+fmt.Sprintf("%q", data.FlowKey.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.SourceDataStoreId.IsNull() && !data.SourceDataStoreId.IsUnknown() {
        filters["sourceDataStoreId"] = data.SourceDataStoreId.ValueString()
        filterNames = append(filterNames, "source_data_store_id = "+fmt.Sprintf("%q", data.SourceDataStoreId.ValueString()))
    }
    if !data.ServiceId.IsNull() && !data.ServiceId.IsUnknown() {
        filters["serviceId"] = data.ServiceId.ValueString()
        filterNames = append(filterNames, "service_id = "+fmt.Sprintf("%q", data.ServiceId.ValueString()))
    }
    if !data.DestinationDataStoreId.IsNull() && !data.DestinationDataStoreId.IsUnknown() {
        filters["destinationDataStoreId"] = data.DestinationDataStoreId.ValueString()
        filterNames = append(filterNames, "destination_data_store_id = "+fmt.Sprintf("%q", data.DestinationDataStoreId.ValueString()))
    }
    if !data.Destination.IsNull() && !data.Destination.IsUnknown() {
        filters["destination"] = data.Destination.ValueString()
        filterNames = append(filterNames, "destination = "+fmt.Sprintf("%q", data.Destination.ValueString()))
    }
    if !data.Sensitivity.IsNull() && !data.Sensitivity.IsUnknown() {
        filters["sensitivity"] = data.Sensitivity.ValueString()
        filterNames = append(filterNames, "sensitivity = "+fmt.Sprintf("%q", data.Sensitivity.ValueString()))
    }
    if !data.IsConfirmed.IsNull() && !data.IsConfirmed.IsUnknown() {
        filters["isConfirmed"] = data.IsConfirmed.ValueBool()
        filterNames = append(filterNames, "is_confirmed = "+fmt.Sprintf("%t", data.IsConfirmed.ValueBool()))
    }
    if !data.OverlapPercent.IsNull() && !data.OverlapPercent.IsUnknown() {
        filters["overlapPercent"] = lookupNumber(data.OverlapPercent)
        filterNames = append(filterNames, "overlap_percent = "+data.OverlapPercent.ValueBigFloat().String())
    }
    if !data.CallCount.IsNull() && !data.CallCount.IsUnknown() {
        filters["callCount"] = lookupNumber(data.CallCount)
        filterNames = append(filterNames, "call_count = "+data.CallCount.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data flow up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data flow up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "flowKey": true,
        "kind": true,
        "sourceDataStoreId": true,
        "serviceId": true,
        "destinationDataStoreId": true,
        "destination": true,
        "dataTypeIds": true,
        "sensitivity": true,
        "isConfirmed": true,
        "overlapPercent": true,
        "callCount": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-flow/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_flow, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data flow found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_flow: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-flow/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_flow, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_flow: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data flow matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data flow matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_flow.")
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
    if obj, ok := item["flowKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FlowKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FlowKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FlowKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FlowKey = types.StringValue(string(jsonBytes))
        } else {
            data.FlowKey = types.StringNull()
        }
    } else if val, ok := item["flowKey"].(string); ok {
        data.FlowKey = types.StringValue(val)
    } else {
        data.FlowKey = types.StringNull()
    }
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := item["sourceDataStoreId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SourceDataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SourceDataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SourceDataStoreId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SourceDataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.SourceDataStoreId = types.StringNull()
        }
    } else if val, ok := item["sourceDataStoreId"].(string); ok {
        data.SourceDataStoreId = types.StringValue(val)
    } else {
        data.SourceDataStoreId = types.StringNull()
    }
    if obj, ok := item["serviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := item["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := item["destinationDataStoreId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationDataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DestinationDataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DestinationDataStoreId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DestinationDataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationDataStoreId = types.StringNull()
        }
    } else if val, ok := item["destinationDataStoreId"].(string); ok {
        data.DestinationDataStoreId = types.StringValue(val)
    } else {
        data.DestinationDataStoreId = types.StringNull()
    }
    if obj, ok := item["destination"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Destination = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Destination = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Destination = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Destination = types.StringValue(string(jsonBytes))
        } else {
            data.Destination = types.StringNull()
        }
    } else if val, ok := item["destination"].(string); ok {
        data.Destination = types.StringValue(val)
    } else {
        data.Destination = types.StringNull()
    }
    if obj, ok := item["dataTypeIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataTypeIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataTypeIds = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeIds = types.StringNull()
        }
    } else if val, ok := item["dataTypeIds"].(string); ok {
        data.DataTypeIds = types.StringValue(val)
    } else {
        data.DataTypeIds = types.StringNull()
    }
    if obj, ok := item["sensitivity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Sensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if val, ok := item["sensitivity"].(string); ok {
        data.Sensitivity = types.StringValue(val)
    } else {
        data.Sensitivity = types.StringNull()
    }
    if val, ok := item["isConfirmed"].(bool); ok {
        data.IsConfirmed = types.BoolValue(val)
    } else {
        data.IsConfirmed = types.BoolNull()
    }
    data.OverlapPercent = numberFromAPI(data.OverlapPercent, item["overlapPercent"])
    data.CallCount = numberFromAPI(data.CallCount, item["callCount"])
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

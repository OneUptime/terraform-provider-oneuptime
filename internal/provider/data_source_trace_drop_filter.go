package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TraceDropFilterDataSource{}

func NewTraceDropFilterDataSource() datasource.DataSource {
    return &TraceDropFilterDataSource{}
}

// TraceDropFilterDataSource defines the data source implementation.
type TraceDropFilterDataSource struct {
    client *Client
}

// TraceDropFilterDataSourceModel describes the data source data model.
type TraceDropFilterDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    FilterQuery types.String `tfsdk:"filter_query"`
    Action types.String `tfsdk:"action"`
    SamplePercentage types.Number `tfsdk:"sample_percentage"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    SortOrder types.Number `tfsdk:"sort_order"`
    DroppedCount types.Number `tfsdk:"dropped_count"`
    LastDroppedAt types.String `tfsdk:"last_dropped_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *TraceDropFilterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_trace_drop_filter"
}

func (d *TraceDropFilterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules to drop or sample spans before storage to reduce volume and cost. Look up an existing trace drop filter by `id`, or by any of its other arguments (`name`, `action`, `created_by_user_id`, ...): each one set must match, and exactly one trace drop filter may match them all.",

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
                MarkdownDescription: "ID of the project this trace drop filter belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this drop filter.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this drop filter does.",
                Optional: true,
                Computed: true,
            },
            "filter_query": schema.StringAttribute{
                MarkdownDescription: "Filter expression that identifies which spans to drop or sample.",
                Optional: true,
                Computed: true,
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "What to do with matching spans: 'drop' to discard entirely, 'sample' to keep a percentage.",
                Optional: true,
                Computed: true,
            },
            "sample_percentage": schema.NumberAttribute{
                MarkdownDescription: "When action is 'sample', the percentage of matching spans to keep (1-99).",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this drop filter is active.",
                Optional: true,
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Where this filter is evaluated among the project's span drop filters, lowest number first. A new filter is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "dropped_count": schema.NumberAttribute{
                MarkdownDescription: "Total number of spans this filter has discarded since it was created.",
                Optional: true,
                Computed: true,
            },
            "last_dropped_at": schema.StringAttribute{
                MarkdownDescription: "When this filter most recently discarded a span. Null means it has never matched anything.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this trace drop filter. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *TraceDropFilterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TraceDropFilterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TraceDropFilterDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.FilterQuery.IsNull() && !data.FilterQuery.IsUnknown() {
        filters["filterQuery"] = data.FilterQuery.ValueString()
        filterNames = append(filterNames, "filter_query = "+fmt.Sprintf("%q", data.FilterQuery.ValueString()))
    }
    if !data.Action.IsNull() && !data.Action.IsUnknown() {
        filters["action"] = data.Action.ValueString()
        filterNames = append(filterNames, "action = "+fmt.Sprintf("%q", data.Action.ValueString()))
    }
    if !data.SamplePercentage.IsNull() && !data.SamplePercentage.IsUnknown() {
        filters["samplePercentage"] = lookupNumber(data.SamplePercentage)
        filterNames = append(filterNames, "sample_percentage = "+data.SamplePercentage.ValueBigFloat().String())
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.SortOrder.IsNull() && !data.SortOrder.IsUnknown() {
        filters["sortOrder"] = lookupNumber(data.SortOrder)
        filterNames = append(filterNames, "sort_order = "+data.SortOrder.ValueBigFloat().String())
    }
    if !data.DroppedCount.IsNull() && !data.DroppedCount.IsUnknown() {
        filters["droppedCount"] = lookupNumber(data.DroppedCount)
        filterNames = append(filterNames, "dropped_count = "+data.DroppedCount.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the trace drop filter up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the trace drop filter up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "filterQuery": true,
        "action": true,
        "samplePercentage": true,
        "isEnabled": true,
        "sortOrder": true,
        "droppedCount": true,
        "lastDroppedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/trace-drop-filter/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trace_drop_filter, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No trace drop filter found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read trace_drop_filter: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/trace-drop-filter/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list trace_drop_filter, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list trace_drop_filter: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No trace drop filter matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one trace drop filter matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for trace_drop_filter.")
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
    if obj, ok := item["filterQuery"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FilterQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FilterQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FilterQuery = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FilterQuery = types.StringValue(string(jsonBytes))
        } else {
            data.FilterQuery = types.StringNull()
        }
    } else if val, ok := item["filterQuery"].(string); ok {
        data.FilterQuery = types.StringValue(val)
    } else {
        data.FilterQuery = types.StringNull()
    }
    if obj, ok := item["action"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := item["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
    }
    if val, ok := item["samplePercentage"].(float64); ok {
        data.SamplePercentage = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["samplePercentage"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SamplePercentage = types.NumberValue(big.NewFloat(val))
        } else {
            data.SamplePercentage = types.NumberNull()
        }
    } else {
        data.SamplePercentage = types.NumberNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["sortOrder"].(float64); ok {
        data.SortOrder = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sortOrder"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SortOrder = types.NumberValue(big.NewFloat(val))
        } else {
            data.SortOrder = types.NumberNull()
        }
    } else {
        data.SortOrder = types.NumberNull()
    }
    if val, ok := item["droppedCount"].(float64); ok {
        data.DroppedCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["droppedCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DroppedCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.DroppedCount = types.NumberNull()
        }
    } else {
        data.DroppedCount = types.NumberNull()
    }
    if obj, ok := item["lastDroppedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastDroppedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastDroppedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastDroppedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastDroppedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastDroppedAt = types.StringNull()
        }
    } else if val, ok := item["lastDroppedAt"].(string); ok {
        data.LastDroppedAt = types.StringValue(val)
    } else {
        data.LastDroppedAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

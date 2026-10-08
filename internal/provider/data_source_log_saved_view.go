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
var _ datasource.DataSource = &LogSavedViewDataSource{}

func NewLogSavedViewDataSource() datasource.DataSource {
    return &LogSavedViewDataSource{}
}

// LogSavedViewDataSource defines the data source implementation.
type LogSavedViewDataSource struct {
    client *Client
}

// LogSavedViewDataSourceModel describes the data source data model.
type LogSavedViewDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Query types.String `tfsdk:"query"`
    Columns types.String `tfsdk:"columns"`
    SortField types.String `tfsdk:"sort_field"`
    SortOrder types.String `tfsdk:"sort_order"`
    PageSize types.Number `tfsdk:"page_size"`
    TimeRange types.String `tfsdk:"time_range"`
    IsDefault types.Bool `tfsdk:"is_default"`
}

func (d *LogSavedViewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log_saved_view"
}

func (d *LogSavedViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Save and reuse log explorer views, including the current filters, columns, sorting, and page size. Look up an existing log saved view by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `is_default`, ...): each one set must match, and exactly one log saved view may match them all.",

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
                MarkdownDescription: "ID of the project this saved log view belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this saved log view.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this saved log view. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "query": schema.StringAttribute{
                MarkdownDescription: "Serialized log query for this saved view. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "columns": schema.StringAttribute{
                MarkdownDescription: "Selected log table columns for this saved view. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "sort_field": schema.StringAttribute{
                MarkdownDescription: "Active sort field for this saved log view.",
                Optional: true,
                Computed: true,
            },
            "sort_order": schema.StringAttribute{
                MarkdownDescription: "Sort order for this saved log view.",
                Optional: true,
                Computed: true,
            },
            "page_size": schema.NumberAttribute{
                MarkdownDescription: "Number of logs per page for this saved view.",
                Optional: true,
                Computed: true,
            },
            "time_range": schema.StringAttribute{
                MarkdownDescription: "Time selection for this saved view — the rolling range token (e.g. Past 1 Hour), or an absolute window when the range is Custom. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_default": schema.BoolAttribute{
                MarkdownDescription: "Whether this saved log view should be applied by default.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *LogSavedViewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogSavedViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LogSavedViewDataSourceModel

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
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.SortField.IsNull() && !data.SortField.IsUnknown() {
        filters["sortField"] = data.SortField.ValueString()
        filterNames = append(filterNames, "sort_field = "+fmt.Sprintf("%q", data.SortField.ValueString()))
    }
    if !data.SortOrder.IsNull() && !data.SortOrder.IsUnknown() {
        filters["sortOrder"] = data.SortOrder.ValueString()
        filterNames = append(filterNames, "sort_order = "+fmt.Sprintf("%q", data.SortOrder.ValueString()))
    }
    if !data.PageSize.IsNull() && !data.PageSize.IsUnknown() {
        filters["pageSize"] = lookupNumber(data.PageSize)
        filterNames = append(filterNames, "page_size = "+data.PageSize.ValueBigFloat().String())
    }
    if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() {
        filters["isDefault"] = data.IsDefault.ValueBool()
        filterNames = append(filterNames, "is_default = "+fmt.Sprintf("%t", data.IsDefault.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the log saved view up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the log saved view up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "createdByUserId": true,
        "query": true,
        "columns": true,
        "sortField": true,
        "sortOrder": true,
        "pageSize": true,
        "timeRange": true,
        "isDefault": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/log-saved-view/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log_saved_view, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No log saved view found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read log_saved_view: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/log-saved-view/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list log_saved_view, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list log_saved_view: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No log saved view matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one log saved view matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for log_saved_view.")
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
    if obj, ok := item["query"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Query = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Query = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Query = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Query = types.StringValue(string(jsonBytes))
        } else {
            data.Query = types.StringNull()
        }
    } else if val, ok := item["query"].(string); ok {
        data.Query = types.StringValue(val)
    } else {
        data.Query = types.StringNull()
    }
    if obj, ok := item["columns"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Columns = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Columns = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Columns = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Columns = types.StringValue(string(jsonBytes))
        } else {
            data.Columns = types.StringNull()
        }
    } else if val, ok := item["columns"].(string); ok {
        data.Columns = types.StringValue(val)
    } else {
        data.Columns = types.StringNull()
    }
    if obj, ok := item["sortField"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SortField = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SortField = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SortField = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SortField = types.StringValue(string(jsonBytes))
        } else {
            data.SortField = types.StringNull()
        }
    } else if val, ok := item["sortField"].(string); ok {
        data.SortField = types.StringValue(val)
    } else {
        data.SortField = types.StringNull()
    }
    if obj, ok := item["sortOrder"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SortOrder = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SortOrder = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SortOrder = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SortOrder = types.StringValue(string(jsonBytes))
        } else {
            data.SortOrder = types.StringNull()
        }
    } else if val, ok := item["sortOrder"].(string); ok {
        data.SortOrder = types.StringValue(val)
    } else {
        data.SortOrder = types.StringNull()
    }
    if val, ok := item["pageSize"].(float64); ok {
        data.PageSize = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["pageSize"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PageSize = types.NumberValue(big.NewFloat(val))
        } else {
            data.PageSize = types.NumberNull()
        }
    } else {
        data.PageSize = types.NumberNull()
    }
    if obj, ok := item["timeRange"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TimeRange = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TimeRange = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TimeRange = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TimeRange = types.StringValue(string(jsonBytes))
        } else {
            data.TimeRange = types.StringNull()
        }
    } else if val, ok := item["timeRange"].(string); ok {
        data.TimeRange = types.StringValue(val)
    } else {
        data.TimeRange = types.StringNull()
    }
    if val, ok := item["isDefault"].(bool); ok {
        data.IsDefault = types.BoolValue(val)
    } else {
        data.IsDefault = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

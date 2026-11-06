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
var _ datasource.DataSource = &DataAssetLineageDataSource{}

func NewDataAssetLineageDataSource() datasource.DataSource {
    return &DataAssetLineageDataSource{}
}

// DataAssetLineageDataSource defines the data source implementation.
type DataAssetLineageDataSource struct {
    client *Client
}

// DataAssetLineageDataSourceModel describes the data source data model.
type DataAssetLineageDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    UpstreamDataAssetId types.String `tfsdk:"upstream_data_asset_id"`
    DownstreamDataAssetId types.String `tfsdk:"downstream_data_asset_id"`
    DataPipelineId types.String `tfsdk:"data_pipeline_id"`
    EdgeKey types.String `tfsdk:"edge_key"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    Source types.String `tfsdk:"source"`
    SeenBy types.String `tfsdk:"seen_by"`
    ColumnMappings types.String `tfsdk:"column_mappings"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataAssetLineageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_asset_lineage"
}

func (d *DataAssetLineageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "One table feeding another, from the pipelines that read the first and write the second, the database's query history, or a view's definition. Look up an existing data asset lineage by `id`, or by any of its other arguments (`created_by_user_id`, `data_pipeline_id`, `downstream_data_asset_id`, ...): each one set must match, and exactly one data asset lineage may match them all.",

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
            "upstream_data_asset_id": schema.StringAttribute{
                MarkdownDescription: "ID of the table the data comes from; empty for a pipeline that writes the table from data OneUptime does not watch. The ID of a `oneuptime_data_asset`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_asset"),
            },
            "downstream_data_asset_id": schema.StringAttribute{
                MarkdownDescription: "ID of the table built from it. The ID of a `oneuptime_data_asset`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_asset"),
            },
            "data_pipeline_id": schema.StringAttribute{
                MarkdownDescription: "ID of the pipeline that moves the data. The ID of a `oneuptime_data_pipeline` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_data_pipeline"),
            },
            "edge_key": schema.StringAttribute{
                MarkdownDescription: "The upstream table, the downstream table and the pipeline together: one row each.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the edge was first seen.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a run or query last moved data along the edge, or a view was last seen defined on it.",
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Where the edge was seen: a pipeline's run, the database's query history, or a view's definition. Empty means a pipeline's run.",
                Optional: true,
                Computed: true,
            },
            "seen_by": schema.StringAttribute{
                MarkdownDescription: "For an edge from query history: a login or service account whose queries moved the data.",
                Optional: true,
                Computed: true,
            },
            "column_mappings": schema.StringAttribute{
                MarkdownDescription: "Which columns of the upstream table feed which columns of the downstream one, when its source names them. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *DataAssetLineageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataAssetLineageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataAssetLineageDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.UpstreamDataAssetId.IsNull() && !data.UpstreamDataAssetId.IsUnknown() {
        filters["upstreamDataAssetId"] = data.UpstreamDataAssetId.ValueString()
        filterNames = append(filterNames, "upstream_data_asset_id = "+fmt.Sprintf("%q", data.UpstreamDataAssetId.ValueString()))
    }
    if !data.DownstreamDataAssetId.IsNull() && !data.DownstreamDataAssetId.IsUnknown() {
        filters["downstreamDataAssetId"] = data.DownstreamDataAssetId.ValueString()
        filterNames = append(filterNames, "downstream_data_asset_id = "+fmt.Sprintf("%q", data.DownstreamDataAssetId.ValueString()))
    }
    if !data.DataPipelineId.IsNull() && !data.DataPipelineId.IsUnknown() {
        filters["dataPipelineId"] = data.DataPipelineId.ValueString()
        filterNames = append(filterNames, "data_pipeline_id = "+fmt.Sprintf("%q", data.DataPipelineId.ValueString()))
    }
    if !data.EdgeKey.IsNull() && !data.EdgeKey.IsUnknown() {
        filters["edgeKey"] = data.EdgeKey.ValueString()
        filterNames = append(filterNames, "edge_key = "+fmt.Sprintf("%q", data.EdgeKey.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.SeenBy.IsNull() && !data.SeenBy.IsUnknown() {
        filters["seenBy"] = data.SeenBy.ValueString()
        filterNames = append(filterNames, "seen_by = "+fmt.Sprintf("%q", data.SeenBy.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data asset lineage up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data asset lineage up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "upstreamDataAssetId": true,
        "downstreamDataAssetId": true,
        "dataPipelineId": true,
        "edgeKey": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "source": true,
        "seenBy": true,
        "columnMappings": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-asset-lineage/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_asset_lineage, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data asset lineage found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_asset_lineage: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-asset-lineage/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_asset_lineage, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_asset_lineage: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data asset lineage matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data asset lineage matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_asset_lineage.")
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
    if obj, ok := item["upstreamDataAssetId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpstreamDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpstreamDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpstreamDataAssetId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpstreamDataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.UpstreamDataAssetId = types.StringNull()
        }
    } else if val, ok := item["upstreamDataAssetId"].(string); ok {
        data.UpstreamDataAssetId = types.StringValue(val)
    } else {
        data.UpstreamDataAssetId = types.StringNull()
    }
    if obj, ok := item["downstreamDataAssetId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DownstreamDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DownstreamDataAssetId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DownstreamDataAssetId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DownstreamDataAssetId = types.StringValue(string(jsonBytes))
        } else {
            data.DownstreamDataAssetId = types.StringNull()
        }
    } else if val, ok := item["downstreamDataAssetId"].(string); ok {
        data.DownstreamDataAssetId = types.StringValue(val)
    } else {
        data.DownstreamDataAssetId = types.StringNull()
    }
    if obj, ok := item["dataPipelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataPipelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataPipelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataPipelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataPipelineId = types.StringValue(string(jsonBytes))
        } else {
            data.DataPipelineId = types.StringNull()
        }
    } else if val, ok := item["dataPipelineId"].(string); ok {
        data.DataPipelineId = types.StringValue(val)
    } else {
        data.DataPipelineId = types.StringNull()
    }
    if obj, ok := item["edgeKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EdgeKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EdgeKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EdgeKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EdgeKey = types.StringValue(string(jsonBytes))
        } else {
            data.EdgeKey = types.StringNull()
        }
    } else if val, ok := item["edgeKey"].(string); ok {
        data.EdgeKey = types.StringValue(val)
    } else {
        data.EdgeKey = types.StringNull()
    }
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
    if obj, ok := item["source"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := item["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := item["seenBy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeenBy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SeenBy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SeenBy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SeenBy = types.StringValue(string(jsonBytes))
        } else {
            data.SeenBy = types.StringNull()
        }
    } else if val, ok := item["seenBy"].(string); ok {
        data.SeenBy = types.StringValue(val)
    } else {
        data.SeenBy = types.StringNull()
    }
    if obj, ok := item["columnMappings"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ColumnMappings = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ColumnMappings = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ColumnMappings = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ColumnMappings = types.StringValue(string(jsonBytes))
        } else {
            data.ColumnMappings = types.StringNull()
        }
    } else if val, ok := item["columnMappings"].(string); ok {
        data.ColumnMappings = types.StringValue(val)
    } else {
        data.ColumnMappings = types.StringNull()
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

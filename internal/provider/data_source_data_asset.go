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
var _ datasource.DataSource = &DataAssetDataSource{}

func NewDataAssetDataSource() datasource.DataSource {
    return &DataAssetDataSource{}
}

// DataAssetDataSource defines the data source implementation.
type DataAssetDataSource struct {
    client *Client
}

// DataAssetDataSourceModel describes the data source data model.
type DataAssetDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    DataStoreName types.String `tfsdk:"data_store_name"`
    DataStoreKind types.String `tfsdk:"data_store_kind"`
    AssetKey types.String `tfsdk:"asset_key"`
    SchemaName types.String `tfsdk:"schema_name"`
    Name types.String `tfsdk:"name"`
    AssetKind types.String `tfsdk:"asset_kind"`
    Description types.String `tfsdk:"description"`
    Labels types.Set `tfsdk:"labels"`
    RowCountChangeTolerancePercent types.Number `tfsdk:"row_count_change_tolerance_percent"`
    FreshnessColumn types.String `tfsdk:"freshness_column"`
    Columns types.String `tfsdk:"columns"`
    DbtNode types.String `tfsdk:"dbt_node"`
    ColumnCount types.Number `tfsdk:"column_count"`
    RowCount types.Number `tfsdk:"row_count"`
    SizeInBytes types.Number `tfsdk:"size_in_bytes"`
    WriteCredits types.Number `tfsdk:"write_credits"`
    WriteBytesBilled types.Number `tfsdk:"write_bytes_billed"`
    LastChangedAt types.String `tfsdk:"last_changed_at"`
    LastReadAt types.String `tfsdk:"last_read_at"`
    UnreadSince types.String `tfsdk:"unread_since"`
    UpdateIntervalInSeconds types.Number `tfsdk:"update_interval_in_seconds"`
    LastCollectedAt types.String `tfsdk:"last_collected_at"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    IsDropped types.Bool `tfsdk:"is_dropped"`
    DroppedAt types.String `tfsdk:"dropped_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataAssetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_asset"
}

func (d *DataAssetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Every table and view of your connected data stores: its columns, how many rows it holds, how big it is and when it last changed. Look up an existing data asset by `id`, or by any of its other arguments (`name`, `asset_key`, `asset_kind`, ...): each one set must match, and exactly one data asset may match them all.",

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
                MarkdownDescription: "ID of the project this data asset belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "ID of the data store this asset is in. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_data_store"),
            },
            "data_store_name": schema.StringAttribute{
                MarkdownDescription: "The name of the data store this asset is in.",
                Optional: true,
                Computed: true,
            },
            "data_store_kind": schema.StringAttribute{
                MarkdownDescription: "What kind of store this asset is in, such as PostgreSQL or Snowflake.",
                Optional: true,
                Computed: true,
            },
            "asset_key": schema.StringAttribute{
                MarkdownDescription: "The asset's name with its schema, such as \"public.orders\": what its metrics and checks name it by.",
                Optional: true,
                Computed: true,
            },
            "schema_name": schema.StringAttribute{
                MarkdownDescription: "The schema (or database, or dataset) the asset is in, such as public.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The table's or view's own name, such as orders.",
                Optional: true,
                Computed: true,
            },
            "asset_kind": schema.StringAttribute{
                MarkdownDescription: "Whether it is a table, a view, a materialized view, a partitioned table or a foreign table.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the asset holds and who relies on it, in your own words.",
                Optional: true,
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter tables by. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "row_count_change_tolerance_percent": schema.NumberAttribute{
                MarkdownDescription: "A change in rows this small counts as none here, raised each time such an alert is marked not an issue.",
                Optional: true,
                Computed: true,
            },
            "freshness_column": schema.StringAttribute{
                MarkdownDescription: "A timestamp column that says when a row was written. Its latest value says when the table last changed, where the database's own counters cannot.",
                Optional: true,
                Computed: true,
            },
            "columns": schema.StringAttribute{
                MarkdownDescription: "The asset's columns at its last collection: each one's name, type, whether it allows empty values, and its position. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "dbt_node": schema.StringAttribute{
                MarkdownDescription: "The dbt model, seed, snapshot or source that is this table, from the project's last uploaded artifacts: its id, materialization, tags, owner, and how its last build went. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "column_count": schema.NumberAttribute{
                MarkdownDescription: "How many columns the asset has.",
                Optional: true,
                Computed: true,
            },
            "row_count": schema.NumberAttribute{
                MarkdownDescription: "How many rows the asset holds, from the engine's own statistics at the last collection.",
                Optional: true,
                Computed: true,
            },
            "size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "How much space the asset and its indexes take.",
                Optional: true,
                Computed: true,
            },
            "write_credits": schema.NumberAttribute{
                MarkdownDescription: "Snowflake credits the jobs that write the table spent in the last day, a job's shared among the tables it writes.",
                Optional: true,
                Computed: true,
            },
            "write_bytes_billed": schema.NumberAttribute{
                MarkdownDescription: "BigQuery bytes billed for the jobs that write the table in the last day, a job's shared among the tables it writes.",
                Optional: true,
                Computed: true,
            },
            "last_changed_at": schema.StringAttribute{
                MarkdownDescription: "When a row was last inserted, updated or deleted, as near as the engine tells.",
                Computed: true,
            },
            "last_read_at": schema.StringAttribute{
                MarkdownDescription: "When the table was last seen read, to the hour: by the engine's read counters, or by a service's database spans.",
                Computed: true,
            },
            "unread_since": schema.StringAttribute{
                MarkdownDescription: "The time from which the table is known to have had no read: its last read, or when the engine's read counters started when it has had none since. Empty when that cannot be said.",
                Computed: true,
            },
            "update_interval_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How often the asset usually changes, learned from the changes OneUptime has seen.",
                Optional: true,
                Computed: true,
            },
            "last_collected_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last read the asset's health.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime first found the asset.",
                Computed: true,
            },
            "is_dropped": schema.BoolAttribute{
                MarkdownDescription: "The asset is gone from its store. It is kept with its history, and comes back if the table does.",
                Optional: true,
                Computed: true,
            },
            "dropped_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime saw that the asset was gone.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this asset, if a person did rather than OneUptime. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *DataAssetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataAssetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataAssetDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DataStoreId.IsNull() && !data.DataStoreId.IsUnknown() {
        filters["dataStoreId"] = data.DataStoreId.ValueString()
        filterNames = append(filterNames, "data_store_id = "+fmt.Sprintf("%q", data.DataStoreId.ValueString()))
    }
    if !data.DataStoreName.IsNull() && !data.DataStoreName.IsUnknown() {
        filters["dataStoreName"] = data.DataStoreName.ValueString()
        filterNames = append(filterNames, "data_store_name = "+fmt.Sprintf("%q", data.DataStoreName.ValueString()))
    }
    if !data.DataStoreKind.IsNull() && !data.DataStoreKind.IsUnknown() {
        filters["dataStoreKind"] = data.DataStoreKind.ValueString()
        filterNames = append(filterNames, "data_store_kind = "+fmt.Sprintf("%q", data.DataStoreKind.ValueString()))
    }
    if !data.AssetKey.IsNull() && !data.AssetKey.IsUnknown() {
        filters["assetKey"] = data.AssetKey.ValueString()
        filterNames = append(filterNames, "asset_key = "+fmt.Sprintf("%q", data.AssetKey.ValueString()))
    }
    if !data.SchemaName.IsNull() && !data.SchemaName.IsUnknown() {
        filters["schemaName"] = data.SchemaName.ValueString()
        filterNames = append(filterNames, "schema_name = "+fmt.Sprintf("%q", data.SchemaName.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.AssetKind.IsNull() && !data.AssetKind.IsUnknown() {
        filters["assetKind"] = data.AssetKind.ValueString()
        filterNames = append(filterNames, "asset_kind = "+fmt.Sprintf("%q", data.AssetKind.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.RowCountChangeTolerancePercent.IsNull() && !data.RowCountChangeTolerancePercent.IsUnknown() {
        filters["rowCountChangeTolerancePercent"] = lookupNumber(data.RowCountChangeTolerancePercent)
        filterNames = append(filterNames, "row_count_change_tolerance_percent = "+data.RowCountChangeTolerancePercent.ValueBigFloat().String())
    }
    if !data.FreshnessColumn.IsNull() && !data.FreshnessColumn.IsUnknown() {
        filters["freshnessColumn"] = data.FreshnessColumn.ValueString()
        filterNames = append(filterNames, "freshness_column = "+fmt.Sprintf("%q", data.FreshnessColumn.ValueString()))
    }
    if !data.ColumnCount.IsNull() && !data.ColumnCount.IsUnknown() {
        filters["columnCount"] = lookupNumber(data.ColumnCount)
        filterNames = append(filterNames, "column_count = "+data.ColumnCount.ValueBigFloat().String())
    }
    if !data.RowCount.IsNull() && !data.RowCount.IsUnknown() {
        filters["rowCount"] = lookupNumber(data.RowCount)
        filterNames = append(filterNames, "row_count = "+data.RowCount.ValueBigFloat().String())
    }
    if !data.SizeInBytes.IsNull() && !data.SizeInBytes.IsUnknown() {
        filters["sizeInBytes"] = lookupNumber(data.SizeInBytes)
        filterNames = append(filterNames, "size_in_bytes = "+data.SizeInBytes.ValueBigFloat().String())
    }
    if !data.WriteCredits.IsNull() && !data.WriteCredits.IsUnknown() {
        filters["writeCredits"] = lookupNumber(data.WriteCredits)
        filterNames = append(filterNames, "write_credits = "+data.WriteCredits.ValueBigFloat().String())
    }
    if !data.WriteBytesBilled.IsNull() && !data.WriteBytesBilled.IsUnknown() {
        filters["writeBytesBilled"] = lookupNumber(data.WriteBytesBilled)
        filterNames = append(filterNames, "write_bytes_billed = "+data.WriteBytesBilled.ValueBigFloat().String())
    }
    if !data.UpdateIntervalInSeconds.IsNull() && !data.UpdateIntervalInSeconds.IsUnknown() {
        filters["updateIntervalInSeconds"] = lookupNumber(data.UpdateIntervalInSeconds)
        filterNames = append(filterNames, "update_interval_in_seconds = "+data.UpdateIntervalInSeconds.ValueBigFloat().String())
    }
    if !data.IsDropped.IsNull() && !data.IsDropped.IsUnknown() {
        filters["isDropped"] = data.IsDropped.ValueBool()
        filterNames = append(filterNames, "is_dropped = "+fmt.Sprintf("%t", data.IsDropped.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data asset up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data asset up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataStoreId": true,
        "dataStoreName": true,
        "dataStoreKind": true,
        "assetKey": true,
        "schemaName": true,
        "name": true,
        "assetKind": true,
        "description": true,
        "labels": true,
        "rowCountChangeTolerancePercent": true,
        "freshnessColumn": true,
        "columns": true,
        "dbtNode": true,
        "columnCount": true,
        "rowCount": true,
        "sizeInBytes": true,
        "writeCredits": true,
        "writeBytesBilled": true,
        "lastChangedAt": true,
        "lastReadAt": true,
        "unreadSince": true,
        "updateIntervalInSeconds": true,
        "lastCollectedAt": true,
        "firstSeenAt": true,
        "isDropped": true,
        "droppedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-asset/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_asset, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data asset found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_asset: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-asset/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_asset, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_asset: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data asset matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data asset matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_asset.")
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
    if obj, ok := item["dataStoreId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataStoreId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataStoreId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataStoreId = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreId = types.StringNull()
        }
    } else if val, ok := item["dataStoreId"].(string); ok {
        data.DataStoreId = types.StringValue(val)
    } else {
        data.DataStoreId = types.StringNull()
    }
    if obj, ok := item["dataStoreName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataStoreName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataStoreName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataStoreName = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreName = types.StringNull()
        }
    } else if val, ok := item["dataStoreName"].(string); ok {
        data.DataStoreName = types.StringValue(val)
    } else {
        data.DataStoreName = types.StringNull()
    }
    if obj, ok := item["dataStoreKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataStoreKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataStoreKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataStoreKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataStoreKind = types.StringValue(string(jsonBytes))
        } else {
            data.DataStoreKind = types.StringNull()
        }
    } else if val, ok := item["dataStoreKind"].(string); ok {
        data.DataStoreKind = types.StringValue(val)
    } else {
        data.DataStoreKind = types.StringNull()
    }
    if obj, ok := item["assetKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssetKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssetKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssetKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssetKey = types.StringValue(string(jsonBytes))
        } else {
            data.AssetKey = types.StringNull()
        }
    } else if val, ok := item["assetKey"].(string); ok {
        data.AssetKey = types.StringValue(val)
    } else {
        data.AssetKey = types.StringNull()
    }
    if obj, ok := item["schemaName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SchemaName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SchemaName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SchemaName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SchemaName = types.StringValue(string(jsonBytes))
        } else {
            data.SchemaName = types.StringNull()
        }
    } else if val, ok := item["schemaName"].(string); ok {
        data.SchemaName = types.StringValue(val)
    } else {
        data.SchemaName = types.StringNull()
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
    if obj, ok := item["assetKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssetKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssetKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssetKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssetKind = types.StringValue(string(jsonBytes))
        } else {
            data.AssetKind = types.StringNull()
        }
    } else if val, ok := item["assetKind"].(string); ok {
        data.AssetKind = types.StringValue(val)
    } else {
        data.AssetKind = types.StringNull()
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
    data.RowCountChangeTolerancePercent = numberFromAPI(data.RowCountChangeTolerancePercent, item["rowCountChangeTolerancePercent"])
    if obj, ok := item["freshnessColumn"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FreshnessColumn = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FreshnessColumn = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FreshnessColumn = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FreshnessColumn = types.StringValue(string(jsonBytes))
        } else {
            data.FreshnessColumn = types.StringNull()
        }
    } else if val, ok := item["freshnessColumn"].(string); ok {
        data.FreshnessColumn = types.StringValue(val)
    } else {
        data.FreshnessColumn = types.StringNull()
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
    if obj, ok := item["dbtNode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DbtNode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DbtNode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DbtNode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DbtNode = types.StringValue(string(jsonBytes))
        } else {
            data.DbtNode = types.StringNull()
        }
    } else if val, ok := item["dbtNode"].(string); ok {
        data.DbtNode = types.StringValue(val)
    } else {
        data.DbtNode = types.StringNull()
    }
    data.ColumnCount = numberFromAPI(data.ColumnCount, item["columnCount"])
    data.RowCount = numberFromAPI(data.RowCount, item["rowCount"])
    data.SizeInBytes = numberFromAPI(data.SizeInBytes, item["sizeInBytes"])
    data.WriteCredits = numberFromAPI(data.WriteCredits, item["writeCredits"])
    data.WriteBytesBilled = numberFromAPI(data.WriteBytesBilled, item["writeBytesBilled"])
    if obj, ok := item["lastChangedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastChangedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastChangedAt = types.StringNull()
        }
    } else if val, ok := item["lastChangedAt"].(string); ok {
        data.LastChangedAt = types.StringValue(val)
    } else {
        data.LastChangedAt = types.StringNull()
    }
    if obj, ok := item["lastReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastReadAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastReadAt = types.StringNull()
        }
    } else if val, ok := item["lastReadAt"].(string); ok {
        data.LastReadAt = types.StringValue(val)
    } else {
        data.LastReadAt = types.StringNull()
    }
    if obj, ok := item["unreadSince"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UnreadSince = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UnreadSince = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UnreadSince = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UnreadSince = types.StringValue(string(jsonBytes))
        } else {
            data.UnreadSince = types.StringNull()
        }
    } else if val, ok := item["unreadSince"].(string); ok {
        data.UnreadSince = types.StringValue(val)
    } else {
        data.UnreadSince = types.StringNull()
    }
    data.UpdateIntervalInSeconds = numberFromAPI(data.UpdateIntervalInSeconds, item["updateIntervalInSeconds"])
    if obj, ok := item["lastCollectedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastCollectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastCollectedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastCollectedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastCollectedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastCollectedAt = types.StringNull()
        }
    } else if val, ok := item["lastCollectedAt"].(string); ok {
        data.LastCollectedAt = types.StringValue(val)
    } else {
        data.LastCollectedAt = types.StringNull()
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
    if val, ok := item["isDropped"].(bool); ok {
        data.IsDropped = types.BoolValue(val)
    } else {
        data.IsDropped = types.BoolNull()
    }
    if obj, ok := item["droppedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DroppedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DroppedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DroppedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DroppedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DroppedAt = types.StringNull()
        }
    } else if val, ok := item["droppedAt"].(string); ok {
        data.DroppedAt = types.StringValue(val)
    } else {
        data.DroppedAt = types.StringNull()
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

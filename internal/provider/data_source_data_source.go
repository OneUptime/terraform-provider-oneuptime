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
var _ datasource.DataSource = &DataSourceDataSource{}

func NewDataSourceDataSource() datasource.DataSource {
    return &DataSourceDataSource{}
}

// DataSourceDataSource defines the data source implementation.
type DataSourceDataSource struct {
    client *Client
}

// DataSourceDataSourceModel describes the data source data model.
type DataSourceDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    DataSourceType types.String `tfsdk:"data_source_type"`
    Url types.String `tfsdk:"url"`
    DatabaseHost types.String `tfsdk:"database_host"`
    DatabasePort types.Number `tfsdk:"database_port"`
    DatabaseName types.String `tfsdk:"database_name"`
    Username types.String `tfsdk:"username"`
    AdditionalOptions types.String `tfsdk:"additional_options"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataSourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_source"
}

func (d *DataSourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Connect external systems — Prometheus, SQL databases, ClickHouse, Loki, Elasticsearch, or REST APIs — and build dashboards on their data. Look up an existing data source by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `data_source_type`, ...): each one set must match, and exactly one data source may match them all.",

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
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "data_source_type": schema.StringAttribute{
                MarkdownDescription: "The kind of external system this data source connects to (Prometheus, PostgreSQL, MySQL, Microsoft SQL Server, ClickHouse, Loki, Elasticsearch, or REST API).",
                Optional: true,
                Computed: true,
            },
            "url": schema.StringAttribute{
                MarkdownDescription: "Base URL for HTTP-based sources (Prometheus, Loki, Elasticsearch, REST API) — e.g. https://prometheus.example.com.",
                Optional: true,
                Computed: true,
            },
            "database_host": schema.StringAttribute{
                MarkdownDescription: "Hostname or IP address for database sources (PostgreSQL, MySQL, SQL Server, ClickHouse).",
                Optional: true,
                Computed: true,
            },
            "database_port": schema.NumberAttribute{
                MarkdownDescription: "Port for database sources. Defaults per engine: PostgreSQL 5432, MySQL 3306, SQL Server 1433, ClickHouse 8123.",
                Optional: true,
                Computed: true,
            },
            "database_name": schema.StringAttribute{
                MarkdownDescription: "Database (or ClickHouse database) to connect to.",
                Optional: true,
                Computed: true,
            },
            "username": schema.StringAttribute{
                MarkdownDescription: "Username for database sources, or HTTP basic-auth username for HTTP sources. Use a READ-ONLY account — dashboards only ever read.",
                Optional: true,
                Computed: true,
            },
            "additional_options": schema.StringAttribute{
                MarkdownDescription: "Per-type options that are not secrets — e.g. { \"sslEnabled\": true, \"elasticsearchIndex\": \"logs-*\" }. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DataSourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataSourceDataSourceModel

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
    if !data.DataSourceType.IsNull() && !data.DataSourceType.IsUnknown() {
        filters["dataSourceType"] = data.DataSourceType.ValueString()
        filterNames = append(filterNames, "data_source_type = "+fmt.Sprintf("%q", data.DataSourceType.ValueString()))
    }
    if !data.Url.IsNull() && !data.Url.IsUnknown() {
        filters["url"] = data.Url.ValueString()
        filterNames = append(filterNames, "url = "+fmt.Sprintf("%q", data.Url.ValueString()))
    }
    if !data.DatabaseHost.IsNull() && !data.DatabaseHost.IsUnknown() {
        filters["databaseHost"] = data.DatabaseHost.ValueString()
        filterNames = append(filterNames, "database_host = "+fmt.Sprintf("%q", data.DatabaseHost.ValueString()))
    }
    if !data.DatabasePort.IsNull() && !data.DatabasePort.IsUnknown() {
        filters["databasePort"] = lookupNumber(data.DatabasePort)
        filterNames = append(filterNames, "database_port = "+data.DatabasePort.ValueBigFloat().String())
    }
    if !data.DatabaseName.IsNull() && !data.DatabaseName.IsUnknown() {
        filters["databaseName"] = data.DatabaseName.ValueString()
        filterNames = append(filterNames, "database_name = "+fmt.Sprintf("%q", data.DatabaseName.ValueString()))
    }
    if !data.Username.IsNull() && !data.Username.IsUnknown() {
        filters["username"] = data.Username.ValueString()
        filterNames = append(filterNames, "username = "+fmt.Sprintf("%q", data.Username.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data source up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data source up by.",
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
        "dataSourceType": true,
        "url": true,
        "databaseHost": true,
        "databasePort": true,
        "databaseName": true,
        "username": true,
        "additionalOptions": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-source/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_source, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data source found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_source: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-source/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_source, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_source: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data source matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data source matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_source.")
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
    if obj, ok := item["dataSourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataSourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataSourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataSourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataSourceType = types.StringValue(string(jsonBytes))
        } else {
            data.DataSourceType = types.StringNull()
        }
    } else if val, ok := item["dataSourceType"].(string); ok {
        data.DataSourceType = types.StringValue(val)
    } else {
        data.DataSourceType = types.StringNull()
    }
    if obj, ok := item["url"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Url = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Url = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Url = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Url = types.StringValue(string(jsonBytes))
        } else {
            data.Url = types.StringNull()
        }
    } else if val, ok := item["url"].(string); ok {
        data.Url = types.StringValue(val)
    } else {
        data.Url = types.StringNull()
    }
    if obj, ok := item["databaseHost"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseHost = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DatabaseHost = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DatabaseHost = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DatabaseHost = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseHost = types.StringNull()
        }
    } else if val, ok := item["databaseHost"].(string); ok {
        data.DatabaseHost = types.StringValue(val)
    } else {
        data.DatabaseHost = types.StringNull()
    }
    if val, ok := item["databasePort"].(float64); ok {
        data.DatabasePort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["databasePort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DatabasePort = types.NumberValue(big.NewFloat(val))
        } else {
            data.DatabasePort = types.NumberNull()
        }
    } else {
        data.DatabasePort = types.NumberNull()
    }
    if obj, ok := item["databaseName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DatabaseName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DatabaseName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DatabaseName = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseName = types.StringNull()
        }
    } else if val, ok := item["databaseName"].(string); ok {
        data.DatabaseName = types.StringValue(val)
    } else {
        data.DatabaseName = types.StringNull()
    }
    if obj, ok := item["username"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if val, ok := item["username"].(string); ok {
        data.Username = types.StringValue(val)
    } else {
        data.Username = types.StringNull()
    }
    if obj, ok := item["additionalOptions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AdditionalOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AdditionalOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AdditionalOptions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AdditionalOptions = types.StringValue(string(jsonBytes))
        } else {
            data.AdditionalOptions = types.StringNull()
        }
    } else if val, ok := item["additionalOptions"].(string); ok {
        data.AdditionalOptions = types.StringValue(val)
    } else {
        data.AdditionalOptions = types.StringNull()
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

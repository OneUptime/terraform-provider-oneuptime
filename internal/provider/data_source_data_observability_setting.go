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
var _ datasource.DataSource = &DataObservabilitySettingDataSource{}

func NewDataObservabilitySettingDataSource() datasource.DataSource {
    return &DataObservabilitySettingDataSource{}
}

// DataObservabilitySettingDataSource defines the data source implementation.
type DataObservabilitySettingDataSource struct {
    client *Client
}

// DataObservabilitySettingDataSourceModel describes the data source data model.
type DataObservabilitySettingDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    MaxGigabytesScannedPerCheck types.Number `tfsdk:"max_gigabytes_scanned_per_check"`
    WaitWhenDatabaseIsBusy types.Bool `tfsdk:"wait_when_database_is_busy"`
    KeepFailingRows types.Bool `tfsdk:"keep_failing_rows"`
}

func (d *DataObservabilitySettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_observability_setting"
}

func (d *DataObservabilitySettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept. Look up an existing data observability setting by `id`, or by any of its other arguments (`keep_failing_rows`, `max_gigabytes_scanned_per_check`, `wait_when_database_is_busy`): each one set must match, and exactly one data observability setting may match them all.",

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
                MarkdownDescription: "ID of the project these settings belong to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "max_gigabytes_scanned_per_check": schema.NumberAttribute{
                MarkdownDescription: "BigQuery and Athena bill by bytes read. A check that would read more is stopped, and says so.",
                Optional: true,
                Computed: true,
            },
            "wait_when_database_is_busy": schema.BoolAttribute{
                MarkdownDescription: "Checks wait for the next collection when a database is close to its connection limit or has many statements running.",
                Optional: true,
                Computed: true,
            },
            "keep_failing_rows": schema.BoolAttribute{
                MarkdownDescription: "A failing check keeps up to 20 of the rows that broke it, with sensitive columns masked. They are stored in OneUptime. Turning this off deletes them.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DataObservabilitySettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataObservabilitySettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataObservabilitySettingDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.MaxGigabytesScannedPerCheck.IsNull() && !data.MaxGigabytesScannedPerCheck.IsUnknown() {
        filters["maxGigabytesScannedPerCheck"] = lookupNumber(data.MaxGigabytesScannedPerCheck)
        filterNames = append(filterNames, "max_gigabytes_scanned_per_check = "+data.MaxGigabytesScannedPerCheck.ValueBigFloat().String())
    }
    if !data.WaitWhenDatabaseIsBusy.IsNull() && !data.WaitWhenDatabaseIsBusy.IsUnknown() {
        filters["waitWhenDatabaseIsBusy"] = data.WaitWhenDatabaseIsBusy.ValueBool()
        filterNames = append(filterNames, "wait_when_database_is_busy = "+fmt.Sprintf("%t", data.WaitWhenDatabaseIsBusy.ValueBool()))
    }
    if !data.KeepFailingRows.IsNull() && !data.KeepFailingRows.IsUnknown() {
        filters["keepFailingRows"] = data.KeepFailingRows.ValueBool()
        filterNames = append(filterNames, "keep_failing_rows = "+fmt.Sprintf("%t", data.KeepFailingRows.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data observability setting up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data observability setting up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "maxGigabytesScannedPerCheck": true,
        "waitWhenDatabaseIsBusy": true,
        "keepFailingRows": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-observability-setting/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_observability_setting, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data observability setting found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_observability_setting: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-observability-setting/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_observability_setting, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_observability_setting: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data observability setting matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data observability setting matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_observability_setting.")
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
    data.MaxGigabytesScannedPerCheck = numberFromAPI(data.MaxGigabytesScannedPerCheck, item["maxGigabytesScannedPerCheck"])
    if val, ok := item["waitWhenDatabaseIsBusy"].(bool); ok {
        data.WaitWhenDatabaseIsBusy = types.BoolValue(val)
    } else {
        data.WaitWhenDatabaseIsBusy = types.BoolNull()
    }
    if val, ok := item["keepFailingRows"].(bool); ok {
        data.KeepFailingRows = types.BoolValue(val)
    } else {
        data.KeepFailingRows = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

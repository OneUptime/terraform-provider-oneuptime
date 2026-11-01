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
var _ datasource.DataSource = &DataSecurityScanDataSource{}

func NewDataSecurityScanDataSource() datasource.DataSource {
    return &DataSecurityScanDataSource{}
}

// DataSecurityScanDataSource defines the data source implementation.
type DataSecurityScanDataSource struct {
    client *Client
}

// DataSecurityScanDataSourceModel describes the data source data model.
type DataSecurityScanDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DataStoreId types.String `tfsdk:"data_store_id"`
    CloudConnectionId types.String `tfsdk:"cloud_connection_id"`
    ProbeId types.String `tfsdk:"probe_id"`
    Kind types.String `tfsdk:"kind"`
    Status types.String `tfsdk:"status"`
    StartedAt types.String `tfsdk:"started_at"`
    FinishedAt types.String `tfsdk:"finished_at"`
    DurationInMs types.Number `tfsdk:"duration_in_ms"`
    StoresRead types.Number `tfsdk:"stores_read"`
    ObjectsRead types.Number `tfsdk:"objects_read"`
    RowsRead types.Number `tfsdk:"rows_read"`
    BytesRead types.Number `tfsdk:"bytes_read"`
    FindingsFound types.Number `tfsdk:"findings_found"`
    ReadLog types.String `tfsdk:"read_log"`
    Warnings types.String `tfsdk:"warnings"`
    ErrorMessage types.String `tfsdk:"error_message"`
}

func (d *DataSecurityScanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_security_scan"
}

func (d *DataSecurityScanDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Every scan Data Security ran: what it read, how much, what it found, how long it took and what went wrong. Look up an existing data security scan by `id`, or by any of its other arguments (`bytes_read`, `cloud_connection_id`, `data_store_id`, ...): each one set must match, and exactly one data security scan may match them all.",

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
                MarkdownDescription: "ID of the project this scan belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "data_store_id": schema.StringAttribute{
                MarkdownDescription: "The store scanned, for a content or telemetry scan. The ID of a `oneuptime_data_store`.",
                Optional: true,
                Computed: true,
            },
            "cloud_connection_id": schema.StringAttribute{
                MarkdownDescription: "The connection read, for a configuration read. The ID of a `oneuptime_cloud_connection`.",
                Optional: true,
                Computed: true,
            },
            "probe_id": schema.StringAttribute{
                MarkdownDescription: "The probe that ran the scan, when one did. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What the scan read: Configuration, Content or Telemetry.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Running, Succeeded, Partly succeeded, Failed or Skipped.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the scan started.",
                Computed: true,
            },
            "finished_at": schema.StringAttribute{
                MarkdownDescription: "When the scan finished.",
                Computed: true,
            },
            "duration_in_ms": schema.NumberAttribute{
                MarkdownDescription: "How long the scan took.",
                Optional: true,
                Computed: true,
            },
            "stores_read": schema.NumberAttribute{
                MarkdownDescription: "How many stores a configuration read covered.",
                Optional: true,
                Computed: true,
            },
            "objects_read": schema.NumberAttribute{
                MarkdownDescription: "How many objects or files the scan sampled.",
                Optional: true,
                Computed: true,
            },
            "rows_read": schema.NumberAttribute{
                MarkdownDescription: "How many rows or records the scan sampled.",
                Optional: true,
                Computed: true,
            },
            "bytes_read": schema.NumberAttribute{
                MarkdownDescription: "How many bytes the scan read.",
                Optional: true,
                Computed: true,
            },
            "findings_found": schema.NumberAttribute{
                MarkdownDescription: "How many findings the scan made or updated.",
                Optional: true,
                Computed: true,
            },
            "read_log": schema.StringAttribute{
                MarkdownDescription: "The objects, tables or services read, by name. Never their contents. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "warnings": schema.StringAttribute{
                MarkdownDescription: "What the scan could not read, and why. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Why the scan failed, if it did.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DataSecurityScanDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSecurityScanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataSecurityScanDataSourceModel

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
    if !data.CloudConnectionId.IsNull() && !data.CloudConnectionId.IsUnknown() {
        filters["cloudConnectionId"] = data.CloudConnectionId.ValueString()
        filterNames = append(filterNames, "cloud_connection_id = "+fmt.Sprintf("%q", data.CloudConnectionId.ValueString()))
    }
    if !data.ProbeId.IsNull() && !data.ProbeId.IsUnknown() {
        filters["probeId"] = data.ProbeId.ValueString()
        filterNames = append(filterNames, "probe_id = "+fmt.Sprintf("%q", data.ProbeId.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.DurationInMs.IsNull() && !data.DurationInMs.IsUnknown() {
        filters["durationInMs"] = lookupNumber(data.DurationInMs)
        filterNames = append(filterNames, "duration_in_ms = "+data.DurationInMs.ValueBigFloat().String())
    }
    if !data.StoresRead.IsNull() && !data.StoresRead.IsUnknown() {
        filters["storesRead"] = lookupNumber(data.StoresRead)
        filterNames = append(filterNames, "stores_read = "+data.StoresRead.ValueBigFloat().String())
    }
    if !data.ObjectsRead.IsNull() && !data.ObjectsRead.IsUnknown() {
        filters["objectsRead"] = lookupNumber(data.ObjectsRead)
        filterNames = append(filterNames, "objects_read = "+data.ObjectsRead.ValueBigFloat().String())
    }
    if !data.RowsRead.IsNull() && !data.RowsRead.IsUnknown() {
        filters["rowsRead"] = lookupNumber(data.RowsRead)
        filterNames = append(filterNames, "rows_read = "+data.RowsRead.ValueBigFloat().String())
    }
    if !data.BytesRead.IsNull() && !data.BytesRead.IsUnknown() {
        filters["bytesRead"] = lookupNumber(data.BytesRead)
        filterNames = append(filterNames, "bytes_read = "+data.BytesRead.ValueBigFloat().String())
    }
    if !data.FindingsFound.IsNull() && !data.FindingsFound.IsUnknown() {
        filters["findingsFound"] = lookupNumber(data.FindingsFound)
        filterNames = append(filterNames, "findings_found = "+data.FindingsFound.ValueBigFloat().String())
    }
    if !data.ErrorMessage.IsNull() && !data.ErrorMessage.IsUnknown() {
        filters["errorMessage"] = data.ErrorMessage.ValueString()
        filterNames = append(filterNames, "error_message = "+fmt.Sprintf("%q", data.ErrorMessage.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data security scan up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data security scan up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "dataStoreId": true,
        "cloudConnectionId": true,
        "probeId": true,
        "kind": true,
        "status": true,
        "startedAt": true,
        "finishedAt": true,
        "durationInMs": true,
        "storesRead": true,
        "objectsRead": true,
        "rowsRead": true,
        "bytesRead": true,
        "findingsFound": true,
        "readLog": true,
        "warnings": true,
        "errorMessage": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-scan/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_scan, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security scan found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_security_scan: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-scan/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_security_scan, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_security_scan: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security scan matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data security scan matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_security_scan.")
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
    if obj, ok := item["cloudConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := item["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := item["probeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeId = types.StringNull()
        }
    } else if val, ok := item["probeId"].(string); ok {
        data.ProbeId = types.StringValue(val)
    } else {
        data.ProbeId = types.StringNull()
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
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartedAt = types.StringNull()
        }
    } else if val, ok := item["startedAt"].(string); ok {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if obj, ok := item["finishedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FinishedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FinishedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FinishedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FinishedAt = types.StringValue(string(jsonBytes))
        } else {
            data.FinishedAt = types.StringNull()
        }
    } else if val, ok := item["finishedAt"].(string); ok {
        data.FinishedAt = types.StringValue(val)
    } else {
        data.FinishedAt = types.StringNull()
    }
    data.DurationInMs = numberFromAPI(data.DurationInMs, item["durationInMs"])
    data.StoresRead = numberFromAPI(data.StoresRead, item["storesRead"])
    data.ObjectsRead = numberFromAPI(data.ObjectsRead, item["objectsRead"])
    data.RowsRead = numberFromAPI(data.RowsRead, item["rowsRead"])
    data.BytesRead = numberFromAPI(data.BytesRead, item["bytesRead"])
    data.FindingsFound = numberFromAPI(data.FindingsFound, item["findingsFound"])
    if obj, ok := item["readLog"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ReadLog = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ReadLog = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ReadLog = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ReadLog = types.StringValue(string(jsonBytes))
        } else {
            data.ReadLog = types.StringNull()
        }
    } else if val, ok := item["readLog"].(string); ok {
        data.ReadLog = types.StringValue(val)
    } else {
        data.ReadLog = types.StringNull()
    }
    if obj, ok := item["warnings"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Warnings = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Warnings = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Warnings = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Warnings = types.StringValue(string(jsonBytes))
        } else {
            data.Warnings = types.StringNull()
        }
    } else if val, ok := item["warnings"].(string); ok {
        data.Warnings = types.StringValue(val)
    } else {
        data.Warnings = types.StringNull()
    }
    if obj, ok := item["errorMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ErrorMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ErrorMessage = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorMessage = types.StringNull()
        }
    } else if val, ok := item["errorMessage"].(string); ok {
        data.ErrorMessage = types.StringValue(val)
    } else {
        data.ErrorMessage = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

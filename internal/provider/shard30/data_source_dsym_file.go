package shard30

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
var _ datasource.DataSource = &DsymFileDataSource{}

func NewDsymFileDataSource() datasource.DataSource {
    return &DsymFileDataSource{}
}

// NewDsymFileLegacyDataSource registers this data source under the
// name it had before type names kept mixed-case words whole,
// oneuptime_d_sym. Deprecated.
func NewDsymFileLegacyDataSource() datasource.DataSource {
    return &DsymFileDataSource{isLegacyAlias: true}
}

// DsymFileDataSource defines the data source implementation.
type DsymFileDataSource struct {
    client *Client
    // Registered under the name this data source had before.
    isLegacyAlias bool
}

// DsymFileDataSourceModel describes the data source data model.
type DsymFileDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    DebugId types.String `tfsdk:"debug_id"`
    BinaryName types.String `tfsdk:"binary_name"`
    Architecture types.String `tfsdk:"architecture"`
    AppName types.String `tfsdk:"app_name"`
    Symbols types.String `tfsdk:"symbols"`
    SymbolCount types.Number `tfsdk:"symbol_count"`
    SizeInBytes types.Number `tfsdk:"size_in_bytes"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DsymFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    if d.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_d_sym"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_dsym_file"
}

func (d *DsymFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The function symbols of one build of an iOS app's binary, from its dSYM. Used to name the frames of the app's crash and hang stacks by their functions, and to group them by function. Matched to a frame by the binary's UUID. Look up an existing dsym file by `id`, or by any of its other arguments (`app_name`, `architecture`, `binary_name`, ...): each one set must match, and exactly one dsym file may match them all.",

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
            "debug_id": schema.StringAttribute{
                MarkdownDescription: "The binary's UUID (its LC_UUID), upper case with dashes: the UUID a crash's or hang's frames name the build by.",
                Optional: true,
                Computed: true,
            },
            "binary_name": schema.StringAttribute{
                MarkdownDescription: "The binary the symbols are of, as its dSYM names it: the app's executable or one of its frameworks.",
                Optional: true,
                Computed: true,
            },
            "architecture": schema.StringAttribute{
                MarkdownDescription: "The architecture this build of the binary is for.",
                Optional: true,
                Computed: true,
            },
            "app_name": schema.StringAttribute{
                MarkdownDescription: "The app the dSYM was uploaded for, as the upload named it (its service.name). Shown beside the dSYM; a frame is matched by the binary's UUID alone.",
                Optional: true,
                Computed: true,
            },
            "symbols": schema.StringAttribute{
                MarkdownDescription: "Every function symbol of this build of the binary: its offset from the binary's load address and its name.",
                Optional: true,
                Computed: true,
            },
            "symbol_count": schema.NumberAttribute{
                MarkdownDescription: "How many function symbols this build of the binary has.",
                Optional: true,
                Computed: true,
            },
            "size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "Size of the stored symbols in bytes.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
    if d.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_d_sym has been renamed to oneuptime_dsym_file. Use the new name; the old one keeps working until then."
    }
}

func (d *DsymFileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DsymFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DsymFileDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.DebugId.IsNull() && !data.DebugId.IsUnknown() {
        filters["debugId"] = data.DebugId.ValueString()
        filterNames = append(filterNames, "debug_id = "+fmt.Sprintf("%q", data.DebugId.ValueString()))
    }
    if !data.BinaryName.IsNull() && !data.BinaryName.IsUnknown() {
        filters["binaryName"] = data.BinaryName.ValueString()
        filterNames = append(filterNames, "binary_name = "+fmt.Sprintf("%q", data.BinaryName.ValueString()))
    }
    if !data.Architecture.IsNull() && !data.Architecture.IsUnknown() {
        filters["architecture"] = data.Architecture.ValueString()
        filterNames = append(filterNames, "architecture = "+fmt.Sprintf("%q", data.Architecture.ValueString()))
    }
    if !data.AppName.IsNull() && !data.AppName.IsUnknown() {
        filters["appName"] = data.AppName.ValueString()
        filterNames = append(filterNames, "app_name = "+fmt.Sprintf("%q", data.AppName.ValueString()))
    }
    if !data.Symbols.IsNull() && !data.Symbols.IsUnknown() {
        filters["symbols"] = data.Symbols.ValueString()
        filterNames = append(filterNames, "symbols = "+fmt.Sprintf("%q", data.Symbols.ValueString()))
    }
    if !data.SymbolCount.IsNull() && !data.SymbolCount.IsUnknown() {
        filters["symbolCount"] = LookupNumber(data.SymbolCount)
        filterNames = append(filterNames, "symbol_count = "+data.SymbolCount.ValueBigFloat().String())
    }
    if !data.SizeInBytes.IsNull() && !data.SizeInBytes.IsUnknown() {
        filters["sizeInBytes"] = LookupNumber(data.SizeInBytes)
        filterNames = append(filterNames, "size_in_bytes = "+data.SizeInBytes.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the dsym file up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the dsym file up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "debugId": true,
        "binaryName": true,
        "architecture": true,
        "appName": true,
        "symbols": true,
        "symbolCount": true,
        "sizeInBytes": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/telemetry-dsym/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dsym_file, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dsym file found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read dsym_file: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/telemetry-dsym/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list dsym_file, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list dsym_file: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No dsym file matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one dsym file matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for dsym_file.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["debugId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DebugId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DebugId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DebugId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DebugId = types.StringValue(string(jsonBytes))
        } else {
            data.DebugId = types.StringNull()
        }
    } else if arr, ok := item["debugId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DebugId = types.StringValue(string(jsonBytes))
        } else {
            data.DebugId = types.StringNull()
        }
    } else if val, ok := item["debugId"].(string); ok {
        data.DebugId = types.StringValue(val)
    } else {
        data.DebugId = types.StringNull()
    }
    if obj, ok := item["binaryName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BinaryName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BinaryName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BinaryName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BinaryName = types.StringValue(string(jsonBytes))
        } else {
            data.BinaryName = types.StringNull()
        }
    } else if arr, ok := item["binaryName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.BinaryName = types.StringValue(string(jsonBytes))
        } else {
            data.BinaryName = types.StringNull()
        }
    } else if val, ok := item["binaryName"].(string); ok {
        data.BinaryName = types.StringValue(val)
    } else {
        data.BinaryName = types.StringNull()
    }
    if obj, ok := item["architecture"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Architecture = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Architecture = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Architecture = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Architecture = types.StringValue(string(jsonBytes))
        } else {
            data.Architecture = types.StringNull()
        }
    } else if arr, ok := item["architecture"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Architecture = types.StringValue(string(jsonBytes))
        } else {
            data.Architecture = types.StringNull()
        }
    } else if val, ok := item["architecture"].(string); ok {
        data.Architecture = types.StringValue(val)
    } else {
        data.Architecture = types.StringNull()
    }
    if obj, ok := item["appName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AppName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AppName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AppName = types.StringValue(string(jsonBytes))
        } else {
            data.AppName = types.StringNull()
        }
    } else if arr, ok := item["appName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AppName = types.StringValue(string(jsonBytes))
        } else {
            data.AppName = types.StringNull()
        }
    } else if val, ok := item["appName"].(string); ok {
        data.AppName = types.StringValue(val)
    } else {
        data.AppName = types.StringNull()
    }
    if obj, ok := item["symbols"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Symbols = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Symbols = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Symbols = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Symbols = types.StringValue(string(jsonBytes))
        } else {
            data.Symbols = types.StringNull()
        }
    } else if arr, ok := item["symbols"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Symbols = types.StringValue(string(jsonBytes))
        } else {
            data.Symbols = types.StringNull()
        }
    } else if val, ok := item["symbols"].(string); ok {
        data.Symbols = types.StringValue(val)
    } else {
        data.Symbols = types.StringNull()
    }
    data.SymbolCount = NumberFromAPI(data.SymbolCount, item["symbolCount"])
    data.SizeInBytes = NumberFromAPI(data.SizeInBytes, item["sizeInBytes"])
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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

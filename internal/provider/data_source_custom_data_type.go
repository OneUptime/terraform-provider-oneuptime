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
var _ datasource.DataSource = &CustomDataTypeDataSource{}

func NewCustomDataTypeDataSource() datasource.DataSource {
    return &CustomDataTypeDataSource{}
}

// CustomDataTypeDataSource defines the data source implementation.
type CustomDataTypeDataSource struct {
    client *Client
}

// CustomDataTypeDataSourceModel describes the data source data model.
type CustomDataTypeDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Category types.String `tfsdk:"category"`
    Sensitivity types.String `tfsdk:"sensitivity"`
    Pattern types.String `tfsdk:"pattern"`
    IgnoreCase types.Bool `tfsdk:"ignore_case"`
    Keywords types.String `tfsdk:"keywords"`
    RequiresKeyword types.Bool `tfsdk:"requires_keyword"`
    NameHints types.String `tfsdk:"name_hints"`
    Validator types.String `tfsdk:"validator"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *CustomDataTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_custom_data_type"
}

func (d *CustomDataTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Your own kinds of sensitive data, such as employee IDs or policy numbers: a pattern, the words and field names that confirm it, and an optional checksum. Look up an existing custom data type by `id`, or by any of its other arguments (`name`, `category`, `created_by_user_id`, ...): each one set must match, and exactly one custom data type may match them all.",

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
                MarkdownDescription: "ID of the project this custom data type belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The data type's name, such as Employee ID.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the data type is.",
                Optional: true,
                Computed: true,
            },
            "category": schema.StringAttribute{
                MarkdownDescription: "The category it belongs to.",
                Optional: true,
                Computed: true,
            },
            "sensitivity": schema.StringAttribute{
                MarkdownDescription: "How sensitive it is: Restricted, Confidential, Internal or Public.",
                Optional: true,
                Computed: true,
            },
            "pattern": schema.StringAttribute{
                MarkdownDescription: "A JavaScript regular expression for one value. Leave empty to find it from field names alone.",
                Optional: true,
                Computed: true,
            },
            "ignore_case": schema.BoolAttribute{
                MarkdownDescription: "Whether the pattern ignores upper and lower case.",
                Optional: true,
                Computed: true,
            },
            "keywords": schema.StringAttribute{
                MarkdownDescription: "Words that, near a value, confirm it is this type. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "requires_keyword": schema.BoolAttribute{
                MarkdownDescription: "Count a value only when a keyword is nearby or the field's name hints at the type.",
                Optional: true,
                Computed: true,
            },
            "name_hints": schema.StringAttribute{
                MarkdownDescription: "Field, column and attribute names that hold this type, such as employee id. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "validator": schema.StringAttribute{
                MarkdownDescription: "A checksum a value must pass: None, Luhn, Mod97 or Verhoeff.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether scans look for this data type.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this custom data type. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *CustomDataTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomDataTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CustomDataTypeDataSourceModel

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
    if !data.Category.IsNull() && !data.Category.IsUnknown() {
        filters["category"] = data.Category.ValueString()
        filterNames = append(filterNames, "category = "+fmt.Sprintf("%q", data.Category.ValueString()))
    }
    if !data.Sensitivity.IsNull() && !data.Sensitivity.IsUnknown() {
        filters["sensitivity"] = data.Sensitivity.ValueString()
        filterNames = append(filterNames, "sensitivity = "+fmt.Sprintf("%q", data.Sensitivity.ValueString()))
    }
    if !data.Pattern.IsNull() && !data.Pattern.IsUnknown() {
        filters["pattern"] = data.Pattern.ValueString()
        filterNames = append(filterNames, "pattern = "+fmt.Sprintf("%q", data.Pattern.ValueString()))
    }
    if !data.IgnoreCase.IsNull() && !data.IgnoreCase.IsUnknown() {
        filters["ignoreCase"] = data.IgnoreCase.ValueBool()
        filterNames = append(filterNames, "ignore_case = "+fmt.Sprintf("%t", data.IgnoreCase.ValueBool()))
    }
    if !data.RequiresKeyword.IsNull() && !data.RequiresKeyword.IsUnknown() {
        filters["requiresKeyword"] = data.RequiresKeyword.ValueBool()
        filterNames = append(filterNames, "requires_keyword = "+fmt.Sprintf("%t", data.RequiresKeyword.ValueBool()))
    }
    if !data.Validator.IsNull() && !data.Validator.IsUnknown() {
        filters["validator"] = data.Validator.ValueString()
        filterNames = append(filterNames, "validator = "+fmt.Sprintf("%q", data.Validator.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the custom data type up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the custom data type up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "category": true,
        "sensitivity": true,
        "pattern": true,
        "ignoreCase": true,
        "keywords": true,
        "requiresKeyword": true,
        "nameHints": true,
        "validator": true,
        "isEnabled": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-custom-data-type/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read custom_data_type, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No custom data type found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read custom_data_type: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-custom-data-type/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list custom_data_type, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list custom_data_type: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No custom data type matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one custom data type matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for custom_data_type.")
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
    if obj, ok := item["category"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Category = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Category = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Category = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Category = types.StringValue(string(jsonBytes))
        } else {
            data.Category = types.StringNull()
        }
    } else if val, ok := item["category"].(string); ok {
        data.Category = types.StringValue(val)
    } else {
        data.Category = types.StringNull()
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
    if obj, ok := item["pattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Pattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Pattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Pattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Pattern = types.StringValue(string(jsonBytes))
        } else {
            data.Pattern = types.StringNull()
        }
    } else if val, ok := item["pattern"].(string); ok {
        data.Pattern = types.StringValue(val)
    } else {
        data.Pattern = types.StringNull()
    }
    if val, ok := item["ignoreCase"].(bool); ok {
        data.IgnoreCase = types.BoolValue(val)
    } else {
        data.IgnoreCase = types.BoolNull()
    }
    if obj, ok := item["keywords"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Keywords = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Keywords = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Keywords = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Keywords = types.StringValue(string(jsonBytes))
        } else {
            data.Keywords = types.StringNull()
        }
    } else if val, ok := item["keywords"].(string); ok {
        data.Keywords = types.StringValue(val)
    } else {
        data.Keywords = types.StringNull()
    }
    if val, ok := item["requiresKeyword"].(bool); ok {
        data.RequiresKeyword = types.BoolValue(val)
    } else {
        data.RequiresKeyword = types.BoolNull()
    }
    if obj, ok := item["nameHints"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NameHints = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NameHints = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NameHints = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NameHints = types.StringValue(string(jsonBytes))
        } else {
            data.NameHints = types.StringNull()
        }
    } else if val, ok := item["nameHints"].(string); ok {
        data.NameHints = types.StringValue(val)
    } else {
        data.NameHints = types.StringNull()
    }
    if obj, ok := item["validator"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Validator = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Validator = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Validator = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Validator = types.StringValue(string(jsonBytes))
        } else {
            data.Validator = types.StringNull()
        }
    } else if val, ok := item["validator"].(string); ok {
        data.Validator = types.StringValue(val)
    } else {
        data.Validator = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
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

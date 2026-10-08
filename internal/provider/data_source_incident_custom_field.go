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
var _ datasource.DataSource = &IncidentCustomFieldDataSource{}

func NewIncidentCustomFieldDataSource() datasource.DataSource {
    return &IncidentCustomFieldDataSource{}
}

// IncidentCustomFieldDataSource defines the data source implementation.
type IncidentCustomFieldDataSource struct {
    client *Client
}

// IncidentCustomFieldDataSourceModel describes the data source data model.
type IncidentCustomFieldDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    CustomFieldType types.String `tfsdk:"custom_field_type"`
    DropdownOptions types.String `tfsdk:"dropdown_options"`
    MapFromResourceType types.String `tfsdk:"map_from_resource_type"`
    MapFromCustomFieldName types.String `tfsdk:"map_from_custom_field_name"`
    IsRequiredOnCreate types.Bool `tfsdk:"is_required_on_create"`
    SortOrder types.Number `tfsdk:"sort_order"`
    ShowOnCreate types.Bool `tfsdk:"show_on_create"`
    IncludeInSubscriberNotifications types.Bool `tfsdk:"include_in_subscriber_notifications"`
    VariableKey types.String `tfsdk:"variable_key"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentCustomFieldDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_custom_field"
}

func (d *IncidentCustomFieldDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage custom fields for your incident. Look up an existing incident custom field by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one incident custom field may match them all.",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description of this custom field that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "custom_field_type": schema.StringAttribute{
                MarkdownDescription: "Is this field Text, Number or Boolean? A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "dropdown_options": schema.StringAttribute{
                MarkdownDescription: "Options and optional colors for dropdown fields. Plain one-per-line values remain supported.",
                Optional: true,
                Computed: true,
            },
            "map_from_resource_type": schema.StringAttribute{
                MarkdownDescription: "Related resource this field copies its value from. Empty means values are entered by hand.",
                Optional: true,
                Computed: true,
            },
            "map_from_custom_field_name": schema.StringAttribute{
                MarkdownDescription: "Name of the custom field on the related resource this field copies its value from.",
                Optional: true,
                Computed: true,
            },
            "is_required_on_create": schema.BoolAttribute{
                MarkdownDescription: "When on, an incident declared from the dashboard cannot be created until this field is filled in (a Boolean field must be ticked). It applies only to fields shown on create. Incidents created by monitors, the API, Slack, Microsoft Teams or AI can leave it empty, and it stays optional when an incident is edited later.",
                Optional: true,
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Where this field appears among the incident's custom fields, lowest number first. A new field is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "show_on_create": schema.BoolAttribute{
                MarkdownDescription: "When on, this field is asked for in a Details step when an incident is declared from the dashboard, and incident templates can fill it in.",
                Optional: true,
                Computed: true,
            },
            "include_in_subscriber_notifications": schema.BoolAttribute{
                MarkdownDescription: "When on, this field and its value appear in the messages status page subscribers get about an incident: the default email, Slack and Microsoft Teams messages, and webhooks (under customFields, by the field's template variable key). The default SMS is kept short and leaves it out. Subscribers are usually people outside your team, so turn this on only for fields that are safe to share with them.",
                Optional: true,
                Computed: true,
            },
            "variable_key": schema.StringAttribute{
                MarkdownDescription: "The key this field is reached by in templates, as {{incident.customFields.<key>}}. Made from the field's name when it is created - lowercase letters, digits and underscores, with _2, _3 and so on added when another field already has it - and never changed afterwards, so renaming the field does not break templates that use it.",
                Optional: true,
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

func (d *IncidentCustomFieldDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentCustomFieldDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentCustomFieldDataSourceModel

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
    if !data.DropdownOptions.IsNull() && !data.DropdownOptions.IsUnknown() {
        filters["dropdownOptions"] = data.DropdownOptions.ValueString()
        filterNames = append(filterNames, "dropdown_options = "+fmt.Sprintf("%q", data.DropdownOptions.ValueString()))
    }
    if !data.MapFromResourceType.IsNull() && !data.MapFromResourceType.IsUnknown() {
        filters["mapFromResourceType"] = data.MapFromResourceType.ValueString()
        filterNames = append(filterNames, "map_from_resource_type = "+fmt.Sprintf("%q", data.MapFromResourceType.ValueString()))
    }
    if !data.MapFromCustomFieldName.IsNull() && !data.MapFromCustomFieldName.IsUnknown() {
        filters["mapFromCustomFieldName"] = data.MapFromCustomFieldName.ValueString()
        filterNames = append(filterNames, "map_from_custom_field_name = "+fmt.Sprintf("%q", data.MapFromCustomFieldName.ValueString()))
    }
    if !data.IsRequiredOnCreate.IsNull() && !data.IsRequiredOnCreate.IsUnknown() {
        filters["isRequiredOnCreate"] = data.IsRequiredOnCreate.ValueBool()
        filterNames = append(filterNames, "is_required_on_create = "+fmt.Sprintf("%t", data.IsRequiredOnCreate.ValueBool()))
    }
    if !data.SortOrder.IsNull() && !data.SortOrder.IsUnknown() {
        filters["sortOrder"] = lookupNumber(data.SortOrder)
        filterNames = append(filterNames, "sort_order = "+data.SortOrder.ValueBigFloat().String())
    }
    if !data.ShowOnCreate.IsNull() && !data.ShowOnCreate.IsUnknown() {
        filters["showOnCreate"] = data.ShowOnCreate.ValueBool()
        filterNames = append(filterNames, "show_on_create = "+fmt.Sprintf("%t", data.ShowOnCreate.ValueBool()))
    }
    if !data.IncludeInSubscriberNotifications.IsNull() && !data.IncludeInSubscriberNotifications.IsUnknown() {
        filters["includeInSubscriberNotifications"] = data.IncludeInSubscriberNotifications.ValueBool()
        filterNames = append(filterNames, "include_in_subscriber_notifications = "+fmt.Sprintf("%t", data.IncludeInSubscriberNotifications.ValueBool()))
    }
    if !data.VariableKey.IsNull() && !data.VariableKey.IsUnknown() {
        filters["variableKey"] = data.VariableKey.ValueString()
        filterNames = append(filterNames, "variable_key = "+fmt.Sprintf("%q", data.VariableKey.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident custom field up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident custom field up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "customFieldType": true,
        "dropdownOptions": true,
        "mapFromResourceType": true,
        "mapFromCustomFieldName": true,
        "isRequiredOnCreate": true,
        "sortOrder": true,
        "showOnCreate": true,
        "includeInSubscriberNotifications": true,
        "variableKey": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-custom-field/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_custom_field, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident custom field found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_custom_field: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-custom-field/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_custom_field, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_custom_field: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident custom field matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident custom field matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_custom_field.")
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
    if obj, ok := item["customFieldType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFieldType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFieldType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFieldType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFieldType = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFieldType = types.StringNull()
        }
    } else if val, ok := item["customFieldType"].(string); ok {
        data.CustomFieldType = types.StringValue(val)
    } else {
        data.CustomFieldType = types.StringNull()
    }
    if obj, ok := item["dropdownOptions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DropdownOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DropdownOptions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DropdownOptions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DropdownOptions = types.StringValue(string(jsonBytes))
        } else {
            data.DropdownOptions = types.StringNull()
        }
    } else if val, ok := item["dropdownOptions"].(string); ok {
        data.DropdownOptions = types.StringValue(val)
    } else {
        data.DropdownOptions = types.StringNull()
    }
    if obj, ok := item["mapFromResourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MapFromResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MapFromResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MapFromResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MapFromResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.MapFromResourceType = types.StringNull()
        }
    } else if val, ok := item["mapFromResourceType"].(string); ok {
        data.MapFromResourceType = types.StringValue(val)
    } else {
        data.MapFromResourceType = types.StringNull()
    }
    if obj, ok := item["mapFromCustomFieldName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MapFromCustomFieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MapFromCustomFieldName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MapFromCustomFieldName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MapFromCustomFieldName = types.StringValue(string(jsonBytes))
        } else {
            data.MapFromCustomFieldName = types.StringNull()
        }
    } else if val, ok := item["mapFromCustomFieldName"].(string); ok {
        data.MapFromCustomFieldName = types.StringValue(val)
    } else {
        data.MapFromCustomFieldName = types.StringNull()
    }
    if val, ok := item["isRequiredOnCreate"].(bool); ok {
        data.IsRequiredOnCreate = types.BoolValue(val)
    } else {
        data.IsRequiredOnCreate = types.BoolNull()
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
    if val, ok := item["showOnCreate"].(bool); ok {
        data.ShowOnCreate = types.BoolValue(val)
    } else {
        data.ShowOnCreate = types.BoolNull()
    }
    if val, ok := item["includeInSubscriberNotifications"].(bool); ok {
        data.IncludeInSubscriberNotifications = types.BoolValue(val)
    } else {
        data.IncludeInSubscriberNotifications = types.BoolNull()
    }
    if obj, ok := item["variableKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VariableKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VariableKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VariableKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VariableKey = types.StringValue(string(jsonBytes))
        } else {
            data.VariableKey = types.StringNull()
        }
    } else if val, ok := item["variableKey"].(string); ok {
        data.VariableKey = types.StringValue(val)
    } else {
        data.VariableKey = types.StringNull()
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

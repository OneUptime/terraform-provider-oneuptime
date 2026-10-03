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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *IncidentCustomFieldDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_custom_field"
}

func (d *IncidentCustomFieldDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage custom fields for your incident. Look up an existing incident_custom_field by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description of this custom field that will help you remember.",
                Computed: true,
            },
            "custom_field_type": schema.StringAttribute{
                MarkdownDescription: "Is this field Text, Number or Boolean?.",
                Computed: true,
            },
            "dropdown_options": schema.StringAttribute{
                MarkdownDescription: "Options and optional colors for dropdown fields. Plain one-per-line values remain supported..",
                Computed: true,
            },
            "map_from_resource_type": schema.StringAttribute{
                MarkdownDescription: "Related resource this field copies its value from. Empty means values are entered by hand..",
                Computed: true,
            },
            "map_from_custom_field_name": schema.StringAttribute{
                MarkdownDescription: "Name of the custom field on the related resource this field copies its value from..",
                Computed: true,
            },
            "is_required_on_create": schema.BoolAttribute{
                MarkdownDescription: "When on, an incident declared from the dashboard cannot be created until this field is filled in (a Boolean field must be ticked). It applies only to fields shown on create. Incidents created by monitors, the API, Slack, Microsoft Teams or AI can leave it empty, and it stays optional when an incident is edited later..",
                Computed: true,
            },
            "sort_order": schema.NumberAttribute{
                MarkdownDescription: "Where this field appears among the incident's custom fields, lowest number first. A new field is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them..",
                Computed: true,
            },
            "show_on_create": schema.BoolAttribute{
                MarkdownDescription: "When on, this field is asked for in a Details step when an incident is declared from the dashboard, and incident templates can fill it in..",
                Computed: true,
            },
            "include_in_subscriber_notifications": schema.BoolAttribute{
                MarkdownDescription: "When on, this field and its value appear in the messages status page subscribers get about an incident: the default email, Slack and Microsoft Teams messages, and webhooks (under customFields, by the field's template variable key). The default SMS is kept short and leaves it out. Subscribers are usually people outside your team, so turn this on only for fields that are safe to share with them..",
                Computed: true,
            },
            "variable_key": schema.StringAttribute{
                MarkdownDescription: "The key this field is reached by in templates, as {{incident.customFields.<key>}}. Made from the field's name when it is created - lowercase letters, digits and underscores, with _2, _3 and so on added when another field already has it - and never changed afterwards, so renaming the field does not break templates that use it..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a incident_custom_field.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
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
        "deletedByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_custom_field found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_custom_field found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident_custom_field matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

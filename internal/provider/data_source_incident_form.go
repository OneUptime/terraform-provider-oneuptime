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
var _ datasource.DataSource = &IncidentFormDataSource{}

func NewIncidentFormDataSource() datasource.DataSource {
    return &IncidentFormDataSource{}
}

// IncidentFormDataSource defines the data source implementation.
type IncidentFormDataSource struct {
    client *Client
}

// IncidentFormDataSourceModel describes the data source data model.
type IncidentFormDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ShareKey types.String `tfsdk:"share_key"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    AllowReporterToChooseSeverity types.Bool `tfsdk:"allow_reporter_to_choose_severity"`
    IncidentTemplateId types.String `tfsdk:"incident_template_id"`
    DescriptionSetting types.String `tfsdk:"description_setting"`
    CustomFieldSettings types.String `tfsdk:"custom_field_settings"`
    IsReporterDetailsRequired types.Bool `tfsdk:"is_reporter_details_required"`
    SuccessMessage types.String `tfsdk:"success_message"`
    IpWhitelist types.String `tfsdk:"ip_whitelist"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentFormDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_form"
}

func (d *IncidentFormDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Forms anyone with the link can fill in to report an incident, without a OneUptime account. Each submission declares an incident in this project. Look up an existing incident_form by `id` or by `name`.",

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
                MarkdownDescription: "Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown..",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the form's link works. While the form is turned off, its public page shows a not-available message and nothing can be submitted..",
                Computed: true,
            },
            "share_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "allow_reporter_to_choose_severity": schema.BoolAttribute{
                MarkdownDescription: "When on, the form asks the reporter to choose a severity from the project's incident severities, with the form's own severity chosen to begin with..",
                Computed: true,
            },
            "incident_template_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description_setting": schema.StringAttribute{
                MarkdownDescription: "Whether the form asks the reporter to describe the incident: Required, Optional or Hidden..",
                Computed: true,
            },
            "custom_field_settings": schema.StringAttribute{
                MarkdownDescription: "The incident custom fields the form asks for, keyed by each field's template variable key (variableKey). Required means the reporter must answer it, Optional that they may leave it empty. Only the fields listed as Required or Optional are asked: a field that is not listed, or is Hidden or Default, is not on the form. The answers become the incident's custom field values..",
                Computed: true,
            },
            "is_reporter_details_required": schema.BoolAttribute{
                MarkdownDescription: "When on, the reporter must give their name and email. When off, they may report anonymously..",
                Computed: true,
            },
            "success_message": schema.StringAttribute{
                MarkdownDescription: "Shown to the reporter after they submit the form, together with the new incident's number. Markdown..",
                Computed: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network..",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncidentFormDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentFormDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentFormDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a incident_form.",
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
        "isEnabled": true,
        "shareKey": true,
        "incidentSeverityId": true,
        "allowReporterToChooseSeverity": true,
        "incidentTemplateId": true,
        "descriptionSetting": true,
        "customFieldSettings": true,
        "isReporterDetailsRequired": true,
        "successMessage": true,
        "ipWhitelist": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-form/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_form, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_form found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_form: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-form/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_form, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_form: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident_form found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident_form matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_form.")
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["shareKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ShareKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ShareKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ShareKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ShareKey = types.StringValue(string(jsonBytes))
        } else {
            data.ShareKey = types.StringNull()
        }
    } else if val, ok := item["shareKey"].(string); ok {
        data.ShareKey = types.StringValue(val)
    } else {
        data.ShareKey = types.StringNull()
    }
    if obj, ok := item["incidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := item["allowReporterToChooseSeverity"].(bool); ok {
        data.AllowReporterToChooseSeverity = types.BoolValue(val)
    } else {
        data.AllowReporterToChooseSeverity = types.BoolNull()
    }
    if obj, ok := item["incidentTemplateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentTemplateId = types.StringNull()
        }
    } else if val, ok := item["incidentTemplateId"].(string); ok {
        data.IncidentTemplateId = types.StringValue(val)
    } else {
        data.IncidentTemplateId = types.StringNull()
    }
    if obj, ok := item["descriptionSetting"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DescriptionSetting = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DescriptionSetting = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DescriptionSetting = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DescriptionSetting = types.StringValue(string(jsonBytes))
        } else {
            data.DescriptionSetting = types.StringNull()
        }
    } else if val, ok := item["descriptionSetting"].(string); ok {
        data.DescriptionSetting = types.StringValue(val)
    } else {
        data.DescriptionSetting = types.StringNull()
    }
    if obj, ok := item["customFieldSettings"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFieldSettings = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFieldSettings = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFieldSettings = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFieldSettings = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFieldSettings = types.StringNull()
        }
    } else if val, ok := item["customFieldSettings"].(string); ok {
        data.CustomFieldSettings = types.StringValue(val)
    } else {
        data.CustomFieldSettings = types.StringNull()
    }
    if val, ok := item["isReporterDetailsRequired"].(bool); ok {
        data.IsReporterDetailsRequired = types.BoolValue(val)
    } else {
        data.IsReporterDetailsRequired = types.BoolNull()
    }
    if obj, ok := item["successMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SuccessMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SuccessMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SuccessMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SuccessMessage = types.StringValue(string(jsonBytes))
        } else {
            data.SuccessMessage = types.StringNull()
        }
    } else if val, ok := item["successMessage"].(string); ok {
        data.SuccessMessage = types.StringValue(val)
    } else {
        data.SuccessMessage = types.StringNull()
    }
    if obj, ok := item["ipWhitelist"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IpWhitelist = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IpWhitelist = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := item["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
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

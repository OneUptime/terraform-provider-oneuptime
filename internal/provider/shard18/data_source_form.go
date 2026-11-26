package shard18

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
var _ datasource.DataSource = &FormDataSource{}

func NewFormDataSource() datasource.DataSource {
    return &FormDataSource{}
}

// FormDataSource defines the data source implementation.
type FormDataSource struct {
    client *Client
}

// FormDataSourceModel describes the data source data model.
type FormDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ShareKey types.String `tfsdk:"share_key"`
    TargetType types.String `tfsdk:"target_type"`
    Fields types.String `tfsdk:"fields"`
    Templates types.String `tfsdk:"templates"`
    TargetSettings types.String `tfsdk:"target_settings"`
    SuccessMessage types.String `tfsdk:"success_message"`
    IpWhitelist types.String `tfsdk:"ip_whitelist"`
    RequireSignIn types.Bool `tfsdk:"require_sign_in"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    LogoFileId types.String `tfsdk:"logo_file_id"`
    LogoAltText types.String `tfsdk:"logo_alt_text"`
    FaviconFileId types.String `tfsdk:"favicon_file_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *FormDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_form"
}

func (d *FormDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Forms anyone with the link can fill in, without a OneUptime account. Each submission creates an incident, an alert or a maintenance event in this project. Look up an existing form by `id`, or by any of its other arguments (`name`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one form may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The form's name, shown as the heading of its public page. Unique within the project.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the form's link works. While it is off, the public page shows a not-available message and nothing can be submitted.",
                Optional: true,
                Computed: true,
            },
            "share_key": schema.StringAttribute{
                MarkdownDescription: "The key in the form's public link, /accounts/form/<shareKey>. Generated when the form is created. Resetting the link in the dashboard replaces it, and the old link stops working.",
                Optional: true,
                Computed: true,
            },
            "target_type": schema.StringAttribute{
                MarkdownDescription: "What each submission creates: Incident, Alert, or ScheduledMaintenance (a maintenance event). A form that creates alerts takes permission to create alerts to set up.",
                Optional: true,
                Computed: true,
            },
            "fields": schema.StringAttribute{
                MarkdownDescription: "The questions the form asks, in order. Each has an id, a source (Question: one of the form's own, answered by type; TargetField: a built-in field of what the form creates, by targetField; TargetCustomField: one of its custom fields, by customFieldId; Submitter: the submitter's Name or Email), a label, optional help text, isRequired and isHidden (not shown on the public form unless the template a submission starts from asks it, and otherwise answered only from the template a submission started from; never required, and never a field the target cannot be created without). isRequired and isHidden are the form's default: each template can make a question Required, Optional or Hidden (its fieldSettings). A new form starts with a title, a description and the submitter's name and email. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "templates": schema.StringAttribute{
                MarkdownDescription: "Named sets of answers a submission can start from, in the order the form lists them. Each has an id, a name (unique within the form), isDefault (the form opens with it; at most one template), answers: an object keyed by question id, each answer as a submission sends it - text, a number, true or false, an option's value, or a list of values for a multi-select - and fieldSettings: an object keyed by question id that makes a question Required, Optional or Hidden when a submission starts from the template; a question it does not list is asked as the form asks it, and a question the form cannot create its record without is always asked and required. The public form lists the templates above its questions and fills in a template's answers when one is chosen, or when its link names one (?template=<id>). A question the submission was not asked - hidden by the form or by the template - is answered only from the template a submission started from. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "target_settings": schema.StringAttribute{
                MarkdownDescription: "What every submission starts with besides the answers. For incidents: defaultTitle, incidentSeverityId, incidentTemplateId, monitorIds, labelIds, onCallDutyPolicyIds, ownerUserIds and ownerTeamIds. For alerts: defaultTitle, alertSeverityId, alertStateId, monitorId, labelIds, onCallDutyPolicyIds, ownerUserIds and ownerTeamIds. For maintenance events: defaultTitle, monitorIds, statusPageIds, labelIds, ownerUserIds, ownerTeamIds, showOnStatusPages and notifySubscribers. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "success_message": schema.StringAttribute{
                MarkdownDescription: "Shown after the form is submitted, together with the number of what the submission created. Markdown.",
                Optional: true,
                Computed: true,
            },
            "ip_whitelist": schema.StringAttribute{
                MarkdownDescription: "The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network.",
                Optional: true,
                Computed: true,
            },
            "require_sign_in": schema.BoolAttribute{
                MarkdownDescription: "Only this form's private users can open and submit it on its public link, once they sign in with a password or single sign-on.",
                Optional: true,
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "While the form asks its visitors to sign in, only its single sign-on lets them in: private users' passwords no longer do.",
                Optional: true,
                Computed: true,
            },
            "logo_file_id": schema.StringAttribute{
                MarkdownDescription: "ID of the file the form's public page shows as its logo: upload the image to /api/file first. Leave it empty to show the OneUptime logo. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_file"),
            },
            "logo_alt_text": schema.StringAttribute{
                MarkdownDescription: "What the logo says, read out by screen readers: usually your organization's name. Leave it empty and screen readers skip the logo.",
                Optional: true,
                Computed: true,
            },
            "favicon_file_id": schema.StringAttribute{
                MarkdownDescription: "ID of the file the form's public page shows as the browser tab's icon: upload the image to /api/file first. Leave it empty to show the OneUptime favicon. The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_file"),
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *FormDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FormDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data FormDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.ShareKey.IsNull() && !data.ShareKey.IsUnknown() {
        filters["shareKey"] = data.ShareKey.ValueString()
        filterNames = append(filterNames, "share_key = "+fmt.Sprintf("%q", data.ShareKey.ValueString()))
    }
    if !data.TargetType.IsNull() && !data.TargetType.IsUnknown() {
        filters["targetType"] = data.TargetType.ValueString()
        filterNames = append(filterNames, "target_type = "+fmt.Sprintf("%q", data.TargetType.ValueString()))
    }
    if !data.SuccessMessage.IsNull() && !data.SuccessMessage.IsUnknown() {
        filters["successMessage"] = data.SuccessMessage.ValueString()
        filterNames = append(filterNames, "success_message = "+fmt.Sprintf("%q", data.SuccessMessage.ValueString()))
    }
    if !data.IpWhitelist.IsNull() && !data.IpWhitelist.IsUnknown() {
        filters["ipWhitelist"] = data.IpWhitelist.ValueString()
        filterNames = append(filterNames, "ip_whitelist = "+fmt.Sprintf("%q", data.IpWhitelist.ValueString()))
    }
    if !data.RequireSignIn.IsNull() && !data.RequireSignIn.IsUnknown() {
        filters["requireSignIn"] = data.RequireSignIn.ValueBool()
        filterNames = append(filterNames, "require_sign_in = "+fmt.Sprintf("%t", data.RequireSignIn.ValueBool()))
    }
    if !data.RequireSsoForLogin.IsNull() && !data.RequireSsoForLogin.IsUnknown() {
        filters["requireSsoForLogin"] = data.RequireSsoForLogin.ValueBool()
        filterNames = append(filterNames, "require_sso_for_login = "+fmt.Sprintf("%t", data.RequireSsoForLogin.ValueBool()))
    }
    if !data.LogoFileId.IsNull() && !data.LogoFileId.IsUnknown() {
        filters["logoFileId"] = data.LogoFileId.ValueString()
        filterNames = append(filterNames, "logo_file_id = "+fmt.Sprintf("%q", data.LogoFileId.ValueString()))
    }
    if !data.LogoAltText.IsNull() && !data.LogoAltText.IsUnknown() {
        filters["logoAltText"] = data.LogoAltText.ValueString()
        filterNames = append(filterNames, "logo_alt_text = "+fmt.Sprintf("%q", data.LogoAltText.ValueString()))
    }
    if !data.FaviconFileId.IsNull() && !data.FaviconFileId.IsUnknown() {
        filters["faviconFileId"] = data.FaviconFileId.ValueString()
        filterNames = append(filterNames, "favicon_file_id = "+fmt.Sprintf("%q", data.FaviconFileId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the form up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the form up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "shareKey": true,
        "targetType": true,
        "fields": true,
        "templates": true,
        "targetSettings": true,
        "successMessage": true,
        "ipWhitelist": true,
        "requireSignIn": true,
        "requireSsoForLogin": true,
        "logoFileId": true,
        "logoAltText": true,
        "faviconFileId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/form/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read form, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No form found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read form: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/form/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list form, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list form: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No form matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one form matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for form.")
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
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["shareKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ShareKey = types.StringValue(string(jsonBytes))
        } else {
            data.ShareKey = types.StringNull()
        }
    } else if val, ok := item["shareKey"].(string); ok {
        data.ShareKey = types.StringValue(val)
    } else {
        data.ShareKey = types.StringNull()
    }
    if obj, ok := item["targetType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetType = types.StringValue(string(jsonBytes))
        } else {
            data.TargetType = types.StringNull()
        }
    } else if arr, ok := item["targetType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetType = types.StringValue(string(jsonBytes))
        } else {
            data.TargetType = types.StringNull()
        }
    } else if val, ok := item["targetType"].(string); ok {
        data.TargetType = types.StringValue(val)
    } else {
        data.TargetType = types.StringNull()
    }
    if obj, ok := item["fields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Fields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Fields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Fields = types.StringValue(string(jsonBytes))
        } else {
            data.Fields = types.StringNull()
        }
    } else if arr, ok := item["fields"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Fields = types.StringValue(string(jsonBytes))
        } else {
            data.Fields = types.StringNull()
        }
    } else if val, ok := item["fields"].(string); ok {
        data.Fields = types.StringValue(val)
    } else {
        data.Fields = types.StringNull()
    }
    if obj, ok := item["templates"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Templates = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Templates = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Templates = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Templates = types.StringValue(string(jsonBytes))
        } else {
            data.Templates = types.StringNull()
        }
    } else if arr, ok := item["templates"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Templates = types.StringValue(string(jsonBytes))
        } else {
            data.Templates = types.StringNull()
        }
    } else if val, ok := item["templates"].(string); ok {
        data.Templates = types.StringValue(val)
    } else {
        data.Templates = types.StringNull()
    }
    if obj, ok := item["targetSettings"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetSettings = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetSettings = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetSettings = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetSettings = types.StringValue(string(jsonBytes))
        } else {
            data.TargetSettings = types.StringNull()
        }
    } else if arr, ok := item["targetSettings"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TargetSettings = types.StringValue(string(jsonBytes))
        } else {
            data.TargetSettings = types.StringNull()
        }
    } else if val, ok := item["targetSettings"].(string); ok {
        data.TargetSettings = types.StringValue(val)
    } else {
        data.TargetSettings = types.StringNull()
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
    } else if arr, ok := item["successMessage"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["ipWhitelist"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IpWhitelist = types.StringValue(string(jsonBytes))
        } else {
            data.IpWhitelist = types.StringNull()
        }
    } else if val, ok := item["ipWhitelist"].(string); ok {
        data.IpWhitelist = types.StringValue(val)
    } else {
        data.IpWhitelist = types.StringNull()
    }
    if val, ok := item["requireSignIn"].(bool); ok {
        data.RequireSignIn = types.BoolValue(val)
    } else {
        data.RequireSignIn = types.BoolNull()
    }
    if val, ok := item["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if obj, ok := item["logoFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LogoFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LogoFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if arr, ok := item["logoFileId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LogoFileId = types.StringValue(string(jsonBytes))
        } else {
            data.LogoFileId = types.StringNull()
        }
    } else if val, ok := item["logoFileId"].(string); ok {
        data.LogoFileId = types.StringValue(val)
    } else {
        data.LogoFileId = types.StringNull()
    }
    if obj, ok := item["logoAltText"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LogoAltText = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LogoAltText = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LogoAltText = types.StringValue(string(jsonBytes))
        } else {
            data.LogoAltText = types.StringNull()
        }
    } else if arr, ok := item["logoAltText"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LogoAltText = types.StringValue(string(jsonBytes))
        } else {
            data.LogoAltText = types.StringNull()
        }
    } else if val, ok := item["logoAltText"].(string); ok {
        data.LogoAltText = types.StringValue(val)
    } else {
        data.LogoAltText = types.StringNull()
    }
    if obj, ok := item["faviconFileId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FaviconFileId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FaviconFileId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if arr, ok := item["faviconFileId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FaviconFileId = types.StringValue(string(jsonBytes))
        } else {
            data.FaviconFileId = types.StringNull()
        }
    } else if val, ok := item["faviconFileId"].(string); ok {
        data.FaviconFileId = types.StringValue(val)
    } else {
        data.FaviconFileId = types.StringNull()
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

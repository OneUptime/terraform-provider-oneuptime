package shard01

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
var _ datasource.DataSource = &CallAndSmsConfigDataSource{}

func NewCallAndSmsConfigDataSource() datasource.DataSource {
    return &CallAndSmsConfigDataSource{}
}

// CallAndSmsConfigDataSource defines the data source implementation.
type CallAndSmsConfigDataSource struct {
    client *Client
}

// CallAndSmsConfigDataSourceModel describes the data source data model.
type CallAndSmsConfigDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    UsedFor types.String `tfsdk:"used_for"`
    TwilioAccountSid types.String `tfsdk:"twilio_account_sid"`
    TwilioAuthToken types.String `tfsdk:"twilio_auth_token"`
    TwilioPrimaryPhoneNumber types.String `tfsdk:"twilio_primary_phone_number"`
    TwilioSecondaryPhoneNumbers types.String `tfsdk:"twilio_secondary_phone_numbers"`
    TwilioSmsPhoneNumbers types.String `tfsdk:"twilio_sms_phone_numbers"`
    IsProjectDefault types.Bool `tfsdk:"is_project_default"`
    IsProjectDefaultForSms types.Bool `tfsdk:"is_project_default_for_sms"`
    IsProjectDefaultForCalls types.Bool `tfsdk:"is_project_default_for_calls"`
}

func (d *CallAndSmsConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_call_and_sms_config"
}

func (d *CallAndSmsConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A project's own Twilio account. The project's default accounts send the SMS and make the calls to its team members instead of OneUptime's own. A status page can name one for the SMS to its subscribers, and an incoming call policy one for its calls. Look up an existing call and sms config by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one call and sms config may match them all.",

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
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create Call and SMS], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Call and SMS], Update: [Project Owner, Project Admin, Edit Call and SMS]",
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "used_for": schema.StringAttribute{
                MarkdownDescription: "What this Twilio config is used for: \"Calls and SMS\" (the default), \"SMS only\" or \"Calls only\". Calls and SMS and Calls only need a Twilio Primary Phone Number; SMS only needs Twilio SMS Phone Numbers. A config used for SMS only makes no calls, and one used for calls only sends no SMS: those go through OneUptime's own Twilio account, as for a project without a config.",
                Optional: true,
                Computed: true,
            },
            "twilio_account_sid": schema.StringAttribute{
                MarkdownDescription: "Account SID for your Twilio Account.",
                Optional: true,
                Computed: true,
            },
            "twilio_auth_token": schema.StringAttribute{
                MarkdownDescription: "Auth Token for your Twilio Account.",
                Computed: true,
                Sensitive: true,
            },
            "twilio_primary_phone_number": schema.StringAttribute{
                MarkdownDescription: "Primary Phone Number for your Twilio account.",
                Computed: true,
            },
            "twilio_secondary_phone_numbers": schema.StringAttribute{
                MarkdownDescription: "Secondary Phone Number for your Twilio account.",
                Optional: true,
                Computed: true,
            },
            "twilio_sms_phone_numbers": schema.StringAttribute{
                MarkdownDescription: "Comma-separated numbers to send SMS from, for phone numbers that can only make calls. Once set, every SMS goes out from the one in the recipient's country, or else the first. Leave it empty to send SMS from the same numbers as calls.",
                Optional: true,
                Computed: true,
            },
            "is_project_default": schema.BoolAttribute{
                MarkdownDescription: "Whether the SMS and calls sent to project team members (on-call notifications, alerts, verification codes, etc.) go through this Twilio config instead of the global config - the ones it is used for. A project has one default for SMS and one for calls, so a config used for SMS only and one used for calls only can both be the default. Setting this to true makes the config the project default for everything it is used for, taking that from the project's other configs; false makes it the default for nothing. It reads true while the config is the default for SMS, calls or both: isProjectDefaultForSms and isProjectDefaultForCalls say which. A project's first Twilio config becomes the project default when this is not set on create; an explicit false is kept. This does not affect status pages — they use the config explicitly assigned to each status page.",
                Optional: true,
                Computed: true,
            },
            "is_project_default_for_sms": schema.BoolAttribute{
                MarkdownDescription: "Whether the SMS sent to project team members go through this Twilio config: it is the project default for SMS. Set by OneUptime from isProjectDefault and what the config is used for; at most one config of a project is the default for SMS.",
                Optional: true,
                Computed: true,
            },
            "is_project_default_for_calls": schema.BoolAttribute{
                MarkdownDescription: "Whether the calls made to project team members go through this Twilio config: it is the project default for calls. Set by OneUptime from isProjectDefault and what the config is used for; at most one config of a project is the default for calls.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *CallAndSmsConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CallAndSmsConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CallAndSmsConfigDataSourceModel

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
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.UsedFor.IsNull() && !data.UsedFor.IsUnknown() {
        filters["usedFor"] = data.UsedFor.ValueString()
        filterNames = append(filterNames, "used_for = "+fmt.Sprintf("%q", data.UsedFor.ValueString()))
    }
    if !data.TwilioAccountSid.IsNull() && !data.TwilioAccountSid.IsUnknown() {
        filters["twilioAccountSID"] = data.TwilioAccountSid.ValueString()
        filterNames = append(filterNames, "twilio_account_sid = "+fmt.Sprintf("%q", data.TwilioAccountSid.ValueString()))
    }
    if !data.TwilioSecondaryPhoneNumbers.IsNull() && !data.TwilioSecondaryPhoneNumbers.IsUnknown() {
        filters["twilioSecondaryPhoneNumbers"] = data.TwilioSecondaryPhoneNumbers.ValueString()
        filterNames = append(filterNames, "twilio_secondary_phone_numbers = "+fmt.Sprintf("%q", data.TwilioSecondaryPhoneNumbers.ValueString()))
    }
    if !data.TwilioSmsPhoneNumbers.IsNull() && !data.TwilioSmsPhoneNumbers.IsUnknown() {
        filters["twilioSmsPhoneNumbers"] = data.TwilioSmsPhoneNumbers.ValueString()
        filterNames = append(filterNames, "twilio_sms_phone_numbers = "+fmt.Sprintf("%q", data.TwilioSmsPhoneNumbers.ValueString()))
    }
    if !data.IsProjectDefault.IsNull() && !data.IsProjectDefault.IsUnknown() {
        filters["isProjectDefault"] = data.IsProjectDefault.ValueBool()
        filterNames = append(filterNames, "is_project_default = "+fmt.Sprintf("%t", data.IsProjectDefault.ValueBool()))
    }
    if !data.IsProjectDefaultForSms.IsNull() && !data.IsProjectDefaultForSms.IsUnknown() {
        filters["isProjectDefaultForSms"] = data.IsProjectDefaultForSms.ValueBool()
        filterNames = append(filterNames, "is_project_default_for_sms = "+fmt.Sprintf("%t", data.IsProjectDefaultForSms.ValueBool()))
    }
    if !data.IsProjectDefaultForCalls.IsNull() && !data.IsProjectDefaultForCalls.IsUnknown() {
        filters["isProjectDefaultForCalls"] = data.IsProjectDefaultForCalls.ValueBool()
        filterNames = append(filterNames, "is_project_default_for_calls = "+fmt.Sprintf("%t", data.IsProjectDefaultForCalls.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the call and sms config up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the call and sms config up by.",
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
        "createdByUserId": true,
        "usedFor": true,
        "twilioAccountSID": true,
        "twilioAuthToken": true,
        "twilioPrimaryPhoneNumber": true,
        "twilioSecondaryPhoneNumbers": true,
        "twilioSmsPhoneNumbers": true,
        "isProjectDefault": true,
        "isProjectDefaultForSms": true,
        "isProjectDefaultForCalls": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/call-sms-config/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read call_and_sms_config, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No call and sms config found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read call_and_sms_config: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/call-sms-config/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list call_and_sms_config, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list call_and_sms_config: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No call and sms config matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one call and sms config matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for call_and_sms_config.")
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
    } else if arr, ok := item["slug"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    if obj, ok := item["usedFor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UsedFor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UsedFor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UsedFor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UsedFor = types.StringValue(string(jsonBytes))
        } else {
            data.UsedFor = types.StringNull()
        }
    } else if arr, ok := item["usedFor"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UsedFor = types.StringValue(string(jsonBytes))
        } else {
            data.UsedFor = types.StringNull()
        }
    } else if val, ok := item["usedFor"].(string); ok {
        data.UsedFor = types.StringValue(val)
    } else {
        data.UsedFor = types.StringNull()
    }
    if obj, ok := item["twilioAccountSID"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TwilioAccountSid = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TwilioAccountSid = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TwilioAccountSid = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TwilioAccountSid = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioAccountSid = types.StringNull()
        }
    } else if arr, ok := item["twilioAccountSID"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TwilioAccountSid = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioAccountSid = types.StringNull()
        }
    } else if val, ok := item["twilioAccountSID"].(string); ok {
        data.TwilioAccountSid = types.StringValue(val)
    } else {
        data.TwilioAccountSid = types.StringNull()
    }
    if obj, ok := item["twilioAuthToken"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TwilioAuthToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TwilioAuthToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TwilioAuthToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TwilioAuthToken = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioAuthToken = types.StringNull()
        }
    } else if arr, ok := item["twilioAuthToken"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TwilioAuthToken = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioAuthToken = types.StringNull()
        }
    } else if val, ok := item["twilioAuthToken"].(string); ok {
        data.TwilioAuthToken = types.StringValue(val)
    } else {
        data.TwilioAuthToken = types.StringNull()
    }
    if obj, ok := item["twilioPrimaryPhoneNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TwilioPrimaryPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TwilioPrimaryPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TwilioPrimaryPhoneNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TwilioPrimaryPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioPrimaryPhoneNumber = types.StringNull()
        }
    } else if arr, ok := item["twilioPrimaryPhoneNumber"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TwilioPrimaryPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioPrimaryPhoneNumber = types.StringNull()
        }
    } else if val, ok := item["twilioPrimaryPhoneNumber"].(string); ok {
        data.TwilioPrimaryPhoneNumber = types.StringValue(val)
    } else {
        data.TwilioPrimaryPhoneNumber = types.StringNull()
    }
    if obj, ok := item["twilioSecondaryPhoneNumbers"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TwilioSecondaryPhoneNumbers = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TwilioSecondaryPhoneNumbers = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TwilioSecondaryPhoneNumbers = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TwilioSecondaryPhoneNumbers = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioSecondaryPhoneNumbers = types.StringNull()
        }
    } else if arr, ok := item["twilioSecondaryPhoneNumbers"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TwilioSecondaryPhoneNumbers = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioSecondaryPhoneNumbers = types.StringNull()
        }
    } else if val, ok := item["twilioSecondaryPhoneNumbers"].(string); ok {
        data.TwilioSecondaryPhoneNumbers = types.StringValue(val)
    } else {
        data.TwilioSecondaryPhoneNumbers = types.StringNull()
    }
    if obj, ok := item["twilioSmsPhoneNumbers"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TwilioSmsPhoneNumbers = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TwilioSmsPhoneNumbers = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TwilioSmsPhoneNumbers = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TwilioSmsPhoneNumbers = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioSmsPhoneNumbers = types.StringNull()
        }
    } else if arr, ok := item["twilioSmsPhoneNumbers"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TwilioSmsPhoneNumbers = types.StringValue(string(jsonBytes))
        } else {
            data.TwilioSmsPhoneNumbers = types.StringNull()
        }
    } else if val, ok := item["twilioSmsPhoneNumbers"].(string); ok {
        data.TwilioSmsPhoneNumbers = types.StringValue(val)
    } else {
        data.TwilioSmsPhoneNumbers = types.StringNull()
    }
    if val, ok := item["isProjectDefault"].(bool); ok {
        data.IsProjectDefault = types.BoolValue(val)
    } else {
        data.IsProjectDefault = types.BoolNull()
    }
    if val, ok := item["isProjectDefaultForSms"].(bool); ok {
        data.IsProjectDefaultForSms = types.BoolValue(val)
    } else {
        data.IsProjectDefaultForSms = types.BoolNull()
    }
    if val, ok := item["isProjectDefaultForCalls"].(bool); ok {
        data.IsProjectDefaultForCalls = types.BoolValue(val)
    } else {
        data.IsProjectDefaultForCalls = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

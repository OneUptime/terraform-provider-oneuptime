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
var _ datasource.DataSource = &SubscriberNotificationTemplateDataSource{}

func NewSubscriberNotificationTemplateDataSource() datasource.DataSource {
    return &SubscriberNotificationTemplateDataSource{}
}

// SubscriberNotificationTemplateDataSource defines the data source implementation.
type SubscriberNotificationTemplateDataSource struct {
    client *Client
}

// SubscriberNotificationTemplateDataSourceModel describes the data source data model.
type SubscriberNotificationTemplateDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    EventType types.String `tfsdk:"event_type"`
    NotificationMethod types.String `tfsdk:"notification_method"`
    EmailSubject types.String `tfsdk:"email_subject"`
    TemplateBody types.String `tfsdk:"template_body"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SubscriberNotificationTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_subscriber_notification_template"
}

func (d *SubscriberNotificationTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage custom notification templates for status page subscribers. These templates can be used to customize the notifications sent to subscribers via Email, SMS, Slack, Microsoft Teams, and Webhooks. Look up an existing subscriber notification template by `id`, or by any of its other arguments (`created_by_user_id`, `email_subject`, `event_type`, ...): each one set must match, and exactly one subscriber notification template may match them all.",

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
            "template_name": schema.StringAttribute{
                MarkdownDescription: "A friendly name for this notification template.",
                Optional: true,
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "A description for this notification template.",
                Optional: true,
                Computed: true,
            },
            "event_type": schema.StringAttribute{
                MarkdownDescription: "The type of event this template is for (e.g., Incident Created, Announcement Created).",
                Optional: true,
                Computed: true,
            },
            "notification_method": schema.StringAttribute{
                MarkdownDescription: "The notification method this template is for (Email, SMS, Slack, Microsoft Teams, Webhook).",
                Optional: true,
                Computed: true,
            },
            "email_subject": schema.StringAttribute{
                MarkdownDescription: "The subject line for email notifications. Only used for Email notification method.",
                Optional: true,
                Computed: true,
            },
            "template_body": schema.StringAttribute{
                MarkdownDescription: "The template body content. For Email: HTML template. For SMS: Plain text. For Slack/Teams: Markdown.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *SubscriberNotificationTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SubscriberNotificationTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SubscriberNotificationTemplateDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.TemplateName.IsNull() && !data.TemplateName.IsUnknown() {
        filters["templateName"] = data.TemplateName.ValueString()
        filterNames = append(filterNames, "template_name = "+fmt.Sprintf("%q", data.TemplateName.ValueString()))
    }
    if !data.TemplateDescription.IsNull() && !data.TemplateDescription.IsUnknown() {
        filters["templateDescription"] = data.TemplateDescription.ValueString()
        filterNames = append(filterNames, "template_description = "+fmt.Sprintf("%q", data.TemplateDescription.ValueString()))
    }
    if !data.EventType.IsNull() && !data.EventType.IsUnknown() {
        filters["eventType"] = data.EventType.ValueString()
        filterNames = append(filterNames, "event_type = "+fmt.Sprintf("%q", data.EventType.ValueString()))
    }
    if !data.NotificationMethod.IsNull() && !data.NotificationMethod.IsUnknown() {
        filters["notificationMethod"] = data.NotificationMethod.ValueString()
        filterNames = append(filterNames, "notification_method = "+fmt.Sprintf("%q", data.NotificationMethod.ValueString()))
    }
    if !data.EmailSubject.IsNull() && !data.EmailSubject.IsUnknown() {
        filters["emailSubject"] = data.EmailSubject.ValueString()
        filterNames = append(filterNames, "email_subject = "+fmt.Sprintf("%q", data.EmailSubject.ValueString()))
    }
    if !data.TemplateBody.IsNull() && !data.TemplateBody.IsUnknown() {
        filters["templateBody"] = data.TemplateBody.ValueString()
        filterNames = append(filterNames, "template_body = "+fmt.Sprintf("%q", data.TemplateBody.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the subscriber notification template up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the subscriber notification template up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "templateName": true,
        "templateDescription": true,
        "eventType": true,
        "notificationMethod": true,
        "emailSubject": true,
        "templateBody": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-subscriber-notification-template/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read subscriber_notification_template, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No subscriber notification template found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read subscriber_notification_template: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-subscriber-notification-template/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list subscriber_notification_template, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list subscriber_notification_template: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No subscriber notification template matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one subscriber notification template matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for subscriber_notification_template.")
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
    if obj, ok := item["templateName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateName = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateName = types.StringNull()
        }
    } else if val, ok := item["templateName"].(string); ok {
        data.TemplateName = types.StringValue(val)
    } else {
        data.TemplateName = types.StringNull()
    }
    if obj, ok := item["templateDescription"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateDescription = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateDescription = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateDescription = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateDescription = types.StringNull()
        }
    } else if val, ok := item["templateDescription"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    } else {
        data.TemplateDescription = types.StringNull()
    }
    if obj, ok := item["eventType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventType = types.StringValue(string(jsonBytes))
        } else {
            data.EventType = types.StringNull()
        }
    } else if val, ok := item["eventType"].(string); ok {
        data.EventType = types.StringValue(val)
    } else {
        data.EventType = types.StringNull()
    }
    if obj, ok := item["notificationMethod"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotificationMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NotificationMethod = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NotificationMethod = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NotificationMethod = types.StringValue(string(jsonBytes))
        } else {
            data.NotificationMethod = types.StringNull()
        }
    } else if val, ok := item["notificationMethod"].(string); ok {
        data.NotificationMethod = types.StringValue(val)
    } else {
        data.NotificationMethod = types.StringNull()
    }
    if obj, ok := item["emailSubject"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EmailSubject = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EmailSubject = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EmailSubject = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EmailSubject = types.StringValue(string(jsonBytes))
        } else {
            data.EmailSubject = types.StringNull()
        }
    } else if val, ok := item["emailSubject"].(string); ok {
        data.EmailSubject = types.StringValue(val)
    } else {
        data.EmailSubject = types.StringNull()
    }
    if obj, ok := item["templateBody"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TemplateBody = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TemplateBody = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TemplateBody = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TemplateBody = types.StringValue(string(jsonBytes))
        } else {
            data.TemplateBody = types.StringNull()
        }
    } else if val, ok := item["templateBody"].(string); ok {
        data.TemplateBody = types.StringValue(val)
    } else {
        data.TemplateBody = types.StringNull()
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

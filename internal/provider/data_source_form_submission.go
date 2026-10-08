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
var _ datasource.DataSource = &FormSubmissionDataSource{}

func NewFormSubmissionDataSource() datasource.DataSource {
    return &FormSubmissionDataSource{}
}

// FormSubmissionDataSource defines the data source implementation.
type FormSubmissionDataSource struct {
    client *Client
}

// FormSubmissionDataSourceModel describes the data source data model.
type FormSubmissionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    FormId types.String `tfsdk:"form_id"`
    Answers types.String `tfsdk:"answers"`
    SubmitterName types.String `tfsdk:"submitter_name"`
    SubmitterEmail types.String `tfsdk:"submitter_email"`
    TargetType types.String `tfsdk:"target_type"`
    IncidentId types.String `tfsdk:"incident_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
}

func (d *FormSubmissionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_form_submission"
}

func (d *FormSubmissionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident or scheduled maintenance event each one created. Look up an existing form submission by `id`, or by any of its other arguments (`form_id`, `incident_id`, `scheduled_maintenance_id`, ...): each one set must match, and exactly one form submission may match them all.",

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
            "form_id": schema.StringAttribute{
                MarkdownDescription: "ID of the form this submission was made through. The ID of a `oneuptime_form`.",
                Optional: true,
                Computed: true,
            },
            "answers": schema.StringAttribute{
                MarkdownDescription: "Every question the submitter answered, in the form's order: the question's id (fieldId), its label when the form was submitted, the stored value, and the value as a person reads it (displayValue). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "submitter_name": schema.StringAttribute{
                MarkdownDescription: "The name the submitter gave, as they typed it. Empty when the form does not ask for it, or lets people leave it out and they did.",
                Optional: true,
                Computed: true,
            },
            "submitter_email": schema.StringAttribute{
                MarkdownDescription: "The email address the submitter gave. It is not verified, and nothing is sent to it. Empty when the form does not ask for it, or lets people leave it out and they did.",
                Computed: true,
            },
            "target_type": schema.StringAttribute{
                MarkdownDescription: "What the submission created: Incident, or ScheduledMaintenance (a scheduled maintenance event).",
                Optional: true,
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident this submission created. Empty for a form that schedules maintenance, and once that incident is deleted. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "ID of the scheduled maintenance event this submission created. Empty for a form that creates incidents, and once that event is deleted. The ID of a `oneuptime_scheduled_maintenance_event`.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *FormSubmissionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FormSubmissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data FormSubmissionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.FormId.IsNull() && !data.FormId.IsUnknown() {
        filters["formId"] = data.FormId.ValueString()
        filterNames = append(filterNames, "form_id = "+fmt.Sprintf("%q", data.FormId.ValueString()))
    }
    if !data.SubmitterName.IsNull() && !data.SubmitterName.IsUnknown() {
        filters["submitterName"] = data.SubmitterName.ValueString()
        filterNames = append(filterNames, "submitter_name = "+fmt.Sprintf("%q", data.SubmitterName.ValueString()))
    }
    if !data.TargetType.IsNull() && !data.TargetType.IsUnknown() {
        filters["targetType"] = data.TargetType.ValueString()
        filterNames = append(filterNames, "target_type = "+fmt.Sprintf("%q", data.TargetType.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.ScheduledMaintenanceId.IsNull() && !data.ScheduledMaintenanceId.IsUnknown() {
        filters["scheduledMaintenanceId"] = data.ScheduledMaintenanceId.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_id = "+fmt.Sprintf("%q", data.ScheduledMaintenanceId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the form submission up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the form submission up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "formId": true,
        "answers": true,
        "submitterName": true,
        "submitterEmail": true,
        "targetType": true,
        "incidentId": true,
        "scheduledMaintenanceId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/form-submission/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read form_submission, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No form submission found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read form_submission: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/form-submission/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list form_submission, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list form_submission: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No form submission matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one form submission matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for form_submission.")
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
    if obj, ok := item["formId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FormId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FormId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FormId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FormId = types.StringValue(string(jsonBytes))
        } else {
            data.FormId = types.StringNull()
        }
    } else if val, ok := item["formId"].(string); ok {
        data.FormId = types.StringValue(val)
    } else {
        data.FormId = types.StringNull()
    }
    if obj, ok := item["answers"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Answers = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Answers = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Answers = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Answers = types.StringValue(string(jsonBytes))
        } else {
            data.Answers = types.StringNull()
        }
    } else if val, ok := item["answers"].(string); ok {
        data.Answers = types.StringValue(val)
    } else {
        data.Answers = types.StringNull()
    }
    if obj, ok := item["submitterName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubmitterName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubmitterName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubmitterName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubmitterName = types.StringValue(string(jsonBytes))
        } else {
            data.SubmitterName = types.StringNull()
        }
    } else if val, ok := item["submitterName"].(string); ok {
        data.SubmitterName = types.StringValue(val)
    } else {
        data.SubmitterName = types.StringNull()
    }
    if obj, ok := item["submitterEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubmitterEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubmitterEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubmitterEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubmitterEmail = types.StringValue(string(jsonBytes))
        } else {
            data.SubmitterEmail = types.StringNull()
        }
    } else if val, ok := item["submitterEmail"].(string); ok {
        data.SubmitterEmail = types.StringValue(val)
    } else {
        data.SubmitterEmail = types.StringNull()
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
    } else if val, ok := item["targetType"].(string); ok {
        data.TargetType = types.StringValue(val)
    } else {
        data.TargetType = types.StringNull()
    }
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := item["scheduledMaintenanceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IncidentResource{}
var _ resource.ResourceWithImportState = &IncidentResource{}

func NewIncidentResource() resource.Resource {
    return &IncidentResource{}
}

// IncidentResource defines the resource implementation.
type IncidentResource struct {
    client *Client
}

// IncidentResourceModel describes the resource data model.
type IncidentResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Title types.String `tfsdk:"title"`
    Description types.String `tfsdk:"description"`
    Monitors types.List `tfsdk:"monitors"`
    OnCallDutyPolicies types.List `tfsdk:"on_call_duty_policies"`
    Labels types.List `tfsdk:"labels"`
    CurrentIncidentStateId types.String `tfsdk:"current_incident_state_id"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    ChangeMonitorStatusToId types.String `tfsdk:"change_monitor_status_to_id"`
    ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated types.Bool `tfsdk:"should_status_page_subscribers_be_notified_on_incident_created"`
    CustomFields types.Map `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    RootCause types.String `tfsdk:"root_cause"`
    RemediationNotes types.String `tfsdk:"remediation_notes"`
    TelemetryQuery types.Map `tfsdk:"telemetry_query"`
    IncidentNumber types.Number `tfsdk:"incident_number"`
    IsVisibleOnStatusPage types.Bool `tfsdk:"is_visible_on_status_page"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsStatusPageSubscribersNotifiedOnIncidentCreated types.Bool `tfsdk:"is_status_page_subscribers_notified_on_incident_created"`
    CreatedStateLog types.Map `tfsdk:"created_state_log"`
    CreatedCriteriaId types.String `tfsdk:"created_criteria_id"`
    CreatedIncidentTemplateId types.String `tfsdk:"created_incident_template_id"`
    CreatedByProbeId types.String `tfsdk:"created_by_probe_id"`
    IsCreatedAutomatically types.Bool `tfsdk:"is_created_automatically"`
}

func (r *IncidentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident"
}

func (r *IncidentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incident resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "title": schema.StringAttribute{
                MarkdownDescription: "Title",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description",
                Optional: true,
            },
            "monitors": schema.ListAttribute{
                MarkdownDescription: "Monitors",
                Optional: true,
                ElementType: types.StringType,
            },
            "on_call_duty_policies": schema.ListAttribute{
                MarkdownDescription: "On-Call Duty Policies",
                Optional: true,
                ElementType: types.StringType,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Labels",
                Optional: true,
                ElementType: types.StringType,
            },
            "current_incident_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "change_monitor_status_to_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "should_status_page_subscribers_be_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Should subscribers be notified?",
                Optional: true,
            },
            "custom_fields": schema.MapAttribute{
                MarkdownDescription: "Custom Fields",
                Optional: true,
                ElementType: types.StringType,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are Owners Notified Of Resource Creation?",
                Optional: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "Root Cause",
                Optional: true,
            },
            "remediation_notes": schema.StringAttribute{
                MarkdownDescription: "Remediation Notes",
                Optional: true,
            },
            "telemetry_query": schema.MapAttribute{
                MarkdownDescription: "Telemetry Query",
                Optional: true,
                ElementType: types.StringType,
            },
            "incident_number": schema.NumberAttribute{
                MarkdownDescription: "Incident Number",
                Optional: true,
            },
            "is_visible_on_status_page": schema.BoolAttribute{
                MarkdownDescription: "Should be visible on status page?",
                Optional: true,
            },
            "created_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "updated_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "deleted_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Incident], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_status_page_subscribers_notified_on_incident_created": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_state_log": schema.MapAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_criteria_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_incident_template_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_probe_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_created_automatically": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incident], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *IncidentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *IncidentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data IncidentResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    incidentRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "title": data.Title.ValueString(),
        "description": data.Description.ValueString(),
        "monitors": r.convertTerraformListToInterface(data.Monitors),
        "onCallDutyPolicies": r.convertTerraformListToInterface(data.OnCallDutyPolicies),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "currentIncidentStateId": data.CurrentIncidentStateId.ValueString(),
        "incidentSeverityId": data.IncidentSeverityId.ValueString(),
        "changeMonitorStatusToId": data.ChangeMonitorStatusToId.ValueString(),
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.ValueBool(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "rootCause": data.RootCause.ValueString(),
        "remediationNotes": data.RemediationNotes.ValueString(),
        "telemetryQuery": r.convertTerraformMapToInterface(data.TelemetryQuery),
        "incidentNumber": data.IncidentNumber.ValueBigFloat(),
        "isVisibleOnStatusPage": data.IsVisibleOnStatusPage.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/incident", incidentRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incident, got error: %s", err))
        return
    }

    var incidentResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incidentResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incidentResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Monitors = listValue
    } else if dataMap["monitors"] == nil {
        data.Monitors = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.OnCallDutyPolicies = listValue
    } else if dataMap["onCallDutyPolicies"] == nil {
        data.OnCallDutyPolicies = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["currentIncidentStateId"].(string); ok && val != "" {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if val, ok := dataMap["incidentSeverityId"].(string); ok && val != "" {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if val, ok := dataMap["remediationNotes"].(string); ok && val != "" {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if val, ok := dataMap["telemetryQuery"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryQuery = mapValue
    } else if dataMap["telemetryQuery"] == nil {
        data.TelemetryQuery = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if dataMap["incidentNumber"] == nil {
        data.IncidentNumber = types.NumberNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else if dataMap["isVisibleOnStatusPage"] == nil {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isStatusPageSubscribersNotifiedOnIncidentCreated"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolValue(val)
    } else if dataMap["isStatusPageSubscribersNotifiedOnIncidentCreated"] == nil {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolNull()
    }
    if val, ok := dataMap["createdStateLog"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedStateLog = mapValue
    } else if dataMap["createdStateLog"] == nil {
        data.CreatedStateLog = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["createdCriteriaId"].(string); ok && val != "" {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if val, ok := dataMap["createdIncidentTemplateId"].(string); ok && val != "" {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
    }
    if val, ok := dataMap["createdByProbeId"].(string); ok && val != "" {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := dataMap["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    } else if dataMap["isCreatedAutomatically"] == nil {
        data.IsCreatedAutomatically = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data IncidentResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "monitors": true,
        "onCallDutyPolicies": true,
        "labels": true,
        "currentIncidentStateId": true,
        "incidentSeverityId": true,
        "changeMonitorStatusToId": true,
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "rootCause": true,
        "remediationNotes": true,
        "telemetryQuery": true,
        "incidentNumber": true,
        "isVisibleOnStatusPage": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "isStatusPageSubscribersNotifiedOnIncidentCreated": true,
        "createdStateLog": true,
        "createdCriteriaId": true,
        "createdIncidentTemplateId": true,
        "createdByProbeId": true,
        "isCreatedAutomatically": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/incident/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var incidentResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := incidentResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = incidentResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Monitors = listValue
    } else if dataMap["monitors"] == nil {
        data.Monitors = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.OnCallDutyPolicies = listValue
    } else if dataMap["onCallDutyPolicies"] == nil {
        data.OnCallDutyPolicies = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["currentIncidentStateId"].(string); ok && val != "" {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if val, ok := dataMap["incidentSeverityId"].(string); ok && val != "" {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if val, ok := dataMap["remediationNotes"].(string); ok && val != "" {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if val, ok := dataMap["telemetryQuery"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryQuery = mapValue
    } else if dataMap["telemetryQuery"] == nil {
        data.TelemetryQuery = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if dataMap["incidentNumber"] == nil {
        data.IncidentNumber = types.NumberNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else if dataMap["isVisibleOnStatusPage"] == nil {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isStatusPageSubscribersNotifiedOnIncidentCreated"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolValue(val)
    } else if dataMap["isStatusPageSubscribersNotifiedOnIncidentCreated"] == nil {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolNull()
    }
    if val, ok := dataMap["createdStateLog"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedStateLog = mapValue
    } else if dataMap["createdStateLog"] == nil {
        data.CreatedStateLog = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["createdCriteriaId"].(string); ok && val != "" {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if val, ok := dataMap["createdIncidentTemplateId"].(string); ok && val != "" {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
    }
    if val, ok := dataMap["createdByProbeId"].(string); ok && val != "" {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := dataMap["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    } else if dataMap["isCreatedAutomatically"] == nil {
        data.IsCreatedAutomatically = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data IncidentResourceModel
    var state IncidentResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    incidentRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "title": data.Title.ValueString(),
        "description": data.Description.ValueString(),
        "monitors": r.convertTerraformListToInterface(data.Monitors),
        "onCallDutyPolicies": r.convertTerraformListToInterface(data.OnCallDutyPolicies),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "currentIncidentStateId": data.CurrentIncidentStateId.ValueString(),
        "incidentSeverityId": data.IncidentSeverityId.ValueString(),
        "changeMonitorStatusToId": data.ChangeMonitorStatusToId.ValueString(),
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated.ValueBool(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "rootCause": data.RootCause.ValueString(),
        "remediationNotes": data.RemediationNotes.ValueString(),
        "telemetryQuery": r.convertTerraformMapToInterface(data.TelemetryQuery),
        "incidentNumber": data.IncidentNumber.ValueBigFloat(),
        "isVisibleOnStatusPage": data.IsVisibleOnStatusPage.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/incident/" + data.Id.ValueString() + "", incidentRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incident, got error: %s", err))
        return
    }

    // Parse the update response
    var incidentResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &incidentResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "title": true,
        "description": true,
        "monitors": true,
        "onCallDutyPolicies": true,
        "labels": true,
        "currentIncidentStateId": true,
        "incidentSeverityId": true,
        "changeMonitorStatusToId": true,
        "shouldStatusPageSubscribersBeNotifiedOnIncidentCreated": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "rootCause": true,
        "remediationNotes": true,
        "telemetryQuery": true,
        "incidentNumber": true,
        "isVisibleOnStatusPage": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "isStatusPageSubscribersNotifiedOnIncidentCreated": true,
        "createdStateLog": true,
        "createdCriteriaId": true,
        "createdIncidentTemplateId": true,
        "createdByProbeId": true,
        "isCreatedAutomatically": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/incident/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incident read response, got error: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["title"].(string); ok && val != "" {
        data.Title = types.StringValue(val)
    } else {
        data.Title = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["monitors"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Monitors = listValue
    } else if dataMap["monitors"] == nil {
        data.Monitors = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.OnCallDutyPolicies = listValue
    } else if dataMap["onCallDutyPolicies"] == nil {
        data.OnCallDutyPolicies = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["currentIncidentStateId"].(string); ok && val != "" {
        data.CurrentIncidentStateId = types.StringValue(val)
    } else {
        data.CurrentIncidentStateId = types.StringNull()
    }
    if val, ok := dataMap["incidentSeverityId"].(string); ok && val != "" {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := dataMap["changeMonitorStatusToId"].(string); ok && val != "" {
        data.ChangeMonitorStatusToId = types.StringValue(val)
    } else {
        data.ChangeMonitorStatusToId = types.StringNull()
    }
    if val, ok := dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"].(bool); ok {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolValue(val)
    } else if dataMap["shouldStatusPageSubscribersBeNotifiedOnIncidentCreated"] == nil {
        data.ShouldStatusPageSubscribersBeNotifiedOnIncidentCreated = types.BoolNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if val, ok := dataMap["remediationNotes"].(string); ok && val != "" {
        data.RemediationNotes = types.StringValue(val)
    } else {
        data.RemediationNotes = types.StringNull()
    }
    if val, ok := dataMap["telemetryQuery"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryQuery = mapValue
    } else if dataMap["telemetryQuery"] == nil {
        data.TelemetryQuery = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["incidentNumber"].(float64); ok {
        data.IncidentNumber = types.NumberValue(big.NewFloat(val))
    } else if dataMap["incidentNumber"] == nil {
        data.IncidentNumber = types.NumberNull()
    }
    if val, ok := dataMap["isVisibleOnStatusPage"].(bool); ok {
        data.IsVisibleOnStatusPage = types.BoolValue(val)
    } else if dataMap["isVisibleOnStatusPage"] == nil {
        data.IsVisibleOnStatusPage = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isStatusPageSubscribersNotifiedOnIncidentCreated"].(bool); ok {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolValue(val)
    } else if dataMap["isStatusPageSubscribersNotifiedOnIncidentCreated"] == nil {
        data.IsStatusPageSubscribersNotifiedOnIncidentCreated = types.BoolNull()
    }
    if val, ok := dataMap["createdStateLog"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedStateLog = mapValue
    } else if dataMap["createdStateLog"] == nil {
        data.CreatedStateLog = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["createdCriteriaId"].(string); ok && val != "" {
        data.CreatedCriteriaId = types.StringValue(val)
    } else {
        data.CreatedCriteriaId = types.StringNull()
    }
    if val, ok := dataMap["createdIncidentTemplateId"].(string); ok && val != "" {
        data.CreatedIncidentTemplateId = types.StringValue(val)
    } else {
        data.CreatedIncidentTemplateId = types.StringNull()
    }
    if val, ok := dataMap["createdByProbeId"].(string); ok && val != "" {
        data.CreatedByProbeId = types.StringValue(val)
    } else {
        data.CreatedByProbeId = types.StringNull()
    }
    if val, ok := dataMap["isCreatedAutomatically"].(bool); ok {
        data.IsCreatedAutomatically = types.BoolValue(val)
    } else if dataMap["isCreatedAutomatically"] == nil {
        data.IsCreatedAutomatically = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IncidentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data IncidentResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/incident/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incident, got error: %s", err))
        return
    }
}


func (r *IncidentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *IncidentResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *IncidentResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

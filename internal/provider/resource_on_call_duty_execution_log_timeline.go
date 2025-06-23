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
var _ resource.Resource = &OnCallDutyExecutionLogTimelineResource{}
var _ resource.ResourceWithImportState = &OnCallDutyExecutionLogTimelineResource{}

func NewOnCallDutyExecutionLogTimelineResource() resource.Resource {
    return &OnCallDutyExecutionLogTimelineResource{}
}

// OnCallDutyExecutionLogTimelineResource defines the resource implementation.
type OnCallDutyExecutionLogTimelineResource struct {
    client *Client
}

// OnCallDutyExecutionLogTimelineResourceModel describes the resource data model.
type OnCallDutyExecutionLogTimelineResourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    OnCallDutyPolicyExecutionLogId types.String `tfsdk:"on_call_duty_policy_execution_log_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    UserNotificationEventType types.String `tfsdk:"user_notification_event_type"`
    AlertSentToUserId types.String `tfsdk:"alert_sent_to_user_id"`
    UserBelongsToTeamId types.String `tfsdk:"user_belongs_to_team_id"`
    OnCallDutyScheduleId types.String `tfsdk:"on_call_duty_schedule_id"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsAcknowledged types.Bool `tfsdk:"is_acknowledged"`
    AcknowledgedAt types.Map `tfsdk:"acknowledged_at"`
    OverridedByUserId types.String `tfsdk:"overrided_by_user_id"`
}

func (r *OnCallDutyExecutionLogTimelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_duty_execution_log_timeline"
}

func (r *OnCallDutyExecutionLogTimelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "on_call_duty_execution_log_timeline resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "triggered_by_alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_execution_log_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_notification_event_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "alert_sent_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_belongs_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_schedule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_acknowledged": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "acknowledged_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "overrided_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (r *OnCallDutyExecutionLogTimelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *OnCallDutyExecutionLogTimelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data OnCallDutyExecutionLogTimelineResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    onCallDutyExecutionLogTimelineRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/on-call-duty-policy-execution-log-timeline/count", onCallDutyExecutionLogTimelineRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create on_call_duty_execution_log_timeline, got error: %s", err))
        return
    }

    var onCallDutyExecutionLogTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &onCallDutyExecutionLogTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_duty_execution_log_timeline response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := onCallDutyExecutionLogTimelineResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = onCallDutyExecutionLogTimelineResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    if val, ok := dataMap["onCallDutyPolicyId"].(string); ok && val != "" {
        data.OnCallDutyPolicyId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyId = types.StringNull()
    }
    if val, ok := dataMap["triggeredByIncidentId"].(string); ok && val != "" {
        data.TriggeredByIncidentId = types.StringValue(val)
    } else {
        data.TriggeredByIncidentId = types.StringNull()
    }
    if val, ok := dataMap["triggeredByAlertId"].(string); ok && val != "" {
        data.TriggeredByAlertId = types.StringValue(val)
    } else {
        data.TriggeredByAlertId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicyExecutionLogId"].(string); ok && val != "" {
        data.OnCallDutyPolicyExecutionLogId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyExecutionLogId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicyEscalationRuleId"].(string); ok && val != "" {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["userNotificationEventType"].(string); ok && val != "" {
        data.UserNotificationEventType = types.StringValue(val)
    } else {
        data.UserNotificationEventType = types.StringNull()
    }
    if val, ok := dataMap["alertSentToUserId"].(string); ok && val != "" {
        data.AlertSentToUserId = types.StringValue(val)
    } else {
        data.AlertSentToUserId = types.StringNull()
    }
    if val, ok := dataMap["userBelongsToTeamId"].(string); ok && val != "" {
        data.UserBelongsToTeamId = types.StringValue(val)
    } else {
        data.UserBelongsToTeamId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyScheduleId"].(string); ok && val != "" {
        data.OnCallDutyScheduleId = types.StringValue(val)
    } else {
        data.OnCallDutyScheduleId = types.StringNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isAcknowledged"].(bool); ok {
        data.IsAcknowledged = types.BoolValue(val)
    } else if dataMap["isAcknowledged"] == nil {
        data.IsAcknowledged = types.BoolNull()
    }
    if val, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.AcknowledgedAt = mapValue
    } else if dataMap["acknowledgedAt"] == nil {
        data.AcknowledgedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["overridedByUserId"].(string); ok && val != "" {
        data.OverridedByUserId = types.StringValue(val)
    } else {
        data.OverridedByUserId = types.StringNull()
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

func (r *OnCallDutyExecutionLogTimelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data OnCallDutyExecutionLogTimelineResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "onCallDutyPolicyId": true,
        "triggeredByIncidentId": true,
        "triggeredByAlertId": true,
        "onCallDutyPolicyExecutionLogId": true,
        "onCallDutyPolicyEscalationRuleId": true,
        "userNotificationEventType": true,
        "alertSentToUserId": true,
        "userBelongsToTeamId": true,
        "onCallDutyScheduleId": true,
        "statusMessage": true,
        "status": true,
        "createdByUserId": true,
        "isAcknowledged": true,
        "acknowledgedAt": true,
        "overridedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/on-call-duty-policy-execution-log-timeline/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_duty_execution_log_timeline, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var onCallDutyExecutionLogTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &onCallDutyExecutionLogTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_duty_execution_log_timeline response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := onCallDutyExecutionLogTimelineResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = onCallDutyExecutionLogTimelineResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    if val, ok := dataMap["onCallDutyPolicyId"].(string); ok && val != "" {
        data.OnCallDutyPolicyId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyId = types.StringNull()
    }
    if val, ok := dataMap["triggeredByIncidentId"].(string); ok && val != "" {
        data.TriggeredByIncidentId = types.StringValue(val)
    } else {
        data.TriggeredByIncidentId = types.StringNull()
    }
    if val, ok := dataMap["triggeredByAlertId"].(string); ok && val != "" {
        data.TriggeredByAlertId = types.StringValue(val)
    } else {
        data.TriggeredByAlertId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicyExecutionLogId"].(string); ok && val != "" {
        data.OnCallDutyPolicyExecutionLogId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyExecutionLogId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicyEscalationRuleId"].(string); ok && val != "" {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["userNotificationEventType"].(string); ok && val != "" {
        data.UserNotificationEventType = types.StringValue(val)
    } else {
        data.UserNotificationEventType = types.StringNull()
    }
    if val, ok := dataMap["alertSentToUserId"].(string); ok && val != "" {
        data.AlertSentToUserId = types.StringValue(val)
    } else {
        data.AlertSentToUserId = types.StringNull()
    }
    if val, ok := dataMap["userBelongsToTeamId"].(string); ok && val != "" {
        data.UserBelongsToTeamId = types.StringValue(val)
    } else {
        data.UserBelongsToTeamId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyScheduleId"].(string); ok && val != "" {
        data.OnCallDutyScheduleId = types.StringValue(val)
    } else {
        data.OnCallDutyScheduleId = types.StringNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["isAcknowledged"].(bool); ok {
        data.IsAcknowledged = types.BoolValue(val)
    } else if dataMap["isAcknowledged"] == nil {
        data.IsAcknowledged = types.BoolNull()
    }
    if val, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.AcknowledgedAt = mapValue
    } else if dataMap["acknowledgedAt"] == nil {
        data.AcknowledgedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["overridedByUserId"].(string); ok && val != "" {
        data.OverridedByUserId = types.StringValue(val)
    } else {
        data.OverridedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OnCallDutyExecutionLogTimelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *OnCallDutyExecutionLogTimelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *OnCallDutyExecutionLogTimelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *OnCallDutyExecutionLogTimelineResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *OnCallDutyExecutionLogTimelineResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

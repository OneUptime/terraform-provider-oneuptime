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
var _ resource.Resource = &OnCallDutyExecutionLogResource{}
var _ resource.ResourceWithImportState = &OnCallDutyExecutionLogResource{}

func NewOnCallDutyExecutionLogResource() resource.Resource {
    return &OnCallDutyExecutionLogResource{}
}

// OnCallDutyExecutionLogResource defines the resource implementation.
type OnCallDutyExecutionLogResource struct {
    client *Client
}

// OnCallDutyExecutionLogResourceModel describes the resource data model.
type OnCallDutyExecutionLogResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    TriggeredByIncidentId types.String `tfsdk:"triggered_by_incident_id"`
    TriggeredByAlertId types.String `tfsdk:"triggered_by_alert_id"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    UserNotificationEventType types.String `tfsdk:"user_notification_event_type"`
    AcknowledgedByUserId types.String `tfsdk:"acknowledged_by_user_id"`
    AcknowledgedAt types.Map `tfsdk:"acknowledged_at"`
    AcknowledgedByTeamId types.String `tfsdk:"acknowledged_by_team_id"`
    LastExecutedEscalationRuleOrder types.Number `tfsdk:"last_executed_escalation_rule_order"`
    LastEscalationRuleExecutedAt types.Map `tfsdk:"last_escalation_rule_executed_at"`
    LastExecutedEscalationRuleId types.String `tfsdk:"last_executed_escalation_rule_id"`
    ExecuteNextEscalationRuleInMinutes types.Number `tfsdk:"execute_next_escalation_rule_in_minutes"`
    OnCallPolicyExecutionRepeatCount types.Number `tfsdk:"on_call_policy_execution_repeat_count"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (r *OnCallDutyExecutionLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_on_call_duty_execution_log"
}

func (r *OnCallDutyExecutionLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "on_call_duty_execution_log resource",

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
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "triggered_by_incident_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "triggered_by_alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status",
                Required: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Status Message",
                Required: true,
            },
            "user_notification_event_type": schema.StringAttribute{
                MarkdownDescription: "Notification Event Type",
                Required: true,
            },
            "acknowledged_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "acknowledged_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "acknowledged_by_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "last_executed_escalation_rule_order": schema.NumberAttribute{
                MarkdownDescription: "Executed Escalation Rule Order",
                Optional: true,
            },
            "last_escalation_rule_executed_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "last_executed_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "execute_next_escalation_rule_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Execute next escalation rule in minutes",
                Optional: true,
            },
            "on_call_policy_execution_repeat_count": schema.NumberAttribute{
                MarkdownDescription: "On-Call Policy Execution Repeat Count",
                Optional: true,
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
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

func (r *OnCallDutyExecutionLogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *OnCallDutyExecutionLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data OnCallDutyExecutionLogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    onCallDutyExecutionLogRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "onCallDutyPolicyId": data.OnCallDutyPolicyId.ValueString(),
        "triggeredByIncidentId": data.TriggeredByIncidentId.ValueString(),
        "triggeredByAlertId": data.TriggeredByAlertId.ValueString(),
        "status": data.Status.ValueString(),
        "statusMessage": data.StatusMessage.ValueString(),
        "userNotificationEventType": data.UserNotificationEventType.ValueString(),
        "acknowledgedByUserId": data.AcknowledgedByUserId.ValueString(),
        "acknowledgedAt": r.convertTerraformMapToInterface(data.AcknowledgedAt),
        "acknowledgedByTeamId": data.AcknowledgedByTeamId.ValueString(),
        "lastExecutedEscalationRuleOrder": data.LastExecutedEscalationRuleOrder.ValueBigFloat(),
        "lastEscalationRuleExecutedAt": r.convertTerraformMapToInterface(data.LastEscalationRuleExecutedAt),
        "lastExecutedEscalationRuleId": data.LastExecutedEscalationRuleId.ValueString(),
        "executeNextEscalationRuleInMinutes": data.ExecuteNextEscalationRuleInMinutes.ValueBigFloat(),
        "onCallPolicyExecutionRepeatCount": data.OnCallPolicyExecutionRepeatCount.ValueBigFloat(),
        "triggeredByUserId": data.TriggeredByUserId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/on-call-duty-policy-execution-log", onCallDutyExecutionLogRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create on_call_duty_execution_log, got error: %s", err))
        return
    }

    var onCallDutyExecutionLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &onCallDutyExecutionLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_duty_execution_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := onCallDutyExecutionLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = onCallDutyExecutionLogResponse
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
    if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["userNotificationEventType"].(string); ok && val != "" {
        data.UserNotificationEventType = types.StringValue(val)
    } else {
        data.UserNotificationEventType = types.StringNull()
    }
    if val, ok := dataMap["acknowledgedByUserId"].(string); ok && val != "" {
        data.AcknowledgedByUserId = types.StringValue(val)
    } else {
        data.AcknowledgedByUserId = types.StringNull()
    }
    if val, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.AcknowledgedAt = mapValue
    } else if dataMap["acknowledgedAt"] == nil {
        data.AcknowledgedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["acknowledgedByTeamId"].(string); ok && val != "" {
        data.AcknowledgedByTeamId = types.StringValue(val)
    } else {
        data.AcknowledgedByTeamId = types.StringNull()
    }
    if val, ok := dataMap["lastExecutedEscalationRuleOrder"].(float64); ok {
        data.LastExecutedEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
    } else if dataMap["lastExecutedEscalationRuleOrder"] == nil {
        data.LastExecutedEscalationRuleOrder = types.NumberNull()
    }
    if val, ok := dataMap["lastEscalationRuleExecutedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastEscalationRuleExecutedAt = mapValue
    } else if dataMap["lastEscalationRuleExecutedAt"] == nil {
        data.LastEscalationRuleExecutedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastExecutedEscalationRuleId"].(string); ok && val != "" {
        data.LastExecutedEscalationRuleId = types.StringValue(val)
    } else {
        data.LastExecutedEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["executeNextEscalationRuleInMinutes"].(float64); ok {
        data.ExecuteNextEscalationRuleInMinutes = types.NumberValue(big.NewFloat(val))
    } else if dataMap["executeNextEscalationRuleInMinutes"] == nil {
        data.ExecuteNextEscalationRuleInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["onCallPolicyExecutionRepeatCount"].(float64); ok {
        data.OnCallPolicyExecutionRepeatCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["onCallPolicyExecutionRepeatCount"] == nil {
        data.OnCallPolicyExecutionRepeatCount = types.NumberNull()
    }
    if val, ok := dataMap["triggeredByUserId"].(string); ok && val != "" {
        data.TriggeredByUserId = types.StringValue(val)
    } else {
        data.TriggeredByUserId = types.StringNull()
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
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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

func (r *OnCallDutyExecutionLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data OnCallDutyExecutionLogResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "onCallDutyPolicyId": true,
        "triggeredByIncidentId": true,
        "triggeredByAlertId": true,
        "status": true,
        "statusMessage": true,
        "userNotificationEventType": true,
        "acknowledgedByUserId": true,
        "acknowledgedAt": true,
        "acknowledgedByTeamId": true,
        "lastExecutedEscalationRuleOrder": true,
        "lastEscalationRuleExecutedAt": true,
        "lastExecutedEscalationRuleId": true,
        "executeNextEscalationRuleInMinutes": true,
        "onCallPolicyExecutionRepeatCount": true,
        "triggeredByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/on-call-duty-policy-execution-log/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read on_call_duty_execution_log, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var onCallDutyExecutionLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &onCallDutyExecutionLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse on_call_duty_execution_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := onCallDutyExecutionLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = onCallDutyExecutionLogResponse
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
    if val, ok := dataMap["status"].(string); ok && val != "" {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if val, ok := dataMap["statusMessage"].(string); ok && val != "" {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := dataMap["userNotificationEventType"].(string); ok && val != "" {
        data.UserNotificationEventType = types.StringValue(val)
    } else {
        data.UserNotificationEventType = types.StringNull()
    }
    if val, ok := dataMap["acknowledgedByUserId"].(string); ok && val != "" {
        data.AcknowledgedByUserId = types.StringValue(val)
    } else {
        data.AcknowledgedByUserId = types.StringNull()
    }
    if val, ok := dataMap["acknowledgedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.AcknowledgedAt = mapValue
    } else if dataMap["acknowledgedAt"] == nil {
        data.AcknowledgedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["acknowledgedByTeamId"].(string); ok && val != "" {
        data.AcknowledgedByTeamId = types.StringValue(val)
    } else {
        data.AcknowledgedByTeamId = types.StringNull()
    }
    if val, ok := dataMap["lastExecutedEscalationRuleOrder"].(float64); ok {
        data.LastExecutedEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
    } else if dataMap["lastExecutedEscalationRuleOrder"] == nil {
        data.LastExecutedEscalationRuleOrder = types.NumberNull()
    }
    if val, ok := dataMap["lastEscalationRuleExecutedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastEscalationRuleExecutedAt = mapValue
    } else if dataMap["lastEscalationRuleExecutedAt"] == nil {
        data.LastEscalationRuleExecutedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastExecutedEscalationRuleId"].(string); ok && val != "" {
        data.LastExecutedEscalationRuleId = types.StringValue(val)
    } else {
        data.LastExecutedEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["executeNextEscalationRuleInMinutes"].(float64); ok {
        data.ExecuteNextEscalationRuleInMinutes = types.NumberValue(big.NewFloat(val))
    } else if dataMap["executeNextEscalationRuleInMinutes"] == nil {
        data.ExecuteNextEscalationRuleInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["onCallPolicyExecutionRepeatCount"].(float64); ok {
        data.OnCallPolicyExecutionRepeatCount = types.NumberValue(big.NewFloat(val))
    } else if dataMap["onCallPolicyExecutionRepeatCount"] == nil {
        data.OnCallPolicyExecutionRepeatCount = types.NumberNull()
    }
    if val, ok := dataMap["triggeredByUserId"].(string); ok && val != "" {
        data.TriggeredByUserId = types.StringValue(val)
    } else {
        data.TriggeredByUserId = types.StringNull()
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
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OnCallDutyExecutionLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *OnCallDutyExecutionLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *OnCallDutyExecutionLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *OnCallDutyExecutionLogResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *OnCallDutyExecutionLogResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

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
var _ resource.Resource = &UsersOnCallDutyEscalationRuleResource{}
var _ resource.ResourceWithImportState = &UsersOnCallDutyEscalationRuleResource{}

func NewUsersOnCallDutyEscalationRuleResource() resource.Resource {
    return &UsersOnCallDutyEscalationRuleResource{}
}

// UsersOnCallDutyEscalationRuleResource defines the resource implementation.
type UsersOnCallDutyEscalationRuleResource struct {
    client *Client
}

// UsersOnCallDutyEscalationRuleResourceModel describes the resource data model.
type UsersOnCallDutyEscalationRuleResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    OnCallDutyPolicyEscalationRuleId types.String `tfsdk:"on_call_duty_policy_escalation_rule_id"`
    UserId types.String `tfsdk:"user_id"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *UsersOnCallDutyEscalationRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_users_on_call_duty_escalation_rule"
}

func (r *UsersOnCallDutyEscalationRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "users_on_call_duty_escalation_rule resource",

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
            "on_call_duty_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "user_id": schema.StringAttribute{
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
        },
    }
}

func (r *UsersOnCallDutyEscalationRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *UsersOnCallDutyEscalationRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data UsersOnCallDutyEscalationRuleResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    usersOnCallDutyEscalationRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "onCallDutyPolicyId": data.OnCallDutyPolicyId.ValueString(),
        "onCallDutyPolicyEscalationRuleId": data.OnCallDutyPolicyEscalationRuleId.ValueString(),
        "userId": data.UserId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/on-call-duty-policy-esclation-rule-user", usersOnCallDutyEscalationRuleRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create users_on_call_duty_escalation_rule, got error: %s", err))
        return
    }

    var usersOnCallDutyEscalationRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &usersOnCallDutyEscalationRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse users_on_call_duty_escalation_rule response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := usersOnCallDutyEscalationRuleResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = usersOnCallDutyEscalationRuleResponse
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
    if val, ok := dataMap["onCallDutyPolicyEscalationRuleId"].(string); ok && val != "" {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
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

func (r *UsersOnCallDutyEscalationRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data UsersOnCallDutyEscalationRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "onCallDutyPolicyId": true,
        "onCallDutyPolicyEscalationRuleId": true,
        "userId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/on-call-duty-policy-esclation-rule-user/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users_on_call_duty_escalation_rule, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var usersOnCallDutyEscalationRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &usersOnCallDutyEscalationRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse users_on_call_duty_escalation_rule response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := usersOnCallDutyEscalationRuleResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = usersOnCallDutyEscalationRuleResponse
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
    if val, ok := dataMap["onCallDutyPolicyEscalationRuleId"].(string); ok && val != "" {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsersOnCallDutyEscalationRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data UsersOnCallDutyEscalationRuleResourceModel
    var state UsersOnCallDutyEscalationRuleResourceModel

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
    usersOnCallDutyEscalationRuleRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "onCallDutyPolicyId": data.OnCallDutyPolicyId.ValueString(),
        "onCallDutyPolicyEscalationRuleId": data.OnCallDutyPolicyEscalationRuleId.ValueString(),
        "userId": data.UserId.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/on-call-duty-policy-esclation-rule-user/" + data.Id.ValueString() + "", usersOnCallDutyEscalationRuleRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update users_on_call_duty_escalation_rule, got error: %s", err))
        return
    }

    // Parse the update response
    var usersOnCallDutyEscalationRuleResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &usersOnCallDutyEscalationRuleResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse users_on_call_duty_escalation_rule response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "onCallDutyPolicyId": true,
        "onCallDutyPolicyEscalationRuleId": true,
        "userId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/on-call-duty-policy-esclation-rule-user/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users_on_call_duty_escalation_rule after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse users_on_call_duty_escalation_rule read response, got error: %s", err))
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
    if val, ok := dataMap["onCallDutyPolicyId"].(string); ok && val != "" {
        data.OnCallDutyPolicyId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyId = types.StringNull()
    }
    if val, ok := dataMap["onCallDutyPolicyEscalationRuleId"].(string); ok && val != "" {
        data.OnCallDutyPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.OnCallDutyPolicyEscalationRuleId = types.StringNull()
    }
    if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsersOnCallDutyEscalationRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data UsersOnCallDutyEscalationRuleResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/on-call-duty-policy-esclation-rule-user/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete users_on_call_duty_escalation_rule, got error: %s", err))
        return
    }
}


func (r *UsersOnCallDutyEscalationRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *UsersOnCallDutyEscalationRuleResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *UsersOnCallDutyEscalationRuleResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

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
var _ resource.Resource = &AlertStateTimelineResource{}
var _ resource.ResourceWithImportState = &AlertStateTimelineResource{}

func NewAlertStateTimelineResource() resource.Resource {
    return &AlertStateTimelineResource{}
}

// AlertStateTimelineResource defines the resource implementation.
type AlertStateTimelineResource struct {
    client *Client
}

// AlertStateTimelineResourceModel describes the resource data model.
type AlertStateTimelineResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    AlertId types.String `tfsdk:"alert_id"`
    AlertStateId types.String `tfsdk:"alert_state_id"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    RootCause types.String `tfsdk:"root_cause"`
    EndsAt types.Map `tfsdk:"ends_at"`
    StartsAt types.Map `tfsdk:"starts_at"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    StateChangeLog types.Map `tfsdk:"state_change_log"`
}

func (r *AlertStateTimelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_alert_state_timeline"
}

func (r *AlertStateTimelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "alert_state_timeline resource",

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
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "alert_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Are Owners Notified",
                Optional: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "Root Cause",
                Optional: true,
            },
            "ends_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "starts_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
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
            "state_change_log": schema.MapAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Alert State Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (r *AlertStateTimelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *AlertStateTimelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data AlertStateTimelineResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    alertStateTimelineRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "alertId": data.AlertId.ValueString(),
        "alertStateId": data.AlertStateId.ValueString(),
        "isOwnerNotified": data.IsOwnerNotified.ValueBool(),
        "rootCause": data.RootCause.ValueString(),
        "endsAt": r.convertTerraformMapToInterface(data.EndsAt),
        "startsAt": r.convertTerraformMapToInterface(data.StartsAt),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/alert-state-timeline", alertStateTimelineRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alert_state_timeline, got error: %s", err))
        return
    }

    var alertStateTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &alertStateTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_state_timeline response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := alertStateTimelineResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = alertStateTimelineResponse
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
    if val, ok := dataMap["alertId"].(string); ok && val != "" {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if val, ok := dataMap["alertStateId"].(string); ok && val != "" {
        data.AlertStateId = types.StringValue(val)
    } else {
        data.AlertStateId = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else if dataMap["isOwnerNotified"] == nil {
        data.IsOwnerNotified = types.BoolNull()
    }
    if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if val, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.EndsAt = mapValue
    } else if dataMap["endsAt"] == nil {
        data.EndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StartsAt = mapValue
    } else if dataMap["startsAt"] == nil {
        data.StartsAt = types.MapNull(types.StringType)
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
    if val, ok := dataMap["stateChangeLog"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StateChangeLog = mapValue
    } else if dataMap["stateChangeLog"] == nil {
        data.StateChangeLog = types.MapNull(types.StringType)
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

func (r *AlertStateTimelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data AlertStateTimelineResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "alertId": true,
        "alertStateId": true,
        "isOwnerNotified": true,
        "rootCause": true,
        "endsAt": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "stateChangeLog": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/alert-state-timeline/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_state_timeline, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var alertStateTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &alertStateTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_state_timeline response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := alertStateTimelineResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = alertStateTimelineResponse
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
    if val, ok := dataMap["alertId"].(string); ok && val != "" {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if val, ok := dataMap["alertStateId"].(string); ok && val != "" {
        data.AlertStateId = types.StringValue(val)
    } else {
        data.AlertStateId = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else if dataMap["isOwnerNotified"] == nil {
        data.IsOwnerNotified = types.BoolNull()
    }
    if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if val, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.EndsAt = mapValue
    } else if dataMap["endsAt"] == nil {
        data.EndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StartsAt = mapValue
    } else if dataMap["startsAt"] == nil {
        data.StartsAt = types.MapNull(types.StringType)
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
    if val, ok := dataMap["stateChangeLog"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StateChangeLog = mapValue
    } else if dataMap["stateChangeLog"] == nil {
        data.StateChangeLog = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertStateTimelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data AlertStateTimelineResourceModel
    var state AlertStateTimelineResourceModel

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
    alertStateTimelineRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "alertId": data.AlertId.ValueString(),
        "alertStateId": data.AlertStateId.ValueString(),
        "isOwnerNotified": data.IsOwnerNotified.ValueBool(),
        "rootCause": data.RootCause.ValueString(),
        "endsAt": r.convertTerraformMapToInterface(data.EndsAt),
        "startsAt": r.convertTerraformMapToInterface(data.StartsAt),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/alert-state-timeline/" + data.Id.ValueString() + "", alertStateTimelineRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update alert_state_timeline, got error: %s", err))
        return
    }

    // Parse the update response
    var alertStateTimelineResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &alertStateTimelineResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_state_timeline response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "alertId": true,
        "alertStateId": true,
        "isOwnerNotified": true,
        "rootCause": true,
        "endsAt": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "stateChangeLog": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/alert-state-timeline/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert_state_timeline after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse alert_state_timeline read response, got error: %s", err))
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
    if val, ok := dataMap["alertId"].(string); ok && val != "" {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if val, ok := dataMap["alertStateId"].(string); ok && val != "" {
        data.AlertStateId = types.StringValue(val)
    } else {
        data.AlertStateId = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    } else if dataMap["isOwnerNotified"] == nil {
        data.IsOwnerNotified = types.BoolNull()
    }
    if val, ok := dataMap["rootCause"].(string); ok && val != "" {
        data.RootCause = types.StringValue(val)
    } else {
        data.RootCause = types.StringNull()
    }
    if val, ok := dataMap["endsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.EndsAt = mapValue
    } else if dataMap["endsAt"] == nil {
        data.EndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StartsAt = mapValue
    } else if dataMap["startsAt"] == nil {
        data.StartsAt = types.MapNull(types.StringType)
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
    if val, ok := dataMap["stateChangeLog"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StateChangeLog = mapValue
    } else if dataMap["stateChangeLog"] == nil {
        data.StateChangeLog = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertStateTimelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data AlertStateTimelineResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/alert-state-timeline/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete alert_state_timeline, got error: %s", err))
        return
    }
}


func (r *AlertStateTimelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *AlertStateTimelineResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *AlertStateTimelineResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

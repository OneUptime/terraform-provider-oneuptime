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
var _ resource.Resource = &ProbeResource{}
var _ resource.ResourceWithImportState = &ProbeResource{}

func NewProbeResource() resource.Resource {
    return &ProbeResource{}
}

// ProbeResource defines the resource implementation.
type ProbeResource struct {
    client *Client
}

// ProbeResourceModel describes the resource data model.
type ProbeResourceModel struct {
    Id types.String `tfsdk:"id"`
    Key types.String `tfsdk:"key"`
    Name types.Map `tfsdk:"name"`
    Description types.Map `tfsdk:"description"`
    ProbeVersion types.Map `tfsdk:"probe_version"`
    LastAlive types.Map `tfsdk:"last_alive"`
    IconFileId types.String `tfsdk:"icon_file_id"`
    ProjectId types.String `tfsdk:"project_id"`
    ShouldAutoEnableProbeOnNewMonitors types.Bool `tfsdk:"should_auto_enable_probe_on_new_monitors"`
    Labels types.List `tfsdk:"labels"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    ConnectionStatus types.String `tfsdk:"connection_status"`
}

func (r *ProbeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_probe"
}

func (r *ProbeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "probe resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Probe], Read: [Project Owner, Project Admin], Update: [Project Owner, Project Admin, Project Member, Edit Probe]",
                Required: true,
            },
            "name": schema.MapAttribute{
                MarkdownDescription: "Name object",
                Required: true,
                ElementType: types.StringType,
            },
            "description": schema.MapAttribute{
                MarkdownDescription: "Name object",
                Optional: true,
                ElementType: types.StringType,
            },
            "probe_version": schema.MapAttribute{
                MarkdownDescription: "Version object",
                Required: true,
                ElementType: types.StringType,
            },
            "last_alive": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "icon_file_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "should_auto_enable_probe_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "Auto Enable Probe",
                Optional: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Labels",
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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Public], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "connection_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *ProbeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ProbeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ProbeResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    probeRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "key": data.Key.ValueString(),
        "name": r.convertTerraformMapToInterface(data.Name),
        "description": r.convertTerraformMapToInterface(data.Description),
        "probeVersion": r.convertTerraformMapToInterface(data.ProbeVersion),
        "lastAlive": r.convertTerraformMapToInterface(data.LastAlive),
        "iconFileId": data.IconFileId.ValueString(),
        "projectId": data.ProjectId.ValueString(),
        "shouldAutoEnableProbeOnNewMonitors": data.ShouldAutoEnableProbeOnNewMonitors.ValueBool(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/probe", probeRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create probe, got error: %s", err))
        return
    }

    var probeResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &probeResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse probe response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := probeResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = probeResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["key"].(string); ok && val != "" {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
    }
    if val, ok := dataMap["name"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Name = mapValue
    } else if dataMap["name"] == nil {
        data.Name = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["description"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Description = mapValue
    } else if dataMap["description"] == nil {
        data.Description = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["probeVersion"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ProbeVersion = mapValue
    } else if dataMap["probeVersion"] == nil {
        data.ProbeVersion = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastAlive"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastAlive = mapValue
    } else if dataMap["lastAlive"] == nil {
        data.LastAlive = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["iconFileId"].(string); ok && val != "" {
        data.IconFileId = types.StringValue(val)
    } else {
        data.IconFileId = types.StringNull()
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
    if val, ok := dataMap["shouldAutoEnableProbeOnNewMonitors"].(bool); ok {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolValue(val)
    } else if dataMap["shouldAutoEnableProbeOnNewMonitors"] == nil {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
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
    if val, ok := dataMap["connectionStatus"].(string); ok && val != "" {
        data.ConnectionStatus = types.StringValue(val)
    } else {
        data.ConnectionStatus = types.StringNull()
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

func (r *ProbeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ProbeResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "key": true,
        "name": true,
        "description": true,
        "probeVersion": true,
        "lastAlive": true,
        "iconFileId": true,
        "projectId": true,
        "shouldAutoEnableProbeOnNewMonitors": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "connectionStatus": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/probe/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read probe, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var probeResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &probeResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse probe response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := probeResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = probeResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["key"].(string); ok && val != "" {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
    }
    if val, ok := dataMap["name"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Name = mapValue
    } else if dataMap["name"] == nil {
        data.Name = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["description"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Description = mapValue
    } else if dataMap["description"] == nil {
        data.Description = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["probeVersion"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ProbeVersion = mapValue
    } else if dataMap["probeVersion"] == nil {
        data.ProbeVersion = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastAlive"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastAlive = mapValue
    } else if dataMap["lastAlive"] == nil {
        data.LastAlive = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["iconFileId"].(string); ok && val != "" {
        data.IconFileId = types.StringValue(val)
    } else {
        data.IconFileId = types.StringNull()
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
    if val, ok := dataMap["shouldAutoEnableProbeOnNewMonitors"].(bool); ok {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolValue(val)
    } else if dataMap["shouldAutoEnableProbeOnNewMonitors"] == nil {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
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
    if val, ok := dataMap["connectionStatus"].(string); ok && val != "" {
        data.ConnectionStatus = types.StringValue(val)
    } else {
        data.ConnectionStatus = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProbeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ProbeResourceModel
    var state ProbeResourceModel

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
    probeRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "key": data.Key.ValueString(),
        "name": r.convertTerraformMapToInterface(data.Name),
        "description": r.convertTerraformMapToInterface(data.Description),
        "probeVersion": r.convertTerraformMapToInterface(data.ProbeVersion),
        "lastAlive": r.convertTerraformMapToInterface(data.LastAlive),
        "iconFileId": data.IconFileId.ValueString(),
        "shouldAutoEnableProbeOnNewMonitors": data.ShouldAutoEnableProbeOnNewMonitors.ValueBool(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/probe/" + data.Id.ValueString() + "", probeRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update probe, got error: %s", err))
        return
    }

    // Parse the update response
    var probeResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &probeResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse probe response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "key": true,
        "name": true,
        "description": true,
        "probeVersion": true,
        "lastAlive": true,
        "iconFileId": true,
        "projectId": true,
        "shouldAutoEnableProbeOnNewMonitors": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "connectionStatus": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/probe/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read probe after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse probe read response, got error: %s", err))
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
    if val, ok := dataMap["key"].(string); ok && val != "" {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
    }
    if val, ok := dataMap["name"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Name = mapValue
    } else if dataMap["name"] == nil {
        data.Name = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["description"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.Description = mapValue
    } else if dataMap["description"] == nil {
        data.Description = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["probeVersion"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ProbeVersion = mapValue
    } else if dataMap["probeVersion"] == nil {
        data.ProbeVersion = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["lastAlive"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.LastAlive = mapValue
    } else if dataMap["lastAlive"] == nil {
        data.LastAlive = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["iconFileId"].(string); ok && val != "" {
        data.IconFileId = types.StringValue(val)
    } else {
        data.IconFileId = types.StringNull()
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
    if val, ok := dataMap["shouldAutoEnableProbeOnNewMonitors"].(bool); ok {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolValue(val)
    } else if dataMap["shouldAutoEnableProbeOnNewMonitors"] == nil {
        data.ShouldAutoEnableProbeOnNewMonitors = types.BoolNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
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
    if val, ok := dataMap["connectionStatus"].(string); ok && val != "" {
        data.ConnectionStatus = types.StringValue(val)
    } else {
        data.ConnectionStatus = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProbeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ProbeResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/probe/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete probe, got error: %s", err))
        return
    }
}


func (r *ProbeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ProbeResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ProbeResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

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
var _ resource.Resource = &SmsLogResource{}
var _ resource.ResourceWithImportState = &SmsLogResource{}

func NewSmsLogResource() resource.Resource {
    return &SmsLogResource{}
}

// SmsLogResource defines the resource implementation.
type SmsLogResource struct {
    client *Client
}

// SmsLogResourceModel describes the resource data model.
type SmsLogResourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ToNumber types.Map `tfsdk:"to_number"`
    FromNumber types.Map `tfsdk:"from_number"`
    SmsText types.String `tfsdk:"sms_text"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
    SmsCostInUsdCents types.Number `tfsdk:"sms_cost_in_usd_cents"`
}

func (r *SmsLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_sms_log"
}

func (r *SmsLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "sms_log resource",

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
            "to_number": schema.MapAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
                ElementType: types.StringType,
            },
            "from_number": schema.MapAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
                ElementType: types.StringType,
            },
            "sms_text": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "sms_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *SmsLogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *SmsLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data SmsLogResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    smsLogRequest := map[string]interface{}{
        "data": map[string]interface{}{

        },
    }

    // Make API call
    httpResp, err := r.client.Post("/sms-log/count", smsLogRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create sms_log, got error: %s", err))
        return
    }

    var smsLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &smsLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse sms_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := smsLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = smsLogResponse
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
    if val, ok := dataMap["toNumber"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ToNumber = mapValue
    } else if dataMap["toNumber"] == nil {
        data.ToNumber = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["fromNumber"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.FromNumber = mapValue
    } else if dataMap["fromNumber"] == nil {
        data.FromNumber = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["smsText"].(string); ok && val != "" {
        data.SmsText = types.StringValue(val)
    } else {
        data.SmsText = types.StringNull()
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
    if val, ok := dataMap["smsCostInUSDCents"].(float64); ok {
        data.SmsCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if dataMap["smsCostInUSDCents"] == nil {
        data.SmsCostInUsdCents = types.NumberNull()
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

func (r *SmsLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data SmsLogResourceModel

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
        "toNumber": true,
        "fromNumber": true,
        "smsText": true,
        "statusMessage": true,
        "status": true,
        "smsCostInUSDCents": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/sms-log/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sms_log, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var smsLogResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &smsLogResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse sms_log response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := smsLogResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = smsLogResponse
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
    if val, ok := dataMap["toNumber"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ToNumber = mapValue
    } else if dataMap["toNumber"] == nil {
        data.ToNumber = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["fromNumber"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.FromNumber = mapValue
    } else if dataMap["fromNumber"] == nil {
        data.FromNumber = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["smsText"].(string); ok && val != "" {
        data.SmsText = types.StringValue(val)
    } else {
        data.SmsText = types.StringNull()
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
    if val, ok := dataMap["smsCostInUSDCents"].(float64); ok {
        data.SmsCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if dataMap["smsCostInUSDCents"] == nil {
        data.SmsCostInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SmsLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *SmsLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *SmsLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *SmsLogResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *SmsLogResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

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
var _ resource.Resource = &MonitorFeedResource{}
var _ resource.ResourceWithImportState = &MonitorFeedResource{}

func NewMonitorFeedResource() resource.Resource {
    return &MonitorFeedResource{}
}

// MonitorFeedResource defines the resource implementation.
type MonitorFeedResource struct {
    client *Client
}

// MonitorFeedResourceModel describes the resource data model.
type MonitorFeedResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    FeedInfoInMarkdown types.String `tfsdk:"feed_info_in_markdown"`
    MoreInformationInMarkdown types.String `tfsdk:"more_information_in_markdown"`
    MonitorFeedEventType types.String `tfsdk:"monitor_feed_event_type"`
    DisplayColor types.Map `tfsdk:"display_color"`
    UserId types.String `tfsdk:"user_id"`
    PostedAt types.Map `tfsdk:"posted_at"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *MonitorFeedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor_feed"
}

func (r *MonitorFeedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "monitor_feed resource",

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
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "feed_info_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Log (in Markdown)",
                Required: true,
            },
            "more_information_in_markdown": schema.StringAttribute{
                MarkdownDescription: "More Information (in Markdown)",
                Optional: true,
            },
            "monitor_feed_event_type": schema.StringAttribute{
                MarkdownDescription: "Monitor Feed Event",
                Required: true,
            },
            "display_color": schema.MapAttribute{
                MarkdownDescription: "Color object",
                Required: true,
                ElementType: types.StringType,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "posted_at": schema.MapAttribute{
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
        },
    }
}

func (r *MonitorFeedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *MonitorFeedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data MonitorFeedResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    monitorFeedRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "monitorId": data.MonitorId.ValueString(),
        "feedInfoInMarkdown": data.FeedInfoInMarkdown.ValueString(),
        "moreInformationInMarkdown": data.MoreInformationInMarkdown.ValueString(),
        "monitorFeedEventType": data.MonitorFeedEventType.ValueString(),
        "displayColor": r.convertTerraformMapToInterface(data.DisplayColor),
        "userId": data.UserId.ValueString(),
        "postedAt": r.convertTerraformMapToInterface(data.PostedAt),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/monitor-feed", monitorFeedRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor_feed, got error: %s", err))
        return
    }

    var monitorFeedResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &monitorFeedResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor_feed response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := monitorFeedResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = monitorFeedResponse
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
    if val, ok := dataMap["monitorId"].(string); ok && val != "" {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if val, ok := dataMap["feedInfoInMarkdown"].(string); ok && val != "" {
        data.FeedInfoInMarkdown = types.StringValue(val)
    } else {
        data.FeedInfoInMarkdown = types.StringNull()
    }
    if val, ok := dataMap["moreInformationInMarkdown"].(string); ok && val != "" {
        data.MoreInformationInMarkdown = types.StringValue(val)
    } else {
        data.MoreInformationInMarkdown = types.StringNull()
    }
    if val, ok := dataMap["monitorFeedEventType"].(string); ok && val != "" {
        data.MonitorFeedEventType = types.StringValue(val)
    } else {
        data.MonitorFeedEventType = types.StringNull()
    }
    if val, ok := dataMap["displayColor"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DisplayColor = mapValue
    } else if dataMap["displayColor"] == nil {
        data.DisplayColor = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["postedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.PostedAt = mapValue
    } else if dataMap["postedAt"] == nil {
        data.PostedAt = types.MapNull(types.StringType)
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

func (r *MonitorFeedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data MonitorFeedResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "monitorId": true,
        "feedInfoInMarkdown": true,
        "moreInformationInMarkdown": true,
        "monitorFeedEventType": true,
        "displayColor": true,
        "userId": true,
        "postedAt": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/monitor-feed/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor_feed, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var monitorFeedResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &monitorFeedResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor_feed response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := monitorFeedResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = monitorFeedResponse
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
    if val, ok := dataMap["monitorId"].(string); ok && val != "" {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if val, ok := dataMap["feedInfoInMarkdown"].(string); ok && val != "" {
        data.FeedInfoInMarkdown = types.StringValue(val)
    } else {
        data.FeedInfoInMarkdown = types.StringNull()
    }
    if val, ok := dataMap["moreInformationInMarkdown"].(string); ok && val != "" {
        data.MoreInformationInMarkdown = types.StringValue(val)
    } else {
        data.MoreInformationInMarkdown = types.StringNull()
    }
    if val, ok := dataMap["monitorFeedEventType"].(string); ok && val != "" {
        data.MonitorFeedEventType = types.StringValue(val)
    } else {
        data.MonitorFeedEventType = types.StringNull()
    }
    if val, ok := dataMap["displayColor"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DisplayColor = mapValue
    } else if dataMap["displayColor"] == nil {
        data.DisplayColor = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["userId"].(string); ok && val != "" {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if val, ok := dataMap["postedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.PostedAt = mapValue
    } else if dataMap["postedAt"] == nil {
        data.PostedAt = types.MapNull(types.StringType)
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

func (r *MonitorFeedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    resp.Diagnostics.AddError(
        "Update Not Implemented",
        "This resource does not support update operations",
    )
}

func (r *MonitorFeedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    resp.Diagnostics.AddError(
        "Delete Not Implemented",
        "This resource does not support delete operations", 
    )
}


func (r *MonitorFeedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *MonitorFeedResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *MonitorFeedResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

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
var _ resource.Resource = &StatusPageResourceResource{}
var _ resource.ResourceWithImportState = &StatusPageResourceResource{}

func NewStatusPageResourceResource() resource.Resource {
    return &StatusPageResourceResource{}
}

// StatusPageResourceResource defines the resource implementation.
type StatusPageResourceResource struct {
    client *Client
}

// StatusPageResourceResourceModel describes the resource data model.
type StatusPageResourceResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    MonitorGroupId types.String `tfsdk:"monitor_group_id"`
    StatusPageGroupId types.String `tfsdk:"status_page_group_id"`
    DisplayName types.String `tfsdk:"display_name"`
    DisplayDescription types.String `tfsdk:"display_description"`
    DisplayTooltip types.String `tfsdk:"display_tooltip"`
    ShowCurrentStatus types.Bool `tfsdk:"show_current_status"`
    ShowUptimePercent types.Bool `tfsdk:"show_uptime_percent"`
    UptimePercentPrecision types.String `tfsdk:"uptime_percent_precision"`
    ShowStatusHistoryChart types.Bool `tfsdk:"show_status_history_chart"`
    Order types.Number `tfsdk:"order"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *StatusPageResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_resource"
}

func (r *StatusPageResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_resource resource",

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
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "monitor_group_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "status_page_group_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "display_name": schema.StringAttribute{
                MarkdownDescription: "Display Name",
                Required: true,
            },
            "display_description": schema.StringAttribute{
                MarkdownDescription: "Display Description",
                Optional: true,
            },
            "display_tooltip": schema.StringAttribute{
                MarkdownDescription: "Display Tooltip",
                Optional: true,
            },
            "show_current_status": schema.BoolAttribute{
                MarkdownDescription: "Show current status",
                Optional: true,
            },
            "show_uptime_percent": schema.BoolAttribute{
                MarkdownDescription: "Show Uptime Percent",
                Optional: true,
            },
            "uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Uptime Percent Precision",
                Optional: true,
            },
            "show_status_history_chart": schema.BoolAttribute{
                MarkdownDescription: "Show History Chart",
                Optional: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order",
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

func (r *StatusPageResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *StatusPageResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageResourceResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    statusPageResourceRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "statusPageId": data.StatusPageId.ValueString(),
        "monitorId": data.MonitorId.ValueString(),
        "monitorGroupId": data.MonitorGroupId.ValueString(),
        "statusPageGroupId": data.StatusPageGroupId.ValueString(),
        "displayName": data.DisplayName.ValueString(),
        "displayDescription": data.DisplayDescription.ValueString(),
        "displayTooltip": data.DisplayTooltip.ValueString(),
        "showCurrentStatus": data.ShowCurrentStatus.ValueBool(),
        "showUptimePercent": data.ShowUptimePercent.ValueBool(),
        "uptimePercentPrecision": data.UptimePercentPrecision.ValueString(),
        "showStatusHistoryChart": data.ShowStatusHistoryChart.ValueBool(),
        "order": data.Order.ValueBigFloat(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/status-page-resource", statusPageResourceRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page_resource, got error: %s", err))
        return
    }

    var statusPageResourceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageResourceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_resource response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageResourceResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageResourceResponse
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
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["monitorId"].(string); ok && val != "" {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if val, ok := dataMap["monitorGroupId"].(string); ok && val != "" {
        data.MonitorGroupId = types.StringValue(val)
    } else {
        data.MonitorGroupId = types.StringNull()
    }
    if val, ok := dataMap["statusPageGroupId"].(string); ok && val != "" {
        data.StatusPageGroupId = types.StringValue(val)
    } else {
        data.StatusPageGroupId = types.StringNull()
    }
    if val, ok := dataMap["displayName"].(string); ok && val != "" {
        data.DisplayName = types.StringValue(val)
    } else {
        data.DisplayName = types.StringNull()
    }
    if val, ok := dataMap["displayDescription"].(string); ok && val != "" {
        data.DisplayDescription = types.StringValue(val)
    } else {
        data.DisplayDescription = types.StringNull()
    }
    if val, ok := dataMap["displayTooltip"].(string); ok && val != "" {
        data.DisplayTooltip = types.StringValue(val)
    } else {
        data.DisplayTooltip = types.StringNull()
    }
    if val, ok := dataMap["showCurrentStatus"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    } else if dataMap["showCurrentStatus"] == nil {
        data.ShowCurrentStatus = types.BoolNull()
    }
    if val, ok := dataMap["showUptimePercent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    } else if dataMap["showUptimePercent"] == nil {
        data.ShowUptimePercent = types.BoolNull()
    }
    if val, ok := dataMap["uptimePercentPrecision"].(string); ok && val != "" {
        data.UptimePercentPrecision = types.StringValue(val)
    } else {
        data.UptimePercentPrecision = types.StringNull()
    }
    if val, ok := dataMap["showStatusHistoryChart"].(bool); ok {
        data.ShowStatusHistoryChart = types.BoolValue(val)
    } else if dataMap["showStatusHistoryChart"] == nil {
        data.ShowStatusHistoryChart = types.BoolNull()
    }
    if val, ok := dataMap["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if dataMap["order"] == nil {
        data.Order = types.NumberNull()
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

func (r *StatusPageResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageResourceResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "statusPageId": true,
        "monitorId": true,
        "monitorGroupId": true,
        "statusPageGroupId": true,
        "displayName": true,
        "displayDescription": true,
        "displayTooltip": true,
        "showCurrentStatus": true,
        "showUptimePercent": true,
        "uptimePercentPrecision": true,
        "showStatusHistoryChart": true,
        "order": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/status-page-resource/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_resource, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var statusPageResourceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageResourceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_resource response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageResourceResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageResourceResponse
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
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["monitorId"].(string); ok && val != "" {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if val, ok := dataMap["monitorGroupId"].(string); ok && val != "" {
        data.MonitorGroupId = types.StringValue(val)
    } else {
        data.MonitorGroupId = types.StringNull()
    }
    if val, ok := dataMap["statusPageGroupId"].(string); ok && val != "" {
        data.StatusPageGroupId = types.StringValue(val)
    } else {
        data.StatusPageGroupId = types.StringNull()
    }
    if val, ok := dataMap["displayName"].(string); ok && val != "" {
        data.DisplayName = types.StringValue(val)
    } else {
        data.DisplayName = types.StringNull()
    }
    if val, ok := dataMap["displayDescription"].(string); ok && val != "" {
        data.DisplayDescription = types.StringValue(val)
    } else {
        data.DisplayDescription = types.StringNull()
    }
    if val, ok := dataMap["displayTooltip"].(string); ok && val != "" {
        data.DisplayTooltip = types.StringValue(val)
    } else {
        data.DisplayTooltip = types.StringNull()
    }
    if val, ok := dataMap["showCurrentStatus"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    } else if dataMap["showCurrentStatus"] == nil {
        data.ShowCurrentStatus = types.BoolNull()
    }
    if val, ok := dataMap["showUptimePercent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    } else if dataMap["showUptimePercent"] == nil {
        data.ShowUptimePercent = types.BoolNull()
    }
    if val, ok := dataMap["uptimePercentPrecision"].(string); ok && val != "" {
        data.UptimePercentPrecision = types.StringValue(val)
    } else {
        data.UptimePercentPrecision = types.StringNull()
    }
    if val, ok := dataMap["showStatusHistoryChart"].(bool); ok {
        data.ShowStatusHistoryChart = types.BoolValue(val)
    } else if dataMap["showStatusHistoryChart"] == nil {
        data.ShowStatusHistoryChart = types.BoolNull()
    }
    if val, ok := dataMap["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if dataMap["order"] == nil {
        data.Order = types.NumberNull()
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

func (r *StatusPageResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StatusPageResourceResourceModel
    var state StatusPageResourceResourceModel

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
    statusPageResourceRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "statusPageId": data.StatusPageId.ValueString(),
        "monitorId": data.MonitorId.ValueString(),
        "monitorGroupId": data.MonitorGroupId.ValueString(),
        "statusPageGroupId": data.StatusPageGroupId.ValueString(),
        "displayName": data.DisplayName.ValueString(),
        "displayDescription": data.DisplayDescription.ValueString(),
        "displayTooltip": data.DisplayTooltip.ValueString(),
        "showCurrentStatus": data.ShowCurrentStatus.ValueBool(),
        "showUptimePercent": data.ShowUptimePercent.ValueBool(),
        "uptimePercentPrecision": data.UptimePercentPrecision.ValueString(),
        "showStatusHistoryChart": data.ShowStatusHistoryChart.ValueBool(),
        "order": data.Order.ValueBigFloat(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/status-page-resource/" + data.Id.ValueString() + "", statusPageResourceRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status_page_resource, got error: %s", err))
        return
    }

    // Parse the update response
    var statusPageResourceResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageResourceResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_resource response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "statusPageId": true,
        "monitorId": true,
        "monitorGroupId": true,
        "statusPageGroupId": true,
        "displayName": true,
        "displayDescription": true,
        "displayTooltip": true,
        "showCurrentStatus": true,
        "showUptimePercent": true,
        "uptimePercentPrecision": true,
        "showStatusHistoryChart": true,
        "order": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/status-page-resource/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_resource after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_resource read response, got error: %s", err))
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
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["monitorId"].(string); ok && val != "" {
        data.MonitorId = types.StringValue(val)
    } else {
        data.MonitorId = types.StringNull()
    }
    if val, ok := dataMap["monitorGroupId"].(string); ok && val != "" {
        data.MonitorGroupId = types.StringValue(val)
    } else {
        data.MonitorGroupId = types.StringNull()
    }
    if val, ok := dataMap["statusPageGroupId"].(string); ok && val != "" {
        data.StatusPageGroupId = types.StringValue(val)
    } else {
        data.StatusPageGroupId = types.StringNull()
    }
    if val, ok := dataMap["displayName"].(string); ok && val != "" {
        data.DisplayName = types.StringValue(val)
    } else {
        data.DisplayName = types.StringNull()
    }
    if val, ok := dataMap["displayDescription"].(string); ok && val != "" {
        data.DisplayDescription = types.StringValue(val)
    } else {
        data.DisplayDescription = types.StringNull()
    }
    if val, ok := dataMap["displayTooltip"].(string); ok && val != "" {
        data.DisplayTooltip = types.StringValue(val)
    } else {
        data.DisplayTooltip = types.StringNull()
    }
    if val, ok := dataMap["showCurrentStatus"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    } else if dataMap["showCurrentStatus"] == nil {
        data.ShowCurrentStatus = types.BoolNull()
    }
    if val, ok := dataMap["showUptimePercent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    } else if dataMap["showUptimePercent"] == nil {
        data.ShowUptimePercent = types.BoolNull()
    }
    if val, ok := dataMap["uptimePercentPrecision"].(string); ok && val != "" {
        data.UptimePercentPrecision = types.StringValue(val)
    } else {
        data.UptimePercentPrecision = types.StringNull()
    }
    if val, ok := dataMap["showStatusHistoryChart"].(bool); ok {
        data.ShowStatusHistoryChart = types.BoolValue(val)
    } else if dataMap["showStatusHistoryChart"] == nil {
        data.ShowStatusHistoryChart = types.BoolNull()
    }
    if val, ok := dataMap["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if dataMap["order"] == nil {
        data.Order = types.NumberNull()
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

func (r *StatusPageResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageResourceResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/status-page-resource/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status_page_resource, got error: %s", err))
        return
    }
}


func (r *StatusPageResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *StatusPageResourceResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *StatusPageResourceResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorTemplateDataDataSource{}

func NewMonitorTemplateDataDataSource() datasource.DataSource {
    return &MonitorTemplateDataDataSource{}
}

// MonitorTemplateDataDataSource defines the data source implementation.
type MonitorTemplateDataDataSource struct {
    client *Client
}

// MonitorTemplateDataDataSourceModel describes the data source data model.
type MonitorTemplateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    TemplateName types.String `tfsdk:"template_name"`
    TemplateDescription types.String `tfsdk:"template_description"`
    Slug types.String `tfsdk:"slug"`
    MonitorName types.String `tfsdk:"monitor_name"`
    MonitorDescription types.String `tfsdk:"monitor_description"`
    MonitorType types.String `tfsdk:"monitor_type"`
    MonitorSteps types.String `tfsdk:"monitor_steps"`
    MonitoringInterval types.String `tfsdk:"monitoring_interval"`
    Labels types.Set `tfsdk:"labels"`
    CustomFields types.String `tfsdk:"custom_fields"`
    MinimumProbeAgreement types.Number `tfsdk:"minimum_probe_agreement"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *MonitorTemplateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor_template_data"
}

func (d *MonitorTemplateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "monitor_template_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "template_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Monitor Template. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "template_description": schema.StringAttribute{
                MarkdownDescription: "Description of the Monitor Template. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "monitor_name": schema.StringAttribute{
                MarkdownDescription: "Default name applied to monitors created from this template. Users can override on creation.. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "monitor_description": schema.StringAttribute{
                MarkdownDescription: "Default description applied to monitors created from this template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "monitor_type": schema.StringAttribute{
                MarkdownDescription: "What is the type of monitor created from this template?. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "monitor_steps": schema.StringAttribute{
                MarkdownDescription: "MonitorSteps object",
                Computed: true,
            },
            "monitoring_interval": schema.StringAttribute{
                MarkdownDescription: "Default monitoring interval for monitors created from this template. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Default labels applied to monitors created from this template.. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
                ElementType: types.StringType,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource.. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "minimum_probe_agreement": schema.NumberAttribute{
                MarkdownDescription: "Default minimum number of probes that must agree on a status before the monitor status changes.. Permissions - Create: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Create Monitor Template], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Template, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Monitor Admin, Monitor Member, Edit Monitor Template]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *MonitorTemplateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorTemplateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MonitorTemplateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "monitor-template" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor_template_data, got error: %s", err))
        return
    }

    var monitorTemplateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &monitorTemplateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor_template_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := monitorTemplateDataResponse["data"].(map[string]interface{}); ok {
        monitorTemplateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := monitorTemplateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := monitorTemplateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["template_name"].(string); ok {
        data.TemplateName = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["template_description"].(string); ok {
        data.TemplateDescription = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["monitor_name"].(string); ok {
        data.MonitorName = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["monitor_description"].(string); ok {
        data.MonitorDescription = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["monitor_type"].(string); ok {
        data.MonitorType = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["monitor_steps"].(string); ok {
        data.MonitorSteps = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["monitoring_interval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["labels"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.Labels = setValue
    }
    if val, ok := monitorTemplateDataResponse["custom_fields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    }
    if val, ok := monitorTemplateDataResponse["minimum_probe_agreement"].(float64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := monitorTemplateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

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
var _ datasource.DataSource = &ProfileDataDataSource{}

func NewProfileDataDataSource() datasource.DataSource {
    return &ProfileDataDataSource{}
}

// ProfileDataDataSource defines the data source implementation.
type ProfileDataDataSource struct {
    client *Client
}

// ProfileDataDataSourceModel describes the data source data model.
type ProfileDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    ProfileId types.String `tfsdk:"profile_id"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    StartTime types.String `tfsdk:"start_time"`
    EndTime types.String `tfsdk:"end_time"`
    StartTimeUnixNano types.String `tfsdk:"start_time_unix_nano"`
    EndTimeUnixNano types.String `tfsdk:"end_time_unix_nano"`
    DurationNano types.String `tfsdk:"duration_nano"`
    ProfileType types.String `tfsdk:"profile_type"`
    Unit types.String `tfsdk:"unit"`
    PeriodType types.String `tfsdk:"period_type"`
    Period types.String `tfsdk:"period"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
    EntityKeys types.Set `tfsdk:"entity_keys"`
    ServiceEntityKey types.String `tfsdk:"service_entity_key"`
    HostEntityKey types.String `tfsdk:"host_entity_key"`
    K8sPodEntityKey types.String `tfsdk:"k8s_pod_entity_key"`
    K8sNodeEntityKey types.String `tfsdk:"k8s_node_entity_key"`
    K8sClusterEntityKey types.String `tfsdk:"k8s_cluster_entity_key"`
    ContainerEntityKey types.String `tfsdk:"container_entity_key"`
    SampleCount types.Number `tfsdk:"sample_count"`
    OriginalPayloadFormat types.String `tfsdk:"original_payload_format"`
}

func (d *ProfileDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_profile_data"
}

func (d *ProfileDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "profile_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Service ID",
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Service Type",
                Computed: true,
            },
            "profile_id": schema.StringAttribute{
                MarkdownDescription: "Profile ID",
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID",
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID",
                Computed: true,
            },
            "start_time": schema.StringAttribute{
                MarkdownDescription: "Start Time",
                Computed: true,
            },
            "end_time": schema.StringAttribute{
                MarkdownDescription: "End Time",
                Computed: true,
            },
            "start_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "Start Time in Unix Nano",
                Computed: true,
            },
            "end_time_unix_nano": schema.StringAttribute{
                MarkdownDescription: "End Time in Unix Nano",
                Computed: true,
            },
            "duration_nano": schema.StringAttribute{
                MarkdownDescription: "Duration in Nanoseconds",
                Computed: true,
            },
            "profile_type": schema.StringAttribute{
                MarkdownDescription: "Profile Type",
                Computed: true,
            },
            "unit": schema.StringAttribute{
                MarkdownDescription: "Unit",
                Computed: true,
            },
            "period_type": schema.StringAttribute{
                MarkdownDescription: "Period Type",
                Computed: true,
            },
            "period": schema.StringAttribute{
                MarkdownDescription: "Period",
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "entity_keys": schema.SetAttribute{
                MarkdownDescription: "Entity Keys",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_entity_key": schema.StringAttribute{
                MarkdownDescription: "Service Entity Key",
                Computed: true,
            },
            "host_entity_key": schema.StringAttribute{
                MarkdownDescription: "Host Entity Key",
                Computed: true,
            },
            "k8s_pod_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Pod Entity Key",
                Computed: true,
            },
            "k8s_node_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Node Entity Key",
                Computed: true,
            },
            "k8s_cluster_entity_key": schema.StringAttribute{
                MarkdownDescription: "Kubernetes Cluster Entity Key",
                Computed: true,
            },
            "container_entity_key": schema.StringAttribute{
                MarkdownDescription: "Container Entity Key",
                Computed: true,
            },
            "sample_count": schema.NumberAttribute{
                MarkdownDescription: "Sample Count",
                Computed: true,
            },
            "original_payload_format": schema.StringAttribute{
                MarkdownDescription: "Original Payload Format",
                Computed: true,
            },
        },
    }
}

func (d *ProfileDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProfileDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProfileDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "profile" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profile_data, got error: %s", err))
        return
    }

    var profileDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &profileDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse profile_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := profileDataResponse["data"].(map[string]interface{}); ok {
        profileDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := profileDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := profileDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := profileDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := profileDataResponse["primary_entity_id"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    }
    if val, ok := profileDataResponse["primary_entity_type"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    }
    if val, ok := profileDataResponse["profile_id"].(string); ok {
        data.ProfileId = types.StringValue(val)
    }
    if val, ok := profileDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := profileDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }
    if val, ok := profileDataResponse["start_time"].(string); ok {
        data.StartTime = types.StringValue(val)
    }
    if val, ok := profileDataResponse["end_time"].(string); ok {
        data.EndTime = types.StringValue(val)
    }
    if val, ok := profileDataResponse["start_time_unix_nano"].(string); ok {
        data.StartTimeUnixNano = types.StringValue(val)
    }
    if val, ok := profileDataResponse["end_time_unix_nano"].(string); ok {
        data.EndTimeUnixNano = types.StringValue(val)
    }
    if val, ok := profileDataResponse["duration_nano"].(string); ok {
        data.DurationNano = types.StringValue(val)
    }
    if val, ok := profileDataResponse["profile_type"].(string); ok {
        data.ProfileType = types.StringValue(val)
    }
    if val, ok := profileDataResponse["unit"].(string); ok {
        data.Unit = types.StringValue(val)
    }
    if val, ok := profileDataResponse["period_type"].(string); ok {
        data.PeriodType = types.StringValue(val)
    }
    if val, ok := profileDataResponse["period"].(string); ok {
        data.Period = types.StringValue(val)
    }
    if val, ok := profileDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }
    if val, ok := profileDataResponse["attribute_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.AttributeKeys = setValue
    }
    if val, ok := profileDataResponse["entity_keys"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        setValue, _ := types.SetValue(types.StringType, elements)
        data.EntityKeys = setValue
    }
    if val, ok := profileDataResponse["service_entity_key"].(string); ok {
        data.ServiceEntityKey = types.StringValue(val)
    }
    if val, ok := profileDataResponse["host_entity_key"].(string); ok {
        data.HostEntityKey = types.StringValue(val)
    }
    if val, ok := profileDataResponse["k8s_pod_entity_key"].(string); ok {
        data.K8sPodEntityKey = types.StringValue(val)
    }
    if val, ok := profileDataResponse["k8s_node_entity_key"].(string); ok {
        data.K8sNodeEntityKey = types.StringValue(val)
    }
    if val, ok := profileDataResponse["k8s_cluster_entity_key"].(string); ok {
        data.K8sClusterEntityKey = types.StringValue(val)
    }
    if val, ok := profileDataResponse["container_entity_key"].(string); ok {
        data.ContainerEntityKey = types.StringValue(val)
    }
    if val, ok := profileDataResponse["sample_count"].(float64); ok {
        data.SampleCount = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := profileDataResponse["original_payload_format"].(string); ok {
        data.OriginalPayloadFormat = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

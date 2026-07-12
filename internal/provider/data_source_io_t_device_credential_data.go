package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IoTDeviceCredentialDataDataSource{}

func NewIoTDeviceCredentialDataDataSource() datasource.DataSource {
    return &IoTDeviceCredentialDataDataSource{}
}

// IoTDeviceCredentialDataDataSource defines the data source implementation.
type IoTDeviceCredentialDataDataSource struct {
    client *Client
}

// IoTDeviceCredentialDataDataSourceModel describes the data source data model.
type IoTDeviceCredentialDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IotFleetId types.String `tfsdk:"iot_fleet_id"`
    ExternalId types.String `tfsdk:"external_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    LastConnectedAt types.String `tfsdk:"last_connected_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SecretKey types.String `tfsdk:"secret_key"`
}

func (d *IoTDeviceCredentialDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_io_t_device_credential_data"
}

func (d *IoTDeviceCredentialDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "io_t_device_credential_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
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
            "iot_fleet_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "external_id": schema.StringAttribute{
                MarkdownDescription: "The device id — must match the device.id label the device stamps on its datapoints. It is also the <device> segment of the device's MQTT topics, so a device that reports directly over MQTT cannot use an id containing '/', '+', or '#' (such devices can still report through a gateway).. Permissions - Create: [Project Owner, Project Admin, Create IoT Device Credential], Read: [Project Owner, Project Admin, Read IoT Device Credential, Read IoT Fleet], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Disabled credentials are rejected at MQTT CONNECT and stop the device's silent-death offline detection.. Permissions - Create: [Project Owner, Project Admin, Create IoT Device Credential], Read: [Project Owner, Project Admin, Read IoT Device Credential, Read IoT Fleet], Update: [Project Owner, Project Admin, Edit IoT Device Credential]",
                Computed: true,
            },
            "last_connected_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IoTDeviceCredentialDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IoTDeviceCredentialDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IoTDeviceCredentialDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "iot-device-credential" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read io_t_device_credential_data, got error: %s", err))
        return
    }

    var ioTDeviceCredentialDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &ioTDeviceCredentialDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse io_t_device_credential_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := ioTDeviceCredentialDataResponse["data"].(map[string]interface{}); ok {
        ioTDeviceCredentialDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := ioTDeviceCredentialDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := ioTDeviceCredentialDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["iot_fleet_id"].(string); ok {
        data.IotFleetId = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["external_id"].(string); ok {
        data.ExternalId = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["last_connected_at"].(string); ok {
        data.LastConnectedAt = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := ioTDeviceCredentialDataResponse["secret_key"].(string); ok {
        data.SecretKey = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

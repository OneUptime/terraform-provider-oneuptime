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
var _ datasource.DataSource = &IncomingCallPolicyDataDataSource{}

func NewIncomingCallPolicyDataDataSource() datasource.DataSource {
    return &IncomingCallPolicyDataDataSource{}
}

// IncomingCallPolicyDataDataSource defines the data source implementation.
type IncomingCallPolicyDataDataSource struct {
    client *Client
}

// IncomingCallPolicyDataDataSourceModel describes the data source data model.
type IncomingCallPolicyDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    RoutingPhoneNumber types.String `tfsdk:"routing_phone_number"`
    CallProviderPhoneNumberId types.String `tfsdk:"call_provider_phone_number_id"`
    PhoneNumberCountryCode types.String `tfsdk:"phone_number_country_code"`
    PhoneNumberAreaCode types.String `tfsdk:"phone_number_area_code"`
    PhoneNumberPurchasedAt types.String `tfsdk:"phone_number_purchased_at"`
    GreetingMessage types.String `tfsdk:"greeting_message"`
    NoAnswerMessage types.String `tfsdk:"no_answer_message"`
    NoOneAvailableMessage types.String `tfsdk:"no_one_available_message"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    RepeatPolicyIfNoOneAnswers types.Bool `tfsdk:"repeat_policy_if_no_one_answers"`
    RepeatPolicyIfNoOneAnswersTimes types.Number `tfsdk:"repeat_policy_if_no_one_answers_times"`
    Labels types.Set `tfsdk:"labels"`
    ProjectCallSmsConfigId types.String `tfsdk:"project_call_sms_config_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncomingCallPolicyDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_policy_data"
}

func (d *IncomingCallPolicyDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_policy_data data source",

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
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "routing_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "call_provider_phone_number_id": schema.StringAttribute{
                MarkdownDescription: "The call provider's ID for the phone number (e.g., Twilio SID). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "phone_number_country_code": schema.StringAttribute{
                MarkdownDescription: "Country code of the phone number (US, GB, etc.). Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "phone_number_area_code": schema.StringAttribute{
                MarkdownDescription: "Area code of the phone number. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "phone_number_purchased_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "greeting_message": schema.StringAttribute{
                MarkdownDescription: "Custom TTS greeting message for incoming calls. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "no_answer_message": schema.StringAttribute{
                MarkdownDescription: "Message when escalation is exhausted and no one answers. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "no_one_available_message": schema.StringAttribute{
                MarkdownDescription: "Message when no one is on-call or reachable. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Enable or disable this incoming call policy. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "repeat_policy_if_no_one_answers": schema.BoolAttribute{
                MarkdownDescription: "Restart from first rule if all fail. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "repeat_policy_if_no_one_answers_times": schema.NumberAttribute{
                MarkdownDescription: "Maximum repeat attempts if no one answers. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in.. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Create Incoming Call Policy], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Incoming Call Policy], Update: [Project Owner, Project Admin, Project Member, Settings Admin, Settings Member, Edit Incoming Call Policy]",
                Computed: true,
                ElementType: types.StringType,
            },
            "project_call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallPolicyDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallPolicyDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallPolicyDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incoming-call-policy" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_policy_data, got error: %s", err))
        return
    }

    var incomingCallPolicyDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incomingCallPolicyDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_policy_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incomingCallPolicyDataResponse["data"].(map[string]interface{}); ok {
        incomingCallPolicyDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incomingCallPolicyDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallPolicyDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["routing_phone_number"].(string); ok {
        data.RoutingPhoneNumber = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["call_provider_phone_number_id"].(string); ok {
        data.CallProviderPhoneNumberId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["phone_number_country_code"].(string); ok {
        data.PhoneNumberCountryCode = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["phone_number_area_code"].(string); ok {
        data.PhoneNumberAreaCode = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["phone_number_purchased_at"].(string); ok {
        data.PhoneNumberPurchasedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["greeting_message"].(string); ok {
        data.GreetingMessage = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["no_answer_message"].(string); ok {
        data.NoAnswerMessage = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["no_one_available_message"].(string); ok {
        data.NoOneAvailableMessage = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["repeat_policy_if_no_one_answers"].(bool); ok {
        data.RepeatPolicyIfNoOneAnswers = types.BoolValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["repeat_policy_if_no_one_answers_times"].(float64); ok {
        data.RepeatPolicyIfNoOneAnswersTimes = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallPolicyDataResponse["labels"].([]interface{}); ok {
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
    if val, ok := incomingCallPolicyDataResponse["project_call_sms_config_id"].(string); ok {
        data.ProjectCallSmsConfigId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

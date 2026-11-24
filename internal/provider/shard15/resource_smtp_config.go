package shard15

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SmtpConfigResource{}
var _ resource.ResourceWithImportState = &SmtpConfigResource{}
var _ resource.ResourceWithModifyPlan = &SmtpConfigResource{}

func NewSmtpConfigResource() resource.Resource {
    return &SmtpConfigResource{}
}

// SmtpConfigResource defines the resource implementation.
type SmtpConfigResource struct {
    client *Client
}

// SmtpConfigResourceModel describes the resource data model.
type SmtpConfigResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Username types.String `tfsdk:"username"`
    Password types.String `tfsdk:"password"`
    Hostname types.String `tfsdk:"hostname"`
    Port types.Number `tfsdk:"port"`
    FromEmail JSONSubsetValue `tfsdk:"from_email"`
    FromName types.String `tfsdk:"from_name"`
    Secure types.Bool `tfsdk:"secure"`
    TransportType types.String `tfsdk:"transport_type"`
    AuthType types.String `tfsdk:"auth_type"`
    ClientId types.String `tfsdk:"client_id"`
    ClientSecret types.String `tfsdk:"client_secret"`
    TokenUrl types.String `tfsdk:"token_url"`
    Scope types.String `tfsdk:"scope"`
    OauthProviderType types.String `tfsdk:"oauth_provider_type"`
    AmazonSesAccessKeyId types.String `tfsdk:"amazon_ses_access_key_id"`
    AmazonSesSecretAccessKey types.String `tfsdk:"amazon_ses_secret_access_key"`
    AmazonSesRegion types.String `tfsdk:"amazon_ses_region"`
    AmazonSesAuthenticationType types.String `tfsdk:"amazon_ses_authentication_type"`
    AmazonSesRoleArn types.String `tfsdk:"amazon_ses_role_arn"`
    SendgridApiKey types.String `tfsdk:"sendgrid_api_key"`
    MailgunApiKey types.String `tfsdk:"mailgun_api_key"`
    MailgunDomain types.String `tfsdk:"mailgun_domain"`
    MailgunRegion types.String `tfsdk:"mailgun_region"`
    PostmarkServerToken types.String `tfsdk:"postmark_server_token"`
    ResendApiKey types.String `tfsdk:"resend_api_key"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *SmtpConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_smtp_config"
}

func (r *SmtpConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *SmtpConfigResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "Manage Custom SMTP Servers for your project",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "username": schema.StringAttribute{
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "password": schema.StringAttribute{
                MarkdownDescription: "The password the SMTP server is signed in to with, or the API key your provider gives you for SMTP. Used when Authentication Type is Username and Password. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "hostname": schema.StringAttribute{
                MarkdownDescription: "SMTP server hostname. Required when Transport is SMTP. Not used by HTTP-API transports like Microsoft Graph.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "port": schema.NumberAttribute{
                MarkdownDescription: "SMTP server port. Required when Transport is SMTP. Not used by HTTP-API transports.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "from_email": schema.StringAttribute{
                CustomType: JSONSubsetType{},
                Required: true,
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "from_name": schema.StringAttribute{
                Required: true,
            },
            "secure": schema.BoolAttribute{
                MarkdownDescription: "Enable secure SMTP connection. Used only for SMTP transport — HTTP-API transports always use HTTPS.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "transport_type": schema.StringAttribute{
                MarkdownDescription: "How OneUptime delivers mail for this config. 'SMTP' uses the hostname and port. 'Microsoft Graph' sends through the Microsoft Graph REST API, for a Microsoft 365 tenant with SMTP AUTH turned off. 'Amazon SES', 'SendGrid', 'Mailgun', 'Postmark' and 'Resend' send through that provider's HTTPS API with its own key, set in the columns named for it.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("SMTP"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("SMTP", "Microsoft Graph", "Amazon SES", "SendGrid", "Mailgun", "Postmark", "Resend"),
                },
            },
            "auth_type": schema.StringAttribute{
                MarkdownDescription: "The type of authentication to use for this SMTP server. Options: Username and Password, OAuth, or None.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Username and Password"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Username and Password", "OAuth", "None"),
                },
            },
            "client_id": schema.StringAttribute{
                MarkdownDescription: "The Client ID from your OAuth application registration. Required for OAuth authentication.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "client_secret": schema.StringAttribute{
                MarkdownDescription: "The Client Secret from your OAuth application registration. Required for OAuth authentication. For Google service accounts, this is the private key from the JSON key file. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "token_url": schema.StringAttribute{
                MarkdownDescription: "The OAuth token endpoint URL. For Microsoft 365: https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token. For Google: https://oauth2.googleapis.com/token.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "scope": schema.StringAttribute{
                MarkdownDescription: "The OAuth scope(s) required for SMTP access. For Microsoft 365: https://outlook.office365.com/.default. For Google: https://mail.google.com/.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "oauth_provider_type": schema.StringAttribute{
                MarkdownDescription: "The OAuth grant type to use. 'Client Credentials' for Microsoft 365 and most providers. 'JWT Bearer' for Google Workspace service accounts.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Client Credentials", "JWT Bearer"),
                },
            },
            "amazon_ses_access_key_id": schema.StringAttribute{
                MarkdownDescription: "The access key ID of an IAM user or role allowed to send with Amazon SES (ses:SendEmail). Required when Transport is Amazon SES.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "amazon_ses_secret_access_key": schema.StringAttribute{
                MarkdownDescription: "The secret access key that goes with the access key ID. Required when Transport is Amazon SES. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "amazon_ses_region": schema.StringAttribute{
                MarkdownDescription: "The AWS Region your sender identity is verified in, such as us-east-1. Required when Transport is Amazon SES.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "amazon_ses_authentication_type": schema.StringAttribute{
                MarkdownDescription: "How OneUptime signs in to Amazon SES: 'Access Key' with the access key ID and secret access key, or 'IAM Role' with an IAM role in your own AWS account (Amazon SES Role ARN), which OneUptime signs in to with this project's external ID. Read only when Transport is Amazon SES.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Access Key"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Access Key", "IAM Role"),
                },
            },
            "amazon_ses_role_arn": schema.StringAttribute{
                MarkdownDescription: "The ARN of the IAM role in your AWS account that OneUptime signs in to, when Amazon SES Authentication Type is 'IAM Role'. Its trust policy must let OneUptime's AWS account in with this project's external ID, and the role must be allowed ses:SendEmail.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "sendgrid_api_key": schema.StringAttribute{
                MarkdownDescription: "A SendGrid API key with the Mail Send permission. Required when Transport is SendGrid. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "mailgun_api_key": schema.StringAttribute{
                MarkdownDescription: "A Mailgun API key, or a sending API key for the domain. Required when Transport is Mailgun. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "mailgun_domain": schema.StringAttribute{
                MarkdownDescription: "The Mailgun sending domain mail goes out from, such as mg.example.com. Required when Transport is Mailgun.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "mailgun_region": schema.StringAttribute{
                MarkdownDescription: "The region the Mailgun domain was created in: US or EU. US when not set.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("US", "EU"),
                },
            },
            "postmark_server_token": schema.StringAttribute{
                MarkdownDescription: "The API token of the Postmark server that sends, from its API Tokens tab. Required when Transport is Postmark. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "resend_api_key": schema.StringAttribute{
                MarkdownDescription: "A Resend API key with sending access. Required when Transport is Resend. Write-only: kept encrypted and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *SmtpConfigResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    // Under Crossplane, what a manifest leaves out is the record's (see crossplane.go).
    KeepLeftOutUnderCrossplane(ctx, req, resp)

    KeepStateWhenNothingChanges(ctx, req, resp)
}

func (r *SmtpConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *SmtpConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data SmtpConfigResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config SmtpConfigResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    smtpConfigRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := smtpConfigRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Username.IsNull() && !data.Username.IsUnknown() {
        requestDataMap["username"] = data.Username.ValueString()
    }
    if !data.Password.IsNull() && !data.Password.IsUnknown() {
        requestDataMap["password"] = data.Password.ValueString()
    }
    if !data.Hostname.IsNull() && !data.Hostname.IsUnknown() {
        requestDataMap["hostname"] = data.Hostname.ValueString()
    }
    if !data.Port.IsNull() && !data.Port.IsUnknown() {
        requestDataMap["port"] = r.bigFloatToFloat64(data.Port.ValueBigFloat())
    }
    if parsedFromEmail := r.parseJSONField(data.FromEmail); parsedFromEmail != nil {
        requestDataMap["fromEmail"] = parsedFromEmail
    }
    if !data.FromName.IsNull() && !data.FromName.IsUnknown() {
        requestDataMap["fromName"] = data.FromName.ValueString()
    }
    if !data.Secure.IsNull() && !data.Secure.IsUnknown() {
        requestDataMap["secure"] = data.Secure.ValueBool()
    }
    if !data.TransportType.IsNull() && !data.TransportType.IsUnknown() {
        requestDataMap["transportType"] = data.TransportType.ValueString()
    }
    if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
        requestDataMap["authType"] = data.AuthType.ValueString()
    }
    if !data.ClientId.IsNull() && !data.ClientId.IsUnknown() {
        requestDataMap["clientId"] = data.ClientId.ValueString()
    }
    if !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown() {
        requestDataMap["clientSecret"] = data.ClientSecret.ValueString()
    }
    if !data.TokenUrl.IsNull() && !data.TokenUrl.IsUnknown() {
        requestDataMap["tokenUrl"] = data.TokenUrl.ValueString()
    }
    if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
        requestDataMap["scope"] = data.Scope.ValueString()
    }
    if !data.OauthProviderType.IsNull() && !data.OauthProviderType.IsUnknown() {
        requestDataMap["oauthProviderType"] = data.OauthProviderType.ValueString()
    }
    if !data.AmazonSesAccessKeyId.IsNull() && !data.AmazonSesAccessKeyId.IsUnknown() {
        requestDataMap["amazonSesAccessKeyId"] = data.AmazonSesAccessKeyId.ValueString()
    }
    if !data.AmazonSesSecretAccessKey.IsNull() && !data.AmazonSesSecretAccessKey.IsUnknown() {
        requestDataMap["amazonSesSecretAccessKey"] = data.AmazonSesSecretAccessKey.ValueString()
    }
    if !data.AmazonSesRegion.IsNull() && !data.AmazonSesRegion.IsUnknown() {
        requestDataMap["amazonSesRegion"] = data.AmazonSesRegion.ValueString()
    }
    if !data.AmazonSesAuthenticationType.IsNull() && !data.AmazonSesAuthenticationType.IsUnknown() {
        requestDataMap["amazonSesAuthenticationType"] = data.AmazonSesAuthenticationType.ValueString()
    }
    if !data.AmazonSesRoleArn.IsNull() && !data.AmazonSesRoleArn.IsUnknown() {
        requestDataMap["amazonSesRoleArn"] = data.AmazonSesRoleArn.ValueString()
    }
    if !data.SendgridApiKey.IsNull() && !data.SendgridApiKey.IsUnknown() {
        requestDataMap["sendgridApiKey"] = data.SendgridApiKey.ValueString()
    }
    if !data.MailgunApiKey.IsNull() && !data.MailgunApiKey.IsUnknown() {
        requestDataMap["mailgunApiKey"] = data.MailgunApiKey.ValueString()
    }
    if !data.MailgunDomain.IsNull() && !data.MailgunDomain.IsUnknown() {
        requestDataMap["mailgunDomain"] = data.MailgunDomain.ValueString()
    }
    if !data.MailgunRegion.IsNull() && !data.MailgunRegion.IsUnknown() {
        requestDataMap["mailgunRegion"] = data.MailgunRegion.ValueString()
    }
    if !data.PostmarkServerToken.IsNull() && !data.PostmarkServerToken.IsUnknown() {
        requestDataMap["postmarkServerToken"] = data.PostmarkServerToken.ValueString()
    }
    if !data.ResendApiKey.IsNull() && !data.ResendApiKey.IsUnknown() {
        requestDataMap["resendApiKey"] = data.ResendApiKey.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/smtp-config", smtpConfigRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create smtp_config, got error: %s", err))
        return
    }

    var smtpConfigResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &smtpConfigResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create smtp_config: %s", err))
        return
    }

    // What OneUptime left out because only it writes it, if anything (see entries.go).
    resp.Diagnostics.Append(FieldsLeftOutWarning(httpResp, "smtp config")...)

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := smtpConfigResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := smtpConfigResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for smtp_config did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)


    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * smtp_config is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "username": true,
        "hostname": true,
        "port": true,
        "fromEmail": true,
        "fromName": true,
        "secure": true,
        "transportType": true,
        "authType": true,
        "clientId": true,
        "tokenUrl": true,
        "scope": true,
        "oauthProviderType": true,
        "amazonSesAccessKeyId": true,
        "amazonSesRegion": true,
        "amazonSesAuthenticationType": true,
        "amazonSesRoleArn": true,
        "mailgunDomain": true,
        "mailgunRegion": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/smtp-config/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created smtp_config but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created smtp_config but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["username"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Username = types.StringValue(string(jsonBytes))
            } else {
                data.Username = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Username = types.StringValue(string(jsonBytes))
            } else {
                data.Username = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if arr, ok := dataMap["username"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if val, ok := dataMap["username"].(string); ok {
        data.Username = types.StringValue(val)
    } else {
        data.Username = types.StringNull()
    }
    if obj, ok := dataMap["hostname"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if arr, ok := dataMap["hostname"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := dataMap["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
    }
    data.Port = NumberFromAPI(data.Port, dataMap["port"])
    if obj, ok := dataMap["fromEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEmail = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromEmail = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.FromEmail = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["fromEmail"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.FromEmail = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["fromEmail"].(string); ok {
        data.FromEmail = NewJSONSubsetValue(val)
    } else {
        data.FromEmail = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["fromName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromName = types.StringValue(string(jsonBytes))
            } else {
                data.FromName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromName = types.StringValue(string(jsonBytes))
            } else {
                data.FromName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if arr, ok := dataMap["fromName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if val, ok := dataMap["fromName"].(string); ok {
        data.FromName = types.StringValue(val)
    } else {
        data.FromName = types.StringNull()
    }
    if val, ok := dataMap["secure"].(bool); ok {
        data.Secure = types.BoolValue(val)
    }
    if obj, ok := dataMap["transportType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TransportType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TransportType = types.StringValue(string(jsonBytes))
            } else {
                data.TransportType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TransportType = types.StringValue(string(jsonBytes))
            } else {
                data.TransportType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if arr, ok := dataMap["transportType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if val, ok := dataMap["transportType"].(string); ok {
        data.TransportType = types.StringValue(val)
    } else {
        data.TransportType = types.StringNull()
    }
    if obj, ok := dataMap["authType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AuthType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AuthType = types.StringValue(string(jsonBytes))
            } else {
                data.AuthType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AuthType = types.StringValue(string(jsonBytes))
            } else {
                data.AuthType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if arr, ok := dataMap["authType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if val, ok := dataMap["authType"].(string); ok {
        data.AuthType = types.StringValue(val)
    } else {
        data.AuthType = types.StringNull()
    }
    if obj, ok := dataMap["clientId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClientId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClientId = types.StringValue(string(jsonBytes))
            } else {
                data.ClientId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClientId = types.StringValue(string(jsonBytes))
            } else {
                data.ClientId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if arr, ok := dataMap["clientId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if val, ok := dataMap["clientId"].(string); ok {
        data.ClientId = types.StringValue(val)
    } else {
        data.ClientId = types.StringNull()
    }
    if obj, ok := dataMap["tokenUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TokenUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TokenUrl = types.StringValue(string(jsonBytes))
            } else {
                data.TokenUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TokenUrl = types.StringValue(string(jsonBytes))
            } else {
                data.TokenUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if arr, ok := dataMap["tokenUrl"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if val, ok := dataMap["tokenUrl"].(string); ok {
        data.TokenUrl = types.StringValue(val)
    } else {
        data.TokenUrl = types.StringNull()
    }
    if obj, ok := dataMap["scope"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Scope = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Scope = types.StringValue(string(jsonBytes))
            } else {
                data.Scope = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Scope = types.StringValue(string(jsonBytes))
            } else {
                data.Scope = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if arr, ok := dataMap["scope"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if val, ok := dataMap["scope"].(string); ok {
        data.Scope = types.StringValue(val)
    } else {
        data.Scope = types.StringNull()
    }
    if obj, ok := dataMap["oauthProviderType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OauthProviderType = types.StringValue(string(jsonBytes))
            } else {
                data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OauthProviderType = types.StringValue(string(jsonBytes))
            } else {
                data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if arr, ok := dataMap["oauthProviderType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if val, ok := dataMap["oauthProviderType"].(string); ok {
        data.OauthProviderType = types.StringValue(val)
    } else {
        data.OauthProviderType = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesAccessKeyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesAccessKeyId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesAccessKeyId"].(string); ok {
        data.AmazonSesAccessKeyId = types.StringValue(val)
    } else {
        data.AmazonSesAccessKeyId = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesRegion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesRegion = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesRegion = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesRegion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesRegion"].(string); ok {
        data.AmazonSesRegion = types.StringValue(val)
    } else {
        data.AmazonSesRegion = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesAuthenticationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAuthenticationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesAuthenticationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAuthenticationType = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesAuthenticationType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAuthenticationType = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesAuthenticationType"].(string); ok {
        data.AmazonSesAuthenticationType = types.StringValue(val)
    } else {
        data.AmazonSesAuthenticationType = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesRoleArn"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRoleArn = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesRoleArn = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRoleArn = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesRoleArn"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRoleArn = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesRoleArn"].(string); ok {
        data.AmazonSesRoleArn = types.StringValue(val)
    } else {
        data.AmazonSesRoleArn = types.StringNull()
    }
    if obj, ok := dataMap["mailgunDomain"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MailgunDomain = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MailgunDomain = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if arr, ok := dataMap["mailgunDomain"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if val, ok := dataMap["mailgunDomain"].(string); ok {
        data.MailgunDomain = types.StringValue(val)
    } else {
        data.MailgunDomain = types.StringNull()
    }
    if obj, ok := dataMap["mailgunRegion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MailgunRegion = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MailgunRegion = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if arr, ok := dataMap["mailgunRegion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if val, ok := dataMap["mailgunRegion"].(string); ok {
        data.MailgunRegion = types.StringValue(val)
    } else {
        data.MailgunRegion = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")


    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SmtpConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data SmtpConfigResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "username": true,
        "hostname": true,
        "port": true,
        "fromEmail": true,
        "fromName": true,
        "secure": true,
        "transportType": true,
        "authType": true,
        "clientId": true,
        "tokenUrl": true,
        "scope": true,
        "oauthProviderType": true,
        "amazonSesAccessKeyId": true,
        "amazonSesRegion": true,
        "amazonSesAuthenticationType": true,
        "amazonSesRoleArn": true,
        "mailgunDomain": true,
        "mailgunRegion": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/smtp-config/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read smtp_config, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var smtpConfigResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &smtpConfigResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse smtp_config response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := smtpConfigResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = smtpConfigResponse
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["username"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Username = types.StringValue(string(jsonBytes))
            } else {
                data.Username = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Username = types.StringValue(string(jsonBytes))
            } else {
                data.Username = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if arr, ok := dataMap["username"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if val, ok := dataMap["username"].(string); ok {
        data.Username = types.StringValue(val)
    } else {
        data.Username = types.StringNull()
    }
    if obj, ok := dataMap["hostname"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if arr, ok := dataMap["hostname"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := dataMap["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
    }
    data.Port = NumberFromAPI(data.Port, dataMap["port"])
    if obj, ok := dataMap["fromEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEmail = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromEmail = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.FromEmail = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["fromEmail"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.FromEmail = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["fromEmail"].(string); ok {
        data.FromEmail = NewJSONSubsetValue(val)
    } else {
        data.FromEmail = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["fromName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromName = types.StringValue(string(jsonBytes))
            } else {
                data.FromName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromName = types.StringValue(string(jsonBytes))
            } else {
                data.FromName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if arr, ok := dataMap["fromName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if val, ok := dataMap["fromName"].(string); ok {
        data.FromName = types.StringValue(val)
    } else {
        data.FromName = types.StringNull()
    }
    if val, ok := dataMap["secure"].(bool); ok {
        data.Secure = types.BoolValue(val)
    }
    if obj, ok := dataMap["transportType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TransportType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TransportType = types.StringValue(string(jsonBytes))
            } else {
                data.TransportType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TransportType = types.StringValue(string(jsonBytes))
            } else {
                data.TransportType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if arr, ok := dataMap["transportType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if val, ok := dataMap["transportType"].(string); ok {
        data.TransportType = types.StringValue(val)
    } else {
        data.TransportType = types.StringNull()
    }
    if obj, ok := dataMap["authType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AuthType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AuthType = types.StringValue(string(jsonBytes))
            } else {
                data.AuthType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AuthType = types.StringValue(string(jsonBytes))
            } else {
                data.AuthType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if arr, ok := dataMap["authType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if val, ok := dataMap["authType"].(string); ok {
        data.AuthType = types.StringValue(val)
    } else {
        data.AuthType = types.StringNull()
    }
    if obj, ok := dataMap["clientId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClientId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClientId = types.StringValue(string(jsonBytes))
            } else {
                data.ClientId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClientId = types.StringValue(string(jsonBytes))
            } else {
                data.ClientId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if arr, ok := dataMap["clientId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if val, ok := dataMap["clientId"].(string); ok {
        data.ClientId = types.StringValue(val)
    } else {
        data.ClientId = types.StringNull()
    }
    if obj, ok := dataMap["tokenUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TokenUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TokenUrl = types.StringValue(string(jsonBytes))
            } else {
                data.TokenUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TokenUrl = types.StringValue(string(jsonBytes))
            } else {
                data.TokenUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if arr, ok := dataMap["tokenUrl"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if val, ok := dataMap["tokenUrl"].(string); ok {
        data.TokenUrl = types.StringValue(val)
    } else {
        data.TokenUrl = types.StringNull()
    }
    if obj, ok := dataMap["scope"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Scope = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Scope = types.StringValue(string(jsonBytes))
            } else {
                data.Scope = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Scope = types.StringValue(string(jsonBytes))
            } else {
                data.Scope = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if arr, ok := dataMap["scope"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if val, ok := dataMap["scope"].(string); ok {
        data.Scope = types.StringValue(val)
    } else {
        data.Scope = types.StringNull()
    }
    if obj, ok := dataMap["oauthProviderType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OauthProviderType = types.StringValue(string(jsonBytes))
            } else {
                data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OauthProviderType = types.StringValue(string(jsonBytes))
            } else {
                data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if arr, ok := dataMap["oauthProviderType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if val, ok := dataMap["oauthProviderType"].(string); ok {
        data.OauthProviderType = types.StringValue(val)
    } else {
        data.OauthProviderType = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesAccessKeyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesAccessKeyId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesAccessKeyId"].(string); ok {
        data.AmazonSesAccessKeyId = types.StringValue(val)
    } else {
        data.AmazonSesAccessKeyId = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesRegion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesRegion = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesRegion = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesRegion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesRegion"].(string); ok {
        data.AmazonSesRegion = types.StringValue(val)
    } else {
        data.AmazonSesRegion = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesAuthenticationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAuthenticationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesAuthenticationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAuthenticationType = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesAuthenticationType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAuthenticationType = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesAuthenticationType"].(string); ok {
        data.AmazonSesAuthenticationType = types.StringValue(val)
    } else {
        data.AmazonSesAuthenticationType = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesRoleArn"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRoleArn = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesRoleArn = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRoleArn = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesRoleArn"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRoleArn = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesRoleArn"].(string); ok {
        data.AmazonSesRoleArn = types.StringValue(val)
    } else {
        data.AmazonSesRoleArn = types.StringNull()
    }
    if obj, ok := dataMap["mailgunDomain"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MailgunDomain = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MailgunDomain = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if arr, ok := dataMap["mailgunDomain"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if val, ok := dataMap["mailgunDomain"].(string); ok {
        data.MailgunDomain = types.StringValue(val)
    } else {
        data.MailgunDomain = types.StringNull()
    }
    if obj, ok := dataMap["mailgunRegion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MailgunRegion = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MailgunRegion = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if arr, ok := dataMap["mailgunRegion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if val, ok := dataMap["mailgunRegion"].(string); ok {
        data.MailgunRegion = types.StringValue(val)
    } else {
        data.MailgunRegion = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
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

func (r *SmtpConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data SmtpConfigResourceModel
    var state SmtpConfigResourceModel

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

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config SmtpConfigResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    smtpConfigRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := smtpConfigRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Username.IsUnknown() && !state.Username.IsUnknown() && !data.Username.Equal(state.Username) {
        requestDataMap["username"] = data.Username.ValueString()
    }
    if !data.Password.IsUnknown() && !state.Password.IsUnknown() && !data.Password.Equal(state.Password) {
        requestDataMap["password"] = data.Password.ValueString()
    }
    if !data.Hostname.IsUnknown() && !state.Hostname.IsUnknown() && !data.Hostname.Equal(state.Hostname) {
        requestDataMap["hostname"] = data.Hostname.ValueString()
    }
    if !data.Port.IsUnknown() && !state.Port.IsUnknown() && !data.Port.Equal(state.Port) {
        requestDataMap["port"] = r.bigFloatToFloat64(data.Port.ValueBigFloat())
    }
    if !data.FromEmail.IsUnknown() && !state.FromEmail.IsUnknown() && !data.FromEmail.Equal(state.FromEmail) {
        var fromemailData interface{}
        if err := json.Unmarshal([]byte(data.FromEmail.ValueString()), &fromemailData); err == nil {
            requestDataMap["fromEmail"] = fromemailData
        } else {
            requestDataMap["fromEmail"] = data.FromEmail.ValueString()
        }
    }
    if !data.FromName.IsUnknown() && !state.FromName.IsUnknown() && !data.FromName.Equal(state.FromName) {
        requestDataMap["fromName"] = data.FromName.ValueString()
    }
    if !data.Secure.IsUnknown() && !state.Secure.IsUnknown() && !data.Secure.Equal(state.Secure) {
        requestDataMap["secure"] = data.Secure.ValueBool()
    }
    if !data.TransportType.IsUnknown() && !state.TransportType.IsUnknown() && !data.TransportType.Equal(state.TransportType) {
        requestDataMap["transportType"] = data.TransportType.ValueString()
    }
    if !data.AuthType.IsUnknown() && !state.AuthType.IsUnknown() && !data.AuthType.Equal(state.AuthType) {
        requestDataMap["authType"] = data.AuthType.ValueString()
    }
    if !data.ClientId.IsUnknown() && !state.ClientId.IsUnknown() && !data.ClientId.Equal(state.ClientId) {
        requestDataMap["clientId"] = data.ClientId.ValueString()
    }
    if !data.ClientSecret.IsUnknown() && !state.ClientSecret.IsUnknown() && !data.ClientSecret.Equal(state.ClientSecret) {
        requestDataMap["clientSecret"] = data.ClientSecret.ValueString()
    }
    if !data.TokenUrl.IsUnknown() && !state.TokenUrl.IsUnknown() && !data.TokenUrl.Equal(state.TokenUrl) {
        requestDataMap["tokenUrl"] = data.TokenUrl.ValueString()
    }
    if !data.Scope.IsUnknown() && !state.Scope.IsUnknown() && !data.Scope.Equal(state.Scope) {
        requestDataMap["scope"] = data.Scope.ValueString()
    }
    if !data.OauthProviderType.IsUnknown() && !state.OauthProviderType.IsUnknown() && !data.OauthProviderType.Equal(state.OauthProviderType) {
        requestDataMap["oauthProviderType"] = data.OauthProviderType.ValueString()
    }
    if !data.AmazonSesAccessKeyId.IsUnknown() && !state.AmazonSesAccessKeyId.IsUnknown() && !data.AmazonSesAccessKeyId.Equal(state.AmazonSesAccessKeyId) {
        requestDataMap["amazonSesAccessKeyId"] = data.AmazonSesAccessKeyId.ValueString()
    }
    if !data.AmazonSesSecretAccessKey.IsUnknown() && !state.AmazonSesSecretAccessKey.IsUnknown() && !data.AmazonSesSecretAccessKey.Equal(state.AmazonSesSecretAccessKey) {
        requestDataMap["amazonSesSecretAccessKey"] = data.AmazonSesSecretAccessKey.ValueString()
    }
    if !data.AmazonSesRegion.IsUnknown() && !state.AmazonSesRegion.IsUnknown() && !data.AmazonSesRegion.Equal(state.AmazonSesRegion) {
        requestDataMap["amazonSesRegion"] = data.AmazonSesRegion.ValueString()
    }
    if !data.AmazonSesAuthenticationType.IsUnknown() && !state.AmazonSesAuthenticationType.IsUnknown() && !data.AmazonSesAuthenticationType.Equal(state.AmazonSesAuthenticationType) {
        requestDataMap["amazonSesAuthenticationType"] = data.AmazonSesAuthenticationType.ValueString()
    }
    if !data.AmazonSesRoleArn.IsUnknown() && !state.AmazonSesRoleArn.IsUnknown() && !data.AmazonSesRoleArn.Equal(state.AmazonSesRoleArn) {
        requestDataMap["amazonSesRoleArn"] = data.AmazonSesRoleArn.ValueString()
    }
    if !data.SendgridApiKey.IsUnknown() && !state.SendgridApiKey.IsUnknown() && !data.SendgridApiKey.Equal(state.SendgridApiKey) {
        requestDataMap["sendgridApiKey"] = data.SendgridApiKey.ValueString()
    }
    if !data.MailgunApiKey.IsUnknown() && !state.MailgunApiKey.IsUnknown() && !data.MailgunApiKey.Equal(state.MailgunApiKey) {
        requestDataMap["mailgunApiKey"] = data.MailgunApiKey.ValueString()
    }
    if !data.MailgunDomain.IsUnknown() && !state.MailgunDomain.IsUnknown() && !data.MailgunDomain.Equal(state.MailgunDomain) {
        requestDataMap["mailgunDomain"] = data.MailgunDomain.ValueString()
    }
    if !data.MailgunRegion.IsUnknown() && !state.MailgunRegion.IsUnknown() && !data.MailgunRegion.Equal(state.MailgunRegion) {
        requestDataMap["mailgunRegion"] = data.MailgunRegion.ValueString()
    }
    if !data.PostmarkServerToken.IsUnknown() && !state.PostmarkServerToken.IsUnknown() && !data.PostmarkServerToken.Equal(state.PostmarkServerToken) {
        requestDataMap["postmarkServerToken"] = data.PostmarkServerToken.ValueString()
    }
    if !data.ResendApiKey.IsUnknown() && !state.ResendApiKey.IsUnknown() && !data.ResendApiKey.Equal(state.ResendApiKey) {
        requestDataMap["resendApiKey"] = data.ResendApiKey.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(smtpConfigRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/smtp-config/" + data.Id.ValueString() + "", smtpConfigRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update smtp_config, got error: %s", err))
            return
        }

        // Parse the update response
        var smtpConfigResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &smtpConfigResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update smtp_config: %s", err))
            return
        }
        _ = smtpConfigResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "username": true,
        "hostname": true,
        "port": true,
        "fromEmail": true,
        "fromName": true,
        "secure": true,
        "transportType": true,
        "authType": true,
        "clientId": true,
        "tokenUrl": true,
        "scope": true,
        "oauthProviderType": true,
        "amazonSesAccessKeyId": true,
        "amazonSesRegion": true,
        "amazonSesAuthenticationType": true,
        "amazonSesRoleArn": true,
        "mailgunDomain": true,
        "mailgunRegion": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/smtp-config/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read smtp_config after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read smtp_config after update: %s", err))
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := dataMap["description"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["username"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Username = types.StringValue(string(jsonBytes))
            } else {
                data.Username = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Username = types.StringValue(string(jsonBytes))
            } else {
                data.Username = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if arr, ok := dataMap["username"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if val, ok := dataMap["username"].(string); ok {
        data.Username = types.StringValue(val)
    } else {
        data.Username = types.StringNull()
    }
    if obj, ok := dataMap["hostname"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Hostname = types.StringValue(string(jsonBytes))
            } else {
                data.Hostname = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if arr, ok := dataMap["hostname"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := dataMap["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
    }
    data.Port = NumberFromAPI(data.Port, dataMap["port"])
    if obj, ok := dataMap["fromEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEmail = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromEmail = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.FromEmail = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.FromEmail = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["fromEmail"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FromEmail = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.FromEmail = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["fromEmail"].(string); ok {
        data.FromEmail = NewJSONSubsetValue(val)
    } else {
        data.FromEmail = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["fromName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromName = types.StringValue(string(jsonBytes))
            } else {
                data.FromName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromName = types.StringValue(string(jsonBytes))
            } else {
                data.FromName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if arr, ok := dataMap["fromName"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if val, ok := dataMap["fromName"].(string); ok {
        data.FromName = types.StringValue(val)
    } else {
        data.FromName = types.StringNull()
    }
    if val, ok := dataMap["secure"].(bool); ok {
        data.Secure = types.BoolValue(val)
    }
    if obj, ok := dataMap["transportType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TransportType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TransportType = types.StringValue(string(jsonBytes))
            } else {
                data.TransportType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TransportType = types.StringValue(string(jsonBytes))
            } else {
                data.TransportType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if arr, ok := dataMap["transportType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if val, ok := dataMap["transportType"].(string); ok {
        data.TransportType = types.StringValue(val)
    } else {
        data.TransportType = types.StringNull()
    }
    if obj, ok := dataMap["authType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AuthType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AuthType = types.StringValue(string(jsonBytes))
            } else {
                data.AuthType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AuthType = types.StringValue(string(jsonBytes))
            } else {
                data.AuthType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if arr, ok := dataMap["authType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if val, ok := dataMap["authType"].(string); ok {
        data.AuthType = types.StringValue(val)
    } else {
        data.AuthType = types.StringNull()
    }
    if obj, ok := dataMap["clientId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ClientId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ClientId = types.StringValue(string(jsonBytes))
            } else {
                data.ClientId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ClientId = types.StringValue(string(jsonBytes))
            } else {
                data.ClientId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if arr, ok := dataMap["clientId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if val, ok := dataMap["clientId"].(string); ok {
        data.ClientId = types.StringValue(val)
    } else {
        data.ClientId = types.StringNull()
    }
    if obj, ok := dataMap["tokenUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TokenUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TokenUrl = types.StringValue(string(jsonBytes))
            } else {
                data.TokenUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TokenUrl = types.StringValue(string(jsonBytes))
            } else {
                data.TokenUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if arr, ok := dataMap["tokenUrl"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if val, ok := dataMap["tokenUrl"].(string); ok {
        data.TokenUrl = types.StringValue(val)
    } else {
        data.TokenUrl = types.StringNull()
    }
    if obj, ok := dataMap["scope"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Scope = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Scope = types.StringValue(string(jsonBytes))
            } else {
                data.Scope = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Scope = types.StringValue(string(jsonBytes))
            } else {
                data.Scope = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if arr, ok := dataMap["scope"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if val, ok := dataMap["scope"].(string); ok {
        data.Scope = types.StringValue(val)
    } else {
        data.Scope = types.StringNull()
    }
    if obj, ok := dataMap["oauthProviderType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.OauthProviderType = types.StringValue(string(jsonBytes))
            } else {
                data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.OauthProviderType = types.StringValue(string(jsonBytes))
            } else {
                data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if arr, ok := dataMap["oauthProviderType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if val, ok := dataMap["oauthProviderType"].(string); ok {
        data.OauthProviderType = types.StringValue(val)
    } else {
        data.OauthProviderType = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesAccessKeyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesAccessKeyId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesAccessKeyId"].(string); ok {
        data.AmazonSesAccessKeyId = types.StringValue(val)
    } else {
        data.AmazonSesAccessKeyId = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesRegion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesRegion = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesRegion = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesRegion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesRegion"].(string); ok {
        data.AmazonSesRegion = types.StringValue(val)
    } else {
        data.AmazonSesRegion = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesAuthenticationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAuthenticationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesAuthenticationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesAuthenticationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAuthenticationType = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesAuthenticationType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesAuthenticationType = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAuthenticationType = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesAuthenticationType"].(string); ok {
        data.AmazonSesAuthenticationType = types.StringValue(val)
    } else {
        data.AmazonSesAuthenticationType = types.StringNull()
    }
    if obj, ok := dataMap["amazonSesRoleArn"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRoleArn = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AmazonSesRoleArn = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
            } else {
                data.AmazonSesRoleArn = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRoleArn = types.StringNull()
        }
    } else if arr, ok := dataMap["amazonSesRoleArn"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AmazonSesRoleArn = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRoleArn = types.StringNull()
        }
    } else if val, ok := dataMap["amazonSesRoleArn"].(string); ok {
        data.AmazonSesRoleArn = types.StringValue(val)
    } else {
        data.AmazonSesRoleArn = types.StringNull()
    }
    if obj, ok := dataMap["mailgunDomain"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MailgunDomain = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MailgunDomain = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if arr, ok := dataMap["mailgunDomain"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if val, ok := dataMap["mailgunDomain"].(string); ok {
        data.MailgunDomain = types.StringValue(val)
    } else {
        data.MailgunDomain = types.StringNull()
    }
    if obj, ok := dataMap["mailgunRegion"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MailgunRegion = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MailgunRegion = types.StringValue(string(jsonBytes))
            } else {
                data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if arr, ok := dataMap["mailgunRegion"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if val, ok := dataMap["mailgunRegion"].(string); ok {
        data.MailgunRegion = types.StringValue(val)
    } else {
        data.MailgunRegion = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := dataMap["slug"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SmtpConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data SmtpConfigResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/smtp-config/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete smtp_config, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete smtp_config: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *SmtpConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *SmtpConfigResource) keepPlannedValues(data *SmtpConfigResourceModel, plan *SmtpConfigResourceModel, config *SmtpConfigResourceModel) {
    if config.Description.IsNull() && !plan.Description.IsUnknown() {
        data.Description = plan.Description
    }
    if config.Username.IsNull() && !plan.Username.IsUnknown() {
        data.Username = plan.Username
    }
    if config.Hostname.IsNull() && !plan.Hostname.IsUnknown() {
        data.Hostname = plan.Hostname
    }
    if config.Port.IsNull() && !plan.Port.IsUnknown() {
        data.Port = plan.Port
    }
    if config.Secure.IsNull() && !plan.Secure.IsUnknown() {
        data.Secure = plan.Secure
    }
    if config.TransportType.IsNull() && !plan.TransportType.IsUnknown() {
        data.TransportType = plan.TransportType
    }
    if config.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
        data.AuthType = plan.AuthType
    }
    if config.ClientId.IsNull() && !plan.ClientId.IsUnknown() {
        data.ClientId = plan.ClientId
    }
    if config.TokenUrl.IsNull() && !plan.TokenUrl.IsUnknown() {
        data.TokenUrl = plan.TokenUrl
    }
    if config.Scope.IsNull() && !plan.Scope.IsUnknown() {
        data.Scope = plan.Scope
    }
    if config.OauthProviderType.IsNull() && !plan.OauthProviderType.IsUnknown() {
        data.OauthProviderType = plan.OauthProviderType
    }
    if config.AmazonSesAccessKeyId.IsNull() && !plan.AmazonSesAccessKeyId.IsUnknown() {
        data.AmazonSesAccessKeyId = plan.AmazonSesAccessKeyId
    }
    if config.AmazonSesRegion.IsNull() && !plan.AmazonSesRegion.IsUnknown() {
        data.AmazonSesRegion = plan.AmazonSesRegion
    }
    if config.AmazonSesAuthenticationType.IsNull() && !plan.AmazonSesAuthenticationType.IsUnknown() {
        data.AmazonSesAuthenticationType = plan.AmazonSesAuthenticationType
    }
    if config.AmazonSesRoleArn.IsNull() && !plan.AmazonSesRoleArn.IsUnknown() {
        data.AmazonSesRoleArn = plan.AmazonSesRoleArn
    }
    if config.MailgunDomain.IsNull() && !plan.MailgunDomain.IsUnknown() {
        data.MailgunDomain = plan.MailgunDomain
    }
    if config.MailgunRegion.IsNull() && !plan.MailgunRegion.IsUnknown() {
        data.MailgunRegion = plan.MailgunRegion
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *SmtpConfigResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *SmtpConfigResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *SmtpConfigResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *SmtpConfigResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *SmtpConfigResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *SmtpConfigResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *SmtpConfigResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *SmtpConfigResource) isValidOneUptimeObjectType(typeStr string) bool {
    return ValidOneUptimeObjectTypes[typeStr]
}

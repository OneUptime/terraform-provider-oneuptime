package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SmtpConfigDataSource{}

func NewSmtpConfigDataSource() datasource.DataSource {
    return &SmtpConfigDataSource{}
}

// SmtpConfigDataSource defines the data source implementation.
type SmtpConfigDataSource struct {
    client *Client
}

// SmtpConfigDataSourceModel describes the data source data model.
type SmtpConfigDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Username types.String `tfsdk:"username"`
    Password types.String `tfsdk:"password"`
    Hostname types.String `tfsdk:"hostname"`
    Port types.Number `tfsdk:"port"`
    FromEmail types.String `tfsdk:"from_email"`
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
    SendgridApiKey types.String `tfsdk:"sendgrid_api_key"`
    MailgunApiKey types.String `tfsdk:"mailgun_api_key"`
    MailgunDomain types.String `tfsdk:"mailgun_domain"`
    MailgunRegion types.String `tfsdk:"mailgun_region"`
    PostmarkServerToken types.String `tfsdk:"postmark_server_token"`
    ResendApiKey types.String `tfsdk:"resend_api_key"`
}

func (d *SmtpConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_smtp_config"
}

func (d *SmtpConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage Custom SMTP Servers for your project Look up an existing smtp config by `id`, or by any of its other arguments (`name`, `amazon_ses_access_key_id`, `amazon_ses_region`, ...): each one set must match, and exactly one smtp config may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "username": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]",
                Optional: true,
                Computed: true,
            },
            "password": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]",
                Computed: true,
                Sensitive: true,
            },
            "hostname": schema.StringAttribute{
                MarkdownDescription: "SMTP server hostname. Required when Transport is SMTP. Not used by HTTP-API transports like Microsoft Graph.",
                Optional: true,
                Computed: true,
            },
            "port": schema.NumberAttribute{
                MarkdownDescription: "SMTP server port. Required when Transport is SMTP. Not used by HTTP-API transports.",
                Optional: true,
                Computed: true,
            },
            "from_email": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]",
                Computed: true,
            },
            "from_name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Create SMTP Config], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read SMTP Config], Update: [Project Owner, Project Admin, Edit SMTP Config]",
                Optional: true,
                Computed: true,
            },
            "secure": schema.BoolAttribute{
                MarkdownDescription: "Enable secure SMTP connection. Used only for SMTP transport — HTTP-API transports always use HTTPS.",
                Optional: true,
                Computed: true,
            },
            "transport_type": schema.StringAttribute{
                MarkdownDescription: "How OneUptime delivers mail for this config. 'SMTP' uses the hostname and port. 'Microsoft Graph' sends through the Microsoft Graph REST API, for a Microsoft 365 tenant with SMTP AUTH turned off. 'Amazon SES', 'SendGrid', 'Mailgun', 'Postmark' and 'Resend' send through that provider's HTTPS API with its own key, set in the columns named for it.",
                Optional: true,
                Computed: true,
            },
            "auth_type": schema.StringAttribute{
                MarkdownDescription: "The type of authentication to use for this SMTP server. Options: Username and Password, OAuth, or None.",
                Optional: true,
                Computed: true,
            },
            "client_id": schema.StringAttribute{
                MarkdownDescription: "The Client ID from your OAuth application registration. Required for OAuth authentication.",
                Optional: true,
                Computed: true,
            },
            "client_secret": schema.StringAttribute{
                MarkdownDescription: "The Client Secret from your OAuth application registration. Required for OAuth authentication. For Google service accounts, this is the private key from the JSON key file.",
                Computed: true,
                Sensitive: true,
            },
            "token_url": schema.StringAttribute{
                MarkdownDescription: "The OAuth token endpoint URL. For Microsoft 365: https://login.microsoftonline.com/{tenant-id}/oauth2/v2.0/token. For Google: https://oauth2.googleapis.com/token.",
                Optional: true,
                Computed: true,
            },
            "scope": schema.StringAttribute{
                MarkdownDescription: "The OAuth scope(s) required for SMTP access. For Microsoft 365: https://outlook.office365.com/.default. For Google: https://mail.google.com/.",
                Optional: true,
                Computed: true,
            },
            "oauth_provider_type": schema.StringAttribute{
                MarkdownDescription: "The OAuth grant type to use. 'Client Credentials' for Microsoft 365 and most providers. 'JWT Bearer' for Google Workspace service accounts.",
                Optional: true,
                Computed: true,
            },
            "amazon_ses_access_key_id": schema.StringAttribute{
                MarkdownDescription: "The access key ID of an IAM user or role allowed to send with Amazon SES (ses:SendEmail). Required when Transport is Amazon SES.",
                Optional: true,
                Computed: true,
            },
            "amazon_ses_secret_access_key": schema.StringAttribute{
                MarkdownDescription: "The secret access key that goes with the access key ID. Required when Transport is Amazon SES. Encrypted at rest.",
                Computed: true,
                Sensitive: true,
            },
            "amazon_ses_region": schema.StringAttribute{
                MarkdownDescription: "The AWS Region your sender identity is verified in, such as us-east-1. Required when Transport is Amazon SES.",
                Optional: true,
                Computed: true,
            },
            "sendgrid_api_key": schema.StringAttribute{
                MarkdownDescription: "A SendGrid API key with the Mail Send permission. Required when Transport is SendGrid. Encrypted at rest.",
                Computed: true,
                Sensitive: true,
            },
            "mailgun_api_key": schema.StringAttribute{
                MarkdownDescription: "A Mailgun API key, or a sending API key for the domain. Required when Transport is Mailgun. Encrypted at rest.",
                Computed: true,
                Sensitive: true,
            },
            "mailgun_domain": schema.StringAttribute{
                MarkdownDescription: "The Mailgun sending domain mail goes out from, such as mg.example.com. Required when Transport is Mailgun.",
                Optional: true,
                Computed: true,
            },
            "mailgun_region": schema.StringAttribute{
                MarkdownDescription: "The region the Mailgun domain was created in: US or EU. US when not set.",
                Optional: true,
                Computed: true,
            },
            "postmark_server_token": schema.StringAttribute{
                MarkdownDescription: "The API token of the Postmark server that sends, from its API Tokens tab. Required when Transport is Postmark. Encrypted at rest.",
                Computed: true,
                Sensitive: true,
            },
            "resend_api_key": schema.StringAttribute{
                MarkdownDescription: "A Resend API key with sending access. Required when Transport is Resend. Encrypted at rest.",
                Computed: true,
                Sensitive: true,
            },
        },
    }
}

func (d *SmtpConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SmtpConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SmtpConfigDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.Username.IsNull() && !data.Username.IsUnknown() {
        filters["username"] = data.Username.ValueString()
        filterNames = append(filterNames, "username = "+fmt.Sprintf("%q", data.Username.ValueString()))
    }
    if !data.Hostname.IsNull() && !data.Hostname.IsUnknown() {
        filters["hostname"] = data.Hostname.ValueString()
        filterNames = append(filterNames, "hostname = "+fmt.Sprintf("%q", data.Hostname.ValueString()))
    }
    if !data.Port.IsNull() && !data.Port.IsUnknown() {
        filters["port"] = lookupNumber(data.Port)
        filterNames = append(filterNames, "port = "+data.Port.ValueBigFloat().String())
    }
    if !data.FromName.IsNull() && !data.FromName.IsUnknown() {
        filters["fromName"] = data.FromName.ValueString()
        filterNames = append(filterNames, "from_name = "+fmt.Sprintf("%q", data.FromName.ValueString()))
    }
    if !data.Secure.IsNull() && !data.Secure.IsUnknown() {
        filters["secure"] = data.Secure.ValueBool()
        filterNames = append(filterNames, "secure = "+fmt.Sprintf("%t", data.Secure.ValueBool()))
    }
    if !data.TransportType.IsNull() && !data.TransportType.IsUnknown() {
        filters["transportType"] = data.TransportType.ValueString()
        filterNames = append(filterNames, "transport_type = "+fmt.Sprintf("%q", data.TransportType.ValueString()))
    }
    if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
        filters["authType"] = data.AuthType.ValueString()
        filterNames = append(filterNames, "auth_type = "+fmt.Sprintf("%q", data.AuthType.ValueString()))
    }
    if !data.ClientId.IsNull() && !data.ClientId.IsUnknown() {
        filters["clientId"] = data.ClientId.ValueString()
        filterNames = append(filterNames, "client_id = "+fmt.Sprintf("%q", data.ClientId.ValueString()))
    }
    if !data.TokenUrl.IsNull() && !data.TokenUrl.IsUnknown() {
        filters["tokenUrl"] = data.TokenUrl.ValueString()
        filterNames = append(filterNames, "token_url = "+fmt.Sprintf("%q", data.TokenUrl.ValueString()))
    }
    if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
        filters["scope"] = data.Scope.ValueString()
        filterNames = append(filterNames, "scope = "+fmt.Sprintf("%q", data.Scope.ValueString()))
    }
    if !data.OauthProviderType.IsNull() && !data.OauthProviderType.IsUnknown() {
        filters["oauthProviderType"] = data.OauthProviderType.ValueString()
        filterNames = append(filterNames, "oauth_provider_type = "+fmt.Sprintf("%q", data.OauthProviderType.ValueString()))
    }
    if !data.AmazonSesAccessKeyId.IsNull() && !data.AmazonSesAccessKeyId.IsUnknown() {
        filters["amazonSesAccessKeyId"] = data.AmazonSesAccessKeyId.ValueString()
        filterNames = append(filterNames, "amazon_ses_access_key_id = "+fmt.Sprintf("%q", data.AmazonSesAccessKeyId.ValueString()))
    }
    if !data.AmazonSesRegion.IsNull() && !data.AmazonSesRegion.IsUnknown() {
        filters["amazonSesRegion"] = data.AmazonSesRegion.ValueString()
        filterNames = append(filterNames, "amazon_ses_region = "+fmt.Sprintf("%q", data.AmazonSesRegion.ValueString()))
    }
    if !data.MailgunDomain.IsNull() && !data.MailgunDomain.IsUnknown() {
        filters["mailgunDomain"] = data.MailgunDomain.ValueString()
        filterNames = append(filterNames, "mailgun_domain = "+fmt.Sprintf("%q", data.MailgunDomain.ValueString()))
    }
    if !data.MailgunRegion.IsNull() && !data.MailgunRegion.IsUnknown() {
        filters["mailgunRegion"] = data.MailgunRegion.ValueString()
        filterNames = append(filterNames, "mailgun_region = "+fmt.Sprintf("%q", data.MailgunRegion.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the smtp config up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the smtp config up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "createdByUserId": true,
        "username": true,
        "password": true,
        "hostname": true,
        "port": true,
        "fromEmail": true,
        "fromName": true,
        "secure": true,
        "transportType": true,
        "authType": true,
        "clientId": true,
        "clientSecret": true,
        "tokenUrl": true,
        "scope": true,
        "oauthProviderType": true,
        "amazonSesAccessKeyId": true,
        "amazonSesSecretAccessKey": true,
        "amazonSesRegion": true,
        "sendgridApiKey": true,
        "mailgunApiKey": true,
        "mailgunDomain": true,
        "mailgunRegion": true,
        "postmarkServerToken": true,
        "resendApiKey": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/smtp-config/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read smtp_config, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No smtp config found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read smtp_config: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/smtp-config/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list smtp_config, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list smtp_config: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No smtp config matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one smtp config matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for smtp_config.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if arr, ok := item["slug"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := item["username"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Username = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Username = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if arr, ok := item["username"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Username = types.StringValue(string(jsonBytes))
        } else {
            data.Username = types.StringNull()
        }
    } else if val, ok := item["username"].(string); ok {
        data.Username = types.StringValue(val)
    } else {
        data.Username = types.StringNull()
    }
    if obj, ok := item["password"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Password = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Password = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Password = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Password = types.StringValue(string(jsonBytes))
        } else {
            data.Password = types.StringNull()
        }
    } else if arr, ok := item["password"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Password = types.StringValue(string(jsonBytes))
        } else {
            data.Password = types.StringNull()
        }
    } else if val, ok := item["password"].(string); ok {
        data.Password = types.StringValue(val)
    } else {
        data.Password = types.StringNull()
    }
    if obj, ok := item["hostname"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Hostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Hostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if arr, ok := item["hostname"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Hostname = types.StringValue(string(jsonBytes))
        } else {
            data.Hostname = types.StringNull()
        }
    } else if val, ok := item["hostname"].(string); ok {
        data.Hostname = types.StringValue(val)
    } else {
        data.Hostname = types.StringNull()
    }
    data.Port = numberFromAPI(data.Port, item["port"])
    if obj, ok := item["fromEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FromEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FromEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FromEmail = types.StringValue(string(jsonBytes))
        } else {
            data.FromEmail = types.StringNull()
        }
    } else if arr, ok := item["fromEmail"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FromEmail = types.StringValue(string(jsonBytes))
        } else {
            data.FromEmail = types.StringNull()
        }
    } else if val, ok := item["fromEmail"].(string); ok {
        data.FromEmail = types.StringValue(val)
    } else {
        data.FromEmail = types.StringNull()
    }
    if obj, ok := item["fromName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FromName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FromName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if arr, ok := item["fromName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FromName = types.StringValue(string(jsonBytes))
        } else {
            data.FromName = types.StringNull()
        }
    } else if val, ok := item["fromName"].(string); ok {
        data.FromName = types.StringValue(val)
    } else {
        data.FromName = types.StringNull()
    }
    if val, ok := item["secure"].(bool); ok {
        data.Secure = types.BoolValue(val)
    } else {
        data.Secure = types.BoolNull()
    }
    if obj, ok := item["transportType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TransportType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TransportType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if arr, ok := item["transportType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TransportType = types.StringValue(string(jsonBytes))
        } else {
            data.TransportType = types.StringNull()
        }
    } else if val, ok := item["transportType"].(string); ok {
        data.TransportType = types.StringValue(val)
    } else {
        data.TransportType = types.StringNull()
    }
    if obj, ok := item["authType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AuthType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AuthType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if arr, ok := item["authType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AuthType = types.StringValue(string(jsonBytes))
        } else {
            data.AuthType = types.StringNull()
        }
    } else if val, ok := item["authType"].(string); ok {
        data.AuthType = types.StringValue(val)
    } else {
        data.AuthType = types.StringNull()
    }
    if obj, ok := item["clientId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClientId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClientId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if arr, ok := item["clientId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ClientId = types.StringValue(string(jsonBytes))
        } else {
            data.ClientId = types.StringNull()
        }
    } else if val, ok := item["clientId"].(string); ok {
        data.ClientId = types.StringValue(val)
    } else {
        data.ClientId = types.StringNull()
    }
    if obj, ok := item["clientSecret"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClientSecret = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClientSecret = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClientSecret = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClientSecret = types.StringValue(string(jsonBytes))
        } else {
            data.ClientSecret = types.StringNull()
        }
    } else if arr, ok := item["clientSecret"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ClientSecret = types.StringValue(string(jsonBytes))
        } else {
            data.ClientSecret = types.StringNull()
        }
    } else if val, ok := item["clientSecret"].(string); ok {
        data.ClientSecret = types.StringValue(val)
    } else {
        data.ClientSecret = types.StringNull()
    }
    if obj, ok := item["tokenUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TokenUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TokenUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if arr, ok := item["tokenUrl"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TokenUrl = types.StringValue(string(jsonBytes))
        } else {
            data.TokenUrl = types.StringNull()
        }
    } else if val, ok := item["tokenUrl"].(string); ok {
        data.TokenUrl = types.StringValue(val)
    } else {
        data.TokenUrl = types.StringNull()
    }
    if obj, ok := item["scope"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Scope = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Scope = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if arr, ok := item["scope"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Scope = types.StringValue(string(jsonBytes))
        } else {
            data.Scope = types.StringNull()
        }
    } else if val, ok := item["scope"].(string); ok {
        data.Scope = types.StringValue(val)
    } else {
        data.Scope = types.StringNull()
    }
    if obj, ok := item["oauthProviderType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OauthProviderType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OauthProviderType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if arr, ok := item["oauthProviderType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.OauthProviderType = types.StringValue(string(jsonBytes))
        } else {
            data.OauthProviderType = types.StringNull()
        }
    } else if val, ok := item["oauthProviderType"].(string); ok {
        data.OauthProviderType = types.StringValue(val)
    } else {
        data.OauthProviderType = types.StringNull()
    }
    if obj, ok := item["amazonSesAccessKeyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AmazonSesAccessKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AmazonSesAccessKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if arr, ok := item["amazonSesAccessKeyId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AmazonSesAccessKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesAccessKeyId = types.StringNull()
        }
    } else if val, ok := item["amazonSesAccessKeyId"].(string); ok {
        data.AmazonSesAccessKeyId = types.StringValue(val)
    } else {
        data.AmazonSesAccessKeyId = types.StringNull()
    }
    if obj, ok := item["amazonSesSecretAccessKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesSecretAccessKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AmazonSesSecretAccessKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AmazonSesSecretAccessKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AmazonSesSecretAccessKey = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesSecretAccessKey = types.StringNull()
        }
    } else if arr, ok := item["amazonSesSecretAccessKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AmazonSesSecretAccessKey = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesSecretAccessKey = types.StringNull()
        }
    } else if val, ok := item["amazonSesSecretAccessKey"].(string); ok {
        data.AmazonSesSecretAccessKey = types.StringValue(val)
    } else {
        data.AmazonSesSecretAccessKey = types.StringNull()
    }
    if obj, ok := item["amazonSesRegion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AmazonSesRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AmazonSesRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if arr, ok := item["amazonSesRegion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AmazonSesRegion = types.StringValue(string(jsonBytes))
        } else {
            data.AmazonSesRegion = types.StringNull()
        }
    } else if val, ok := item["amazonSesRegion"].(string); ok {
        data.AmazonSesRegion = types.StringValue(val)
    } else {
        data.AmazonSesRegion = types.StringNull()
    }
    if obj, ok := item["sendgridApiKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendgridApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SendgridApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SendgridApiKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SendgridApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.SendgridApiKey = types.StringNull()
        }
    } else if arr, ok := item["sendgridApiKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SendgridApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.SendgridApiKey = types.StringNull()
        }
    } else if val, ok := item["sendgridApiKey"].(string); ok {
        data.SendgridApiKey = types.StringValue(val)
    } else {
        data.SendgridApiKey = types.StringNull()
    }
    if obj, ok := item["mailgunApiKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MailgunApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MailgunApiKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MailgunApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunApiKey = types.StringNull()
        }
    } else if arr, ok := item["mailgunApiKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MailgunApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunApiKey = types.StringNull()
        }
    } else if val, ok := item["mailgunApiKey"].(string); ok {
        data.MailgunApiKey = types.StringValue(val)
    } else {
        data.MailgunApiKey = types.StringNull()
    }
    if obj, ok := item["mailgunDomain"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MailgunDomain = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MailgunDomain = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if arr, ok := item["mailgunDomain"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MailgunDomain = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunDomain = types.StringNull()
        }
    } else if val, ok := item["mailgunDomain"].(string); ok {
        data.MailgunDomain = types.StringValue(val)
    } else {
        data.MailgunDomain = types.StringNull()
    }
    if obj, ok := item["mailgunRegion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MailgunRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MailgunRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if arr, ok := item["mailgunRegion"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MailgunRegion = types.StringValue(string(jsonBytes))
        } else {
            data.MailgunRegion = types.StringNull()
        }
    } else if val, ok := item["mailgunRegion"].(string); ok {
        data.MailgunRegion = types.StringValue(val)
    } else {
        data.MailgunRegion = types.StringNull()
    }
    if obj, ok := item["postmarkServerToken"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PostmarkServerToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PostmarkServerToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PostmarkServerToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PostmarkServerToken = types.StringValue(string(jsonBytes))
        } else {
            data.PostmarkServerToken = types.StringNull()
        }
    } else if arr, ok := item["postmarkServerToken"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PostmarkServerToken = types.StringValue(string(jsonBytes))
        } else {
            data.PostmarkServerToken = types.StringNull()
        }
    } else if val, ok := item["postmarkServerToken"].(string); ok {
        data.PostmarkServerToken = types.StringValue(val)
    } else {
        data.PostmarkServerToken = types.StringNull()
    }
    if obj, ok := item["resendApiKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResendApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResendApiKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResendApiKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResendApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.ResendApiKey = types.StringNull()
        }
    } else if arr, ok := item["resendApiKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ResendApiKey = types.StringValue(string(jsonBytes))
        } else {
            data.ResendApiKey = types.StringNull()
        }
    } else if val, ok := item["resendApiKey"].(string); ok {
        data.ResendApiKey = types.StringValue(val)
    } else {
        data.ResendApiKey = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

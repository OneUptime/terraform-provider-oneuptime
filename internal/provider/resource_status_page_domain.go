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
var _ resource.Resource = &StatusPageDomainResource{}
var _ resource.ResourceWithImportState = &StatusPageDomainResource{}

func NewStatusPageDomainResource() resource.Resource {
    return &StatusPageDomainResource{}
}

// StatusPageDomainResource defines the resource implementation.
type StatusPageDomainResource struct {
    client *Client
}

// StatusPageDomainResourceModel describes the resource data model.
type StatusPageDomainResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    DomainId types.String `tfsdk:"domain_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    Subdomain types.String `tfsdk:"subdomain"`
    FullDomain types.String `tfsdk:"full_domain"`
    CnameVerificationToken types.String `tfsdk:"cname_verification_token"`
    CustomCertificate types.String `tfsdk:"custom_certificate"`
    CustomCertificateKey types.String `tfsdk:"custom_certificate_key"`
    IsCustomCertificate types.Bool `tfsdk:"is_custom_certificate"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsCnameVerified types.Bool `tfsdk:"is_cname_verified"`
    IsSslOrdered types.Bool `tfsdk:"is_ssl_ordered"`
    IsSslProvisioned types.Bool `tfsdk:"is_ssl_provisioned"`
}

func (r *StatusPageDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_domain"
}

func (r *StatusPageDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_domain resource",

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
            "domain_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "subdomain": schema.StringAttribute{
                MarkdownDescription: "Sumdomain",
                Required: true,
            },
            "full_domain": schema.StringAttribute{
                MarkdownDescription: "Full Domain",
                Required: true,
            },
            "cname_verification_token": schema.StringAttribute{
                MarkdownDescription: "CNAME Verification Token",
                Required: true,
            },
            "custom_certificate": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
                Optional: true,
            },
            "custom_certificate_key": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
                Optional: true,
            },
            "is_custom_certificate": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Domain], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Domain]",
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
            "is_cname_verified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_ssl_ordered": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_ssl_provisioned": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Status Page Domain], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *StatusPageDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *StatusPageDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageDomainResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    statusPageDomainRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "domainId": data.DomainId.ValueString(),
        "statusPageId": data.StatusPageId.ValueString(),
        "subdomain": data.Subdomain.ValueString(),
        "fullDomain": data.FullDomain.ValueString(),
        "cnameVerificationToken": data.CnameVerificationToken.ValueString(),
        "customCertificate": data.CustomCertificate.ValueString(),
        "customCertificateKey": data.CustomCertificateKey.ValueString(),
        "isCustomCertificate": data.IsCustomCertificate.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/status-page-domain", statusPageDomainRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page_domain, got error: %s", err))
        return
    }

    var statusPageDomainResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageDomainResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_domain response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageDomainResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageDomainResponse
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
    if val, ok := dataMap["domainId"].(string); ok && val != "" {
        data.DomainId = types.StringValue(val)
    } else {
        data.DomainId = types.StringNull()
    }
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["subdomain"].(string); ok && val != "" {
        data.Subdomain = types.StringValue(val)
    } else {
        data.Subdomain = types.StringNull()
    }
    if val, ok := dataMap["fullDomain"].(string); ok && val != "" {
        data.FullDomain = types.StringValue(val)
    } else {
        data.FullDomain = types.StringNull()
    }
    if val, ok := dataMap["cnameVerificationToken"].(string); ok && val != "" {
        data.CnameVerificationToken = types.StringValue(val)
    } else {
        data.CnameVerificationToken = types.StringNull()
    }
    if val, ok := dataMap["customCertificate"].(string); ok && val != "" {
        data.CustomCertificate = types.StringValue(val)
    } else {
        data.CustomCertificate = types.StringNull()
    }
    if val, ok := dataMap["customCertificateKey"].(string); ok && val != "" {
        data.CustomCertificateKey = types.StringValue(val)
    } else {
        data.CustomCertificateKey = types.StringNull()
    }
    if val, ok := dataMap["isCustomCertificate"].(bool); ok {
        data.IsCustomCertificate = types.BoolValue(val)
    } else if dataMap["isCustomCertificate"] == nil {
        data.IsCustomCertificate = types.BoolNull()
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
    if val, ok := dataMap["isCnameVerified"].(bool); ok {
        data.IsCnameVerified = types.BoolValue(val)
    } else if dataMap["isCnameVerified"] == nil {
        data.IsCnameVerified = types.BoolNull()
    }
    if val, ok := dataMap["isSslOrdered"].(bool); ok {
        data.IsSslOrdered = types.BoolValue(val)
    } else if dataMap["isSslOrdered"] == nil {
        data.IsSslOrdered = types.BoolNull()
    }
    if val, ok := dataMap["isSslProvisioned"].(bool); ok {
        data.IsSslProvisioned = types.BoolValue(val)
    } else if dataMap["isSslProvisioned"] == nil {
        data.IsSslProvisioned = types.BoolNull()
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

func (r *StatusPageDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageDomainResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "domainId": true,
        "statusPageId": true,
        "subdomain": true,
        "fullDomain": true,
        "cnameVerificationToken": true,
        "customCertificate": true,
        "customCertificateKey": true,
        "isCustomCertificate": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "isCnameVerified": true,
        "isSslOrdered": true,
        "isSslProvisioned": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/status-page-domain/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_domain, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var statusPageDomainResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageDomainResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_domain response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageDomainResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageDomainResponse
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
    if val, ok := dataMap["domainId"].(string); ok && val != "" {
        data.DomainId = types.StringValue(val)
    } else {
        data.DomainId = types.StringNull()
    }
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["subdomain"].(string); ok && val != "" {
        data.Subdomain = types.StringValue(val)
    } else {
        data.Subdomain = types.StringNull()
    }
    if val, ok := dataMap["fullDomain"].(string); ok && val != "" {
        data.FullDomain = types.StringValue(val)
    } else {
        data.FullDomain = types.StringNull()
    }
    if val, ok := dataMap["cnameVerificationToken"].(string); ok && val != "" {
        data.CnameVerificationToken = types.StringValue(val)
    } else {
        data.CnameVerificationToken = types.StringNull()
    }
    if val, ok := dataMap["customCertificate"].(string); ok && val != "" {
        data.CustomCertificate = types.StringValue(val)
    } else {
        data.CustomCertificate = types.StringNull()
    }
    if val, ok := dataMap["customCertificateKey"].(string); ok && val != "" {
        data.CustomCertificateKey = types.StringValue(val)
    } else {
        data.CustomCertificateKey = types.StringNull()
    }
    if val, ok := dataMap["isCustomCertificate"].(bool); ok {
        data.IsCustomCertificate = types.BoolValue(val)
    } else if dataMap["isCustomCertificate"] == nil {
        data.IsCustomCertificate = types.BoolNull()
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
    if val, ok := dataMap["isCnameVerified"].(bool); ok {
        data.IsCnameVerified = types.BoolValue(val)
    } else if dataMap["isCnameVerified"] == nil {
        data.IsCnameVerified = types.BoolNull()
    }
    if val, ok := dataMap["isSslOrdered"].(bool); ok {
        data.IsSslOrdered = types.BoolValue(val)
    } else if dataMap["isSslOrdered"] == nil {
        data.IsSslOrdered = types.BoolNull()
    }
    if val, ok := dataMap["isSslProvisioned"].(bool); ok {
        data.IsSslProvisioned = types.BoolValue(val)
    } else if dataMap["isSslProvisioned"] == nil {
        data.IsSslProvisioned = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StatusPageDomainResourceModel
    var state StatusPageDomainResourceModel

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
    statusPageDomainRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "domainId": data.DomainId.ValueString(),
        "statusPageId": data.StatusPageId.ValueString(),
        "subdomain": data.Subdomain.ValueString(),
        "fullDomain": data.FullDomain.ValueString(),
        "cnameVerificationToken": data.CnameVerificationToken.ValueString(),
        "customCertificate": data.CustomCertificate.ValueString(),
        "customCertificateKey": data.CustomCertificateKey.ValueString(),
        "isCustomCertificate": data.IsCustomCertificate.ValueBool(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/status-page-domain/" + data.Id.ValueString() + "", statusPageDomainRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status_page_domain, got error: %s", err))
        return
    }

    // Parse the update response
    var statusPageDomainResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageDomainResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_domain response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "domainId": true,
        "statusPageId": true,
        "subdomain": true,
        "fullDomain": true,
        "cnameVerificationToken": true,
        "customCertificate": true,
        "customCertificateKey": true,
        "isCustomCertificate": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "isCnameVerified": true,
        "isSslOrdered": true,
        "isSslProvisioned": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/status-page-domain/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_domain after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_domain read response, got error: %s", err))
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
    if val, ok := dataMap["domainId"].(string); ok && val != "" {
        data.DomainId = types.StringValue(val)
    } else {
        data.DomainId = types.StringNull()
    }
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["subdomain"].(string); ok && val != "" {
        data.Subdomain = types.StringValue(val)
    } else {
        data.Subdomain = types.StringNull()
    }
    if val, ok := dataMap["fullDomain"].(string); ok && val != "" {
        data.FullDomain = types.StringValue(val)
    } else {
        data.FullDomain = types.StringNull()
    }
    if val, ok := dataMap["cnameVerificationToken"].(string); ok && val != "" {
        data.CnameVerificationToken = types.StringValue(val)
    } else {
        data.CnameVerificationToken = types.StringNull()
    }
    if val, ok := dataMap["customCertificate"].(string); ok && val != "" {
        data.CustomCertificate = types.StringValue(val)
    } else {
        data.CustomCertificate = types.StringNull()
    }
    if val, ok := dataMap["customCertificateKey"].(string); ok && val != "" {
        data.CustomCertificateKey = types.StringValue(val)
    } else {
        data.CustomCertificateKey = types.StringNull()
    }
    if val, ok := dataMap["isCustomCertificate"].(bool); ok {
        data.IsCustomCertificate = types.BoolValue(val)
    } else if dataMap["isCustomCertificate"] == nil {
        data.IsCustomCertificate = types.BoolNull()
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
    if val, ok := dataMap["isCnameVerified"].(bool); ok {
        data.IsCnameVerified = types.BoolValue(val)
    } else if dataMap["isCnameVerified"] == nil {
        data.IsCnameVerified = types.BoolNull()
    }
    if val, ok := dataMap["isSslOrdered"].(bool); ok {
        data.IsSslOrdered = types.BoolValue(val)
    } else if dataMap["isSslOrdered"] == nil {
        data.IsSslOrdered = types.BoolNull()
    }
    if val, ok := dataMap["isSslProvisioned"].(bool); ok {
        data.IsSslProvisioned = types.BoolValue(val)
    } else if dataMap["isSslProvisioned"] == nil {
        data.IsSslProvisioned = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageDomainResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/status-page-domain/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status_page_domain, got error: %s", err))
        return
    }
}


func (r *StatusPageDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *StatusPageDomainResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *StatusPageDomainResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}

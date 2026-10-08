package provider

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
    return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
    client *Client
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Name types.String `tfsdk:"name"`
    Email types.String `tfsdk:"email"`
    NewUnverifiedTemporaryEmail types.String `tfsdk:"new_unverified_temporary_email"`
    Password types.String `tfsdk:"password"`
    IsEmailVerified types.Bool `tfsdk:"is_email_verified"`
    CompanyName types.String `tfsdk:"company_name"`
    JobRole types.String `tfsdk:"job_role"`
    CompanySize types.String `tfsdk:"company_size"`
    Referral types.String `tfsdk:"referral"`
    CompanyPhoneNumber types.String `tfsdk:"company_phone_number"`
    ProfilePictureId types.String `tfsdk:"profile_picture_id"`
    TwoFactorAuthEnabled types.Bool `tfsdk:"two_factor_auth_enabled"`
    Timezone types.String `tfsdk:"timezone"`
    IsDisabled types.Bool `tfsdk:"is_disabled"`
    IsBlocked types.Bool `tfsdk:"is_blocked"`
    EnableTwoFactorAuth types.Bool `tfsdk:"enable_two_factor_auth"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A signed up or invited OneUptime user. Look up an existing user by `id`, or by any of its other arguments (`name`, `company_name`, `company_size`, ...): each one set must match, and exactly one user may match them all.",

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
            "name": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "email": schema.StringAttribute{
                MarkdownDescription: "Email.",
                Computed: true,
            },
            "new_unverified_temporary_email": schema.StringAttribute{
                Computed: true,
            },
            "password": schema.StringAttribute{
                MarkdownDescription: "Password.",
                Computed: true,
                Sensitive: true,
            },
            "is_email_verified": schema.BoolAttribute{
                Optional: true,
                Computed: true,
            },
            "company_name": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "job_role": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "company_size": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "referral": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "company_phone_number": schema.StringAttribute{
                Computed: true,
            },
            "profile_picture_id": schema.StringAttribute{
                MarkdownDescription: "The ID of a `oneuptime_file`.",
                Optional: true,
                Computed: true,
            },
            "two_factor_auth_enabled": schema.BoolAttribute{
                Optional: true,
                Computed: true,
            },
            "timezone": schema.StringAttribute{
                Optional: true,
                Computed: true,
            },
            "is_disabled": schema.BoolAttribute{
                Optional: true,
                Computed: true,
            },
            "is_blocked": schema.BoolAttribute{
                Optional: true,
                Computed: true,
            },
            "enable_two_factor_auth": schema.BoolAttribute{
                MarkdownDescription: "Is two factor authentication enabled?",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data UserDataSourceModel

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
    if !data.IsEmailVerified.IsNull() && !data.IsEmailVerified.IsUnknown() {
        filters["isEmailVerified"] = data.IsEmailVerified.ValueBool()
        filterNames = append(filterNames, "is_email_verified = "+fmt.Sprintf("%t", data.IsEmailVerified.ValueBool()))
    }
    if !data.CompanyName.IsNull() && !data.CompanyName.IsUnknown() {
        filters["companyName"] = data.CompanyName.ValueString()
        filterNames = append(filterNames, "company_name = "+fmt.Sprintf("%q", data.CompanyName.ValueString()))
    }
    if !data.JobRole.IsNull() && !data.JobRole.IsUnknown() {
        filters["jobRole"] = data.JobRole.ValueString()
        filterNames = append(filterNames, "job_role = "+fmt.Sprintf("%q", data.JobRole.ValueString()))
    }
    if !data.CompanySize.IsNull() && !data.CompanySize.IsUnknown() {
        filters["companySize"] = data.CompanySize.ValueString()
        filterNames = append(filterNames, "company_size = "+fmt.Sprintf("%q", data.CompanySize.ValueString()))
    }
    if !data.Referral.IsNull() && !data.Referral.IsUnknown() {
        filters["referral"] = data.Referral.ValueString()
        filterNames = append(filterNames, "referral = "+fmt.Sprintf("%q", data.Referral.ValueString()))
    }
    if !data.ProfilePictureId.IsNull() && !data.ProfilePictureId.IsUnknown() {
        filters["profilePictureId"] = data.ProfilePictureId.ValueString()
        filterNames = append(filterNames, "profile_picture_id = "+fmt.Sprintf("%q", data.ProfilePictureId.ValueString()))
    }
    if !data.TwoFactorAuthEnabled.IsNull() && !data.TwoFactorAuthEnabled.IsUnknown() {
        filters["twoFactorAuthEnabled"] = data.TwoFactorAuthEnabled.ValueBool()
        filterNames = append(filterNames, "two_factor_auth_enabled = "+fmt.Sprintf("%t", data.TwoFactorAuthEnabled.ValueBool()))
    }
    if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
        filters["timezone"] = data.Timezone.ValueString()
        filterNames = append(filterNames, "timezone = "+fmt.Sprintf("%q", data.Timezone.ValueString()))
    }
    if !data.IsDisabled.IsNull() && !data.IsDisabled.IsUnknown() {
        filters["isDisabled"] = data.IsDisabled.ValueBool()
        filterNames = append(filterNames, "is_disabled = "+fmt.Sprintf("%t", data.IsDisabled.ValueBool()))
    }
    if !data.IsBlocked.IsNull() && !data.IsBlocked.IsUnknown() {
        filters["isBlocked"] = data.IsBlocked.ValueBool()
        filterNames = append(filterNames, "is_blocked = "+fmt.Sprintf("%t", data.IsBlocked.ValueBool()))
    }
    if !data.EnableTwoFactorAuth.IsNull() && !data.EnableTwoFactorAuth.IsUnknown() {
        filters["enableTwoFactorAuth"] = data.EnableTwoFactorAuth.ValueBool()
        filterNames = append(filterNames, "enable_two_factor_auth = "+fmt.Sprintf("%t", data.EnableTwoFactorAuth.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the user up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the user up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "name": true,
        "email": true,
        "newUnverifiedTemporaryEmail": true,
        "password": true,
        "isEmailVerified": true,
        "companyName": true,
        "jobRole": true,
        "companySize": true,
        "referral": true,
        "companyPhoneNumber": true,
        "profilePictureId": true,
        "twoFactorAuthEnabled": true,
        "timezone": true,
        "isDisabled": true,
        "isBlocked": true,
        "enableTwoFactorAuth": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        // No get endpoint: find it in the list by id.
        filters["_id"] = data.Id.ValueString()
        filterNames = append(filterNames, fmt.Sprintf("id = %q", data.Id.ValueString()))
    }
    if item == nil {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/user/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list user, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list user: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No user matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one user matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for user.")
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
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
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
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["email"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Email = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Email = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Email = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Email = types.StringValue(string(jsonBytes))
        } else {
            data.Email = types.StringNull()
        }
    } else if val, ok := item["email"].(string); ok {
        data.Email = types.StringValue(val)
    } else {
        data.Email = types.StringNull()
    }
    if obj, ok := item["newUnverifiedTemporaryEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NewUnverifiedTemporaryEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NewUnverifiedTemporaryEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NewUnverifiedTemporaryEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NewUnverifiedTemporaryEmail = types.StringValue(string(jsonBytes))
        } else {
            data.NewUnverifiedTemporaryEmail = types.StringNull()
        }
    } else if val, ok := item["newUnverifiedTemporaryEmail"].(string); ok {
        data.NewUnverifiedTemporaryEmail = types.StringValue(val)
    } else {
        data.NewUnverifiedTemporaryEmail = types.StringNull()
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
    } else if val, ok := item["password"].(string); ok {
        data.Password = types.StringValue(val)
    } else {
        data.Password = types.StringNull()
    }
    if val, ok := item["isEmailVerified"].(bool); ok {
        data.IsEmailVerified = types.BoolValue(val)
    } else {
        data.IsEmailVerified = types.BoolNull()
    }
    if obj, ok := item["companyName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompanyName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompanyName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompanyName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompanyName = types.StringValue(string(jsonBytes))
        } else {
            data.CompanyName = types.StringNull()
        }
    } else if val, ok := item["companyName"].(string); ok {
        data.CompanyName = types.StringValue(val)
    } else {
        data.CompanyName = types.StringNull()
    }
    if obj, ok := item["jobRole"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.JobRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.JobRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.JobRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.JobRole = types.StringValue(string(jsonBytes))
        } else {
            data.JobRole = types.StringNull()
        }
    } else if val, ok := item["jobRole"].(string); ok {
        data.JobRole = types.StringValue(val)
    } else {
        data.JobRole = types.StringNull()
    }
    if obj, ok := item["companySize"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompanySize = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompanySize = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompanySize = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompanySize = types.StringValue(string(jsonBytes))
        } else {
            data.CompanySize = types.StringNull()
        }
    } else if val, ok := item["companySize"].(string); ok {
        data.CompanySize = types.StringValue(val)
    } else {
        data.CompanySize = types.StringNull()
    }
    if obj, ok := item["referral"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Referral = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Referral = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Referral = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Referral = types.StringValue(string(jsonBytes))
        } else {
            data.Referral = types.StringNull()
        }
    } else if val, ok := item["referral"].(string); ok {
        data.Referral = types.StringValue(val)
    } else {
        data.Referral = types.StringNull()
    }
    if obj, ok := item["companyPhoneNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompanyPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompanyPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompanyPhoneNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompanyPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CompanyPhoneNumber = types.StringNull()
        }
    } else if val, ok := item["companyPhoneNumber"].(string); ok {
        data.CompanyPhoneNumber = types.StringValue(val)
    } else {
        data.CompanyPhoneNumber = types.StringNull()
    }
    if obj, ok := item["profilePictureId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProfilePictureId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProfilePictureId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProfilePictureId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProfilePictureId = types.StringValue(string(jsonBytes))
        } else {
            data.ProfilePictureId = types.StringNull()
        }
    } else if val, ok := item["profilePictureId"].(string); ok {
        data.ProfilePictureId = types.StringValue(val)
    } else {
        data.ProfilePictureId = types.StringNull()
    }
    if val, ok := item["twoFactorAuthEnabled"].(bool); ok {
        data.TwoFactorAuthEnabled = types.BoolValue(val)
    } else {
        data.TwoFactorAuthEnabled = types.BoolNull()
    }
    if obj, ok := item["timezone"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Timezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Timezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Timezone = types.StringValue(string(jsonBytes))
        } else {
            data.Timezone = types.StringNull()
        }
    } else if val, ok := item["timezone"].(string); ok {
        data.Timezone = types.StringValue(val)
    } else {
        data.Timezone = types.StringNull()
    }
    if val, ok := item["isDisabled"].(bool); ok {
        data.IsDisabled = types.BoolValue(val)
    } else {
        data.IsDisabled = types.BoolNull()
    }
    if val, ok := item["isBlocked"].(bool); ok {
        data.IsBlocked = types.BoolValue(val)
    } else {
        data.IsBlocked = types.BoolNull()
    }
    if val, ok := item["enableTwoFactorAuth"].(bool); ok {
        data.EnableTwoFactorAuth = types.BoolValue(val)
    } else {
        data.EnableTwoFactorAuth = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

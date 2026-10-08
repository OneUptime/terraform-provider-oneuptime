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
var _ datasource.DataSource = &IncomingCallPolicyPhoneNumberDataSource{}

func NewIncomingCallPolicyPhoneNumberDataSource() datasource.DataSource {
    return &IncomingCallPolicyPhoneNumberDataSource{}
}

// IncomingCallPolicyPhoneNumberDataSource defines the data source implementation.
type IncomingCallPolicyPhoneNumberDataSource struct {
    client *Client
}

// IncomingCallPolicyPhoneNumberDataSourceModel describes the data source data model.
type IncomingCallPolicyPhoneNumberDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallPolicyId types.String `tfsdk:"incoming_call_policy_id"`
    ProjectCallSmsConfigId types.String `tfsdk:"project_call_sms_config_id"`
    PhoneNumber types.String `tfsdk:"phone_number"`
    CallProviderPhoneNumberId types.String `tfsdk:"call_provider_phone_number_id"`
    CountryCode types.String `tfsdk:"country_code"`
    AreaCode types.String `tfsdk:"area_code"`
    PhoneNumberPurchasedAt types.String `tfsdk:"phone_number_purchased_at"`
}

func (d *IncomingCallPolicyPhoneNumberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_policy_phone_number"
}

func (d *IncomingCallPolicyPhoneNumberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Phone numbers that route incoming calls to an incoming call policy. Look up an existing incoming call policy phone number by `id`, or by any of its other arguments (`area_code`, `call_provider_phone_number_id`, `country_code`, ...): each one set must match, and exactly one incoming call policy phone number may match them all.",

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
                MarkdownDescription: "ID of the project that owns this phone number. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "incoming_call_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of the policy that receives calls to this number. The ID of a `oneuptime_incoming_call_policy`.",
                Optional: true,
                Computed: true,
            },
            "project_call_sms_config_id": schema.StringAttribute{
                MarkdownDescription: "ID of the call provider configuration that owns this number.",
                Optional: true,
                Computed: true,
            },
            "phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone number that routes calls to the policy.",
                Computed: true,
            },
            "call_provider_phone_number_id": schema.StringAttribute{
                MarkdownDescription: "The call provider identifier for this phone number.",
                Optional: true,
                Computed: true,
            },
            "country_code": schema.StringAttribute{
                MarkdownDescription: "Country code associated with this phone number.",
                Optional: true,
                Computed: true,
            },
            "area_code": schema.StringAttribute{
                MarkdownDescription: "Area code associated with this phone number.",
                Optional: true,
                Computed: true,
            },
            "phone_number_purchased_at": schema.StringAttribute{
                MarkdownDescription: "When this phone number was attached to the policy.",
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallPolicyPhoneNumberDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallPolicyPhoneNumberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallPolicyPhoneNumberDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IncomingCallPolicyId.IsNull() && !data.IncomingCallPolicyId.IsUnknown() {
        filters["incomingCallPolicyId"] = data.IncomingCallPolicyId.ValueString()
        filterNames = append(filterNames, "incoming_call_policy_id = "+fmt.Sprintf("%q", data.IncomingCallPolicyId.ValueString()))
    }
    if !data.ProjectCallSmsConfigId.IsNull() && !data.ProjectCallSmsConfigId.IsUnknown() {
        filters["projectCallSMSConfigId"] = data.ProjectCallSmsConfigId.ValueString()
        filterNames = append(filterNames, "project_call_sms_config_id = "+fmt.Sprintf("%q", data.ProjectCallSmsConfigId.ValueString()))
    }
    if !data.CallProviderPhoneNumberId.IsNull() && !data.CallProviderPhoneNumberId.IsUnknown() {
        filters["callProviderPhoneNumberId"] = data.CallProviderPhoneNumberId.ValueString()
        filterNames = append(filterNames, "call_provider_phone_number_id = "+fmt.Sprintf("%q", data.CallProviderPhoneNumberId.ValueString()))
    }
    if !data.CountryCode.IsNull() && !data.CountryCode.IsUnknown() {
        filters["countryCode"] = data.CountryCode.ValueString()
        filterNames = append(filterNames, "country_code = "+fmt.Sprintf("%q", data.CountryCode.ValueString()))
    }
    if !data.AreaCode.IsNull() && !data.AreaCode.IsUnknown() {
        filters["areaCode"] = data.AreaCode.ValueString()
        filterNames = append(filterNames, "area_code = "+fmt.Sprintf("%q", data.AreaCode.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incoming call policy phone number up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incoming call policy phone number up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incomingCallPolicyId": true,
        "projectCallSMSConfigId": true,
        "phoneNumber": true,
        "callProviderPhoneNumberId": true,
        "countryCode": true,
        "areaCode": true,
        "phoneNumberPurchasedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incoming-call-policy-phone-number/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_policy_phone_number, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incoming call policy phone number found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incoming_call_policy_phone_number: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incoming-call-policy-phone-number/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incoming_call_policy_phone_number, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incoming_call_policy_phone_number: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incoming call policy phone number matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incoming call policy phone number matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incoming_call_policy_phone_number.")
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
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["incomingCallPolicyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if val, ok := item["incomingCallPolicyId"].(string); ok {
        data.IncomingCallPolicyId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyId = types.StringNull()
    }
    if obj, ok := item["projectCallSMSConfigId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectCallSmsConfigId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectCallSmsConfigId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectCallSmsConfigId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectCallSmsConfigId = types.StringNull()
        }
    } else if val, ok := item["projectCallSMSConfigId"].(string); ok {
        data.ProjectCallSmsConfigId = types.StringValue(val)
    } else {
        data.ProjectCallSmsConfigId = types.StringNull()
    }
    if obj, ok := item["phoneNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PhoneNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.PhoneNumber = types.StringNull()
        }
    } else if val, ok := item["phoneNumber"].(string); ok {
        data.PhoneNumber = types.StringValue(val)
    } else {
        data.PhoneNumber = types.StringNull()
    }
    if obj, ok := item["callProviderPhoneNumberId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallProviderPhoneNumberId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallProviderPhoneNumberId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallProviderPhoneNumberId = types.StringValue(string(jsonBytes))
        } else {
            data.CallProviderPhoneNumberId = types.StringNull()
        }
    } else if val, ok := item["callProviderPhoneNumberId"].(string); ok {
        data.CallProviderPhoneNumberId = types.StringValue(val)
    } else {
        data.CallProviderPhoneNumberId = types.StringNull()
    }
    if obj, ok := item["countryCode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CountryCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CountryCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CountryCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CountryCode = types.StringValue(string(jsonBytes))
        } else {
            data.CountryCode = types.StringNull()
        }
    } else if val, ok := item["countryCode"].(string); ok {
        data.CountryCode = types.StringValue(val)
    } else {
        data.CountryCode = types.StringNull()
    }
    if obj, ok := item["areaCode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AreaCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AreaCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AreaCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AreaCode = types.StringValue(string(jsonBytes))
        } else {
            data.AreaCode = types.StringNull()
        }
    } else if val, ok := item["areaCode"].(string); ok {
        data.AreaCode = types.StringValue(val)
    } else {
        data.AreaCode = types.StringNull()
    }
    if obj, ok := item["phoneNumberPurchasedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PhoneNumberPurchasedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PhoneNumberPurchasedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PhoneNumberPurchasedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PhoneNumberPurchasedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PhoneNumberPurchasedAt = types.StringNull()
        }
    } else if val, ok := item["phoneNumberPurchasedAt"].(string); ok {
        data.PhoneNumberPurchasedAt = types.StringValue(val)
    } else {
        data.PhoneNumberPurchasedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

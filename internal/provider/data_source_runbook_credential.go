package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RunbookCredentialDataSource{}

func NewRunbookCredentialDataSource() datasource.DataSource {
    return &RunbookCredentialDataSource{}
}

// RunbookCredentialDataSource defines the data source implementation.
type RunbookCredentialDataSource struct {
    client *Client
}

// RunbookCredentialDataSourceModel describes the data source data model.
type RunbookCredentialDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    CredentialType types.String `tfsdk:"credential_type"`
    SshHostname types.String `tfsdk:"ssh_hostname"`
    SshPort types.Number `tfsdk:"ssh_port"`
    SshUsername types.String `tfsdk:"ssh_username"`
    KubernetesApiServerUrl types.String `tfsdk:"kubernetes_api_server_url"`
    KubernetesCaCertificate types.String `tfsdk:"kubernetes_ca_certificate"`
    Runners types.Set `tfsdk:"runners"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *RunbookCredentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook_credential"
}

func (d *RunbookCredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Access to a system a runbook needs to act on — an SSH host, or a Kubernetes cluster. Secret material is encrypted at rest and can never be read back through the API; it is decrypted only when handed to an assigned Runner as it claims a step. Look up an existing runbook credential by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `credential_type`, ...): each one set must match, and exactly one runbook credential may match them all.",

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
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "credential_type": schema.StringAttribute{
                MarkdownDescription: "SSH, or Kubernetes.",
                Optional: true,
                Computed: true,
            },
            "ssh_hostname": schema.StringAttribute{
                MarkdownDescription: "Hostname or IP address the Runner connects to.",
                Optional: true,
                Computed: true,
            },
            "ssh_port": schema.NumberAttribute{
                MarkdownDescription: "Defaults to 22 when unset.",
                Optional: true,
                Computed: true,
            },
            "ssh_username": schema.StringAttribute{
                MarkdownDescription: "The user the Runner authenticates as.",
                Optional: true,
                Computed: true,
            },
            "kubernetes_api_server_url": schema.StringAttribute{
                MarkdownDescription: "For example https://10.0.0.1:6443.",
                Optional: true,
                Computed: true,
            },
            "kubernetes_ca_certificate": schema.StringAttribute{
                MarkdownDescription: "PEM certificate authority for the API server. Leave empty only if the API server presents a certificate your Runner already trusts.",
                Optional: true,
                Computed: true,
            },
            "runners": schema.SetAttribute{
                MarkdownDescription: "The Runners allowed to use this credential. A step referencing it must target one of them. IDs of `oneuptime_runner` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *RunbookCredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookCredentialDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CredentialType.IsNull() && !data.CredentialType.IsUnknown() {
        filters["credentialType"] = data.CredentialType.ValueString()
        filterNames = append(filterNames, "credential_type = "+fmt.Sprintf("%q", data.CredentialType.ValueString()))
    }
    if !data.SshHostname.IsNull() && !data.SshHostname.IsUnknown() {
        filters["sshHostname"] = data.SshHostname.ValueString()
        filterNames = append(filterNames, "ssh_hostname = "+fmt.Sprintf("%q", data.SshHostname.ValueString()))
    }
    if !data.SshPort.IsNull() && !data.SshPort.IsUnknown() {
        filters["sshPort"] = lookupNumber(data.SshPort)
        filterNames = append(filterNames, "ssh_port = "+data.SshPort.ValueBigFloat().String())
    }
    if !data.SshUsername.IsNull() && !data.SshUsername.IsUnknown() {
        filters["sshUsername"] = data.SshUsername.ValueString()
        filterNames = append(filterNames, "ssh_username = "+fmt.Sprintf("%q", data.SshUsername.ValueString()))
    }
    if !data.KubernetesApiServerUrl.IsNull() && !data.KubernetesApiServerUrl.IsUnknown() {
        filters["kubernetesApiServerUrl"] = data.KubernetesApiServerUrl.ValueString()
        filterNames = append(filterNames, "kubernetes_api_server_url = "+fmt.Sprintf("%q", data.KubernetesApiServerUrl.ValueString()))
    }
    if !data.KubernetesCaCertificate.IsNull() && !data.KubernetesCaCertificate.IsUnknown() {
        filters["kubernetesCaCertificate"] = data.KubernetesCaCertificate.ValueString()
        filterNames = append(filterNames, "kubernetes_ca_certificate = "+fmt.Sprintf("%q", data.KubernetesCaCertificate.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the runbook credential up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the runbook credential up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "credentialType": true,
        "sshHostname": true,
        "sshPort": true,
        "sshUsername": true,
        "kubernetesApiServerUrl": true,
        "kubernetesCaCertificate": true,
        "runners": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/runbook-credential/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook_credential, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runbook credential found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read runbook_credential: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/runbook-credential/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list runbook_credential, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list runbook_credential: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runbook credential matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one runbook credential matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for runbook_credential.")
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
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["credentialType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CredentialType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CredentialType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CredentialType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CredentialType = types.StringValue(string(jsonBytes))
        } else {
            data.CredentialType = types.StringNull()
        }
    } else if val, ok := item["credentialType"].(string); ok {
        data.CredentialType = types.StringValue(val)
    } else {
        data.CredentialType = types.StringNull()
    }
    if obj, ok := item["sshHostname"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SshHostname = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SshHostname = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SshHostname = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SshHostname = types.StringValue(string(jsonBytes))
        } else {
            data.SshHostname = types.StringNull()
        }
    } else if val, ok := item["sshHostname"].(string); ok {
        data.SshHostname = types.StringValue(val)
    } else {
        data.SshHostname = types.StringNull()
    }
    if val, ok := item["sshPort"].(float64); ok {
        data.SshPort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["sshPort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SshPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.SshPort = types.NumberNull()
        }
    } else {
        data.SshPort = types.NumberNull()
    }
    if obj, ok := item["sshUsername"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SshUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SshUsername = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SshUsername = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SshUsername = types.StringValue(string(jsonBytes))
        } else {
            data.SshUsername = types.StringNull()
        }
    } else if val, ok := item["sshUsername"].(string); ok {
        data.SshUsername = types.StringValue(val)
    } else {
        data.SshUsername = types.StringNull()
    }
    if obj, ok := item["kubernetesApiServerUrl"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesApiServerUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KubernetesApiServerUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KubernetesApiServerUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KubernetesApiServerUrl = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesApiServerUrl = types.StringNull()
        }
    } else if val, ok := item["kubernetesApiServerUrl"].(string); ok {
        data.KubernetesApiServerUrl = types.StringValue(val)
    } else {
        data.KubernetesApiServerUrl = types.StringNull()
    }
    if obj, ok := item["kubernetesCaCertificate"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KubernetesCaCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KubernetesCaCertificate = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KubernetesCaCertificate = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KubernetesCaCertificate = types.StringValue(string(jsonBytes))
        } else {
            data.KubernetesCaCertificate = types.StringNull()
        }
    } else if val, ok := item["kubernetesCaCertificate"].(string); ok {
        data.KubernetesCaCertificate = types.StringValue(val)
    } else {
        data.KubernetesCaCertificate = types.StringNull()
    }
    if val, ok := item["runners"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Runners = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Runners = types.SetNull(types.StringType)
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
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

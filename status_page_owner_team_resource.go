package oneuptime

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StatusPageOwnerTeamResource{}
var _ resource.ResourceWithImportState = &StatusPageOwnerTeamResource{}

func NewStatusPageOwnerTeamResource() resource.Resource {
    return &StatusPageOwnerTeamResource{}
}

// StatusPageOwnerTeamResource defines the resource implementation.
type StatusPageOwnerTeamResource struct{}

// Metadata returns the resource type name.
func (r *StatusPageOwnerTeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_owner_team"
}

// Schema defines the schema for the resource.
func (r *StatusPageOwnerTeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = StatusPageOwnerTeamResourceSchema(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *StatusPageOwnerTeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    // Add client configuration here when API client is implemented
}

// Create creates the resource and sets the initial Terraform state.
func (r *StatusPageOwnerTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageOwnerTeamModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // TODO: Implement API call to create resource
    // For now, set a placeholder ID
    data.Id = types.StringValue("placeholder-id")

    // Write logs using the tflog package
    // Documentation: https://terraform.io/plugin/log
    // tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *StatusPageOwnerTeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageOwnerTeamModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // TODO: Implement API call to read resource

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *StatusPageOwnerTeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StatusPageOwnerTeamModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // TODO: Implement API call to update resource

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *StatusPageOwnerTeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageOwnerTeamModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // TODO: Implement API call to delete resource
}

// ImportState imports the resource into Terraform state.
func (r *StatusPageOwnerTeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // TODO: Implement resource import
    resp.Diagnostics.AddError(
        "Import Not Implemented",
        "Import is not yet implemented for this resource.",
    )
}

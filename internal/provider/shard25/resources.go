package shard25

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard25", []OrderedResource{
        {Order: 1, New: NewProjectResource},
        {Order: 17, New: NewLogPipelineResource},
        {Order: 30, New: NewSecurityCasePullRequestResource},
        {Order: 90, New: NewOnCallPolicyResource},
        {Order: 185, New: NewAlertMeasurementResource},
        {Order: 241, New: NewStatusPageCustomFieldResource},
        {Order: 242, New: NewStatusPageResourceResource},
        {Order: 258, New: NewWorkflowResource},
        {Order: 388, New: NewNetworkSiteTypeResource},
        {Order: 396, New: NewNetworkDeviceNamingRuleResource},
        {Order: 413, New: NewIotFleetTeamOwnerResource},
        {Order: 415, New: NewIotDeviceCredentialResource},
        {Order: 432, New: NewServerlessFunctionLabelRuleResource},
        {Order: 441, New: NewCloudResourceOwnerRuleResource},
        {Order: 480, New: NewIotFleetTeamOwnerLegacyResource},
        {Order: 482, New: NewIotDeviceCredentialLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 413, Name: "iot_fleet_team_owner", LegacyName: "io_t_fleet_team_owner", New: NewIotFleetTeamOwnerResource, NewLegacy: NewIotFleetTeamOwnerLegacyResource},
        {Order: 415, Name: "iot_device_credential", LegacyName: "io_t_device_credential", New: NewIotDeviceCredentialResource, NewLegacy: NewIotDeviceCredentialLegacyResource},
    })
}

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
        {Order: 86, New: NewOnCallPolicyResource},
        {Order: 181, New: NewAlertMeasurementResource},
        {Order: 237, New: NewStatusPageCustomFieldResource},
        {Order: 238, New: NewStatusPageResourceResource},
        {Order: 254, New: NewWorkflowResource},
        {Order: 384, New: NewNetworkSiteTypeResource},
        {Order: 392, New: NewNetworkDeviceNamingRuleResource},
        {Order: 409, New: NewIotFleetTeamOwnerResource},
        {Order: 411, New: NewIotDeviceCredentialResource},
        {Order: 428, New: NewServerlessFunctionLabelRuleResource},
        {Order: 437, New: NewCloudResourceOwnerRuleResource},
        {Order: 476, New: NewIotFleetTeamOwnerLegacyResource},
        {Order: 478, New: NewIotDeviceCredentialLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 409, Name: "iot_fleet_team_owner", LegacyName: "io_t_fleet_team_owner", New: NewIotFleetTeamOwnerResource, NewLegacy: NewIotFleetTeamOwnerLegacyResource},
        {Order: 411, Name: "iot_device_credential", LegacyName: "io_t_device_credential", New: NewIotDeviceCredentialResource, NewLegacy: NewIotDeviceCredentialLegacyResource},
    })
}

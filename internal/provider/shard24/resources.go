package shard24

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard24", []OrderedResource{
        {Order: 13, New: NewLogSavedViewResource},
        {Order: 25, New: NewSecurityEventConnectionResource},
        {Order: 59, New: NewPipelineTeamOwnerResource},
        {Order: 87, New: NewOnCallPolicyCustomFieldResource},
        {Order: 106, New: NewIncidentStateTimelineResource},
        {Order: 125, New: NewFormOidcResource},
        {Order: 203, New: NewHostLabelRuleResource},
        {Order: 246, New: NewStatusPagePrivateUserResource},
        {Order: 247, New: NewStatusPageHistoryChartBarColorResource},
        {Order: 335, New: NewDashboardPrivateUserResource},
        {Order: 373, New: NewNetworkDeviceAutoImportRuleResource},
        {Order: 385, New: NewNetworkEndpointResource},
        {Order: 393, New: NewAddressNameResource},
        {Order: 422, New: NewHostTeamOwnerResource},
        {Order: 424, New: NewCollectorConfigResource},
        {Order: 433, New: NewCloudResourceResource},
    }, nil)
}

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
        {Order: 60, New: NewPipelineTeamOwnerResource},
        {Order: 91, New: NewOnCallPolicyCustomFieldResource},
        {Order: 110, New: NewIncidentStateTimelineResource},
        {Order: 129, New: NewFormOidcResource},
        {Order: 207, New: NewHostLabelRuleResource},
        {Order: 250, New: NewStatusPagePrivateUserResource},
        {Order: 251, New: NewStatusPageHistoryChartBarColorResource},
        {Order: 339, New: NewDashboardPrivateUserResource},
        {Order: 377, New: NewNetworkDeviceAutoImportRuleResource},
        {Order: 389, New: NewNetworkEndpointResource},
        {Order: 397, New: NewAddressNameResource},
        {Order: 426, New: NewHostTeamOwnerResource},
        {Order: 428, New: NewCollectorConfigResource},
        {Order: 437, New: NewCloudResourceResource},
    }, nil)
}

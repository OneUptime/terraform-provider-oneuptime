package shard21

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard21", []OrderedResource{
        {Order: 23, New: NewWatchlistValueResource},
        {Order: 63, New: NewDataObservabilitySettingResource},
        {Order: 65, New: NewBlockResource},
        {Order: 124, New: NewPostmortemActionItemResource},
        {Order: 194, New: NewSloFeedResource},
        {Order: 216, New: NewProxmoxClusterOwnerRuleResource},
        {Order: 227, New: NewStorageArrayLabelRuleResource},
        {Order: 249, New: NewStatusPageHeaderLinkResource},
        {Order: 254, New: NewScheduledMaintenanceStateTimelineResource},
        {Order: 274, New: NewTeamUserOwnerResource},
        {Order: 275, New: NewRunbookVariableResource},
        {Order: 368, New: NewKubernetesClusterResource},
        {Order: 374, New: NewNetworkDeviceUserOwnerResource},
        {Order: 391, New: NewNetworkSiteLinkResource},
        {Order: 423, New: NewVcenterTeamOwnerResource},
        {Order: 484, New: NewVcenterTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 423, Name: "vcenter_team_owner", LegacyName: "v_center_team_owner", New: NewVcenterTeamOwnerResource, NewLegacy: NewVcenterTeamOwnerLegacyResource},
    })
}

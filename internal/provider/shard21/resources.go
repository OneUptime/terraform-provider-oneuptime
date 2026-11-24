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
        {Order: 62, New: NewDataObservabilitySettingResource},
        {Order: 64, New: NewBlockResource},
        {Order: 120, New: NewPostmortemActionItemResource},
        {Order: 190, New: NewSloFeedResource},
        {Order: 212, New: NewProxmoxClusterOwnerRuleResource},
        {Order: 223, New: NewStorageArrayLabelRuleResource},
        {Order: 245, New: NewStatusPageHeaderLinkResource},
        {Order: 250, New: NewScheduledMaintenanceStateTimelineResource},
        {Order: 270, New: NewTeamUserOwnerResource},
        {Order: 271, New: NewRunbookVariableResource},
        {Order: 364, New: NewKubernetesClusterResource},
        {Order: 370, New: NewNetworkDeviceUserOwnerResource},
        {Order: 387, New: NewNetworkSiteLinkResource},
        {Order: 419, New: NewVcenterTeamOwnerResource},
        {Order: 480, New: NewVcenterTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 419, Name: "vcenter_team_owner", LegacyName: "v_center_team_owner", New: NewVcenterTeamOwnerResource, NewLegacy: NewVcenterTeamOwnerLegacyResource},
    })
}

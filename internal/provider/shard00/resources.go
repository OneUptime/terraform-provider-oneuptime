package shard00

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard00", []OrderedResource{
        {Order: 55, New: NewTableUserOwnerResource},
        {Order: 121, New: NewFormResource},
        {Order: 144, New: NewAlertOnCallRuleResource},
        {Order: 160, New: NewIncidentEpisodeUserOwnerResource},
        {Order: 198, New: NewMonitorLabelRuleResource},
        {Order: 222, New: NewStorageArrayOwnerRuleResource},
        {Order: 276, New: NewMonitorTeamOwnerResource},
        {Order: 315, New: NewProbeTeamOwnerResource},
        {Order: 329, New: NewInventoryItemResource},
        {Order: 346, New: NewWorkspaceNotificationSummaryResource},
        {Order: 352, New: NewCephClusterFeedResource},
        {Order: 375, New: NewNetworkDeviceRoleResource},
        {Order: 376, New: NewNetworkAlertPolicyResource},
        {Order: 377, New: NewSnmpCredentialProfileResource},
        {Order: 388, New: NewNetworkDeviceLinkResource},
        {Order: 421, New: NewHostResource},
        {Order: 467, New: NewProbeTeamOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 315, Name: "probe_team_owner", LegacyName: "probe_owner_team", New: NewProbeTeamOwnerResource, NewLegacy: NewProbeTeamOwnerLegacyResource},
    })
}

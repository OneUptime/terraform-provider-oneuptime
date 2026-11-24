package shard23

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard23", []OrderedResource{
        {Order: 48, New: NewDataAssetResource},
        {Order: 65, New: NewEntityCriticalityResource},
        {Order: 74, New: NewAiConnectionResource},
        {Order: 104, New: NewIncidentFeedResource},
        {Order: 107, New: NewIncidentInternalNoteResource},
        {Order: 108, New: NewIncidentPublicNoteResource},
        {Order: 199, New: NewStatusPageOwnerRuleResource},
        {Order: 230, New: NewDashboardOwnerRuleResource},
        {Order: 269, New: NewTeamTeamOwnerResource},
        {Order: 307, New: NewIncomingCallPolicyOwnerRuleResource},
        {Order: 309, New: NewCallAndSmsConfigResource},
        {Order: 356, New: NewOpenstackCloudFeedResource},
        {Order: 369, New: NewNetworkDeviceTeamOwnerResource},
        {Order: 372, New: NewNetworkDeviceLabelRuleResource},
        {Order: 408, New: NewIotFleetResource},
        {Order: 410, New: NewIotFleetUserOwnerResource},
        {Order: 431, New: NewCloudBudgetResource},
        {Order: 439, New: NewDatabaseEndpointResource},
        {Order: 470, New: NewOpenstackCloudFeedLegacyResource},
        {Order: 475, New: NewIotFleetLegacyResource},
        {Order: 477, New: NewIotFleetUserOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 356, Name: "openstack_cloud_feed", LegacyName: "open_stack_cloud_feed", New: NewOpenstackCloudFeedResource, NewLegacy: NewOpenstackCloudFeedLegacyResource},
        {Order: 408, Name: "iot_fleet", LegacyName: "io_t_fleet", New: NewIotFleetResource, NewLegacy: NewIotFleetLegacyResource},
        {Order: 410, Name: "iot_fleet_user_owner", LegacyName: "io_t_fleet_user_owner", New: NewIotFleetUserOwnerResource, NewLegacy: NewIotFleetUserOwnerLegacyResource},
    })
}

package shard23

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard23", []OrderedResource{
        {Order: 49, New: NewDataAssetResource},
        {Order: 69, New: NewEntityCriticalityResource},
        {Order: 78, New: NewAiConnectionResource},
        {Order: 108, New: NewIncidentFeedResource},
        {Order: 111, New: NewIncidentInternalNoteResource},
        {Order: 112, New: NewIncidentPublicNoteResource},
        {Order: 203, New: NewStatusPageOwnerRuleResource},
        {Order: 234, New: NewDashboardOwnerRuleResource},
        {Order: 273, New: NewTeamTeamOwnerResource},
        {Order: 311, New: NewIncomingCallPolicyOwnerRuleResource},
        {Order: 313, New: NewCallAndSmsConfigResource},
        {Order: 360, New: NewOpenstackCloudFeedResource},
        {Order: 373, New: NewNetworkDeviceTeamOwnerResource},
        {Order: 376, New: NewNetworkDeviceLabelRuleResource},
        {Order: 412, New: NewIotFleetResource},
        {Order: 414, New: NewIotFleetUserOwnerResource},
        {Order: 435, New: NewCloudBudgetResource},
        {Order: 443, New: NewDatabaseEndpointResource},
        {Order: 474, New: NewOpenstackCloudFeedLegacyResource},
        {Order: 479, New: NewIotFleetLegacyResource},
        {Order: 481, New: NewIotFleetUserOwnerLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 360, Name: "openstack_cloud_feed", LegacyName: "open_stack_cloud_feed", New: NewOpenstackCloudFeedResource, NewLegacy: NewOpenstackCloudFeedLegacyResource},
        {Order: 412, Name: "iot_fleet", LegacyName: "io_t_fleet", New: NewIotFleetResource, NewLegacy: NewIotFleetLegacyResource},
        {Order: 414, Name: "iot_fleet_user_owner", LegacyName: "io_t_fleet_user_owner", New: NewIotFleetUserOwnerResource, NewLegacy: NewIotFleetUserOwnerLegacyResource},
    })
}

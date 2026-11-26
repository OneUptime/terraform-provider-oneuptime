package shard06

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard06", []OrderedResource{
        {Order: 7, New: NewCustomRoleResource},
        {Order: 21, New: NewRuleExceptionResource},
        {Order: 35, New: NewSecurityCaseCustomFieldResource},
        {Order: 36, New: NewSecurityPlaybookResource},
        {Order: 50, New: NewDataAssetFeedResource},
        {Order: 83, New: NewTraceScrubRuleResource},
        {Order: 127, New: NewFormPrivateUserResource},
        {Order: 138, New: NewAlertSeverityResource},
        {Order: 162, New: NewIncidentEpisodeRoleMemberResource},
        {Order: 196, New: NewSloLabelRuleResource},
        {Order: 256, New: NewScheduledMaintenanceCustomFieldResource},
        {Order: 293, New: NewMonitorGroupUserOwnerResource},
        {Order: 327, New: NewDataSourceResource},
        {Order: 382, New: NewNetworkDeviceDiscoveryScanResource},
        {Order: 390, New: NewNetworkSiteStatusEventResource},
        {Order: 400, New: NewPodmanHostResource},
    }, nil)
}

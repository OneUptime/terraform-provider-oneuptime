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
        {Order: 34, New: NewSecurityCaseCustomFieldResource},
        {Order: 35, New: NewSecurityPlaybookResource},
        {Order: 49, New: NewDataAssetFeedResource},
        {Order: 79, New: NewTraceScrubRuleResource},
        {Order: 123, New: NewFormPrivateUserResource},
        {Order: 134, New: NewAlertSeverityResource},
        {Order: 158, New: NewIncidentEpisodeRoleMemberResource},
        {Order: 192, New: NewSloLabelRuleResource},
        {Order: 252, New: NewScheduledMaintenanceCustomFieldResource},
        {Order: 289, New: NewMonitorGroupUserOwnerResource},
        {Order: 323, New: NewDataSourceResource},
        {Order: 378, New: NewNetworkDeviceDiscoveryScanResource},
        {Order: 386, New: NewNetworkSiteStatusEventResource},
        {Order: 396, New: NewPodmanHostResource},
    }, nil)
}

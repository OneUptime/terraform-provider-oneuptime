package shard17

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard17", []OrderedResource{
        {Order: 111, New: NewIncidentTemplateUserOwnerResource},
        {Order: 141, New: NewAlertEpisodeInternalNoteResource},
        {Order: 164, New: NewIncidentEpisodePublicNoteResource},
        {Order: 168, New: NewIncidentOwnerRuleResource},
        {Order: 205, New: NewServiceLabelRuleResource},
        {Order: 213, New: NewOpenstackCloudOwnerRuleResource},
        {Order: 236, New: NewStatusPageDomainResource},
        {Order: 248, New: NewScheduledMaintenanceStateResource},
        {Order: 251, New: NewScheduledEventPublicNoteResource},
        {Order: 274, New: NewStatusPageScimResource},
        {Order: 334, New: NewDashboardDomainResource},
        {Order: 343, New: NewVideoCallConnectionResource},
        {Order: 354, New: NewPodmanHostFeedResource},
        {Order: 362, New: NewMetricTypeResource},
        {Order: 379, New: NewNetworkDeviceDiagnosticResource},
        {Order: 414, New: NewCephClusterUserOwnerResource},
        {Order: 426, New: NewServerlessFunctionTeamOwnerResource},
        {Order: 450, New: NewRumApplicationResource},
        {Order: 459, New: NewOpenstackCloudOwnerRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 213, Name: "openstack_cloud_owner_rule", LegacyName: "open_stack_cloud_owner_rule", New: NewOpenstackCloudOwnerRuleResource, NewLegacy: NewOpenstackCloudOwnerRuleLegacyResource},
    })
}

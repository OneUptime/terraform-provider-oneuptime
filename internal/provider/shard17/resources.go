package shard17

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard17", []OrderedResource{
        {Order: 115, New: NewIncidentTemplateUserOwnerResource},
        {Order: 145, New: NewAlertEpisodeInternalNoteResource},
        {Order: 168, New: NewIncidentEpisodePublicNoteResource},
        {Order: 172, New: NewIncidentOwnerRuleResource},
        {Order: 209, New: NewServiceLabelRuleResource},
        {Order: 217, New: NewOpenstackCloudOwnerRuleResource},
        {Order: 240, New: NewStatusPageDomainResource},
        {Order: 252, New: NewScheduledMaintenanceStateResource},
        {Order: 255, New: NewScheduledEventPublicNoteResource},
        {Order: 278, New: NewStatusPageScimResource},
        {Order: 338, New: NewDashboardDomainResource},
        {Order: 347, New: NewVideoCallConnectionResource},
        {Order: 358, New: NewPodmanHostFeedResource},
        {Order: 366, New: NewMetricTypeResource},
        {Order: 383, New: NewNetworkDeviceDiagnosticResource},
        {Order: 418, New: NewCephClusterUserOwnerResource},
        {Order: 430, New: NewServerlessFunctionTeamOwnerResource},
        {Order: 454, New: NewRumApplicationResource},
        {Order: 463, New: NewOpenstackCloudOwnerRuleLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 217, Name: "openstack_cloud_owner_rule", LegacyName: "open_stack_cloud_owner_rule", New: NewOpenstackCloudOwnerRuleResource, NewLegacy: NewOpenstackCloudOwnerRuleLegacyResource},
    })
}

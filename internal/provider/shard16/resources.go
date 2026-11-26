package shard16

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard16", []OrderedResource{
        {Order: 3, New: NewTeamCustomFieldResource},
        {Order: 10, New: NewTeamComplianceSettingResource},
        {Order: 70, New: NewSecurityEventDestinationResource},
        {Order: 73, New: NewSecurityTelemetrySourceResource},
        {Order: 107, New: NewIncidentResource},
        {Order: 121, New: NewIncidentNoteTemplateResource},
        {Order: 228, New: NewVcenterOwnerRuleResource},
        {Order: 237, New: NewFileResource},
        {Order: 247, New: NewStatusPageSubscriberNotificationTemplateResource},
        {Order: 265, New: NewRunbookUserOwnerResource},
        {Order: 270, New: NewRunnerResource},
        {Order: 318, New: NewCodeRepositoryResource},
        {Order: 323, New: NewAiAgentUserOwnerResource},
        {Order: 332, New: NewDsymFileResource},
        {Order: 353, New: NewKubernetesClusterFeedResource},
        {Order: 387, New: NewNetworkSiteUserOwnerResource},
        {Order: 394, New: NewNetworkTopologySuppressionResource},
        {Order: 404, New: NewOpenstackCloudResource},
        {Order: 438, New: NewCloudResourceTeamOwnerResource},
        {Order: 440, New: NewCloudResourceLabelRuleResource},
        {Order: 467, New: NewVcenterOwnerRuleLegacyResource},
        {Order: 473, New: NewDsymFileLegacyResource},
        {Order: 476, New: NewOpenstackCloudLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 228, Name: "vcenter_owner_rule", LegacyName: "v_center_owner_rule", New: NewVcenterOwnerRuleResource, NewLegacy: NewVcenterOwnerRuleLegacyResource},
        {Order: 332, Name: "dsym_file", LegacyName: "d_sym", New: NewDsymFileResource, NewLegacy: NewDsymFileLegacyResource},
        {Order: 404, Name: "openstack_cloud", LegacyName: "open_stack_cloud", New: NewOpenstackCloudResource, NewLegacy: NewOpenstackCloudLegacyResource},
    })
}

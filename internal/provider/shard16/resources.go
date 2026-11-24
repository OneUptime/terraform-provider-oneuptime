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
        {Order: 66, New: NewSecurityEventDestinationResource},
        {Order: 69, New: NewSecurityTelemetrySourceResource},
        {Order: 103, New: NewIncidentResource},
        {Order: 117, New: NewIncidentNoteTemplateResource},
        {Order: 224, New: NewVcenterOwnerRuleResource},
        {Order: 233, New: NewFileResource},
        {Order: 243, New: NewStatusPageSubscriberNotificationTemplateResource},
        {Order: 261, New: NewRunbookUserOwnerResource},
        {Order: 266, New: NewRunnerResource},
        {Order: 314, New: NewCodeRepositoryResource},
        {Order: 319, New: NewAiAgentUserOwnerResource},
        {Order: 328, New: NewDsymFileResource},
        {Order: 349, New: NewKubernetesClusterFeedResource},
        {Order: 383, New: NewNetworkSiteUserOwnerResource},
        {Order: 390, New: NewNetworkTopologySuppressionResource},
        {Order: 400, New: NewOpenstackCloudResource},
        {Order: 434, New: NewCloudResourceTeamOwnerResource},
        {Order: 436, New: NewCloudResourceLabelRuleResource},
        {Order: 463, New: NewVcenterOwnerRuleLegacyResource},
        {Order: 469, New: NewDsymFileLegacyResource},
        {Order: 472, New: NewOpenstackCloudLegacyResource},
    }, []LegacyResourceAlias{
        {Order: 224, Name: "vcenter_owner_rule", LegacyName: "v_center_owner_rule", New: NewVcenterOwnerRuleResource, NewLegacy: NewVcenterOwnerRuleLegacyResource},
        {Order: 328, Name: "dsym_file", LegacyName: "d_sym", New: NewDsymFileResource, NewLegacy: NewDsymFileLegacyResource},
        {Order: 400, Name: "openstack_cloud", LegacyName: "open_stack_cloud", New: NewOpenstackCloudResource, NewLegacy: NewOpenstackCloudLegacyResource},
    })
}

package shard19

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard19", []OrderedDataSource{
        {Order: 20, New: NewLogDropFilterDataSource},
        {Order: 24, New: NewWatchlistValueDataSource},
        {Order: 37, New: NewSecurityCaseEvidenceDataSource},
        {Order: 98, New: NewAuditLogExportDataSource},
        {Order: 117, New: NewOnCallPolicyDataSource},
        {Order: 235, New: NewMonitorOwnerRuleDataSource},
        {Order: 250, New: NewProxmoxClusterOwnerRuleDataSource},
        {Order: 274, New: NewStatusPageDomainDataSource},
        {Order: 283, New: NewStatusPageHeaderLinkDataSource},
        {Order: 292, New: NewWorkflowDataSource},
        {Order: 392, New: NewInventoryItemRelationshipDataSource},
        {Order: 395, New: NewDashboardDomainDataSource},
        {Order: 475, New: NewIotDeviceCredentialDataSource},
        {Order: 515, New: NewQueueLabelRuleDataSource},
        {Order: 532, New: NewProfileDataSource},
        {Order: 560, New: NewIotDeviceCredentialLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 475, Name: "iot_device_credential", LegacyName: "io_t_device_credential", New: NewIotDeviceCredentialDataSource, NewLegacy: NewIotDeviceCredentialLegacyDataSource},
    })
}

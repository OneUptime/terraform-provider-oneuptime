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
        {Order: 34, New: NewSecurityCaseEvidenceDataSource},
        {Order: 88, New: NewAuditLogExportDataSource},
        {Order: 107, New: NewOnCallPolicyDataSource},
        {Order: 225, New: NewMonitorOwnerRuleDataSource},
        {Order: 240, New: NewProxmoxClusterOwnerRuleDataSource},
        {Order: 264, New: NewStatusPageDomainDataSource},
        {Order: 273, New: NewStatusPageHeaderLinkDataSource},
        {Order: 282, New: NewWorkflowDataSource},
        {Order: 382, New: NewInventoryItemRelationshipDataSource},
        {Order: 385, New: NewDashboardDomainDataSource},
        {Order: 465, New: NewIotDeviceCredentialDataSource},
        {Order: 505, New: NewQueueLabelRuleDataSource},
        {Order: 522, New: NewProfileDataSource},
        {Order: 550, New: NewIotDeviceCredentialLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 465, Name: "iot_device_credential", LegacyName: "io_t_device_credential", New: NewIotDeviceCredentialDataSource, NewLegacy: NewIotDeviceCredentialLegacyDataSource},
    })
}

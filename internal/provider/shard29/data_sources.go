package shard29

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard29", []OrderedDataSource{
        {Order: 95, New: NewEntityRiskDataSource},
        {Order: 103, New: NewLlmCostBudgetDataSource},
        {Order: 134, New: NewMonitorCustomFieldDataSource},
        {Order: 192, New: NewIncidentAlertDataSource},
        {Order: 279, New: NewStatusPageSubscriberDataSource},
        {Order: 424, New: NewOnCallTimeLogDataSource},
        {Order: 481, New: NewStorageArrayUserOwnerDataSource},
        {Order: 492, New: NewServerlessFunctionUserOwnerDataSource},
        {Order: 510, New: NewDatabaseLabelRuleDataSource},
    }, nil)
}

package events

// ============================================================================
// COMPLETE BY ID
// ============================================================================
type TableFuzzySearch struct{ Match string }

func DecodeTableFuzzySearch(e *Event) *TableFuzzySearch { return e.Data.(*TableFuzzySearch) }
func NewTableFuzzySearch(match string) *Event {
	return &Event{
		Type: EventTableFuzzySearch,
		Data: &TableFuzzySearch{Match: match},
	}
}

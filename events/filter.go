package events

// ============================================================================
// COMPLETE BY ID
// ============================================================================

// TableFuzzySearch Is an event that is used to search a table by a fuzzy search
type TableFuzzySearch struct{ Match string }

// DecodeTableFuzzySearch will decode the event to search a table by a fuzzy search
func DecodeTableFuzzySearch(e *Event) *TableFuzzySearch { return e.Data.(*TableFuzzySearch) }

// NewTableFuzzySearch will create a new event to search a table by a fuzzy search
func NewTableFuzzySearch(match string) *Event {
	return &Event{
		Type: EventTableFuzzySearch,
		Data: &TableFuzzySearch{Match: match},
	}
}

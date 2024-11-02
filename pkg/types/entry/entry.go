package entry

type EntryArray struct {
	EntryType    int8
	EntryHash    string
	EntryName    string
	EntryContent []byte
	EntriesUnder []EntryArray
}

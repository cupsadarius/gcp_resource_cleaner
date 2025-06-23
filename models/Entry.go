package models

type EntryType int

const (
	EntryTypeProject EntryType = iota
	EntryTypeFolder
)

type Entry struct {
	Type EntryType
	Id   string
}

var EntryTypes = map[EntryType]string{
	EntryTypeProject: "project",
	EntryTypeFolder:  "folder",
}

func NewEntry(id string, entryType EntryType) *Entry {
	return &Entry{
		Type: entryType,
		Id:   id,
	}
}

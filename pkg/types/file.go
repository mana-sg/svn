package types

type FileNode struct {
	Id       int8       `json:"id"`
	Name     string     `json:"name"`
	Type     int8       `json:"type"`
	Children []FileNode `json:"children"`
	Content  string     `json:"content"`
}

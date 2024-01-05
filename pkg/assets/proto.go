package assets

type ProtoSet[T any] struct {
	Type      string
	TableName string
	Signature string
	DataArray []T
}

type Proto struct {
	Name string
	ID   int32
	SID  string
}

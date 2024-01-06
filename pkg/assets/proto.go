package assets

type ProtoSet[T any] struct {
	Type      string
	TableName string
	Signature string
	DataArray []Proto[T]
}

type Proto[T any] struct {
	Name  string
	ID    int32
	SID   string
	Proto T
}

func (p *ProtoSet[T]) Map() map[int32]Proto[T] {
	ret := make(map[int32]Proto[T])
	for _, v := range p.DataArray {
		ret[v.ID] = v
	}
	return ret
}

package blueprint

func Decode(data string) (*Blueprint, error) {
	b := &Blueprint{}
	return b.Decode(data)
}

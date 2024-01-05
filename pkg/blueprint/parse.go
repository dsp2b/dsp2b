package blueprint

func Decode(data string) (Blueprint, error) {
	b := Blueprint{}
	if err := b.Decode(data); err != nil {
		return Blueprint{}, err
	}
	return b, nil
}

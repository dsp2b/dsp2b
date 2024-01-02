package blueprint

import (
	"testing"
)

func TestDecode(t *testing.T) {
	b, err := Decode("BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,%E6%96%B0%E8%93%9D%E5%9B%BE,\"H4sIAAAAAAAAC2NkQAWMUAxh/2dgOAFlMsKFEWoPSG7Dxv7HYcfwHwpQTWZgAAB4dngncAAAAA==\"2881F7A76BAF3A19C17C948A5C773D72")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", b)
}

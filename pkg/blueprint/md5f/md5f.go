package md5f

import (
	"encoding/binary"
)

type MD5F struct {
	A, B, C, D uint32
}

func (m *MD5F) F(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

func (m *MD5F) G(x, y, z uint32) uint32 {
	return (x & z) | (y & ^z)
}

func (m *MD5F) H(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func (m *MD5F) I(x, y, z uint32) uint32 {
	return y ^ (x | ^z)
}

func (m *MD5F) FF(a *uint32, b, c, d, mj, s uint32, ti uint32) {
	*a = *a + m.F(b, c, d) + mj + ti
	*a = *a<<s | *a>>(32-s)
	*a += b
}

func (m *MD5F) GG(a *uint32, b, c, d, mj, s uint32, ti uint32) {
	*a = *a + m.G(b, c, d) + mj + ti
	*a = *a<<s | *a>>(32-s)
	*a += b
}

func (m *MD5F) HH(a *uint32, b, c, d, mj, s uint32, ti uint32) {
	*a = *a + m.H(b, c, d) + mj + ti
	*a = *a<<s | *a>>(32-s)
	*a += b
}

func (m *MD5F) II(a *uint32, b, c, d, mj, s uint32, ti uint32) {
	*a = *a + m.I(b, c, d) + mj + ti
	*a = *a<<s | *a>>(32-s)
	*a += b
}

func (m *MD5F) MD5_Init() {
	m.A = 1732584193
	m.B = 4024216457
	m.C = 2562383102
	m.D = 271734598
}

func (m *MD5F) MD5_Append(input []byte) []uint32 {
	num := 1
	num2 := len(input)
	num3 := num2 % 64
	var num4, num5 int

	if num3 < 56 {
		num4 = 55 - num3
		num5 = num2 - num3 + 64
	} else if num3 == 56 {
		num4 = 63
		num = 1
		num5 = num2 + 8 + 64
	} else {
		num4 = 63 - num3 + 56
		num5 = num2 + 64 - num3 + 64
	}

	var arrayList []byte
	arrayList = append(arrayList, input...)

	if num == 1 {
		arrayList = append(arrayList, 128)
	}

	for i := 0; i < num4; i++ {
		arrayList = append(arrayList, 0)
	}

	num6 := int64(num2) * 8
	b := byte(num6 & 255)
	b2 := byte((num6 >> 8) & 255)
	b3 := byte((num6 >> 16) & 255)
	b4 := byte((num6 >> 24) & 255)
	b5 := byte((num6 >> 32) & 255)
	b6 := byte((num6 >> 40) & 255)
	b7 := byte((num6 >> 48) & 255)
	b8 := byte(num6 >> 56)

	arrayList = append(arrayList, b, b2, b3, b4, b5, b6, b7, b8)

	array := make([]byte, len(arrayList))
	copy(array, arrayList)

	array2 := make([]uint32, num5/4)

	var num7, num8 int64
	num7 = 0
	num8 = 0

	for num7 < int64(num5) {
		array2[num8] = binary.LittleEndian.Uint32(array[int(num7):])
		num8++
		num7 += 4
	}

	return array2
}

func (m *MD5F) MD5_Trasform(x []uint32) []uint32 {
	for i := 0; i < len(x); i += 16 {
		a := m.A
		b := m.B
		c := m.C
		d := m.D
		m.FF(&a, b, c, d, x[i], 7, 3614090360)
		m.FF(&d, a, b, c, x[i+1], 12, 3906451286)
		m.FF(&c, d, a, b, x[i+2], 17, 606105819)
		m.FF(&b, c, d, a, x[i+3], 22, 3250441966)
		m.FF(&a, b, c, d, x[i+4], 7, 4118548399)
		m.FF(&d, a, b, c, x[i+5], 12, 1200080426)
		m.FF(&c, d, a, b, x[i+6], 17, 2821735971)
		m.FF(&b, c, d, a, x[i+7], 22, 4249261313)
		m.FF(&a, b, c, d, x[i+8], 7, 1770035416)
		m.FF(&d, a, b, c, x[i+9], 12, 2336552879)
		m.FF(&c, d, a, b, x[i+10], 17, 4294925233)
		m.FF(&b, c, d, a, x[i+11], 22, 2304563134)
		m.FF(&a, b, c, d, x[i+12], 7, 1805586722)
		m.FF(&d, a, b, c, x[i+13], 12, 4254626195)
		m.FF(&c, d, a, b, x[i+14], 17, 2792965006)
		m.FF(&b, c, d, a, x[i+15], 22, 968099873)

		m.GG(&a, b, c, d, x[i+1], 5, 4129170786)
		m.GG(&d, a, b, c, x[i+6], 9, 3225465664)
		m.GG(&c, d, a, b, x[i+11], 14, 643717713)
		m.GG(&b, c, d, a, x[i], 20, 3384199082)
		m.GG(&a, b, c, d, x[i+5], 5, 3593408605)
		m.GG(&d, a, b, c, x[i+10], 9, 38024275)
		m.GG(&c, d, a, b, x[i+15], 14, 3634488961)
		m.GG(&b, c, d, a, x[i+4], 20, 3889429448)
		m.GG(&a, b, c, d, x[i+9], 5, 569495014)
		m.GG(&d, a, b, c, x[i+14], 9, 3275163606)
		m.GG(&c, d, a, b, x[i+3], 14, 4107603335)
		m.GG(&b, c, d, a, x[i+8], 20, 1197085933)
		m.GG(&a, b, c, d, x[i+13], 5, 2850285829)
		m.GG(&d, a, b, c, x[i+2], 9, 4243563512)
		m.GG(&c, d, a, b, x[i+7], 14, 1735328473)
		m.GG(&b, c, d, a, x[i+12], 20, 2368359562)

		m.HH(&a, b, c, d, x[i+5], 4, 4294588738)
		m.HH(&d, a, b, c, x[i+8], 11, 2272392833)
		m.HH(&c, d, a, b, x[i+11], 16, 1839030562)
		m.HH(&b, c, d, a, x[i+14], 23, 4259657740)
		m.HH(&a, b, c, d, x[i+1], 4, 2763975236)
		m.HH(&d, a, b, c, x[i+4], 11, 1272893353)
		m.HH(&c, d, a, b, x[i+7], 16, 4139469664)
		m.HH(&b, c, d, a, x[i+10], 23, 3200236656)
		m.HH(&a, b, c, d, x[i+13], 4, 681279174)
		m.HH(&d, a, b, c, x[i], 11, 3936430074)
		m.HH(&c, d, a, b, x[i+3], 16, 3572445317)
		m.HH(&b, c, d, a, x[i+6], 23, 76029189)
		m.HH(&a, b, c, d, x[i+9], 4, 3654602809)
		m.HH(&d, a, b, c, x[i+12], 11, 3873151461)
		m.HH(&c, d, a, b, x[i+15], 16, 530742520)
		m.HH(&b, c, d, a, x[i+2], 23, 3299628645)

		m.II(&a, b, c, d, x[i], 6, 4096336452)
		m.II(&d, a, b, c, x[i+7], 10, 1126891415)
		m.II(&c, d, a, b, x[i+14], 15, 2878612391)
		m.II(&b, c, d, a, x[i+5], 21, 4237533241)
		m.II(&a, b, c, d, x[i+12], 6, 1700485571)
		m.II(&d, a, b, c, x[i+3], 10, 2399980690)
		m.II(&c, d, a, b, x[i+10], 15, 4293915773)
		m.II(&b, c, d, a, x[i+1], 21, 2240044497)
		m.II(&a, b, c, d, x[i+8], 6, 1873313359)
		m.II(&d, a, b, c, x[i+15], 10, 4264355552)
		m.II(&c, d, a, b, x[i+6], 15, 2734768916)
		m.II(&b, c, d, a, x[i+13], 21, 1309151649)
		m.II(&a, b, c, d, x[i+4], 6, 4149444226)
		m.II(&d, a, b, c, x[i+11], 10, 3174756917)
		m.II(&c, d, a, b, x[i+2], 15, 718787259)
		m.II(&b, c, d, a, x[i+9], 21, 3951481745)

		m.A += a
		m.B += b
		m.C += c
		m.D += d
	}
	return []uint32{m.A, m.B, m.C, m.D}
}

func MD5Hash(input string) []byte {
	md5f := MD5F{}
	md5f.MD5_Init()
	data := []byte(input)
	x := md5f.MD5_Append(data)
	result := md5f.MD5_Trasform(x)

	output := make([]byte, 16)
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(output[i*4:], result[i])
	}

	return output
}

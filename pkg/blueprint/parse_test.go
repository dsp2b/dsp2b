package blueprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	str := "BLUEPRINT:0,10,1106,0,0,0,0,0,638423980300910000,0.9.26.13026,%E9%92%9B%E5%9D%97-60.000%2Fmin,%E9%92%9B%E5%9D%97-60.000%20%2Fmin\"H4sIAAAAAAAAA4WRMQ7CMAxFnaY0pRyDibmMCBJxA07DyEE4EAMMzCw9ABs7Ia7T2CVFWHL7E+X95MsKxqVik/YAlygVLPjUaRc+NtdU73oLPhbtuOSqGOpB4D/Weo/9MCuYxx1viDmUqhcFLtr2lSDS0gQADbAqbE8GT00GOt5qRQSh/0coGUoRBDyOoCcizHBxu54TRDqPoH9EqOhIJyCp6QV3s5QRhqGCYWDILDXDeAvCSsA1A9bmmmGTwUV4PsAxqGG0TmR23/nRpJl4QcNAGpnUCc6zF2ECAJugPqLB+pD2AgAA\"25AB8B80EEF56D21AFDF13D2C7E89756"
	b, err := Decode(str)
	assert.NoError(t, err)
	actual, err := b.Encode()
	assert.NoError(t, err)
	t.Log(actual)
	_, err = Decode(actual)
	assert.NoError(t, err)

}

func TestRename(t *testing.T) {
	replace, err := Rename("BLUEPRINT:0,10,1106,0,0,0,0,0,638423980300910000,0.9.26.13026,%E9%92%9B%E5%9D%97-60.000%2Fmin,%E9%92%9B%E5%9D%97-60.000%20%2Fmin\"H4sIAAAAAAAAA4WRMQ7CMAxFnaY0pRyDibmMCBJxA07DyEE4EAMMzCw9ABs7Ia7T2CVFWHL7E+X95MsKxqVik/YAlygVLPjUaRc+NtdU73oLPhbtuOSqGOpB4D/Weo/9MCuYxx1viDmUqhcFLtr2lSDS0gQADbAqbE8GT00GOt5qRQSh/0coGUoRBDyOoCcizHBxu54TRDqPoH9EqOhIJyCp6QV3s5QRhqGCYWDILDXDeAvCSsA1A9bmmmGTwUV4PsAxqGG0TmR23/nRpJl4QcNAGpnUCc6zF2ECAJugPqLB+pD2AgAA\"25AB8B80EEF56D21AFDF13D2C7E89756",
		"替换名字")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(replace)
}

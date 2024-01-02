package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func main() {
	b, err := os.ReadFile("RecipeProtoSet.dat")
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewReader(b)
	var l int32
	var num18 int32
	binary.Read(r, binary.LittleEndian, &l)
	binary.Read(r, binary.LittleEndian, &num18)
	for m := int32(0); m < num18; m++ {
		binary.Read(r, binary.LittleEndian, &l)
		var num19 int32
		binary.Read(r, binary.LittleEndian, &num19)
		var num20 int32
		if num19 < 12000 && num19 > 0 {
			log.Printf("ignore")
		} else {
			num20 = 0
		}
		log.Printf("%d", num19)
		for n := 0; n < 40; n++ {
			var num21 float32
			binary.Read(r, binary.LittleEndian, &num21)
			if num20 > 0 {
				log.Printf("%f", num21)
			}
		}
	}
	os.Exit(0)
}

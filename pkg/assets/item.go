package assets

import (
	"encoding/binary"
	"os"
)

type ItemProtoSet struct {
	ProtoSet[ItemProto]
}

func (i *ItemProtoSet) Load(datFile string) error {
	f, err := os.Open(datFile)
	if err != nil {
		return err
	}
	var flag [20]byte
	_ = binary.Read(f, binary.LittleEndian, &flag)
	if flag != [20]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 0, 0,
	} {
		return err
	}
	ReadInt64(f)
	i.ProtoSet.Type = ReadString(f)
	i.TableName = ReadString(f)
	i.Signature = ReadString(f)
	i01 := ReadInt32(f)
	i.DataArray = make([]ItemProto, i01)
	for n := 0; int32(n) < i01; n++ {
		_ = ReadStruct(f, &i.DataArray[n])
	}
	return nil
}

type ItemProto struct {
	Proto
	Type            EItemType
	SubID           int32
	MiningFrom      string
	ProduceFrom     string
	StackSize       int32
	Grade           int32
	Upgrades        []int32
	IsFluid         bool
	IsEntity        bool
	CanBuild        bool
	BuildInGas      bool
	IconPath        string //15
	ModelIndex      int32
	ModelCount      int32
	HpMax           int32
	Ability         int32
	HeatValue       int64
	Potential       int64
	ReactorInc      float32
	FuelType        int32
	AmmoType        EAmmoType
	BombType        int32 //25
	CraftType       int32
	BuildIndex      int32
	BuildMode       int32
	GridIndex       int32
	UnlockKey       int32
	PreTechOverride int32
	Productive      bool
	MechaMaterialID int32
	DropRate        float32
	EnemyDropLevel  int32 //35
	EnemyDropRange  Vector2
	EnemyDropCount  float32
	EnemyDropMask   int32
	DescFields      []int32
	Description     string
}

type Vector2 struct {
	X float32
	Y float32
}

type EItemType int32

const (
	Unknown EItemType = iota
	Resource
	Material
	Component
	Product
	Logistics
	Production
	Decoration
	Turret
	Defense
	DarkFog
	Matrix
)

type EAmmoType int32

const (
	EAmmoType_None EAmmoType = iota
	Bullet
	Laser
	Cannon
	Plasma
	Missile
)

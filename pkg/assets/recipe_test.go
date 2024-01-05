package assets

import "testing"

func TestRecipeProtoSet_Load(t *testing.T) {
	//r := &RecipeProtoSet{}
	//r.Load("RecipeProtoSet.dat")
	//t.Logf("%v", r)

	item := &ItemProtoSet{}
	item.Load("ItemProtoSet.dat")
	t.Logf("%v", item)
}

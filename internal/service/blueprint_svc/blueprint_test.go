package blueprint_svc

import (
	"context"
	"testing"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/configs/memory"
	"github.com/codfrm/cago/database/cache"
	api "github.com/dsp2b/dsp2b-go/internal/api/blueprint"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cfg, err := configs.NewConfig("test", configs.WithSource(memory.NewSource(map[string]interface{}{
		"cache": map[string]interface{}{
			"type": "memory",
		},
	})))
	if err != nil {
		panic(err)
	}
	err = cache.Cache().Start(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	err = InitBlueprint("../../../data/itemProtoSet.json",
		"../../../data/recipeProtoSet.json")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func Test_blueprintSvc_Parse(t *testing.T) {
	resp, err := Blueprint().Parse(context.Background(), &api.ParseRequest{
		Blueprint: "BLUEPRINT:0,30,1101,1104,1105,0,0,0,638417736554170000,0.9.26.13026,%E9%93%81%E5%9D%97-120.000%2Fmin,%E9%93%81%E5%9D%97-120.000%20%2Fmin%0A%E9%93%9C%E5%9D%97-120.000%20%2Fmin%0A%E9%AB%98%E7%BA%AF%E7%A1%85%E5%9D%97-60.000%20%2Fmin\"H4sIAAAAAAAAA43WTU8TURQG4NuC1ALSAQXlSwvl+3OgfBlNuBPXJuUf9Cd0J0t+iIv+lIkrF7owLo1Jt8LSLUmd4dzb9+B9m8xNyn1p+sw952Qybck8XiX3ktw35quLJfMWn7q5zP7YMMuaqH4xfbf+v2oJ6AEa7Pk6+ZC/fld2zSv3Tr8i5uNo6SGU83+azb8DJFlfxJj8AvmqZq+xvlzgz4hcYMSdalULKg9rYcRNwZhRoEELCqOFBcd8Cy3XwpP8n+/fPg+Q5LCFHNdUC7euhTH5SE+VrfOwFt4ZP/8KkD+1F7Twq7JtltxpvoUr18JTKftetXBPW8jxrGrhzrVQlY/UEyCdpYIflYZuwd+XZhzA96wzcH5KjksKTwBYG2bgcYInVak2zMATBD8D6NowA08GuJzdvsa0szRVpPdq9XowuKlabcpXUJMtVpOOC089AvCnxXTqNdL7NIC1YQaOCJ5RpdowA08T/Byga8MMPEOmfptN/TpLL4r0Pmzqs7JZNWlbeOpzAP40S6c+S3p/CWBtmIHnSO93rnf3JG6pknWWlV9knlQwD+BL1hk47L2cPfOlggV5o61ObdMKFkkFiwD+1HbhClqugiV5o6NO7dAKlkkFywD+1E7hCq6yCt5n6bW8kT+m/aNaZ1nD7sA3AFESZlk/K+umHpRfNp/cF14d0CRhxkVWyAxWipSf41WCVwFSG2bgBsENgK4NM/AawWuypQqkFK8TvA6Q2jADbxC8AfAwpMvHGXiT4E1g425TnYG3CN4CjhSOArxN8LZsPQV6FO8QvOP2RIGE4V2Cd4GNwmHPewTvAQ9+FCWPfyAJ3id4HzhVOA3wAcEHGJgHOgMfEnwoW6pASnFMcAzcVbgb4COCj4BvFL4J8DHBx+jZgx7FTYKbmPYAJAyfEHwiW6RARPEpwafAVmEb4DOCz2SrK1Cn+Jzgc9liBWKKLwi+kM0qYCn23xx9hf8BlMxa5yAOAAA=\"64627B4DFD142D3E53368FE51F673B06",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

}

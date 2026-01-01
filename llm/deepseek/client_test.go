package deepseek

import (
	"context"
	"fmt"
	"testing"

	"github.com/neura-flow/common/llm/types"
	"github.com/neura-flow/common/util"
)

const prompt = `
# 你的角色：一个可以从文本提取数据的助手

# 示例输入:
张三是北京科技大学计算机科学专业的毕业生，今年25岁

# 示例提取字段和格式:
[
  { "description": "姓名", "type": "string", "var_name": "name"},
  { "description": "年龄", "type": "int", "var_name": "age"},
  { "description": "学历", "type": "string", "var_name": "qualification"},
  { "description": "毕业学校", "type": "string", "var_name": "school"},
  { "description": "专业", "type": "string", "var_name": "major"}
]

# 示例输出:
{
  "name": "张三",
  "age": 25,
  "qualification": "",
  "major": "计算机科学专业",
  "school": "北京科技大学"
}


# 输入文本:
A beautiful 3200 square foot house in Santa Barbara, California for rent. $5,300 a month.
Featuring a spacious master bedroom with en-suite bathroom at the ground floor, 2 kids rooms and one more bathroom on the 2nd floor.
The living room is perfect for entertaining guests or relaxing with family, with plenty of natural light and a cozy fireplace.
The house also includes a fully equipped kitchen with modern appliances, laundry facilities, a parking for two cars, and a large backyard for outdoor activities.
This is the perfect place to call home.

# 提取字段和格式如下:
[
  {
    "description": "price of property", "type": "float",
    "var_name": "price"
  },
  {
    "description": "price currency, format: ISO currency code",
    "type": "string",
    "var_name": "price_currency"
  },
  {
    "description": "US state, 2 letters abbreviation",
    "type": "string",
    "var_name": "state"
  },
  {
    "description": "number of bedrooms",
    "type": "float",
    "var_name": "num_bedrooms"
  },
  {
    "description": "number of bathrooms",
    "type": "float",
    "var_name": "num_bathrooms"
  },
  {
    "description": "amenities in property",
    "type": "array[string]",
    "var_name": "amenities",
    "valid_values": [
      "pool",
      "parking",
      "balcony",
      "backyard",
      "elevator"
    ]
  },
  {
    "description": "property area",
    "type": "integer",
    "var_name": "property_area"
  },
  {
    "description": "property area units",
    "type": "string",
    "var_name": "property_area_units",
    "valid_values": [
      "sqm",
      "sqft"
    ]
  },
  {
    "description": "city where property is located",
    "type": "string",
    "var_name": "city"
  },
  {
    "description": "number of parking spots",
    "type": "integer",
    "var_name": "num_parkings"
  },
  {
    "description": "are pets allowed (true) not allowed (false), or unknown (null)",
    "type": "boolean",
    "var_name": "pets_allowed"
  }
]

请按照指定格式提取数据，并以 json 格式返回

`

func TestCall(t *testing.T) {
	client := New(&Config{
		ServerUrl: "https://api.deepseek.com/chat/completions",
		Secret:    "sk-xxx",
	})

	completion := types.NewDeepSeekCompletionRequest("deepseek-reasoner")
	completion.Messages = []any{
		types.SystemMessage{
			Role: "system", Content: "You are a helpful assistant.",
		},
		types.UserMessage{
			Role: "user", Content: prompt,
		},
	}

	resp, err := client.Call(context.TODO(), completion)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response:")
	fmt.Printf("%s\n", util.ToJson(resp))
}

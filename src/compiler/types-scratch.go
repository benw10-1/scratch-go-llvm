package compiler

import (
	"encoding/json"
	"scratch-llvm/globalutil"
)

type Project struct {
	Targets  []Target  `json:"targets"`
	Moniters []Monitor `json:"monitors"`
	// Extensions []Extension `json:"extensions"`
	// MetaData MetaData `json:"metaData"`
}

type MonitorMode string

const (
	ModeDefault MonitorMode = "default"
	ModeLarge   MonitorMode = "large"
	ModeSlider  MonitorMode = "slider"
	ModeList    MonitorMode = "list"
)

type Monitor struct {
	ID         string                 `json:"id"`
	Mode       MonitorMode            `json:"mode"`
	Opcode     string                 `json:"opcode"`
	Params     map[string]interface{} `json:"params"`
	SpriteName string                 `json:"spriteName"`
	X          int                    `json:"x"`
	Y          int                    `json:"y"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Visible    bool                   `json:"visible"`
	Value      interface{}
	SliderMin  float64 `json:"sliderMin"`
	SliderMax  float64 `json:"sliderMax"`
	IsDiscrete bool    `json:"isDiscrete"`
}

type Target struct {
	Name                string              `json:"name"`
	IsStage             bool                `json:"isStage"`
	Variables           map[string]Variable `json:"variables"`
	Blocks              map[string]Block    `json:"blocks"`
	Lists               map[string]List     `json:"lists"`
	Broadcasts          map[string]string   `json:"broadcasts"`
	Comments            map[string]string   `json:"comments"`
	CurrentCostumeIndex int                 `json:"currentCostume"`
	Costumes            []Costume           `json:"costumes"`
	Sounds              []Sound             `json:"sounds"`
	Volume              float64             `json:"volume"`
	LayerOrder          int                 `json:"layerOrder"`
}

type Costume struct {
	BitmapResolution int `json:"bitmapResolution"`
	RotationCenterX  int `json:"rotationCenterX"`
	RotationCenterY  int `json:"rotationCenterY"`
}

type Sound struct {
	Rate        float64 `json:"rate"`
	SampleCount int     `json:"sampleCount"`
}

type List []interface{}

func (l List) Name() string {
	if len(l) == 0 {
		return ""
	}

	return l[0].(string)
}

func (l List) List() []interface{} {
	if len(l) < 2 {
		return nil
	}

	return l[1].([]interface{})
}

type Variable []interface{}

func (v Variable) Name() string {
	return v[0].(string)
}

func (v Variable) ID() string {
	return v[1].(string)
}

func (v Variable) IsCloud() bool {
	return v[2].(bool)
}

type BlockTypeInt int

const (
	BlockTypeNumber BlockTypeInt = (iota + 4)
	BlockTypePostiveNumber
	BlockTypePostiveInteger
	BlockTypeInteger
	BlockTypeAngle
	BlockTypeColor
	BlockTypeString
	BlockTypeBroadcast
	BlockTypeVariable
	BlockTypeList
)

type BlockTypeStringEnum string

const (
	BlockTypeArray  BlockTypeStringEnum = "array"
	BlockTypeObject BlockTypeStringEnum = "object"
)

type Block struct {
	underlying interface{}
	blockType  BlockTypeStringEnum
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.underlying)
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var blockObj BlockObj
	err := json.Unmarshal(data, &blockObj)
	if err == nil {
		b.underlying = blockObj
		b.blockType = BlockTypeObject
		return nil
	}

	var blockArr BlockArray
	err = json.Unmarshal(data, &blockArr)
	if err == nil {
		b.underlying = blockArr
		b.blockType = BlockTypeArray
		return nil
	}

	return globalutil.Errorf(err, "Failed to unmarshal block: %s", string(data))
}

// https://en.scratch-wiki.info/wiki/Scratch_File_Format#Blocks
type BlockArray []interface{}

func (b BlockArray) BlockType() int {
	if len(b) < 1 {
		return 0
	}

	return b[0].(int)
}

func (b BlockArray) IsVariable() bool {
	if len(b) < 1 {
		return false
	}

	return b.BlockType() == int(BlockTypeVariable)
}

func (b BlockArray) Name() string {

	switch b.BlockType() {
	case int(BlockTypeVariable):
		if len(b) < 2 {
			return ""
		}

		return b[1].(string)

	case int(BlockTypeList):
		if len(b) < 2 {
			return ""
		}

		return b[1].(string)
	default:
		return ""
	}
}

func (b BlockArray) ID() interface{} {

	switch b.BlockType() {
	case int(BlockTypeVariable):
		if len(b) < 3 {
			return ""
		}

		return b[2]

	case int(BlockTypeList):
		if len(b) < 3 {
			return ""
		}

		return b[2]
	default:
		return ""
	}
}

func (b BlockArray) IsList() bool {
	if len(b) < 1 {
		return false
	}

	return b.BlockType() == int(BlockTypeList)
}

type BlockObj struct {
	Opcode   string                `json:"opcode"`
	Next     string                `json:"next"`
	Parent   string                `json:"parent,omitempty"`
	Inputs   map[string]BlockInput `json:"inputs"`
	Fields   map[string]BlockField `json:"fields"`
	Shadow   bool                  `json:"shadow"`
	TopLevel bool                  `json:"topLevel"`
}

type BlockField []interface{}

func (b BlockField) Value() interface{} {
	if len(b) < 1 {
		return nil
	}

	return b[0]
}

type BlockInput []interface{}

func (b BlockInput) ShadowStatus() int {
	if len(b) < 1 {
		return 0
	}

	return b[0].(int)
}

func (b BlockInput) BlockID() string {
	if len(b) < 2 {
		return ""
	}

	return b[1].(string)
}

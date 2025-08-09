package ao

import (
	goarTypes "github.com/everFinance/goar/types"
)


type Message struct {
	Id string `json:"Id"`
	Target string `json:"Target"`
	Anchor string `json:"Anchor"`
	Data string `json:"Data"`
	Tags []goarTypes.Tag `json:"Tags,omitempty"`
}

type Result struct {
	Messages []Message `json:"Messages"`
	Assignments []any `json:"Assignments"`
	Spawns []any `json:"Spawns"`
	Output map[string]any `json:"Output"`
	Error any `json:"Error"`
	GasUsed any `json:"GasUsed"`
}

type DryRunInput struct {
	Id string `json:"Id"`
	Owner string `json:"Owner"`
	From string `json:"From"`
	Anchor string `json:"Anchor"`
	Data string `json:"Data"`
	Tags []goarTypes.Tag `json:"Tags,omitempty"`
}

type WriteInput struct {
	Process string 
	Anchor string
	Data string 
	Tags []goarTypes.Tag
}

type SpawnInput struct {
	Module string
	Authority string
	Scheduler string
	Tags []goarTypes.Tag 
	Data any // can be string, number, bytes 
	Target string
	Anchor string
}

type AOClient interface {
	Read(input DryRunInput) (result *Result, err error)
	Write(message WriteInput) (id string, result *Result, err error) // the id, result of the write operation, and an error if any
	Spawn(input SpawnInput) (id string, result *Result, err error) // the id, result of the spawn operation, and an error if any
}


package model

type Instruction struct {
	ID          string `json:"id"`
	StepNumber  int    `json:"stepNumber"`
	Description string `json:"description"`
}

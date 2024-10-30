package blocks

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type AnswerBlock struct {
	BaseBlock
	Prompt string `json:"prompt"`
	Answer string `json:"answer"`
	Fuzzy  bool   `json:"fuzzy"`
}

type answerBlockData struct {
	Attempts int      `json:"attempts"`
	Guesses  []string `json:"guesses"`
}

// Basic Attributes Getters

func (b *AnswerBlock) GetID() string         { return b.ID }
func (b *AnswerBlock) GetType() string       { return "answer" }
func (b *AnswerBlock) GetLocationID() string { return b.LocationID }
func (b *AnswerBlock) GetName() string       { return "Answer" }
func (b *AnswerBlock) GetDescription() string {
	return "Players must enter the correct answer to a prompt."
}
func (b *AnswerBlock) GetOrder() int  { return b.Order }
func (b *AnswerBlock) GetPoints() int { return b.Points }
func (b *AnswerBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-pencil-line"><path d="M12 20h9"/><path d="M16.376 3.622a1 1 0 0 1 3.002 3.002L7.368 18.635a2 2 0 0 1-.855.506l-2.872.838a.5.5 0 0 1-.62-.62l.838-2.872a2 2 0 0 1 .506-.854z"/><path d="m15 5 3 3"/></svg>`
}
func (b *AnswerBlock) GetAdminData() interface{} {
	return &b
}
func (b *AnswerBlock) GetData() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

// Data Operations

func (b *AnswerBlock) ParseData() error {
	return json.Unmarshal(b.Data, b)
}

func (b *AnswerBlock) UpdateBlockData(input map[string][]string) error {
	// Points
	if input["points"] != nil {
		points, err := strconv.Atoi(input["points"][0])
		if err != nil {
			return fmt.Errorf("points must be an integer")
		}
		b.Points = points
	}
	// Prompt and Answer
	if input["prompt"] == nil || input["answer"] == nil {
		return fmt.Errorf("prompt and answer are required fields")
	}
	b.Prompt = input["prompt"][0]
	b.Answer = input["answer"][0]
	if input["fuzzy"] != nil {
		b.Fuzzy = input["fuzzy"][0] == "on"
	}
	return nil
}

// Validation and Points Calculation

func (b *AnswerBlock) RequiresValidation() bool { return true }

func (b *AnswerBlock) ValidatePlayerInput(state PlayerState, input map[string][]string) (PlayerState, error) {
	if input["answer"] == nil {
		return state, fmt.Errorf("answer is a required field")
	}

	var err error
	newPlayerData := answerBlockData{}
	if state.GetPlayerData() != nil {
		json.Unmarshal(state.GetPlayerData(), &newPlayerData)
	}

	// Increment the number of attempts and save guesses
	newPlayerData.Attempts++
	newPlayerData.Guesses = append(newPlayerData.Guesses, input["answer"][0])

	if input["answer"][0] != b.Answer {
		// Incorrect answer, save player data and return an error
		playerData, err := json.Marshal(newPlayerData)
		if err != nil {
			return state, fmt.Errorf("Error saving player data")
		}
		state.SetPlayerData(playerData)
		return state, nil
	}

	// Correct answer, update state to complete
	playerData, err := json.Marshal(newPlayerData)
	if err != nil {
		return state, fmt.Errorf("Error saving player data")
	}
	state.SetPlayerData(playerData)
	state.SetComplete(true)
	state.SetPointsAwarded(b.Points)
	return state, nil
}

func (b *AnswerBlock) CalculatePoints(input map[string][]string) (int, error) {
	return b.Points, nil
}

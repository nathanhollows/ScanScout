package blocks

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PincodeBlock struct {
	BaseBlock
	Prompt  string `json:"prompt"`
	Pincode string `json:"pincode"`
}

type pincodeBlockData struct {
	Attempts int      `json:"attempts"`
	Guesses  []string `json:"guesses"`
}

// Basic Attributes Getters

func (b *PincodeBlock) GetID() string         { return b.ID }
func (b *PincodeBlock) GetType() string       { return "pincode" }
func (b *PincodeBlock) GetLocationID() string { return b.LocationID }
func (b *PincodeBlock) GetName() string       { return "Pincode" }
func (b *PincodeBlock) GetDescription() string {
	return "Players must enter the correct pincode to a prompt."
}
func (b *PincodeBlock) GetOrder() int  { return b.Order }
func (b *PincodeBlock) GetPoints() int { return b.Points }
func (b *PincodeBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-file-digit"><path d="M4 22h14a2 2 0 0 0 2-2V7l-5-5H6a2 2 0 0 0-2 2v4"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><rect width="4" height="6" x="2" y="12" rx="2"/><path d="M10 12h2v6"/><path d="M10 18h4"/></svg>`
}
func (b *PincodeBlock) GetAdminData() interface{} {
	return &b
}
func (b *PincodeBlock) GetData() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

// Data Operations

func (b *PincodeBlock) ParseData() error {
	return json.Unmarshal(b.Data, b)
}

func (b *PincodeBlock) UpdateBlockData(input map[string][]string) error {
	// Points
	if input["points"] != nil {
		points, err := strconv.Atoi(input["points"][0])
		if err != nil {
			return fmt.Errorf("points must be an integer")
		}
		b.Points = points
	}
	// Prompt and Pincode
	if input["prompt"] == nil || input["pincode"] == nil {
		return fmt.Errorf("prompt and pincode are required fields")
	}
	b.Prompt = input["prompt"][0]
	b.Pincode = input["pincode"][0]
	return nil
}

// Validation and Points Calculation

func (b *PincodeBlock) RequiresValidation() bool { return true }

func (b *PincodeBlock) ValidatePlayerInput(state PlayerState, input map[string][]string) (PlayerState, error) {
	if input["pincode"] == nil {
		return state, fmt.Errorf("pincode is a required field")
	}

	var err error
	newPlayerData := pincodeBlockData{}
	if state.GetPlayerData() != nil {
		json.Unmarshal(state.GetPlayerData(), &newPlayerData)
	}

	// Increment the number of attempts and save guesses
	newPlayerData.Attempts++
	newPlayerData.Guesses = append(newPlayerData.Guesses, input["pincode"][0])

	if input["pincode"][0] != b.Pincode {
		// Incorrect pincode, save player data and return an error
		playerData, err := json.Marshal(newPlayerData)
		if err != nil {
			return state, fmt.Errorf("Error saving player data")
		}
		state.SetPlayerData(playerData)
		return state, nil
	}

	// Correct pincode, update state to complete
	playerData, err := json.Marshal(newPlayerData)
	if err != nil {
		return state, fmt.Errorf("Error saving player data")
	}
	state.SetPlayerData(playerData)
	state.SetComplete(true)
	state.SetPointsAwarded(b.Points)
	return state, nil
}

func (b *PincodeBlock) CalculatePoints(input map[string][]string) (int, error) {
	return b.Points, nil
}

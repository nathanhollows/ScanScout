package blocks

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PasswordBlock struct {
	BaseBlock
	Content  string `json:"content"`
	Password string `json:"password"`
	Fuzzy    bool   `json:"fuzzy"`
}

type passwordBlockData struct {
	Attempts int      `json:"attempts"`
	Guesses  []string `json:"guesses"`
}

// Basic Attributes Getters

func (b *PasswordBlock) GetID() string         { return b.ID }
func (b *PasswordBlock) GetType() string       { return "password" }
func (b *PasswordBlock) GetLocationID() string { return b.LocationID }
func (b *PasswordBlock) GetName() string       { return "Password" }
func (b *PasswordBlock) GetDescription() string {
	return "Players must enter the correct password."
}
func (b *PasswordBlock) GetOrder() int  { return b.Order }
func (b *PasswordBlock) GetPoints() int { return b.Points }
func (b *PasswordBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-key-round"><path d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z"/><circle cx="16.5" cy="7.5" r=".5" fill="currentColor"/></svg>`
}
func (b *PasswordBlock) GetAdminData() interface{} {
	return &b
}
func (b *PasswordBlock) GetData() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

// Data Operations

func (b *PasswordBlock) ParseData() error {
	return json.Unmarshal(b.Data, b)
}

func (b *PasswordBlock) UpdateBlockData(input map[string][]string) error {
	// Points
	if input["points"] != nil {
		points, err := strconv.Atoi(input["points"][0])
		if err != nil {
			return fmt.Errorf("points must be an integer")
		}
		b.Points = points
	}
	// Content and Password
	if input["content"] == nil || input["block-passphrase"] == nil {
		return fmt.Errorf("content and block-passphrase are required fields")
	}
	b.Content = input["content"][0]
	b.Password = input["block-passphrase"][0]
	if input["fuzzy"] != nil {
		b.Fuzzy = input["fuzzy"][0] == "on"
	}
	return nil
}

// Validation and Points Calculation

func (b *PasswordBlock) RequiresValidation() bool { return true }

func (b *PasswordBlock) ValidatePlayerInput(state PlayerState, input map[string][]string) (PlayerState, error) {
	if input["password"] == nil {
		return state, fmt.Errorf("password is a required field")
	}

	var err error
	newPlayerData := passwordBlockData{}
	if state.GetPlayerData() != nil {
		json.Unmarshal(state.GetPlayerData(), &newPlayerData)
	}

	// Increment the number of attempts and save guesses
	newPlayerData.Attempts++
	newPlayerData.Guesses = append(newPlayerData.Guesses, input["password"][0])

	if input["password"][0] != b.Password {
		// Incorrect password, save player data and return an error
		playerData, err := json.Marshal(newPlayerData)
		if err != nil {
			return state, fmt.Errorf("Error saving player data")
		}
		state.SetPlayerData(playerData)
		return state, nil
	}

	// Correct password, update state to complete
	playerData, err := json.Marshal(newPlayerData)
	if err != nil {
		return state, fmt.Errorf("Error saving player data")
	}
	state.SetPlayerData(playerData)
	state.SetComplete(true)
	state.SetPointsAwarded(b.Points)
	return state, nil
}

func (b *PasswordBlock) CalculatePoints(input map[string][]string) (int, error) {
	return b.Points, nil
}

package blocks

import (
	"encoding/json"
)

type mockPlayerState struct {
	blockID       string
	playerID      string
	playerData    json.RawMessage
	isComplete    bool
	pointsAwarded int
}

func (m *mockPlayerState) GetBlockID() string {
	return m.blockID
}

func (m *mockPlayerState) GetPlayerID() string {
	return m.playerID
}

func (m *mockPlayerState) GetPlayerData() json.RawMessage {
	return m.playerData
}

func (m *mockPlayerState) SetPlayerData(data json.RawMessage) {
	m.playerData = data
}

func (m *mockPlayerState) IsComplete() bool {
	return m.isComplete
}

func (m *mockPlayerState) SetComplete(complete bool) {
	m.isComplete = complete
}

func (m *mockPlayerState) GetPointsAwarded() int {
	return m.pointsAwarded
}

func (m *mockPlayerState) SetPointsAwarded(points int) {
	m.pointsAwarded = points
}

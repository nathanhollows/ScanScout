package blocks

import (
	"encoding/json"
)

type PasswordBlock struct {
	BaseBlock
	Content  string `json:"content"`
	Password string `json:"password"`
	Fuzzy    bool   `json:"fuzzy"`
}

func (b *PasswordBlock) Validate(userID string, input map[string]string) error {
	// No validation required
	return nil
}

func (b *PasswordBlock) GetName() string { return "Password" }
func (b *PasswordBlock) GetDescription() string {
	return "Players must enter the correct password."
}
func (b *PasswordBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-key-round"><path d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z"/><circle cx="16.5" cy="7.5" r=".5" fill="currentColor"/></svg>`
}
func (b *PasswordBlock) GetType() string       { return "password" }
func (b *PasswordBlock) GetID() string         { return b.ID }
func (b *PasswordBlock) GetOrder() int         { return b.Order }
func (b *PasswordBlock) GetLocationID() string { return b.LocationID }
func (b *PasswordBlock) GetAdminData() interface{} {
	return &b
}
func (b *PasswordBlock) GetData() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

func (b *PasswordBlock) ParseData() error {
	return json.Unmarshal(b.Data, b)
}

func (b *PasswordBlock) UpdateData(data map[string]string) error {
	b.Content = data["content"]
	b.Password = data["block-passphrase"]
	b.Fuzzy = data["fuzzy"] == "on"
	return nil
}

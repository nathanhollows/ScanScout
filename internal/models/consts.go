package models

import "errors"

type NavigationMode int
type NavigationMethod int
type CompletionMethod int

type NavigationModes []NavigationMode
type NavigationMethods []NavigationMethod
type CompletionMethods []CompletionMethod

const (
	RandomNav NavigationMode = iota
	FreeRoamNav
	OrderedNav
)

const (
	ShowMap NavigationMethod = iota
	ShowMapAndNames
	ShowNames
	ShowClues
)

const (
	CheckInOnly CompletionMethod = iota
	CheckInAndOut
	SubmitContent
	Password
	ClickButton
)

// GetNavigationModes returns a list of navigation modes
func GetNavigationModes() NavigationModes {
	return []NavigationMode{RandomNav, FreeRoamNav, OrderedNav}
}

// GetNavigationMethods returns a list of navigation methods
func GetNavigationMethods() NavigationMethods {
	return []NavigationMethod{ShowMap, ShowMapAndNames, ShowNames, ShowClues}
}

// GetCompletionMethods returns a list of completion methods
func GetCompletionMethods() CompletionMethods {
	return []CompletionMethod{CheckInOnly, CheckInAndOut, SubmitContent, Password, ClickButton}
}

// String returns the string representation of the NavigationMode
func (n NavigationMode) String() string {
	return [...]string{"Random", "Free Roam", "Ordered"}[n]
}

// String returns the string representation of the NavigationMethod
func (n NavigationMethod) String() string {
	return [...]string{"Show Map", "Show Map and Names", "Show Location Names", "Show Clues"}[n]
}

// String returns the string representation of the CompletionMethod
func (c CompletionMethod) String() string {
	return [...]string{"Check In Only", "Check In And Out", "Submit Content", "Password", "Click Button"}[c]
}

// Description returns the description of the NavigationMode
func (n NavigationMode) Description() string {
	return [...]string{
		"The game will randomly select locations for players to visit. Good for large groups as it disperses players.",
		"Players can visit locations in any order. This mode shows all locations and is good for exploration.",
		"Players must visit locations in a specific order. Good for narrative experiences.",
	}[n]
}

// Description returns the description of the NavigationMethod
func (n NavigationMethod) Description() string {
	return [...]string{
		"Players are shown a map.",
		"Players are shown a map with location names.",
		"Players are shown a list of locations by name.",
		"Players are shown clues but not the location or name.",
	}[n]
}

// Description returns the description of the CompletionMethod
func (c CompletionMethod) Description() string {
	return [...]string{
		"Players must check in to a location but do not need to check out.",
		"Players must check in and out of a location.",
		"Players must submit content to a location, i.e., a photo or text.",
		"Players must enter a password to a location, i.e., a code or phrase.",
		"Players must click the correct button for the location, i.e., a quick quiz.",
	}[c]
}

// Parse NavigationMode
func ParseNavigationMode(s string) (NavigationMode, error) {
	switch s {
	case "Random":
		return RandomNav, nil
	case "Free Roam":
		return FreeRoamNav, nil
	case "Ordered":
		return OrderedNav, nil
	default:
		return 0, errors.New("invalid NavigationMode")
	}
}

// Parse NavigationMethod
func ParseNavigationMethod(s string) (NavigationMethod, error) {
	switch s {
	case "Show Map":
		return ShowMap, nil
	case "Show Map and Names":
		return ShowMapAndNames, nil
	case "Show Location Names":
		return ShowNames, nil
	case "Show Clues":
		return ShowClues, nil
	default:
		return ShowMap, errors.New("invalid NavigationMethod")
	}
}

// Parse CompletionMethod
func ParseCompletionMethod(s string) (CompletionMethod, error) {
	switch s {
	case "Check In Only":
		return CheckInOnly, nil
	case "Check In And Out":
		return CheckInAndOut, nil
	case "Submit Content":
		return SubmitContent, nil
	case "Password":
		return Password, nil
	case "Click Button":
		return ClickButton, nil
	default:
		return CheckInOnly, errors.New("invalid CompletionMethod")
	}
}

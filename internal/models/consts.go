package models

// NavigationMode represents the mode of navigation for an instance
type NavigationMode int

const (
	// FreeRoamShowAllNavigation is a navigation mode where users can scan in and out of locations in any order
	FreeRoamShowAllNavigation = iota
	// FreeRoamShowNNavigation is a navigation mode where users must scan in and out of locations in a random orde
	FreeRoamShowNNavigation
	// FreeRoamShowNoneNavigation is a navigation mode where users must scan in and out of locations in a random order
	FreeRoamShowNoneNavigation
	// OrderedNavigation is a navigation mode where users must scan in and out of locations in a specific order
	OrderedNavigation
)

func (n NavigationMode) String() string {
	return [...]string{"FreeRoamShowAllNavigation", "FreeRoamShowNNavigation", "FreeRoamShowNoneNavigation", "OrderedNavigation"}[n]
}

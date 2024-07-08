package models

type CompletionCriteria struct {
	baseModel

	ID   string `bun:",pk,type:varchar(36)" json:"id"`
	Name string `bun:",type:varchar(255)" json:"name"`
}

type CompletionCriterias []CompletionCriteria

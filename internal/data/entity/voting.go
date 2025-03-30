package entity

type VotingEntity struct {
	Id        string
	CreatorId string
	IsEnded   bool
	Title     string
	Options   map[string]int
	Users     []string
}

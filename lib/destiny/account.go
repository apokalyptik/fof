package destiny

type Account struct {
	id string
	p  *Platform
}

func (a *Account) AccountSummary() (*Request, error) {
	return a.p.AccountSummary(a.id)
}

func (a *Account) AccountItems() (*Request, error) {
	return a.p.AccountItems(a.id)
}

func (a *Account) ActivityHistory(cid string) (*Request, error) {
	return a.p.ActivityHistory(a.id, cid)
}

func (a *Account) CharacterSummary(cid string) (*Request, error) {
	return a.p.CharacterSummary(a.id, cid)
}

func (a *Account) CharacterActivities(cid string) (*Request, error) {
	return a.p.CharacterActivities(a.id, cid)
}

func (a *Account) CharacterInventory(cid string) (*Request, error) {
	return a.p.CharacterInventory(a.id, cid)
}

func (a *Account) CharacterProgression(cid string) (*Request, error) {
	return a.p.CharacterProgression(a.id, cid)
}

func (a *Account) AggregateActivityStats(cid string) (*Request, error) {
	return a.p.AggregateActivityStats(a.id, cid)
}

func (a *Account) UniqueWeapons(cid string) (*Request, error) {
	return a.p.UniqueWeapons(a.id, cid)
}

func (a *Account) Character(id string) *Character {
	return &Character{
		a:  a,
		id: id,
	}
}

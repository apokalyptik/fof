package destiny

type Platform struct {
	c  *Client
	id int
}

func (p *Platform) SearchDestinyPlayer(gamertag string) (*Request, error) {
	return p.c.SearchDestinyPlayer(p.id, gamertag)
}

func (p *Platform) AccountSummary(id string) (*Request, error) {
	return p.c.AccountSummary(p.id, id)
}

func (p *Platform) AccountItems(id string) (*Request, error) {
	return p.c.AccountItems(p.id, id)
}

func (p *Platform) ActivityHistory(id string, cid string) (*Request, error) {
	return p.c.ActivityHistory(p.id, id, cid)
}

func (p *Platform) CharacterSummary(id string, cid string) (*Request, error) {
	return p.c.CharacterSummary(p.id, id, cid)
}

func (p *Platform) CharacterActivities(id string, cid string) (*Request, error) {
	return p.c.CharacterActivities(p.id, id, cid)
}

func (p *Platform) CharacterInventory(id string, cid string) (*Request, error) {
	return p.c.CharacterInventory(p.id, id, cid)
}

func (p *Platform) CharacterProgression(id string, cid string) (*Request, error) {
	return p.c.CharacterProgression(p.id, id, cid)
}

func (p *Platform) AggregateActivityStats(id string, cid string) (*Request, error) {
	return p.c.AggregateActivityStats(p.id, id, cid)
}

func (p *Platform) UniqueWeapons(id string, cid string) (*Request, error) {
	return p.c.UniqueWeapons(p.id, id, cid)
}

func (p *Platform) Account(id string) *Account {
	return &Account{
		id: id,
		p:  p,
	}
}

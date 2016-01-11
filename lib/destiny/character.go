package destiny

type Character struct {
	id string
	a  *Account
}

func (c *Character) ActivityHistory() (*Request, error) {
	return c.a.ActivityHistory(c.id)
}

func (c *Character) CharacterSummary() (*Request, error) {
	return c.a.CharacterSummary(c.id)
}

func (c *Character) CharacterActivities() (*Request, error) {
	return c.a.CharacterActivities(c.id)
}

func (c *Character) CharacterInventory() (*Request, error) {
	return c.a.CharacterInventory(c.id)
}

func (c *Character) CharacterProgression() (*Request, error) {
	return c.a.CharacterProgression(c.id)
}

func (c *Character) AggregateActivityStats() (*Request, error) {
	return c.a.AggregateActivityStats(c.id)
}

func (c *Character) UniqueWeapons() (*Request, error) {
	return c.a.UniqueWeapons(c.id)
}

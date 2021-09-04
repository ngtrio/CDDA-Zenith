package data

type Mod struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Dependencies []string
	Data         map[string][]string
	TempData     map[string][]string
	Loaded       bool
}

func (mod *Mod) getById(id string) []string {
	if v, ok := mod.Data[id]; ok {
		return v
	}
	return nil
}

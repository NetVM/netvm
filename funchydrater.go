package netvm

// FuncHydrater can hydrate deadheads with funcs
type FuncHydrater map[string]HydratedFunc

// Hydrate hydrates the id with the func
func (f FuncHydrater) Hydrate(id string) (HydratedFunc, error) {
	return f[id], nil
}

package mgollective

type Factsource struct {
}

func (Factsource) GetFact(name string) string {
	return "yes"
}

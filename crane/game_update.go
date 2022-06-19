package crane

func (g *Game) Update() error {
	g.crane.Update()
	g.level.Update()

	g.space.Step(1. / 60.)
	return nil
}

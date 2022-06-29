package crane

func (g *Game) Update() error {
	g.crane.Update()
	g.level.Update()

	g.space.Step(1. / 180.)
	g.space.Step(1. / 180.)
	g.space.Step(1. / 180.)
	return nil
}

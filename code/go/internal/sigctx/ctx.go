package sigctx

type (
	Ctx <-chan struct{}
)

func (c Ctx) Done() <-chan struct{} { return c }

package command

type Command interface {
	Ctx() Context
	isCommand()
}

type Context struct {
	ConnID string
}

func (c Context) Ctx() Context { return c }
func (c Context) isCommand()   {}

package checkpoint

type Context struct {
	AddrMap map[string]string
	SvcMap  map[string]string
}

func NewContext() *Context {
	return &Context{
		AddrMap: make(map[string]string),
		SvcMap:  make(map[string]string),
	}
}

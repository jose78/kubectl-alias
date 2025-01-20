package alias

import "context"


// Command pattern to encapsulate an execution
type Command interface {
	Execute(context.Context)
}
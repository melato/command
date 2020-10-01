package command

/** Provides no-op implementation for Command methods */
type Base struct {
	commands map[string]Command
}

func (t *Base) Run(args []string) error { return nil }
func (t *Base) Init() error             { return nil }
func (t *Base) Configured() error       { return nil }
func (t *Base) Usage() *Usage           { return nil }
func (t *Base) Commands() map[string]Command {
	if t.commands == nil {
		t.commands = make(map[string]Command)
	}
	return t.commands
}

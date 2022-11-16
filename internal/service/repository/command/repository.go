package command

type CommandRepositoryInterface interface {
	Construct() CommandRepositoryInterface
	GetAddCommands() map[string]any
	GetShowCommands() map[string]any
}

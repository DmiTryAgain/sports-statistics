package command

type RepositoryInterface interface {
	Construct() RepositoryInterface
	GetAddCommands() map[string]any
	GetShowCommands() map[string]any
	GetHelpCommands() map[string]any
}

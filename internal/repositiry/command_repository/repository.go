package command_repository

import (
	"sports-statistics/internal/service/helpers"
	"sports-statistics/internal/service/repository/command"
)

type CommandRepository struct {
	addCommands, showCommands, helpCommands map[string]any
}

func (r *CommandRepository) Construct() command.RepositoryInterface {
	r.addCommands = map[string]any{
		"сделал":   helpers.StructStub{},
		"добавь":   helpers.StructStub{},
		"добавить": helpers.StructStub{},
	}

	r.showCommands = map[string]any{
		"покажи": helpers.StructStub{},
		"выведи": helpers.StructStub{},
	}

	r.helpCommands = map[string]any{
		"помощь": helpers.StructStub{},
		"помоги": helpers.StructStub{},
		"help":   helpers.StructStub{},
	}

	return r
}

func (r *CommandRepository) GetAddCommands() map[string]any {
	return r.addCommands
}

func (r *CommandRepository) GetShowCommands() map[string]any {
	return r.showCommands
}

func (r *CommandRepository) GetHelpCommands() map[string]any {
	return r.helpCommands
}

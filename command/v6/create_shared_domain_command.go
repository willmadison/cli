package v6

import (
	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v6/shared"
)

//go:generate counterfeiter . CreateSharedDomainActor
type CreateSharedDomainActor interface {
	GetRouterGroupByName(string, v2action.RouterClient) (v2action.RouterGroup, error)
}

type CreateSharedDomainCommand struct {
	RequiredArgs    flag.Domain `positional-args:"yes"`
	RouterGroup     string      `long:"router-group" description:"Routes for this domain will be configured only on the specified router group"`
	usage           interface{} `usage:"CF_NAME create-shared-domain DOMAIN [--router-group ROUTER_GROUP]"`
	relatedCommands interface{} `related_commands:"create-domain, domains, router-groups"`

	UI           command.UI
	Config       command.Config
	Actor        CreateSharedDomainActor
	SharedActor  command.SharedActor
	RouterClient v2action.RouterClient
}

func (cmd *CreateSharedDomainCommand) Setup(config command.Config, ui command.UI) error {
	ccClient, uaaClient, err := shared.NewClients(config, ui, true)
	routerClient, _ := shared.NewRouterClient(config, ui) // TODO handle the error

	if err != nil {
		return err
	}

	cmd.Actor = v2action.NewActor(ccClient, uaaClient, config)
	cmd.RouterClient = routerClient
	cmd.SharedActor = sharedaction.NewActor(config)
	cmd.Config = config
	cmd.UI = ui
	return nil
}

func (cmd CreateSharedDomainCommand) Execute(args []string) error {
	username, err := cmd.SharedActor.RequireCurrentUser()
	cmd.UI.DisplayTextWithFlavor("Creating shared domain {{.Domain}} as {{.User}}...",
		map[string]interface{}{
			"Domain": cmd.RequiredArgs.Domain,
			"User":   username,
		})

	if err != nil {
		return err
	}
	_, routerGroupErr := cmd.Actor.GetRouterGroupByName("", cmd.RouterClient)

	if routerGroupErr != nil {
		return routerGroupErr
	}
	return nil
}

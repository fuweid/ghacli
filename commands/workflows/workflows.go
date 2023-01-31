package workflows

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/v50/github"
	"github.com/urfave/cli"

	"github.com/fuweid/ghacli/commands"
)

// Command represents workflows related commands.
var Command = cli.Command{
	Name:    "workflows",
	Aliases: []string{"w"},
	Usage:   "View workflows for a repository in GitHub Actions",
	Subcommands: []cli.Command{
		listCommand,
	},
}

var listCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "Lists the workflows in a repository.",
	Action: func(cliCtx *cli.Context) error {
		ctx, client, err := commands.NewGithubClient(cliCtx)
		if err != nil {
			return fmt.Errorf("failed to init github client: %w", err)
		}

		owner, repo := commands.TargetRepo(cliCtx)

		workflows, err := listAllWorkflows(ctx, client.Actions, owner, repo)
		if err != nil {
			return fmt.Errorf("failed to list workflows: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATE\tPATH\t")
		for _, workflow := range workflows {
			_, err := fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n",
				*workflow.ID,
				*workflow.Name,
				*workflow.State,
				*workflow.Path,
			)
			if err != nil {
				return fmt.Errorf("failed to print result: %w", err)
			}
		}
		return w.Flush()
	},
}

func listAllWorkflows(ctx context.Context, action *github.ActionsService, owner, repo string) ([]*github.Workflow, error) {
	handler := func(ctx context.Context, opt *github.ListOptions) ([]*github.Workflow, *github.Response, error) {
		workflows, resp, err := action.ListWorkflows(ctx, owner, repo, opt)
		if err != nil {
			return nil, nil, err
		}
		return workflows.Workflows, resp, nil
	}
	return commands.ListAllItems[*github.Workflow](ctx, handler, 0)
}

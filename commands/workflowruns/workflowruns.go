package workflowruns

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/urfave/cli"

	"github.com/fuweid/ghacli/commands"
)

// Command represents workflowruns related commands.
var Command = cli.Command{
	Name:  "run",
	Usage: "About workflow runs for a repository in GitHub Actions",
	Subcommands: []cli.Command{
		listCommand,
	},
}

var listCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List workflow runs",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "workflow-id",
			Usage: "The ID of the workflow. You can also pass the workflow file name as a string.",
		},
		cli.StringFlag{
			Name: "status",
			Usage: `
Returns workflow runs with the check run status or conclusion that you specify. Can be one of: completed, action_required, cancelled, failure, neutral, skipped, stale, success, timed_out, in_progress, queued, requested, waiting, pending
`,
		},
		cli.StringFlag{
			Name:  "created",
			Usage: "Returns workflow runs created within the given date-time range. For more information on the syntax",
			// https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax#query-for-dates
		},
		cli.StringFlag{
			Name:  "event",
			Usage: "Returns workflow run triggered by the event you specify. For example, push, pull_request or issue.",
			// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows
		},
		cli.Uint64Flag{
			Name:  "limit",
			Usage: "Max number of runs to be fetched",
			Value: 100,
		},
	},
	Action: func(cliCtx *cli.Context) error {
		ctx, client, err := commands.NewGithubClient(cliCtx)
		if err != nil {
			return fmt.Errorf("failed to init github client: %w", err)
		}

		owner, repo := commands.TargetRepo(cliCtx)
		workflowID := cliCtx.String("workflow-id")

		filterOpt := &github.ListWorkflowRunsOptions{
			Created: cliCtx.String("created"),
			Status:  cliCtx.String("status"),
			Event:   cliCtx.String("event"),
		}
		limit := int(cliCtx.Uint64("limit"))

		runs, err := listAllWorkflowRuns(ctx, client.Actions,
			owner, repo, workflowID, filterOpt, limit)
		if err != nil {
			return fmt.Errorf("failed to list workflow runs for %s: %w",
				workflowID, err)
		}

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tEVENT\tSTATUS\tHEAD MESSAGE\tCREATED")
		for _, run := range runs {
			_, err := fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t\n",
				*run.ID,
				*run.Name,
				*run.Event,
				*run.Conclusion, // NOTE: Status is completed even if failure
				subjectLineOfGitMsg(*run.HeadCommit.Message),
				run.CreatedAt.UTC().Format(time.UnixDate),
			)
			if err != nil {
				return fmt.Errorf("failed to print result: %w", err)
			}
		}
		return w.Flush()
	},
}

func listAllWorkflowRuns(
	ctx context.Context,
	action *github.ActionsService,
	owner, repo, workflowID string,
	filterOpt *github.ListWorkflowRunsOptions,
	limit int,
) ([]*github.WorkflowRun, error) {
	handler := func(ctx context.Context, pageOpt *github.ListOptions) ([]*github.WorkflowRun, *github.Response, error) {
		filterOpt.ListOptions = *pageOpt

		var err error
		var runs *github.WorkflowRuns
		var resp *github.Response

		if len(workflowID) == 0 {
			runs, resp, err = action.ListRepositoryWorkflowRuns(ctx, owner, repo, filterOpt)
		} else {
			runs, resp, err = action.ListWorkflowRunsByFileName(ctx, owner, repo, workflowID, filterOpt)
		}

		if err != nil {
			return nil, nil, err
		}
		return runs.WorkflowRuns, resp, nil
	}
	return commands.ListAllItems[*github.WorkflowRun](ctx, handler, limit)
}

// subjectLineOfGitMsg limits message in 80 chars.
func subjectLineOfGitMsg(gitMsg string) string {
	limit := 80

	var msg string
	msgs := strings.Split(gitMsg, "\n")
	if len(msgs) == 0 {
		msg = gitMsg
	} else {
		msg = msgs[0]
	}

	if len(msg) <= limit {
		return msg
	}
	return msg[:limit]
}

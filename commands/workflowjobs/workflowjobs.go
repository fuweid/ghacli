package workflowjobs

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/urfave/cli"

	"github.com/fuweid/ghacli/commands"
)

// Command represents workflowruns related commands.
var Command = cli.Command{
	Name:  "job",
	Usage: "About workflow jobs for a repository in GitHub Actions",
	Subcommands: []cli.Command{
		downloadLogCommand,
	},
}

var downloadLogCommand = cli.Command{
	Name:  "logs",
	Usage: "Download a job's log",
	Action: func(cliCtx *cli.Context) error {
		if !cliCtx.Args().Present() {
			return fmt.Errorf("job ID is required")
		}
		if cliCtx.NArg() > 1 {
			return fmt.Errorf("only support one job")
		}

		jobID, err := strconv.ParseInt(cliCtx.Args().First(), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse %s into int64: %w",
				cliCtx.Args().First(), jobID)
		}

		ctx, client, err := commands.NewGithubClient(cliCtx)
		if err != nil {
			return fmt.Errorf("failed to init github client: %w", err)
		}
		owner, repo := commands.TargetRepo(cliCtx)

		_, resp, err := client.Actions.GetWorkflowJobLogs(ctx, owner, repo, jobID, true)
		if err != nil {
			return fmt.Errorf("failed to locate job %v log: %w", jobID, err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read job %v log: %w", jobID, err)
		}
		return nil
	},
}

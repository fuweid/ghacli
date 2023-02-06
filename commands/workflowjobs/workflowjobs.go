package workflowjobs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/urfave/cli"

	"github.com/fuweid/ghacli/commands"
)

// Command represents workflowruns related commands.
var Command = cli.Command{
	Name:  "job",
	Usage: "About workflow jobs for a repository in GitHub Actions",
	Subcommands: []cli.Command{
		downloadLogCommand,
		listJobsCommand,
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
				cliCtx.Args().First(), err)
		}

		ctx, client, err := commands.NewGithubClient(cliCtx)
		if err != nil {
			return fmt.Errorf("failed to init github client: %w", err)
		}
		owner, repo := commands.TargetRepo(cliCtx)

		url, _, err := client.Actions.GetWorkflowJobLogs(ctx, owner, repo, jobID, true)
		if err != nil {
			return fmt.Errorf("failed to locate job %v log: %w", jobID, err)
		}

		resp, err := http.Get(url.String())
		if err != nil {
			return fmt.Errorf("failed to GET %s: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code %v", resp.StatusCode)
		}

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read job %v log: %w", jobID, err)
		}
		return nil
	},
}

var listJobsCommand = cli.Command{
	Name:  "list",
	Usage: "List workflow jobs with a given workflow run ID",
	Action: func(cliCtx *cli.Context) error {
		if !cliCtx.Args().Present() {
			return fmt.Errorf("run ID is required")
		}
		if cliCtx.NArg() > 1 {
			return fmt.Errorf("only support one run")
		}

		runID, err := strconv.ParseInt(cliCtx.Args().First(), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse %s into int64: %w",
				cliCtx.Args().First(), err)
		}

		ctx, client, err := commands.NewGithubClient(cliCtx)
		if err != nil {
			return fmt.Errorf("failed to init github client: %w", err)
		}
		owner, repo := commands.TargetRepo(cliCtx)

		jobs, err := commands.ListAllJobsForARun(ctx, client.Actions, owner, repo, runID)
		if err != nil {
			return fmt.Errorf("failed to list jobs for run %v: %w", runID, err)
		}

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 2, ' ', 0)

		header := "ID\tNAME\tSTATUS\tRUN_ATTEMPT\tWORKFLOW_NAME"
		fmt.Fprintln(w, header)
		for _, job := range jobs {
			status := *job.Status
			if job.Conclusion != nil {
				status = *job.Conclusion
			}
			_, err := fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%s\n",
				*job.ID,
				*job.Name,
				status,
				*job.RunAttempt,
				*job.WorkflowName,
			)
			if err != nil {
				return fmt.Errorf("failed to print result: %w", err)
			}
		}
		return w.Flush()
	},
}

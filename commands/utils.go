package commands

import (
	"context"
	"net/http"

	"github.com/google/go-github/v50/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

// NewGithubClient returns github client with global flag "token".
func NewGithubClient(cliCtx *cli.Context) (context.Context, *github.Client, error) {
	ctx := context.Background()
	token := cliCtx.GlobalString("token")
	var tc *http.Client

	if len(token) > 0 {
		tc = oauth2.NewClient(ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
				},
			),
		)
	}
	return ctx, github.NewClient(tc), nil
}

// TargetRepo returns target owner and repo from global "owner" and "repo" flags.
func TargetRepo(cliCtx *cli.Context) (owner string, repo string) {
	return cliCtx.GlobalString("owner"), cliCtx.GlobalString("repo")
}

// ListAllItems is used to list items page by page.
func ListAllItems[ItemType any](
	ctx context.Context,
	handler func(ctx context.Context, opt *github.ListOptions) ([]ItemType, *github.Response, error),
	limit int, /* zero means unlimit */
) ([]ItemType, error) {

	opt := &github.ListOptions{
		PerPage: 100,
	}

	var res []ItemType
	for {
		items, resp, err := handler(ctx, opt)
		if err != nil {
			return nil, err
		}

		res = append(res, items...)

		if limit > 0 && len(res) >= limit {
			res = res[:limit]
			break
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return res, nil
}

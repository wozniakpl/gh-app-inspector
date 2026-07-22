package inspector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v89/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Config struct {
	AppID          int64
	InstallationID int64
	PEM            []byte
}

func Run(ctx context.Context, out io.Writer, cfg Config) error {
	appTr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, cfg.PEM)
	if err != nil {
		return fmt.Errorf("apps transport: %w", err)
	}
	appClient, err := github.NewClient(github.WithHTTPClient(&http.Client{Transport: appTr}))
	if err != nil {
		return fmt.Errorf("app client: %w", err)
	}

	app, _, err := appClient.Apps.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("get app: %w", err)
	}
	renderApp(out, app)

	install, _, err := appClient.Apps.GetInstallation(ctx, cfg.InstallationID)
	if err != nil {
		return fmt.Errorf("get installation: %w", err)
	}
	renderInstallation(out, install)
	renderPermissions(out, install.GetPermissions())

	instTr := ghinstallation.NewFromAppsTransport(appTr, cfg.InstallationID)
	instClient, err := github.NewClient(github.WithHTTPClient(&http.Client{Transport: instTr}))
	if err != nil {
		return fmt.Errorf("installation client: %w", err)
	}

	if err := renderRepos(ctx, out, instClient); err != nil {
		return err
	}

	return renderRateLimit(ctx, out, instClient)
}

func renderApp(out io.Writer, a *github.App) {
	t := newTable(out, "GitHub App")
	t.AppendRow(table.Row{"Name", a.GetName()})
	t.AppendRow(table.Row{"Slug", a.GetSlug()})
	t.AppendRow(table.Row{"ID", a.GetID()})
	t.AppendRow(table.Row{"Owner", a.GetOwner().GetLogin()})
	t.AppendRow(table.Row{"HTML URL", a.GetHTMLURL()})
	t.AppendRow(table.Row{"Created", fmtTime(a.GetCreatedAt().Time)})
	t.Render()
	fmt.Fprintln(out)
}

func renderInstallation(out io.Writer, i *github.Installation) {
	t := newTable(out, "Installation")
	t.AppendRow(table.Row{"ID", i.GetID()})
	t.AppendRow(table.Row{"Account", i.GetAccount().GetLogin()})
	t.AppendRow(table.Row{"Account type", i.GetAccount().GetType()})
	t.AppendRow(table.Row{"Target type", i.GetTargetType()})
	t.AppendRow(table.Row{"Repository selection", i.GetRepositorySelection()})
	t.AppendRow(table.Row{"Events", fmt.Sprintf("%v", i.Events)})
	t.AppendRow(table.Row{"Created", fmtTime(i.GetCreatedAt().Time)})
	t.AppendRow(table.Row{"Updated", fmtTime(i.GetUpdatedAt().Time)})
	t.Render()
	fmt.Fprintln(out)
}

func renderPermissions(out io.Writer, p *github.InstallationPermissions) {
	rows := permissionRows(p)
	t := newTable(out, fmt.Sprintf("Permissions (%d)", len(rows)))
	t.AppendHeader(table.Row{"Scope", "Access"})
	for _, r := range rows {
		t.AppendRow(table.Row{r[0], r[1]})
	}
	t.Render()
	fmt.Fprintln(out)
}

func renderRepos(ctx context.Context, out io.Writer, c *github.Client) error {
	opts := &github.ListOptions{PerPage: 100}
	var all []*github.Repository
	for {
		page, resp, err := c.Apps.ListRepos(ctx, opts)
		if err != nil {
			return fmt.Errorf("list repos: %w", err)
		}
		all = append(all, page.Repositories...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	sort.Slice(all, func(i, j int) bool { return all[i].GetFullName() < all[j].GetFullName() })

	t := newTable(out, fmt.Sprintf("Accessible repositories (%d)", len(all)))
	t.AppendHeader(table.Row{"Repository", "Private", "Default branch", "Pushed at"})
	for _, r := range all {
		t.AppendRow(table.Row{
			r.GetFullName(),
			r.GetPrivate(),
			r.GetDefaultBranch(),
			fmtTime(r.GetPushedAt().Time),
		})
	}
	t.Render()
	fmt.Fprintln(out)
	return nil
}

func renderRateLimit(ctx context.Context, out io.Writer, c *github.Client) error {
	rl, _, err := c.RateLimit.Get(ctx)
	if err != nil {
		return fmt.Errorf("rate limit: %w", err)
	}
	core := rl.GetCore()
	t := newTable(out, "Rate limit (core)")
	t.AppendRow(table.Row{"Limit", core.Limit})
	t.AppendRow(table.Row{"Remaining", core.Remaining})
	t.AppendRow(table.Row{"Resets", fmtTime(core.Reset.Time)})
	t.Render()
	fmt.Fprintln(out)
	return nil
}

func newTable(out io.Writer, title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(out)
	t.SetTitle(title)
	t.SetStyle(table.StyleRounded)
	return t
}

func fmtTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.UTC().Format(time.RFC3339)
}

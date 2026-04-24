// preflight.go — `movie preflight` command.
//
// Verifies the local git clone is on the expected branch and in sync with
// origin BEFORE any updater, build, or rest action runs. When stale, prints
// the exact recovery commands documented in
// .lovable/pending-issues/01-local-repo-stale.md.
//
// Read-only: never runs checkout/reset/clean — that boundary belongs to
// the user and run.ps1 (mem://constraints/updater-scope).
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v6/gitcheck"
)

const (
	preflightExitOK    = 0
	preflightExitStale = 4
	preflightExitError = 1
)

var (
	preflightRepoPath string
	preflightRemote   string
	preflightBranch   string
	preflightFetch    bool
)

var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Verify local repo is on expected branch and in sync with origin",
	Long: `Inspects the local git clone and reports whether it is in sync
with the expected upstream (default: origin/main). When the local clone is
stale, on a different branch, or has uncommitted changes, prints the exact
git commands required to recover.

No mutating git command is ever run by preflight itself — fetch is read-only
and only when --fetch is passed.

Exit codes:
  0 = repo is clean and up to date
  4 = repo is stale, dirty, or on the wrong branch
  1 = preflight could not run (no git, bad path, etc.)`,
	Run: runPreflightCmd,
}

func runPreflightCmd(cmd *cobra.Command, args []string) {
	repoPath := resolvePreflightRepo()
	status, err := gitcheck.Check(gitcheck.Options{
		RepoPath: repoPath,
		Remote:   preflightRemote,
		Branch:   preflightBranch,
		DoFetch:  preflightFetch,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "preflight: %v\n", err)
		os.Exit(preflightExitError)
	}
	printPreflight(status)
	os.Exit(preflightExitFor(status))
}

func resolvePreflightRepo() string {
	if preflightRepoPath != "" {
		abs, err := filepath.Abs(preflightRepoPath)
		if err == nil {
			return abs
		}
		return preflightRepoPath
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

func printPreflight(s *gitcheck.Status) {
	fmt.Println("==> movie preflight (repo sync check)")
	fmt.Println("  --------------------------------------------------")
	fmt.Printf("    repo          : %s\n", orDashStr(s.RepoPath))
	fmt.Printf("    expected      : %s/%s\n", s.Remote, s.Branch)
	fmt.Printf("    current branch: %s\n", orDashStr(s.CurrentBranch))
	fmt.Printf("    local commit  : %s\n", orDashStr(s.LocalCommit))
	fmt.Printf("    remote commit : %s\n", orDashStr(s.RemoteCommit))
	fmt.Printf("    ahead/behind  : %s/%s\n", orDashStr(s.Ahead), orDashStr(s.Behind))
	fmt.Println("  --------------------------------------------------")
	printPreflightVerdict(s)
}

func printPreflightVerdict(s *gitcheck.Status) {
	if !s.IsRepo {
		fmt.Println("    [ERR ] not a git repository")
		fmt.Println("      hint: run preflight inside the movie-cli-v6 clone")
		return
	}
	if !s.IsStale() {
		fmt.Println("    [ OK ] repo is on expected branch and up to date")
		return
	}
	printPreflightIssues(s)
	printPreflightRecovery(s)
}

func printPreflightIssues(s *gitcheck.Status) {
	if !s.IsOnBranch {
		fmt.Printf("    [WARN] on '%s', expected '%s'\n", s.CurrentBranch, s.Branch)
	}
	if !s.IsClean {
		fmt.Println("    [WARN] working tree has uncommitted changes")
	}
	if !s.IsUpToDate {
		fmt.Printf("    [WARN] local is behind %s/%s by %s commit(s)\n",
			s.Remote, s.Branch, orDashStr(s.Behind))
	}
}

func printPreflightRecovery(s *gitcheck.Status) {
	fmt.Println("  --------------------------------------------------")
	fmt.Println("    To force-sync the local clone, run:")
	for _, line := range s.RecoveryCommands() {
		fmt.Printf("      %s\n", line)
	}
	fmt.Println("    Warning: `git reset --hard` discards local commits and changes.")
}

func preflightExitFor(s *gitcheck.Status) int {
	if !s.IsRepo {
		return preflightExitError
	}
	if s.IsStale() {
		return preflightExitStale
	}
	return preflightExitOK
}

func orDashStr(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func init() {
	preflightCmd.Flags().StringVar(&preflightRepoPath, "repo-path", "", "Path to the movie-cli-v6 repo (default: current directory)")
	preflightCmd.Flags().StringVar(&preflightRemote, "remote", gitcheck.DefaultRemote, "Expected git remote name")
	preflightCmd.Flags().StringVar(&preflightBranch, "branch", gitcheck.DefaultBranch, "Expected git branch")
	preflightCmd.Flags().BoolVar(&preflightFetch, "fetch", false, "Run `git fetch <remote> <branch>` before checking (read-only)")
}
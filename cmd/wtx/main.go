package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/darkLord19/wtx/internal/config"
	"github.com/darkLord19/wtx/internal/editor"
	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
	"github.com/darkLord19/wtx/internal/tui"
)

var (
	cfg        *config.Config
	gitMgr     *git.Manager
	metaStore  *metadata.Store
	edDetector *editor.Detector
	fullTUI    bool
	isFirstRun bool
)

var rootCmd = &cobra.Command{
	Use:   "wtx",
	Short: "Git worktree workspace manager",
	Long:  `wtx makes Git worktrees feel like instant "workspace tabs" across editors.`,
	RunE:  runInteractive,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add flags
	rootCmd.Flags().BoolVarP(&fullTUI, "tui", "t", false, "Launch full TUI with tabs (worktrees, manage, settings)")

	// Add commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(pruneCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(manageCmd)
	rootCmd.AddCommand(setupCmd)
}

func initConfig() {
	var err error

	// Check if git is installed
	if !git.IsGitInstalled() {
		fmt.Fprintf(os.Stderr, "Error: git is not installed or not in PATH\n")
		os.Exit(1)
	}

	// Load config with first run check
	loadResult, err := config.LoadWithFirstRunCheck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	cfg = loadResult.Config
	isFirstRun = loadResult.IsFirstRun

	// Find git repo
	repoPath, err := git.GetRootPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	repo, err := git.FindRepo(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize managers
	gitMgr = git.NewManager(repo)
	metaStore, err = metadata.Load(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading metadata: %v\n", err)
		os.Exit(1)
	}

	edDetector = editor.NewDetector(cfg)
}

func runInteractive(cmd *cobra.Command, args []string) error {
	// Check for first run and launch setup wizard
	if isFirstRun {
		completed, err := tui.RunSetup(cfg, edDetector)
		if err != nil {
			return fmt.Errorf("setup wizard failed: %w", err)
		}
		if !completed {
			fmt.Println("Setup cancelled. Run 'wtx' again to complete setup.")
			return nil
		}
		// Reload editor detector with new config
		edDetector = editor.NewDetector(cfg)
	}

	var selected *tui.WorktreeItem
	var err error

	if fullTUI {
		// Run full TUI manager with tabs
		selected, err = tui.RunManager(gitMgr, metaStore, cfg, edDetector)
	} else {
		// Run simple TUI selector
		selected, err = tui.Run(gitMgr, metaStore)
	}

	if err != nil {
		return err
	}

	if selected == nil {
		return nil // User cancelled
	}

	// Open in editor
	ed, err := edDetector.GetPreferred()
	if err != nil {
		return err
	}

	fmt.Printf("Opening %s in %s...\n", selected.Name, ed.Name())

	if err := ed.Open(selected.Path, cfg.ReuseWindow); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Update metadata
	metaStore.Touch(selected.Name)
	return metaStore.Save()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

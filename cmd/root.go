package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dodirepository/warehouse-svc/cmd/http"
	"github.com/dodirepository/warehouse-svc/cmd/migration"
	"github.com/spf13/cobra"
)

// Start handler registering service command
func Start() {

	rootCmd := &cobra.Command{}
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	migrateCmd := &cobra.Command{
		Use:   "db:migrate",
		Short: "database migration",
		Run: func(c *cobra.Command, args []string) {
			migration.DatabaseMigration()
		},
	}

	migrateCmd.Flags().BoolP("version", "", false, "print version")
	migrateCmd.Flags().StringP("dir", "", "infrastructure/database/migration/", "directory with migration files")
	migrateCmd.Flags().StringP("table", "", "db", "migrations table name")
	migrateCmd.Flags().BoolP("verbose", "", false, "enable verbose mode")
	migrateCmd.Flags().BoolP("guide", "", false, "print help")

	httpRunner := &cobra.Command{
		Use:   "serve",
		Short: "Run HTTP Server",
		Run: func(cmd *cobra.Command, args []string) {
			http.Start(ctx)
		},
	}

	cmd := []*cobra.Command{
		httpRunner,
		migrateCmd,
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

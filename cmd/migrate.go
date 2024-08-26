package cmd

import (
	"context"

	"github.com/uniwise/parrot/migrations"
	"github.com/uniwise/parrot/pkg/connectors/database"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate the database",
	Run:  migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	dbConfig := database.Config{
		DSN:   viper.GetString("database.dsn"),
		Debug: viper.GetBool("database.debug"),
	}

	l := instantiateLogger()

	db, err := database.NewClient(ctx, dbConfig)
	if err != nil {
		l.WithError(err).Fatal("failed to connect to database")
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.Migration00001Init,
		migrations.Migration00002AddStorageKeyToVersionsTable,
	})

	if err := m.Migrate(); err != nil {
		logrus.WithError(err).Fatal("Migration failed")
	}

	logrus.Println("Migration did run successfully")
}
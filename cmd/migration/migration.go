package migration

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

func init() {
	godotenv.Load()
}

var (
	flags   = flag.NewFlagSet("db:migrate", flag.ExitOnError)
	dir     = flags.String("dir", "infrastructure/database/migration", "directory with migration files")
	table   = flags.String("table", "db_migration", "migrations table name")
	verbose = flags.Bool("verbose", false, "enable verbose mode")
	help    = flags.Bool("guide", false, "print help")
	version = flags.Bool("version", false, "print version")
)

var (
	usageCommands = `
  --dir string     directory with migration files (default "database/migration")
  --guide          print help
  --table string   migrations table name (default "db_migration")
  --verbose        enable verbose mode
  --version        print version

Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)

// DatabaseMigration :nodoc:
func DatabaseMigration() {

	flags.Usage = usage
	flags.Parse(os.Args[2:])

	if *version {
		fmt.Println(goose.VERSION)
		return
	}
	if *verbose {
		goose.SetVerbose(true)
	}

	goose.SetTableName(*table)

	args := flags.Args()

	if len(args) == 0 || *help {
		flags.Usage()
		return
	}
	fmt.Printf("args: %v\n", dir)
	switch args[0] {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	port, _ := strconv.Atoi(os.Getenv("DB_MYSQL_PORT"))
	dbstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", os.Getenv("DB_MYSQL_USERNAME"), os.Getenv("DB_MYSQL_PASSWORD"), os.Getenv("DB_MYSQL_HOST"), port, os.Getenv("DB_MYSQL_DATABASE"))
	db, err := goose.OpenDBWithDriver("mysql", normalizeDBString("mysql", dbstring))

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("db migrate run: %v", err)
	}
}

func usage() {
	fmt.Println(usageCommands)
}

func normalizeDBString(driver string, str string) string {
	if driver == "mysql" {
		var err error
		str, err = normalizeMySQLDSN(str)
		if err != nil {
			log.Fatalf("failed to normalize MySQL connection string: %v", err)
		}
	}
	return str
}

func normalizeMySQLDSN(dsn string) (string, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", err
	}
	config.ParseTime = true
	return config.FormatDSN(), nil
}

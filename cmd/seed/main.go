package main

import (
	"errors"
	"fmt"
	"liga-app/cmd/seed/court_time_slots"
	"liga-app/cmd/seed/leagues"
	"liga-app/cmd/seed/matches"
	"liga-app/cmd/seed/round_schedule"
	"liga-app/cmd/seed/teams"
	"liga-app/cmd/seed/user_ratings"
	"liga-app/config"
	"liga-app/entities"
	"liga-app/infra/gsheet"
	"liga-app/infra/mysql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	gsheetAPI *gsheet.API
)

func LoadEnv() error {
	_, err := os.Stat(".env")
	if errors.Is(err, os.ErrNotExist) {
		log.Println("[WARNING] no .env file found")
		return nil
	}

	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return nil
}

func init() {
	err := LoadEnv()
	if err != nil {
		log.Fatalf("[cmd.app.main] LoadEnv(): failed to load env: %v", err)
	}

	// load and validate config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config.NewConfig(): failed to load config: %v", err)
	}

	// connect to database
	db, err = mysql.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("mysql.NewDB(): failed to connect to database: %v", err)
	}

	gsheetAPI, err = gsheet.NewAPI(cfg)
	if err != nil {
		log.Fatalf("mysql.NewDB(): failed to connect to database: %v", err)
	}
}

func resetDB(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		_ = c

		err := db.Migrator().DropTable(
			&entities.User{},
			&entities.MagicLinkLogin{},
			&entities.MagicLinkRegistration{},
			// &entities.Wilayah{},
			&entities.League{},
			&entities.LeagueRegistration{},
			&entities.Team{},
			&entities.TeamPlayer{},
			&entities.LeaguePlayer{},
			&entities.LeagueTeam{},
			&entities.UserRating{},
			&entities.RuleBook{},
			&entities.Match{},
			&entities.MatchSet{},
		)
		if err != nil {
			return fmt.Errorf("db.Migrator().DropTable(...): %w", err)
		}

		err = db.AutoMigrate(
			&entities.User{},
			&entities.MagicLinkLogin{},
			&entities.MagicLinkRegistration{},
			// &entities.Wilayah{},
			&entities.League{},
			&entities.LeagueRegistration{},
			&entities.Team{},
			&entities.TeamPlayer{},
			&entities.LeaguePlayer{},
			&entities.LeagueTeam{},
			&entities.UserRating{},
			&entities.RuleBook{},
			&entities.Match{},
			&entities.MatchSet{},
		)
		if err != nil {
			return fmt.Errorf("db.AutoMigrate(...): %w", err)
		}

		/*
			err = wilayah.SeedWilayah(db)
			if err != nil {
				return fmt.Errorf("wilayah.SeedWilayah(db): %w", err)
			}
		*/

		return nil
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "reset-db",
				Description: "Reset DB",
				Action:      resetDB(db),
			},
			{
				Name:        "reset-leagues",
				Description: "Reset leagues data",
				Action:      leagues.ResetLeagues(db),
			},
			{
				Name:        "seed-league-teams",
				Description: "Seed teams into leagues",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "team-count",
						Value: 12,
					},
				},
				Action: teams.SeedLeagueTeams(db),
			},
			{
				Name:        "generate-pool-matches",
				Description: "Generate pool matches",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "league-uuid",
					},
					&cli.BoolFlag{
						Name:  "home-and-away",
						Value: false,
					},
				},
				Action: matches.GeneratePoolMatches(db, gsheetAPI),
			},
			{
				Name:        "reset-round-schedules",
				Description: "Reset round schedules",
				Action:      round_schedule.ResetRoundSchedules(db, gsheetAPI),
			},
			{
				Name:        "reset-court-time-slots",
				Description: "Reset Court Time Slots",
				Action:      court_time_slots.ResetCourtTimeSlots(db, gsheetAPI),
			},
			{
				Name:        "update-court-time-slots",
				Description: "Update Court Time Slots",
				Action:      court_time_slots.UpdateCourtTimeSlots(db, gsheetAPI),
			},
			{
				Name:        "update-user-ratings",
				Description: "Update user ratings from gsheet",
				Action:      user_ratings.UpdateUserRatings(db, gsheetAPI),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("seed cli error: %s", err.Error())
	}
}

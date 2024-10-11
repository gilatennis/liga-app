package gsheet

import (
	"context"
	"fmt"
	"liga-app/config"
	"liga-app/infra/gsheet/court_time_slot"
	"liga-app/infra/gsheet/drawings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type API struct {
	CourtTimeSlot *court_time_slot.API
	Drawings      *drawings.API
}

func NewAPI(cfg *config.Config) (*API, error) {
	srv, err := setupService(cfg.Google.ServiceAccountCreds)
	if err != nil {
		return nil, fmt.Errorf("setupService(cfg.Google.ServiceAccountCreds): %w", err)
	}

	apiCourtTimeSlot := court_time_slot.NewAPI(srv)
	apiDrawings := drawings.NewAPI(srv)

	return &API{
		CourtTimeSlot: apiCourtTimeSlot,
		Drawings:      apiDrawings,
	}, nil
}

func setupService(creds string) (*sheets.Service, error) {
	ctx := context.Background()

	b := []byte(creds)

	jwtConf, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fmt.Errorf("google.JWTConfigFromJSON(b, https://www.googleapis.com/auth/spreadsheets.readonly)")
	}

	client := jwtConf.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("sheets.NewService(): %w", err)
	}

	return srv, nil
}

package league

import (
	"context"
	"fmt"
	"liga-app/config"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	"liga-app/interfaces/third_party/duitku"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	cfg           *config.Config
	gormRepo      repositories.IGormRepo
	leagueRepo    repositories.ILeagueRepository
	leagueRegRepo repositories.ILeagueRegistrationRepository
	wilayahRepo   repositories.IWilayahRepository
	duitkuAPI     duitku.IDuitkuAPI
}

func NewService(
	cfg *config.Config,
	gormRepo repositories.IGormRepo,
	leagueRepo repositories.ILeagueRepository,
	leagueRegRepo repositories.ILeagueRegistrationRepository,
	wilayahRepo repositories.IWilayahRepository,
	duitkuAPI duitku.IDuitkuAPI,
) *Service {
	return &Service{
		cfg:           cfg,
		gormRepo:      gormRepo,
		leagueRepo:    leagueRepo,
		leagueRegRepo: leagueRegRepo,
		wilayahRepo:   wilayahRepo,
		duitkuAPI:     duitkuAPI,
	}
}

func (s *Service) Get(ctx context.Context, p *common.Pagination, f *repositories.LeagueFilter) ([]entities.League, int64, error) {
	leagues, count, err := s.leagueRepo.Get(ctx, p, f)
	if err != nil {
		return nil, 0, err
	}

	return leagues, count, err
}

func (s *Service) GetByUUID(ctx context.Context, leagueUUID uuid.UUID) (*entities.League, error) {
	league, err := s.leagueRepo.GetByUUID(ctx, leagueUUID)
	if err != nil {
		return nil, err
	}

	return league, nil
}

func (s *Service) Join(ctx context.Context, leagueUUID uuid.UUID, user *entities.User) (string, error) {
	league, err := s.leagueRepo.GetByUUID(ctx, leagueUUID)
	if err != nil {
		return "", err
	}

	leagueRegistration := entities.LeagueRegistration{
		LeagueUUID:    leagueUUID,
		UserUUID:      user.UUID,
		PaymentStatus: enums.PAYMENT_STATUS_UNPAID,
	}

	leagueRegUUID, err := s.gormRepo.Save(ctx, &leagueRegistration, user)
	if err != nil {
		return "", err
	}

	duitkuPayload := s.createDuitkuPayload(ctx, league, user, *leagueRegUUID)

	resp, err := s.duitkuAPI.CreateInvoice(ctx, duitkuPayload)
	if err != nil {
		return "", err
	}

	return resp.PaymentURL, nil
}

func (s *Service) createDuitkuPayload(ctx context.Context, league *entities.League, user *entities.User, leagueRegUUID uuid.UUID) *duitku.CreateInvoicePayload {
	paymentAmount := int(league.RegistrationFee)
	customerAddress := s.createDuitkuCustomerAddress(ctx, user)
	fullName := customerAddress.FirstName + " " + customerAddress.LastName
	returnURL := s.createDuitkuReturnURL(ctx, league.UUID)

	merchantOrderID := leagueRegUUID.String()

	duitkuPayload := &duitku.CreateInvoicePayload{
		AdditionalParam:  "league_registration",
		PaymentAmount:    paymentAmount,
		MerchantOrderID:  merchantOrderID,
		ProductDetails:   fmt.Sprintf("Pendaftaran Liga - %s", league.Name),
		MerchantUserInfo: user.Email,
		CustomerVAName:   fullName,
		PhoneNumber:      customerAddress.Phone,
		CustomerDetail: duitku.CreateInvoiceCustomer{
			FirstName:       customerAddress.FirstName,
			LastName:        customerAddress.LastName,
			Email:           user.Email,
			PhoneNumber:     customerAddress.Phone,
			BillingAddress:  customerAddress,
			ShippingAddress: customerAddress,
		},
		ItemDetails: []duitku.CreateInvoiceItem{{
			Name:     fmt.Sprintf("Pendaftaran Liga - %s", league.Name),
			Quantity: 1,
			Price:    paymentAmount,
		}},
		CallbackURL:  s.cfg.App.BaseURL + "/api/v1/duitku/callback",
		ReturnURL:    returnURL,
		ExpiryPeriod: 60,
	}

	return duitkuPayload
}

func (s *Service) createDuitkuCustomerAddress(c context.Context, user *entities.User) duitku.CreateInvoiceCustomerAddress {
	firstName := strings.Split(user.Name, " ")[0]
	lastName := " "

	if len(strings.Split(user.Name, " ")) > 1 {
		lastName = strings.Split(user.Name, " ")[1]
	}

	address := "Jakarta"

	province, _ := s.wilayahRepo.Get(
		c,
		&common.Pagination{Page: 1, Limit: 1},
		&repositories.WilayahFilter{IDProvince: user.Province, Type: 1},
	)
	if province != nil {
		address = province[0].Name
	}

	customerAddress := duitku.CreateInvoiceCustomerAddress{
		FirstName:   firstName,
		LastName:    lastName,
		Address:     address,
		City:        " ", // INFO: we don't have city info for now
		PostalCode:  " ", // TODO: we don't have postal code info for now
		Phone:       user.Phone,
		CountryCode: "ID",
	}

	return customerAddress
}

func (s *Service) createDuitkuReturnURL(c context.Context, leagueUUID uuid.UUID) string {
	return s.cfg.App.BaseURL + fmt.Sprintf("/dashboard/leagues/%s/registration/success", leagueUUID)
}

func (s *Service) CompleteLeagueRegistrationPayment(
	ctx context.Context,
	leagueRegUUID uuid.UUID,
	paymentID string,
	paymentAmount float64,
) (*uuid.UUID, error) {
	leagueReg, err := s.leagueRegRepo.GetByUUID(ctx, leagueRegUUID)
	if err != nil {
		return nil, err
	}

	leagueReg.PaymentAmount = paymentAmount
	leagueReg.PaymentID = paymentID
	leagueReg.PaymentStatus = enums.PAYMENT_STATUS_PAID

	updatedUUID, err := s.leagueRegRepo.Update(ctx, leagueReg)
	if err != nil {
		return nil, err
	}

	return updatedUUID, nil
}

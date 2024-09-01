package duitku

import (
	"context"
)

type IDuitkuAPI interface {
	CreateInvoice(ctx context.Context, payload *CreateInvoicePayload) (*CreateInvoiceResponse, error)
}

type CreateInvoicePayload struct {
	PaymentAmount    int                   `json:"paymentAmount"`
	MerchantOrderID  string                `json:"merchantOrderId"`
	ProductDetails   string                `json:"productDetails"`
	AdditionalParam  string                `json:"additionalParam"`
	MerchantUserInfo string                `json:"merchantUserInfo"`
	CustomerVAName   string                `json:"customerVaName"` // nama yg akan muncul pada halaman konfirmasi pembayaran
	PhoneNumber      string                `json:"phoneNumber"`
	ItemDetails      []CreateInvoiceItem   `json:"itemDetails"`
	CustomerDetail   CreateInvoiceCustomer `json:"customerDetail"`
	CallbackURL      string                `json:"callbackUrl"`
	ReturnURL        string                `json:"returnUrl"`
	ExpiryPeriod     int                   `json:"expiryPeriod"`
}

type CreateInvoiceCustomer struct {
	FirstName       string                       `json:"firstName"`
	LastName        string                       `json:"lastName"`
	Email           string                       `json:"email"`
	PhoneNumber     string                       `json:"phoneNumber"`
	BillingAddress  CreateInvoiceCustomerAddress `json:"billingAddress"`
	ShippingAddress CreateInvoiceCustomerAddress `json:"shippingAddress"`
}

type CreateInvoiceCustomerAddress struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postalCode"`
	Phone       string `json:"phone"`
	CountryCode string `json:"countryCode"`
}

type CreateInvoiceItem struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type CreateInvoiceResponse struct {
	MerchantCode  string `json:"merchantCode"`
	Reference     string `json:"reference"` // need to be saved to the database
	PaymentURL    string `json:"paymentUrl"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

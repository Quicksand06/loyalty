package domain

import "time"

type IdentifierType string

const (
	IdentifierTypeLoyaltyCard  IdentifierType = "loyalty_card"
	IdentifierTypeMembershipID IdentifierType = "membership_id"
	IdentifierTypeEmailAddress IdentifierType = "email_address"
)

func (e IdentifierType) IsValid() bool {
	switch e {
	case IdentifierTypeLoyaltyCard, IdentifierTypeMembershipID, IdentifierTypeEmailAddress:
		return true
	}
	return false
}

type Customer struct {
	CustomerID     string
	IdentifierType IdentifierType
}

type Transaction struct {
	TransactionID string
	Amount        float64
	StoreID       string
	Timestamp     time.Time
	CustomerID    string
}

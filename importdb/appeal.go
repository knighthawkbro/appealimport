package importdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Appeal is the main table that links to most of the others
type Appeal struct {
	gorm.Model

	AppealID       string
	Appellant      Person
	IsOnline       bool
	UNDPRU         bool
	NoticeDate     time.Time
	EnrollmentYear string
	DateEntered    time.Time
	DateReceived   time.Time
	Comments       string `gorm:"size:MAX"`
	Reason         Reason
	IDR            IDR
	Outreach       []Outreach
	Action         Action
	NoticeComments string `gorm:"size:MAX"`
}

// Reason struct is the initial reason an appellant may apply for an appeal
type Reason struct {
	ID       uint
	AppealID uint

	Income         bool
	FamilySize     bool
	PremiumWavier  bool
	Residency      bool
	LawfulPresence bool
	OtherInsurance bool
	Incarceration  bool
	Other          bool
	OtherReason    string
}

// IDR is the internal dispute resolution the Appeals team takes when dealing with a case.
type IDR struct {
	ID       uint
	AppealID uint

	AidPendingAppplied bool
	AidPendingRemoved  bool
	InternalReview     bool
	CaseHolder         string
	Expedite           bool
	BusinessEvent      string
	Issue              Issue
}

// TableName makes sure IDR table is called 'idr' and not 'id_rs'
func (IDR) TableName() string {
	return "idr"
}

// Issue table is similar to the reason table because this is the official issue reason determined by the appeals team
type Issue struct {
	ID    uint
	IDRID uint `gorm:"column:idr_id"`

	Income              bool
	FamilySize          bool
	PWDenial            bool
	PublicMEC           bool
	PublicMECText       string
	Residency           bool
	LawfulPresence      bool
	Other               bool
	OtherReason         string
	SEP                 bool
	ESI                 bool
	IncarcerationStatus bool
	TaxFilling          bool
}

// Action table is for
type Action struct {
	ID       uint
	AppealID uint

	DismissedInvalid string
	Dismissed        string
	WDDate           time.Time
	FinalLeterSent   time.Time
	LetterSentBy     string
	Hearing          []Hearing
	Accessibility    Accessibility
}

// Hearing table is for
type Hearing struct {
	ID       uint
	ActionID uint

	Action         string
	HearingDate    time.Time
	HearingOfficer string
	CCARep         string
}

// Accessibility table is for
type Accessibility struct {
	ID       uint
	ActionID uint

	Interpreter   string
	Device        string
	Accommodation string
}

// Outreach table is for
type Outreach struct {
	gorm.Model
	AppealID uint

	Notes         string `gorm:"size:MAX"`
	ContactMethod string
	ContactMade   bool
	Outcome       string
	TimeSpent     uint
}

// Rep table is for
type Rep struct {
	ID       uint
	PersonID uint

	FirstName string
	LastName  string
	Address   Address `gorm:"polymorphic:resident"`
}

// Person table is for
type Person struct {
	ID       uint
	AppealID uint

	FirstName string
	LastName  string
	DOB       string
	MemberID  string
	Email     string
	Phone     string

	Rep     Rep
	Address Address `gorm:"polymorphic:resident"`
}

// Address table is for
type Address struct {
	ID uint

	Street string
	City   string
	State  string
	Zip    string

	ResidentID   uint
	ResidentType string
}

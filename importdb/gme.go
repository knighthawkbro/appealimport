package importdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
0,				1,		2,
"CC_Person_ID",	"Type",	"FNAME",
3,			4,		5,
"LNAME",	"DOB",	"MemberID",
5,			6,		7,
"Address",	"City",	"State",
8,		9,			10,
"Zip",	"Phone",	"Email",
11,				12,			13,
"Appeals_ID",	"isOnline",	"Notice_Date",
14,				15,				16,
"REP_FNAME",	"REP_LNAME",	"REP_Address",
17,				18,				19,
"REP_City",		"REP_State",	"REP_Zip",
20,				21,					22,
"REP_Phone",	"Notice_Number",	"Date_Entered",
23,					24,			25,
"Date_Recieved",	"Income",	"AppealReason",
26,					27,			28,
"Business_Event",	"Expedite",	"IssueGeneratingAppealER",
29,							30,					31,
"IssueGeneratingAppealEE",	"Review_Outreach",	"Hearing_Date",
32,			33,					34,
"Hearing",	"Outreach_Notes",	"Dismissed_Invalid",
35,				36,				37,
"Dismissed",	"Pending_info",	"WD_date",
38,							39,					40,
"Date_Final_Leter_Sent",	"Letter_Sent_by",	"Internal_review",
41,						42,					43,
"Internal_Case_holder",	"Enrollment_Year",	"Death",
44,				45,					46,
"Hearing_time",	"Hearing_Officer",	"Assistive_InterpreterLanguage",
47,					48,								49,
"Assistive_Device",	"Assistive_AccomDisability",	"Connector_Representative",	"Hearing_action",
50,				51,					52,
"Re-Hearing",	"Re-hearing_Date",	"Re-hearing_time",
53,						54,							55,
"Re-Hearing_officer",	"Re-Hearing_Connector_rep",	"Re-Hearing_action",
56,			57,		58,
"QSHIP",	"NPP",	"Medicare",
59,					60,			61
"Merge_Comments",	"Comments",	"UND-PRU" */

// GME (Public) -
type GME struct {
	gorm.Model

	Appellant      GMEPerson
	IsOnline       bool
	UNDPRU         bool
	NoticeDate     time.Time
	EnrollmentYear string
	DateEntered    time.Time
	DateReceived   time.Time
	Comments       string `gorm:"size:MAX"`
	Reason         GMEReason
	IDR            GMEIDR
	Outreach       []GMEOutreach
	Action         GMEAction
	NoticeComments string `gorm:"size:MAX"`
}

// GMEPerson (Public) -
type GMEPerson struct {
}

// GMEReason (Public) -
type GMEReason struct {
}

// GMEIDR (Public) -
type GMEIDR struct {
}

// GMEOutreach (Public) -
type GMEOutreach struct {
}

// GMEAction (Public) -
type GMEAction struct {
}

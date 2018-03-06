package importdb

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// LoadData (Public) - Reads a datafile and loads it into a variable.
func LoadData() []Appeal {
	file, err := os.Open("appeals/tblApealsData.txt")
	if err != nil {
		return nil
	}
	var result []Appeal
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'
	layout := "1/2/2006 15:04:05"
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		if line[0] == "CC_Person_ID" {
			continue
		}

		var appeal Appeal
		id, err := strconv.Atoi(line[0])
		checkErr(err)
		appeal = Appeal{
			Model:          gorm.Model{ID: uint(id)},
			EnrollmentYear: line[63], Comments: line[82],
			NoticeComments: line[81],
		}
		if line[63] != "" {
			appeal.AppealID = fmt.Sprintf("ACA%v-%v", appeal.EnrollmentYear[2:], appeal.ID)
		}
		appeal.IsOnline, err = strconv.ParseBool(line[12])
		checkErr(err)
		appeal.UNDPRU, err = strconv.ParseBool(line[83])
		checkErr(err)
		if len(line[21]) > 0 {
			appeal.DateEntered, err = time.Parse(layout, line[21])
			checkErr(err)
		}
		appeal.DateReceived, err = time.Parse(layout, line[22])
		checkErr(err)
		appeal.Appellant = Person{
			FirstName: line[1], LastName: line[2], DOB: line[3],
			MemberID: line[4], Email: line[10], Phone: line[9],
			Rep: Rep{
				FirstName: line[14], LastName: line[15],
				Address: Address{
					Street: line[16], City: line[17], State: line[18],
					Zip: line[19],
				},
			},
			Address: Address{
				Street: line[5], City: line[6], State: line[7],
				Zip: line[8],
			},
		}
		appeal.Reason.Income, err = strconv.ParseBool(line[23])
		checkErr(err)
		appeal.Reason.FamilySize, err = strconv.ParseBool(line[24])
		checkErr(err)
		appeal.Reason.PremiumWavier, err = strconv.ParseBool(line[25])
		checkErr(err)
		appeal.Reason.Residency, err = strconv.ParseBool(line[26])
		checkErr(err)
		appeal.Reason.LawfulPresence, err = strconv.ParseBool(line[27])
		checkErr(err)
		appeal.Reason.OtherInsurance, err = strconv.ParseBool(line[28])
		checkErr(err)
		appeal.Reason.Incarceration, err = strconv.ParseBool(line[29])
		checkErr(err)
		appeal.Reason.Other, err = strconv.ParseBool(line[30])
		checkErr(err)
		appeal.Reason.OtherReason = line[31]

		appeal.IDR.AidPendingAppplied, err = strconv.ParseBool(line[47])
		checkErr(err)
		appeal.IDR.AidPendingRemoved, err = strconv.ParseBool(line[49])
		checkErr(err)
		appeal.IDR.InternalReview, err = strconv.ParseBool(line[61])
		checkErr(err)
		appeal.IDR.CaseHolder = line[62]
		appeal.IDR.Expedite, err = strconv.ParseBool(line[46])
		checkErr(err)
		appeal.IDR.BusinessEvent = line[32]

		appeal.IDR.Issue.Income, err = strconv.ParseBool(line[33])
		checkErr(err)
		appeal.IDR.Issue.FamilySize, err = strconv.ParseBool(line[38])
		checkErr(err)
		appeal.IDR.Issue.PWDenial, err = strconv.ParseBool(line[40])
		checkErr(err)
		appeal.IDR.Issue.PublicMEC, err = strconv.ParseBool(line[43])
		checkErr(err)
		appeal.IDR.Issue.PublicMECText = line[44]
		appeal.IDR.Issue.Residency, err = strconv.ParseBool(line[34])
		checkErr(err)
		appeal.IDR.Issue.LawfulPresence, err = strconv.ParseBool(line[37])
		checkErr(err)
		appeal.IDR.Issue.Other, err = strconv.ParseBool(line[41])
		checkErr(err)
		appeal.IDR.Issue.OtherReason = line[42]
		appeal.IDR.Issue.SEP, err = strconv.ParseBool(line[35])
		checkErr(err)
		appeal.IDR.Issue.ESI, err = strconv.ParseBool(line[36])
		checkErr(err)
		appeal.IDR.Issue.IncarcerationStatus, err = strconv.ParseBool(line[39])
		checkErr(err)
		appeal.IDR.Issue.TaxFilling, err = strconv.ParseBool(line[45])
		checkErr(err)

		appeal.Action.DismissedInvalid = line[65]
		appeal.Action.Dismissed = line[66]
		if len(line[58]) > 0 {
			if strings.Contains(line[58], ".") {
				appeal.Action.WDDate, err = time.Parse("3.4.2006", line[58])
			} else {
				appeal.Action.WDDate, err = time.Parse("3/4/2006", line[58])
			}
			checkErr(err)
		}
		if len(line[59]) > 0 {
			appeal.Action.FinalLeterSent, err = time.Parse(layout, line[59])
			checkErr(err)
		}
		appeal.Action.LetterSentBy = line[60]
		hearing, err := strconv.ParseBool(line[53])
		checkErr(err)
		if hearing {
			appeal.Action.Hearing = append(appeal.Action.Hearing, Hearing{
				Action: line[71], HearingOfficer: line[66], CCARep: line[70],
			})
			if len(line[52]) > 0 || len(line[65]) > 0 {
				d := strings.Split(line[52], " ")
				t := strings.Split(line[65], " ")
				hd, err := time.Parse(layout, d[0]+" "+t[1])
				checkErr(err)
				appeal.Action.Hearing[0].HearingDate = hd
			}
		}
		rehearing, err := strconv.ParseBool(line[72])
		checkErr(err)
		if rehearing {
			appeal.Action.Hearing = append(appeal.Action.Hearing, Hearing{
				Action: line[77], HearingOfficer: line[75], CCARep: line[76],
			})
			if len(line[73]) > 0 || len(line[74]) > 0 {
				d := strings.Split(line[73], " ")
				t := strings.Split(line[74], " ")
				hd, err := time.Parse(layout, d[0]+" "+t[1])
				checkErr(err)
				if len(appeal.Action.Hearing) > 1 {
					appeal.Action.Hearing[1].HearingDate = hd
				} else {
					appeal.Action.Hearing[0].HearingDate = hd
				}
			}
		}
		appeal.Action.Accessibility = Accessibility{
			Interpreter: line[67], Device: line[68], Accommodation: line[69],
		}
		result = append(result, appeal)
	}

	return result
}

// LoadOutreach (Public) -
func LoadOutreach() []Outreach {
	var result []Outreach
	file, err := os.Open("appeals/tblOutreach.txt")
	checkErr(err)
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'
	layout := "1/2/2006 15:04:05"
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		if line[0] == "Outreach_ID" {
			continue
		}
		outreachID, err := strconv.Atoi(line[0])
		checkErr(err)
		appealID, err := strconv.Atoi(line[1])
		checkErr(err)
		timeSpent, err := strconv.Atoi(line[6])
		checkErr(err)

		outreach := Outreach{
			Model:    gorm.Model{ID: uint(outreachID)},
			AppealID: uint(appealID), Notes: line[2],
			ContactMethod: line[3], Outcome: line[5],
			TimeSpent: uint(timeSpent),
		}
		outreach.ContactMade, err = strconv.ParseBool(line[4])
		checkErr(err)
		if len(line[7]) > 0 {
			outreach.CreatedAt, err = time.Parse(layout, line[7])
			checkErr(err)
		}
		result = append(result, outreach)
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

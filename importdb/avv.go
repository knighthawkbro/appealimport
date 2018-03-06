package importdb

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Need to have SFTP connections to Optum,
// Download Production Data to server
// Transform data from CSV to Struct
// Return Struct

// SFTPConf (Public) -
type SFTPConf struct {
	Address  string
	User     string
	Password string
	Folder   string
}

func AppealsExtract(c SFTPConf) error {
	conf := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		Config: ssh.Config{
			Ciphers: []string{"3des-cbc", "aes256-cbc", "aes192-cbc", "aes128-cbc"},
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", c.Address, conf)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		panic("Failed to create client: " + err.Error())
	}
	defer client.Close()
	files, err := client.ReadDir(c.Folder)
	if err != nil {
		panic(err.Error())
	}
	if len(files) != 0 {
		for _, file := range files {
			fmt.Printf("%v%v\n", c.Folder, file.Name())
			r, err := client.Open(fmt.Sprintf("%v%v", c.Folder, file.Name()))
			if err != nil {
				panic(err)
			}
			defer r.Close()
			f, err := os.Create("./appeals" + file.Name())
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = r.WriteTo(f)
			if err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("No files found.")
	}

	return nil
}

func AppealsTransform() []Appeal {
	files, err := ioutil.ReadDir("appeals")
	if err != nil {
		panic(err)
	}
	var result []Appeal
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("appeals/%v", file.Name()))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		reader := csv.NewReader(bufio.NewReader(f))
		reader.Comma = '|'
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			var appeal Appeal
			appeal.Appellant.FirstName = line[0]
			appeal.Appellant.LastName = line[1]
			DOB, err := time.Parse("2006-01-02", line[2])
			if err != nil {
				panic(err)
			}
			appeal.Appellant.DOB = DOB.Format("1/2/2006 3:04:05")
			appeal.Appellant.MemberID = line[3]
			appeal.Appellant.Address.Street = line[4]
			appeal.Appellant.Address.City = line[5]
			appeal.Appellant.Address.State = line[6]
			appeal.Appellant.Address.Zip = "0" + line[7]
			appeal.Appellant.Phone = line[8]
			appeal.Appellant.Email = line[9]

			appeal.Appellant.Rep.FirstName = line[10]
			appeal.Appellant.Rep.LastName = line[11]
			appeal.Appellant.Rep.Address.Street = line[12]
			appeal.Appellant.Rep.Address.City = line[13]
			appeal.Appellant.Rep.Address.State = line[14]
			appeal.Appellant.Rep.Address.Zip = line[15]

			if line[16] == "Y" {
				appeal.Reason.Income = true
			}
			if line[17] == "Y" {
				appeal.Reason.Residency = true
			}
			if line[18] == "Y" {
				appeal.Reason.OtherInsurance = true
			}
			if line[19] == "Y" {
				appeal.Reason.FamilySize = true
			}
			if line[20] == "Y" {
				appeal.Reason.LawfulPresence = true
			}
			if line[21] == "Y" {
				appeal.Reason.Incarceration = true
			}
			if line[22] == "Y" {
				appeal.Reason.PremiumWavier = true
			}
			if line[23] == "Y" {
				appeal.Reason.Other = true
				appeal.Reason.OtherReason = line[24]
			}

			// Skip line 25 and 26 for some reason
			var additionalComments string
			if line[27] == "1" {
				additionalComments += " - Interpreter Required - Language: " + line[28]
				appeal.Action.Accessibility.Interpreter = line[28]
			}
			if line[29] == "1" {
				additionalComments += " - Assist Device Required - Device: " + line[30]
				appeal.Action.Accessibility.Device = line[30]
			}
			if line[31] == "1" {
				additionalComments += " - Other Accomodation Required - Accomodation: " + line[32]
				appeal.Action.Accessibility.Accommodation = line[32]
			}

			appeal.Comments = line[33] + additionalComments
			// 2018-02-12 16:51:58.037
			appeal.DateReceived, err = time.Parse("2006-01-02 15:04:05.000", line[34])
			if err != nil {
				panic(err)
			}

			appeal.DateEntered = time.Now()
			appeal.IsOnline = true

			//appeal.Action.WDDate, err = time.Parse("3/4/2006", line[58])
			result = append(result, appeal)
		}
	}
	return result
}

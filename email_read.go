package emailsorter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type ImapConfig struct {
	ImapHost string `json:"imap_host"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImapPort int    `json:"imap_port"`
}

type MsgHeader struct {
	Subject  string
	SentDate time.Time
	Sender   string
}

func ReadConfig(filename string) (ImapConfig, error) {

	var imapconfig ImapConfig

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return imapconfig, err
	}

	err = json.Unmarshal(data, &imapconfig)

	if err != nil {
		return imapconfig, err
	}

	return imapconfig, nil

}

func deleteEmail(c *client.Client, uids []uint32) error {
	seqset := new(imap.SeqSet)

	for _, num := range uids {
		seqset.AddNum(num)
	}
	err := c.Store(seqset, "+FLAGS.SILENT", []interface{}{imap.DeletedFlag}, nil)
	if err != nil {
		return err
	}
	err = c.Expunge(nil)
	if err != nil {
		return err
	}
	return nil
}

func NewImapClient(imapconfig ImapConfig) (*client.Client, error) {
	imapclient, err := client.DialTLS(fmt.Sprintf("%s:%d",
		imapconfig.ImapHost,
		imapconfig.ImapPort),
		nil)
	if err != nil {
		return nil, err
	}
	err = imapclient.Login(imapconfig.Email, imapconfig.Password)
	if err != nil {
		return imapclient, err
	}
	return imapclient, nil

}

func copyTo(imapclient *client.Client, uids []uint32, folder string) error {
	seqset := new(imap.SeqSet)

	for _, num := range uids {
		seqset.AddNum(num)
	}
	return imapclient.Copy(seqset, folder)
}

func GetAttachmentNamesAndDownload(imapclient *client.Client,uids []uint32,destDir string) ([]string,error) {
	var attachmentNames []string
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}
    if len(uids) == 0 {
		return nil, errors.New("uids slice of zero length")
	}

	seqset := new(imap.SeqSet)

	for _, num := range uids {
		seqset.AddNum(num)
	}

	messages := make(chan *imap.Message)
	done := make(chan error, 1)

	go func() {

		done <- imapclient.Fetch(seqset, items, messages)
	}()

	

	for msg := range messages {

		r := msg.GetBody(&section)
		
		if r== nil {
			return nil,errors.New("Unable to open the message body")
		}

		mr,err := mail.CreateReader(r)

		if err != nil {
			return nil,err
		}

		for {
			p,err := mr.NextPart()

			if err == io.EOF {
				break
			}

			if err != nil && err != io.EOF {
				return nil,err
			}

			switch h := p.Header.(type) {
			case *mail.AttachmentHeader:
				filename,_ := h.Filename()
				attachmentNames = append(attachmentNames, filename)
				destFileName := destDir +"/"+filename
				w,err := os.Create(destFileName)
				if err != nil {
					return nil,err
				}
				_,err = io.Copy(w,p.Body)
				if err != nil {
					return nil,err
				}
				err = w.Close()
				if err != nil {
					return nil,err
				}
				
			default:
				continue
			}


		}
		
	}
	return attachmentNames,nil
}

func GetMessageHeaders(imapclient *client.Client, uids []uint32) ([]MsgHeader, error) {
	var messageHeaders []MsgHeader

	if len(uids) == 0 {
		return nil, errors.New("uids slice of zero length")
	}

	seqset := new(imap.SeqSet)

	for _, num := range uids {
		seqset.AddNum(num)
	}

	messages := make(chan *imap.Message)
	done := make(chan error, 1)

	go func() {

		done <- imapclient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for msg := range messages {

		msghdr := MsgHeader{
			Subject:  msg.Envelope.Subject,
			SentDate: msg.Envelope.Date,
			Sender:   msg.Envelope.From[0].MailboxName,
		}

		messageHeaders = append(messageHeaders, msghdr)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return messageHeaders, nil

}

func GetEmailUIDs(imapclient *client.Client, params CmdParams) ([]uint32, error) {
	criteria := imap.NewSearchCriteria()
	if len(params.From) != 0 {
		criteria.Header.Add("FROM", params.From)
	}

	if len(params.Subject) != 0 {
		criteria.Header.Add("SUBJECT", params.Subject)
	}

	if len(params.Since) != 0 {
		from_date_filter, err := time.Parse("2006-01-02", params.Since)
		if err != nil {
			return nil, err
		}
		year, month, day := from_date_filter.Date()
		filterFrom_Date := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
		criteria.Since = filterFrom_Date
		
	}

	if len(params.Before) != 0 {
		before_date_filter,err:= time.Parse("2006-01-02", params.Before)
		if err != nil {
			return nil, err
		}
		year, month, day := before_date_filter.Date()
		filterFrom_Date := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
		criteria.Before = filterFrom_Date 
	}

	mbox, err := imapclient.Select(params.SrcFolder, false)

	if err != nil {
		return nil, err
	}
	_ = mbox

	uids, err := imapclient.Search(criteria)
	if err != nil {

		return nil, err
	}

	return uids, nil

}

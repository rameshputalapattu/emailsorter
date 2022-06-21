package emailsorter

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

const attachmentsHelp = `list attachment names and download them in the destination directory for the mails in the source folder meeting search criteria`

func (cmd *AttachmentsCommand) Name() string      { return "list_and_download_attachments" }
func (cmd *AttachmentsCommand) Args() string      { return "" }
func (cmd *AttachmentsCommand) ShortHelp() string { return attachmentsHelp }
func (cmd *AttachmentsCommand) LongHelp() string  { return attachmentsHelp }
func (cmd *AttachmentsCommand) Hidden() bool      { return false }

func (cmd *AttachmentsCommand) Register(fs *flag.FlagSet) {}

type AttachmentsCommand struct {
	Params     *CmdParams
	Imapconfig *ImapConfig
}

func (cmd *AttachmentsCommand) Run(ctx context.Context, args []string) error {

	client, err := NewImapClient(*cmd.Imapconfig)

	if err != nil {

		return err
	}
	uids, err := GetEmailUIDs(client, *cmd.Params)
	if err != nil {

		return err
	}

	if len(uids) == 0 {
		return errors.New("no emails found with the search criteria")
	}

	attachmentsLists,err := GetAttachmentNamesAndDownload(client,uids,cmd.Params.DestDirectory)

	if err != nil {
		return err
	}
	for _,attachName := range attachmentsLists {
		fmt.Println(attachName)
	}

	return nil
}

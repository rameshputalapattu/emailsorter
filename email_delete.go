package emailsorter

import (
	"context"
	"errors"
	"flag"
)

const deleteHelp = `delete emails in the source folder meeting search criteria`

func (cmd *DeleteCommand) Name() string      { return "delete" }
func (cmd *DeleteCommand) Args() string      { return "" }
func (cmd *DeleteCommand) ShortHelp() string { return deleteHelp }
func (cmd *DeleteCommand) LongHelp() string  { return deleteHelp }
func (cmd *DeleteCommand) Hidden() bool      { return false }

func (cmd *DeleteCommand) Register(fs *flag.FlagSet) {}

type DeleteCommand struct {
	Params     *CmdParams
	Imapconfig *ImapConfig
}

func (cmd *DeleteCommand) Run(ctx context.Context, args []string) error {

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

	err = deleteEmail(client, uids)

	return err
}

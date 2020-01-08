package emailsorter

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

const showHelp = `print the subjects to the console`

func (cmd *ShowCommand) Name() string      { return "show" }
func (cmd *ShowCommand) Args() string      { return "" }
func (cmd *ShowCommand) ShortHelp() string { return showHelp }
func (cmd *ShowCommand) LongHelp() string  { return showHelp }
func (cmd *ShowCommand) Hidden() bool      { return false }

func (cmd *ShowCommand) Register(fs *flag.FlagSet) {

}

type ShowCommand struct {
	Params     *CmdParams
	Imapconfig *ImapConfig
}

func (cmd *ShowCommand) Run(ctx context.Context, args []string) error {

	client, err := NewImapClient(*cmd.Imapconfig)

	if err != nil {
		return err
	}

	uids, err := GetEmailUIDs(client, *cmd.Params)
	if err != nil {
		return err
	}

	if len(uids) == 0 {
		return errors.New("No emails retrieved for the search criteria")
	}

	msgheaders, err := GetMessageHeaders(client, uids)

	if err != nil {
		return err
	}

	for _, msg := range msgheaders {
		fmt.Printf("%s|%s|%s\n", msg.Sender, msg.SentDate.Format("2006-01-02"), msg.Subject)
	}

	return nil
}

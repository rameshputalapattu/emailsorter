package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

const showHelp = `print the subjects to the console`

func (cmd *showCommand) Name() string      { return "show" }
func (cmd *showCommand) Args() string      { return "" }
func (cmd *showCommand) ShortHelp() string { return showHelp }
func (cmd *showCommand) LongHelp() string  { return showHelp }
func (cmd *showCommand) Hidden() bool      { return false }

func (cmd *showCommand) Register(fs *flag.FlagSet) {

}

type showCommand struct {
	params     *cmdParams
	imapconfig *ImapConfig
}

func (cmd *showCommand) Run(ctx context.Context, args []string) error {

	client, err := NewImapClient(*cmd.imapconfig)

	if err != nil {
		return err
	}

	uids, err := GetEmailUIDs(client, *cmd.params)
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

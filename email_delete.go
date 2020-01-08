package main

import (
	"context"
	"errors"
	"flag"
)

const deleteHelp = `delete emails in the source folder meeting search criteria`

func (cmd *deleteCommand) Name() string      { return "delete" }
func (cmd *deleteCommand) Args() string      { return "" }
func (cmd *deleteCommand) ShortHelp() string { return deleteHelp }
func (cmd *deleteCommand) LongHelp() string  { return deleteHelp }
func (cmd *deleteCommand) Hidden() bool      { return false }

func (cmd *deleteCommand) Register(fs *flag.FlagSet) {}

type deleteCommand struct {
	params     *cmdParams
	imapconfig *ImapConfig
}

func (cmd *deleteCommand) Run(ctx context.Context, args []string) error {

	client, err := NewImapClient(*cmd.imapconfig)

	if err != nil {

		return err
	}
	uids, err := GetEmailUIDs(client, *cmd.params)
	if err != nil {

		return err
	}

	if len(uids) == 0 {
		return errors.New("no emails found with the search criteria")
	}

	err = deleteEmail(client, uids)

	return err
}

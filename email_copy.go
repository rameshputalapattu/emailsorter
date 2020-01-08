package main

import (
	"context"
	"errors"
	"flag"
)

const copyHelp = `copy emails to the destfolder`

func (cmd *copyCommand) Name() string      { return "copy" }
func (cmd *copyCommand) Args() string      { return "" }
func (cmd *copyCommand) ShortHelp() string { return copyHelp }
func (cmd *copyCommand) LongHelp() string  { return copyHelp }
func (cmd *copyCommand) Hidden() bool      { return false }

func (cmd *copyCommand) Register(fs *flag.FlagSet) {

}

type copyCommand struct {
	params     *cmdParams
	imapconfig *ImapConfig
}

func (cmd *copyCommand) Run(ctx context.Context, args []string) error {

	client, err := NewImapClient(*cmd.imapconfig)

	if err != nil {
		return err
	}

	if len(cmd.params.DestFolder) == 0 {
		return errors.New("destfolder should be provided for copy")
	}

	uids, err := GetEmailUIDs(client, *cmd.params)
	if err != nil {
		return err
	}

	err = copyTo(client, uids, cmd.params.DestFolder)
	if err != nil {
		return err
	}

	err = deleteEmail(client, uids)
	return err
}

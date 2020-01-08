package emailsorter

import (
	"context"
	"errors"
	"flag"
)

const copyHelp = `copy emails to the destfolder`

func (cmd *CopyCommand) Name() string      { return "copy" }
func (cmd *CopyCommand) Args() string      { return "" }
func (cmd *CopyCommand) ShortHelp() string { return copyHelp }
func (cmd *CopyCommand) LongHelp() string  { return copyHelp }
func (cmd *CopyCommand) Hidden() bool      { return false }

func (cmd *CopyCommand) Register(fs *flag.FlagSet) {

}

type CopyCommand struct {
	Params     *CmdParams
	Imapconfig *ImapConfig
}

func (cmd *CopyCommand) Run(ctx context.Context, args []string) error {

	client, err := NewImapClient(*cmd.Imapconfig)

	if err != nil {
		return err
	}

	if len(cmd.Params.DestFolder) == 0 {
		return errors.New("destfolder should be provided for copy")
	}

	uids, err := GetEmailUIDs(client, *cmd.Params)
	if err != nil {
		return err
	}

	err = copyTo(client, uids, cmd.Params.DestFolder)
	if err != nil {
		return err
	}

	err = deleteEmail(client, uids)
	return err
}

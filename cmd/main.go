package main

import (
	"context"
	"errors"
	"flag"
	"time"
        "github.com/rameshputalapattu/emailsorter"
	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)


func main() {
	p := cli.NewProgram()
	p.Name = "emailsorter"
	p.Description = "email sorter bot"

	params := emailsorter.CmdParams{}
	var imapconfig emailsorter.ImapConfig

	p.Commands = []cli.Command{
		&emailsorter.CopyCommand{&params, &imapconfig},
		&emailsorter.DeleteCommand{&params, &imapconfig},
		&emailsorter.ShowCommand{&params, &imapconfig},
	}

	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&params.From, "from", "", "from address filter")
	p.FlagSet.StringVar(&params.Subject, "subject", "", "subject filter")
	p.FlagSet.StringVar(&params.Body, "body", "", "email body filter")
	p.FlagSet.StringVar(&params.SrcFolder, "srcfolder", "INBOX", "source folder")
	p.FlagSet.StringVar(&params.DestFolder, "destfolder", "", "destination folder")
	p.FlagSet.StringVar(&params.Since, "since", "", "since")
	p.FlagSet.StringVar(&params.ConfigFile, "config", "email_config.json", "email config file name (including full path)")

	p.Before = func(ctx context.Context) error {

		if len(params.From) == 0 && len(params.Subject) == 0 && len(params.Since) == 0 {
			return errors.New("Atleast one of from or subject or since must be set")
		}

		var err error

		imapconfig, err = emailsorter.ReadConfig(params.ConfigFile)
		if err != nil {
			return err
		}

		if len(params.Since) != 0 {
			_, err := time.Parse("2006-01-02", params.Since)

			if err != nil {
				return err
			}
		}

		return nil

	}

	p.Run()
	logrus.Info("executed the command successfully")

}

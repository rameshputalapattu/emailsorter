# emailsorter

A simple command line email client with limited features to manage your email.

Email credentials and smtp configuration must be stored in a  config file with a default name called email_config.json. (The name can be overridden with --config flag)

## The config file format

{
        "email":"{{email_address}}",
        "password":"{{password}}",
        "imap_host":"{{imap_host}}",
        "imap_port":{{imap_port}}



}

## Usage



Usage: emailsorter <command>

Flags:

  --body        email body filter (default: <none>)
  --config      email config file name (including full path) (default: email_config.json)
  --destfolder  destination folder (default: <none>)
  --from        from address filter (default: <none>)
  --since       since (default: <none>)
  --srcfolder   source folder (default: INBOX)
  --subject     subject filter (default: <none>)

Commands:

  copy     copy emails to the destfolder
  delete   delete emails in the source folder meeting search criteria
  show     print the subjects to the console
  version  Show the version information.

SeedDMS Client
===============

This program is for accessing the SeedDMS RestAPI on the command line.
It currently has a very limited set of operations to

* upload a document
* download a document
* list the documents and subfolders of a folder

There is also a special upload mode which monitors a local directory
and uploads new files as soon as they appear.

Configuration
-------------

The programm can be configured with command line options, but usually
uses a configuration file named .seeddms-client.yaml which must be
located in the user's home directory and must follow the yaml syntax.
This configuration file may
have several main sections, each containing the options for running the
client. The main section name may be either `default` or any other name
without spaces. Each section may have the following keys

  * Url (url of RestAPI, e.g. http://seeddms-domain/restapi/index.php)
  * User: (username to be used)
  * Password: (password or user)
  * ApiKey: (api key instead of User and Password)

and further sub sections for each of the commands `Upload`,
`Autoupload`, `Ls`.  Command sections for `Upload` and `Autoupload`
may have the keys:

 * Folder
 * Comment

The command section for `Ls` may have only the key

 * Folder

Depending on which command is called. Any parameter in the
configuration can be overriden by a command line option.

  * `--folder` overrides `Folder`
  * `--comment` overrides `Comment`

The main section used for running the client can be selected with the
command line option `--profile` or by setting the environment variable
`SEEDDMS_CLIENT`.  If it is not specified either way, the main section
`default` is used. If that is not configured, then all parameters must
be passed on the command line or seeddms-clients quits with an error
message.

Auto uploading
---------------

This mode is for automatically uploading files into SeedDMS once they
have been created in a local directory. The initial idea was to
upload scanned documents or screenshots without further interaction.
You just need to setup your scanner or screen shot program to store
each new file into a predifined directory and monitor this directory
with `seeddms-client`. Once it is uploaded the file is deleted locally
unless --keep=true was passed.

This mode is not for syncing a local directory with a folder in SeedDMS. 

Adding desktop entries
-----------------------

A very convenient way to upload files into SeedDMS when using a grafical
user interface is by creating a desktop entry and placing it into
the directory `~/.local/share/applications`. seeddms-client already
comes with a desktop entry being installed in
/usr/share/applications/org.seeddms.client.upload.desktop
This could be used as a template for your own more specialized uploads.

  [Desktop Entry]
  Name=SeedDMS Client (default profile)
  Comment=Client for uploading files into SeedDMS using the default profile
  Categories=Network;System;
  Type=Application
  Terminal=true
  Icon=seeddms
  Exec=seeddms-client upload %f

If for example you would like to upload your invoices into a dedicated
folder with id=4711, then just make a copy of the above entry into e.g.
`~/.local/share/applications/org.seeddms.client.upload.invoices.desktop`
Than change the name of the new entry to `SeedDMS Client (invoices)`
and the command to `seeddms-client upload --folder=4711 %f` or create
another section in your `~/.seeddms-client.yaml` titled `invoices` and
change the command to `seeddms-client --profile=invoices upload %f`


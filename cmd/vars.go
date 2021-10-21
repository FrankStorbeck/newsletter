package main

const man = `
SYNOPSIS
%s [-h | -usage | -version | -test] [-auth authfile] [-recipients recipientsfile] \
          [-selectors selfile] [-subject subject] [-max maxRcpnts] \
          [-skip skipRcpnts] mailtemplate [alternativetemplate]

DESCRIPTION
The program sends a newsletter to a number of recipients. These recipients are
read from a '.cvs' file where the fields are separated by semi-colons (';'). The
default file used is 'model/recipients.csv'. It can be changed by providing a
'-recipients' flag. Everything on a line in this file starting from a hash
character ('#') will be disregarded. When a line is empty or only contain white
space characters it will be completely disregarded.

Note: recipients must have a valid e-mail address in a field (collumn) named
'EMail'.

So a file might start as follows:

  #ID;ORGANISATIE;NAME;EMAIL_ADDRESS;REMARKS;IS_VOLUNTEER;WANTS_NEWSLETTER
  23317;;F.W. Johnson;fjohnson@gmail.com;;Y:Y
  23328;;S. Janzen;jsuzanjanzen@gmail.com;;N;Y
  23337;Dot Com;Q. Johannson;quinten@johansson.nl;;Y;N

Recipients are selected for recieving a newsletter by defining values that must
be present in the fields of the recipients file. These are read from a file
'model/selectors.txt'. This file must hold a number of lines each holding the
name of the field, then a equal sign ('=') and then the value that the field
must hold for a recipient to be selected. However, when the value is '*', any
value is accepted. The order of the lines must be the same as the order of the
fields in the recipientsfile. Based on the example above it could be:

  id=*
  organisation=*
  name=*
  EMail=*
  notes=*
  volunteer=Y
  newsletter=Y

The file with the selectors can be changed by providing the '-selectors sel'
arguments.
Note: the fieldnames can be used in the templates for replacement by its value.
In the case of the example above a template could have the follown line:

  Email {{.Get "EMail"}} belongs to {{.Get "name"}}.

It will generate the following result:

  Email fjohnson@gmail.com belongs to F.W. Johnson.

By providing a '-skip' flag the indicated number of selected recipients will be
skipped before starting sending the newsletters.

By providing a '-max' flag no more then the indicated number of newsletters
will be sent.

By using the '-skip' and '-max' flags one can send the newsletters in batches.

When the '-test' flag is provided, newsletters will only be sent to a number of
selected recipients.

To gain access to some SMTP server to send the newsletters the program reads
an authorisation file. It must contain a number of values to access the
server. Each line should contain a keyword ("from", "hostname", "port",
"password", "username"), then a colon followed by an appropiate value.
The default value for the filepath is ".auth.txt". It can be changed by
providing a '-auth' flag. So it should look like:

  hostname: smtp.somehost.com
  port:     587
  username: someusername
  password: somepassword
  from:     noreply@somedomain.com

Note: this file should only be readable by the owner of the program.

The subject for the e-mail holding the newsletter is set to "Newsletter".
It can be changed by providing a '-subject' flag.

By providing a '-usage' flag the program shows this summary about the usage and
quits (alternavely one can provide the '-h' flag to display the use of the
flags).

By providing a '-version' flag the program shows a the version number and quits.

The program expects one or two arguments holding the paths to files with the
templates for construction of the newsletter. A path ending with a ".txt"
extension will be used to construct a plain text version of the newsletter.
When the path ends with a ".html" extension a HTML version will be
constructed. Mails can hold both newsletters. Most modern e-mail clients will
then only display the HTML version. Oldfashioned e-mail clients display the
plain version.
Note: when two paths have the same extension, only the last path will be
used!

When in these files some text in the form {{.Get "fieldname"}} is found, it will
be replaced by the value in the field named 'fieldname' in the file with the
recipients. So {{.Get "EMail"}} will be replaced by the e-mail address of the
recipient.

EXIT STATUS
The program exits 0 on success, and > 0 if an error occurs.
`

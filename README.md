## NAME
newsletter – send a newsletter to a (sub)group of subscribers.

## SYNOPSIS
newsletter [`-auth` _authfile_] [`-dry`] [`-from` _address_] [`-max` _m_]
[`-quota` _q_] [`-selectors` _selfile_] [`-skip` _s_] [`-subject` _sbjct_]
[`-subscribers` _subfile_] [`-version`] _template_


## DESCRIPTION
Send a formatted newsletter to a group of subscribers based on the template
given in the file _template_. It logs the actions on `stdout` and `stderr`.

The following options are available:

`-auth` _autfile_: read the authorisation data for accessing an SMTP server
  from _authfile_. The default is `.auth.txt`.

`-dry`: do a dry run, i.e. no newsletters will be sent.

`-from` _address_: use _address_ as the sender of the newsletters. This overrides the `sender` value as defined in _authfile_.

`-max` _m_: set the maximum number of newsletters to be sent to _m_.

`-quota` _q_: set the maximum number of newsletters sent per hour to _q_.

`-selectors` _selfile_: read the selection criteria for the sub group that must
  receive the newsletter from _selfile_. the default is `selectors.txt`.

`-skip` _s_: first skip _s_ records in the file with the subscribers before
  starting sending newsletters.

`-subject` _sbjct_: use _sbjct_ as the subject for the email with the newsletter. The default value is `Newsletter`.

`-subscribers` _subfile_: read the data for the subscribers from _subfile_. The default is `subscribers.csv`.

`-version`: show the version number and then exit.

## FILES

Data are read from three files.

### _subfile_
  A `csv` file holding the data about the subscribers. The values in each record
  must be separated by semi-colons (`;`). The first record must hold
  the names of the columns. **Capitals in these names are changed into their
  lower case counterparts. Leading and trailing spaces are removed.**
  Empty lines and lines starting with a `#` characters are skipped.
  The column names can be used in the template to format the newsletter.  
  E.g. if a column name is `name`, then the text element `{{.Get "name"}}` in
  the template file will replaced by the field value. There must be one column
  named _email_ holding a valid email address of the subscribers (if this field
  is empty no email will be sent).

  As an example:
```
  id;first name;middle names;family name;email               ;wants newsletter
  1 ;Bob       ;            ;Crypt      ;bob@crpt.com        ;y
  2 ;John      ;O'          ;Doe        ;john@company.com    ;Yes
  3 ;Joyce     ;O'          ;Doe        ;joyce@doe-family.com;Yes
  4 ;Johny     ;O'          ;Doe        ;johny@doe-family.com;y
  5 ;Karin     ;            ;Sheppard   ;                    ;y
  6 ;Daisy     ;O'          ;Doe        ;daisy.doe-family.com;Y
  7 ;Alice     ;            ;Crypt      ;alice@crpt.com      ;n
```

###_selfile_
  A text file providing the criteria for selecting a sub group from the
  subscribers that must receive the newsletter. Each line must hold the column
  name as used in the subscribers file (_subfile_) followed by an equal sign
  (`=`) and then a regular expression placed between quotation marks (`"`).
  **Capitals for the column name are changed into their lower case
  counterparts. Leading and trailing spaces are removed**. Only if all values
  of the record's fields in the named columns are matched by the regular
  expression, the  newsletter will be sent to that subscriber.

  As an example, to send the newsletter only to all members of the family `O'
  Doe`, but only those who want to receive the newsletter:
```
middle names     = "^O'$"
family name      = "^Doe$"
wants newsletter = "^(Y|y)+(es){0,1}$"
```
### _authfile_:
  A text file providing the authorisation data for accessing an SMTP server.
  Each line must hold a key value, then a colon followed by a value. The
  recognised key values are _hostname_ for the name of the SMTP server, _port_
  for the port number to be used to deliver the email, _username_ /_password_
  for the login credentials, _sender_ for the originator of the email and
  _from_ for the reply address. The contents of this file should be kept secret.

  As an example:
```
hostname: smtp.somehosting.com
port:     587
username: some.user.name@somehosting.com
password: some_secret_password
sender:   some.sender@domain.com
from:     noreply@domain.com
```

### _template_
  A text file holding the template for the newsletter. It can define two
  versions for the newsletter, a plain version and an `HTML` version. The
  program can combine the two versions into one email (the receiving mail
  client is then responsible for displaying one of the versions). The plain
  version can be defined between the text elements `{{define "plainBody"}}` and
  `{{end}}`. For the `HTML` version this is `{{define "htmlBody"}}` and
  `{{end}}`. There is no need to provide both versions. If none of the versions
  is defined the contents of the file will be put unformatted into the body of
  the email.

  As an example:

```
{{define "plainBody"}}
Hello {{.Get "first name"}} {{.Get "middle names"}} {{.Get "family name"}}.
{{end}}

{{define "htmlBody"}}
<!DOCTYPE html>
<html>
  <head>
    <title>Newsletter</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0 " />
    <style type="text/css">
      body{
        font-family: sans-serif;
      }
    </style>
  </head>
  <body>
    <p>
    Hello {{.Get "first name"}} {{.Get "middle names"}} {{.Get "family name"}}.
    </p>
  </body>
</html>
{{end}}
```

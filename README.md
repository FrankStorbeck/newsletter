# newsletter

_newsletter_ is `cli` program to send a newsletter to a (sub)group of
  subscribers.

Usage:

newsletter [-auth _authfile_] [-from _from-address_] [-max _m_] [-quota _q_]  \\ \
           [-recipients _recfile_]  [-selectors _selfile_] [-skip _s_] \\ \
           [-subject _sbjct_] [-version] _template_

 where:

 - _autfile_ is the path to the file with the authorisation data (default
  `./.auth.txt`).
 - _from-address_ is the reply address.
 - _m_ is the maximum number of newsletters to be sent (default `100`).
 - _q_ is the maximum number of newsletters that the SMTP server accepts per
   hour (default unlimited).
 - _recfile_ is the path to the file with the subscribers (default
   `./recipients.csv`).
 - _selfile_ is the path to the file with the selection criteria (default
   `./selectors.txt`).
 - _s_ is the number of entries in the subscriber file that must be skipped
   before starting sending the newsletters (default `0`). This can be used when
   sending the newsletter has been terminated too early for some reason and the
   program has to be restarted.
 - _sbjct_ is the subject for the email holding the newsletter (default `Newsletter`).
 - if the `-version` flag is given, the program shows its version number and
   then exits.
- _template_ is the path to the template file.

It uses a number files to perform its task:

- a csv file holding the data for the subscribers. The column names can be used
  in the template. **The column name must start with a capital.** If the name is
  `Name`, then the text element `{{.Get "Name"}}` in the template will
  replaced by its corresponding field value in the newsletter.
- a text file providing the criteria for selecting a sub group that must receive
  the newsletter. Each line must hold the column name in the recipients file
  followed by a `=` and a value. Only if all the field values are equal to the
  provided values in this file, the newsletter will be sent. The wildcard `*`
  matches all field values.
- a text file providing the authorisation data for accessing the SMTP server.
  Each line must hold a key value, then a colon followed by a value. The key
  values are _hostname_, _port_, _username_, _password_, _sender_ and _from_.
  The content of this file should be kept secret.
- a text file holding the template for the newsletter. It can define two
  versions of the newsletter, a plain version and an `HTML` version. The program
  can combine the two versions into one email with the newsletter (the receiving
  mail client is responsible for displaying one of the versions).
  The plain version must be defined between the text elements
  `{{define "plainBody"}}` and `{{end}}`. For the `HTML` version this is
  `{{define "htmlBody"}}` and `{{end}}`. As an example:

```
{{define "plainBody"}}
  Dear {{.Get "FirstName"}} {{.Get "MiddleNames"}} {{.Get "FamilyName"}},

  hello!
{{end}}

{{define "htmlBody"}}
  <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "https://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
  <html xmlns="https://www.w3.org/1999/xhtml">
    <head>
      <title>Hello</title>
      <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
      <style type="text/css">
        body{
          font-family: sans-serif;
        }
      </style>
    </head>
    <body>
      Dear {{.Get "FirstName"}} {{.Get "MiddleNames"}} {{.Get "FamilyName"}},
      <br>
      hello!
    </body>
  </html>
{{end}}
```

# mp4duration

A command-line program for showing the duration of an MPEG-4 file.  I have
tested it with both .mp4 and .m4a files.

Normal output is of the form hh:mm:ss.  If you specify the -m/--millis flag,
you get hh:mm:ss.mmm (where mmm is millseconds).  If you specify the -t/--total
flag, you get the total, either in seconds (without the --millis flag) or milliseconds
(with the --millis flag).

There are two other flags: -h/--no-filename and -H/--with-filename.  These function
as they do for the grep command.  If you only specify one file, by default the filename
is not printed.  If you specify more than one file, the filenames are printed.  If you
specify both, the filenames are printed (not sure why you would do this, though).

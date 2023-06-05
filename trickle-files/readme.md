# Trickle files

This is a small utility to trickle feed files from one folder to another.
The next file will only be copied when the previous file has been removed from
the destination folder.

The main use case is for systems where dumping lots of files into a "watch folder"
would result in processing all of them in parallel, but this behavior is not adjustable
or must be preserved for parallel processing of other files that are put into the same folder.



GO /dev/null
============

A sample upload / download service written in go

**Use at your own risk**

This application receives files on /upload with a POST and serves those files
under /fetch/<filename>.

Internally the files are stored in a tmp directory and each file lives in it's own pseudo random directory. The user can not access the files but has to provide the directory names.

For ease of use a little html form is provided for uploading.

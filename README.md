GO /dev/null
============

A sample upload / download service written in go

**Use at your own risk**

This application receives files on /upload with a POST and serves those files
under /fetch/< dir name >.

Internally the files are stored in a tmp directory and each file lives in it's own pseudo random directory. The user can not access the files but has to provide the directory names.

For ease of use a little html form is provided for uploading.


Docker - from 0 to hero
========================

This application is being used as a getting started guide to docker containers


Installation *docker for mac*
------------------------------

Installation of *docker for mac*: https://docs.docker.com/docker-for-mac/
Docker for mac is a native tool set for using docker under mac. Do not use docker-toolbox anymore.


Create Dockerfile
-----------------

The Dockerfile tells docker how to build / provision a docker container
 * FROM - Which base image should the container use?
 * ADD - Any files to copy?
 * ENTRYPOINT - Which executable to run?
 
    go build
    docker build . # remember the id
    docker run <id>


** stack trace .... **
 
    env GOOS=linux go build # <- we are running the app under linux!
    docker build . # remember the id
    docker run <id>
    
** success ! **


Networking 
----------

Sofar we do not expose the port to the os. 

    docker run -p 80:8080 <id>
    
```-p <target port>:<origin port>``` will expose the origin port from within the container to the target port

Now we can access the upload service under port 80

** success ! ** 


Volumes
--------

So far our uploaded data is being lost on every restart... 

    docker build -t=go-dev-null .
    docker run -p 80:8080 -v `echo $(pwd)`/tmp:/tmp go-dev-null
    
```-v <absolute/path/to/source>:<target/dir>``` will map the source dir path to target dir inside container

Now restart the application, upload a file, stop the container and restart. The file is not lost!

** success ! **


Container size
--------------

To see all local available images run ```docker images```. On my machine the image is about 100mb big.


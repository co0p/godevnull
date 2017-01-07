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

    FROM google/debian:jessie
    MAINTAINER Julian Godesa <julian.godesa@googlemail.com>
    ADD goDevNull goDevNull
    ENTRYPOINT "goDevNull"

* FROM - Which base image should the container use?
* ADD - Any files to copy?
* ENTRYPOINT - Which executable to run?

And now build and run the container: 

    go build
    docker build . # remember the id
    docker run <id>


** OOOPPPS stack trace .... **
 
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

So far our uploaded data is being lost on every restart. Let's attach a local directory as a volume inside the container

    docker build -t=go-dev-null .
    docker run -p 80:8080 -v `echo $(pwd)`/tmp:/tmp go-dev-null
    
```-v <absolute/path/to/source>:<target/dir>``` will map the source dir path to target dir inside container

Now restart the application, upload a file, stop the container and restart. The file is not lost!

** success ! **

Login to container
------------------

Currently we use entrypoint in the Dockerfile. This starts the goDevNull binary right away and the
output is directly send to the std output of the terminal. 

 adjust the Dockerfile to put the binary under /root
 
    FROM google/debian:jessie
    MAINTAINER Julian Godesa <julian.godesa@googlemail.com>
    ADD goDevNull /root/goDevNull
    WORKDIR "/root" # <---- this is new !!
    ENTRYPOINT "/root/goDevNull"
 
 * build the container ```docker build -t=go-dev-null .```
 * run the container image ```docker run -p 80:8080 -v `echo $(pwd)`/tmp:/root/tmp go-dev-null```
 
See what containers are running: 

    docker ps
    
Now connect to the running container via name and execute the bash:

    docker exec -i -it <name> /bin/bash

```-it``` is short hand for interactive (-i) and getting a new tty (-t)

Yeah, now you are inside the container. Go to /root/ and inspect the content

** success: inception baby! **


Further ideas
--------------

 * put a nginx load balancer in front of multiple instances or use a ha-proxy
 * deploy to a cloud provider
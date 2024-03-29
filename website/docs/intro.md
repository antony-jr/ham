---
title: Introduction
sidebar_position: 1
---

:::info

This is alpha stage software, even though the builder and client program works well to a certain degree, we
still lack recipes for different devices. Mostly legacy devices need recipes since LOS does not support them.

I ask help from the community to create git repositories and follow the HAM recipe syntax to create LineageOS
builds for legacy devices that LOS is supported. No need to port it to the latest LOS version, that is the work
done by LOS community, all we need to do is create recipes for such devices and make it build to a stable LOS
version with good platform security for Android. Vendor security cannot be improved since it's closed source.
You can also port to the latest version but stability over features.

Also recipes for OnePlus Devices to make ROMs which can run on locked bootloader. Until we get atleast 10 recipes
this project is considered a "Work in Progress".

After you create a git repository with a stable ham recipe, request to move the repo to https://github.com/ham-community

:::

HAM (Hetzner Android Make) is a Simple tool written in GO which can build LineageOS (or AOSP) from Source using 
Hetzner Cloud. **Build your Own Flavor of Android Under €1.** (Run Directly from your Android Phone too..)

Please install HAM for your Platform and Architecture and follow the Tutorial


Everything is Automated, So you can start a build and just forget about it, the server destroys itself when
the job is done (i.e The build finishes and the assets are uploaded to some other server). The program makes 
sure that even if the build errors out, it destroys itself. This makes each AOSP build economical and faster.

It also runs on **Termux** so you can build AOSP right from your Android Phone.


## How it Works

Ham has two programs, namely ```ham``` and ```ham-build```, **ham-build** program can only run on linux, this is by design
since it is not intended to run the user directly, this is the program which runs on Hetzner's cloud instance that we 
create. **ham** is the client program which can run on all platforms and architecture supported by golang, this program 
reads a recipe, ask question from the user when the build starts, creates a new cloud instance at Hetzner and then 
starts the build. It is also responsible to track the progress of **ham-build** which will be running on the Hetzner's
cloud instance.

The client program asks the user for required arguments for the build like Android Certificates if it's a signed build or 
a API key for the recipe to upload the output to Github or some server. These variables and files provided by the user 
are transported to the build server securely over SSH and SFTP protocols using a temporary SSH key created by HAM client
program. Every communication between ham client program and the ham-build program is done through SSH only (Using EdDSA and not RSA for security reasons).

## HAM Recipes

Ham recipes are simple directories which has **ham.yaml** or **ham.yml** file at the root of these simple directories.
This simple directory can contain anything the recipe author wants to have, this includes a version control too, so ham
recipes can be a git repository and it is recommended to be that way. **The requirement for ham to see it as a recipe 
is that it has a valid ham.yml file and it follows the syntax that ham defines.**

HAM YAML syntax example,

```
# Some sensible title for this recipe
title: "Lineage OS 19.1 (Enchilada/OP6) (Signed)"
# Semver or any string to trigger change for Ham
version: "0.0.1"

args:
  # id: This will be the name of the env variable in the build
  # server, no spaces are allowed, no hypens, only underscores.
  # 
  # prompt: This will the string displayed to the user when ham
  # client program asks questions to the user before the build.
  # 
  # required: true or false, by default it's false if it is not
  # defined by the recipe author. If true, build will not start
  # if user does not provide value for this.
  #
  # type: file, secret or value, file type represents a file path
  # which the user needs to give, secret is some type of secret
  # variable which will be handled with care, value are simple
  # string which will be set as env variables on the build server
  # for you to use.

  # There will be a env var on the build server called
  # ANDROID_CERTS_ZIP, and this env var will contain the
  # path to the uploaded file from the user on the build
  # server, so you can move use it as you like.
  - id: android_certs_zip
    prompt: "Path to Android Certificates in Un-Encrypted ZIP"
    required: true 
    type: file

  # This field will not be required, so user will 
  # have the option to skip
  # 
  # If user gives a value here, you will have a env
  # variable in your build server named TELEGRAM_KEY
  - id: telegram_key
    prompt: "Telegram API Key"
    type: secret
  
  # The id is capatilized and then set as the environmental
  # variable in the build server which you can use on all
  # your build scripts


# The actual list of commands to execute for the build.
# By default ham installs all the required basic requirement
# for building LineageOS, so you don't need to do that.
# It also sets the CCACHE env vars and also installs the
# repo command directly to the system for you to use.
# You might want to install your other deps using apt
# package manager, make sure to pass -y and -qq args to
# apt install

# On each ham build, ham creates /ham-build, /ham-recipe directories by
# default. /ham-recipe directory contains the entire directory of the ham
# recipe, if it's a git repo, then it is cloned into that destination on 
# the build server.
# 
# By default you are cd-ed into /ham-build which is an empty directory, you
# are expected to use this directory as so called home directory, since using
# absolute directories can be helpful when building AOSP and ~ are not really 
# parsed well by AOSP makefiles. 
build:
  # IMPORTANT: Set this for every recipe
  # otherwise your build will fail since repo
  # command needs a default python version,
  # not set by ham since we don't know if you are
  # building for legacy devices.
  - name: Set Python Default Version
    run: apt install -y -qq python-is-python3
 
  - name: Making Directory
    run: mkdir lineage

  # Github Action Style name and 
  # run
  - name: Repo Init
    run: |
        cd lineage
        repo sync -j20 -c < /dev/null
        mkdir -p .repo/local_manifests
  
  - name: Execute Bash Scripts
    run: /ham-recipe/scripts/yourscripts.sh

  - name: Do Patches
    run: patch /ham-build/lineage/Makefile /ham-recipe/patches/yourpatch.patch

  - name: Use Args given by User
    run: |
      echo "$TELEGRAM_KEY" > key.txt
      sleep 20
      echo "Something"
      
  - name: Use Files Given by Users
    run: |
      cd /root/ # Which is the home dir
      cp $ANDROID_CERTS_ZIP certs.zip
      mkdir .android-certs
      cd .android-certs
      mv ../certs.zip .
      unzip certs.zip
      rm -rf certs.zip 

  - name: Make sure to change directory
    run: cd /ham-build/lineage

# Automatically cd-ed into the /ham-build directory
post_build:
  - echo "Finished" > lineage/build.txt
  # Upload your files here.
```

## Why Only Hetzner and not Cloud Provider X

Hetzner is the only cloud provider which has predictable pricing and good bandwidth. Bandwidth is not the only thing that
makes Hetzner perfect, all other cloud providers don't provide **16 vCPUs AMD, 32 GB RAM, 10 GBit Internet and 320GB Storage** for the price point they give. Also all Hetzner Cloud instance has **20 TB** bandwidth.

To summarize,

* Hetzner has clear Pricing

* They have really good and stable Cloud API

* Easy to Understand API

* Official Open Source GO library for Hetzner Cloud API

* They Support Open Source Work

* They are based on Germany thus follows GDPR which means your Build Server's Data is Protected

* They are the Cheapest and Most Reliable Cloud Provider

* Referral Program gives new users free 20 euros cloud credit

* Hetzner is a Big Company as AWS and GCP but not that Popular


Also GCP, AWS and Azure are more focused on enterprise users and thus have different billing system which means a single
mistake can cost you thousands of dollars (search YouTube for incidents like this). That's the reason to go with Hetzner
and it is recommended.

**Also just focusing on one thing will make the program easy to maintain and code, also gives good quality. Just like the
Unix philosophy that a program should only do one and one thing only, and do that correctly.**


## Why not CI/CD self-hosted?

Yes there are a lot of options and solutions for this problem, but none of it offers the lower cost per build like this 
tool, Thanks to Hetzner, cloud is very cheap and powerful at the same time. Cloud's ultimate power comes to it's scaling 
powers, but in CI/CD build systems, there is no way to scale down to zero (which is the ultimate scalability). We waste 
a lot of computing resources doing noting but waiting for a cron job to actually do the work. 

With **HAM**, we can scale the cloud down to zero, ham creates a temporary server, reads a recipe and setups the 
environment and securely transfers required files and variables over SSH, starts the build and tracks it. Even if the 
client program closes for some reason, the server is still running and building Android. Server destroys itself when 
the work is finished without wasting costly computing resources.

**LineageOS 19.1 Signed Build for a Single Device(OnePlus 6) cost me about €0.30, the runtime of the build is around 
2-3 hours.** I don't need to stay awake for the build, it just runs over night and the server destroys itself when the 
job is done. **And the best part is, the client program can be run from Termux too, so I can just use my Android Phone 
to build a new Android OS for itself (remotely).**

Also I have seen some Github Actions which can create you Hetzner Server and Destroy it after running some command,
This is good too, but Ham gives other features like getting Variables and Files from the user and uploads it securely
over to the build server. It also tracks in realtime. Running SSH is not wise, if the line breaks, the build is stopped,
this is not good, that's why Ham uses a daemon on the build server to actually run the build so even if any connection
fails, the build is not affected. Also not everyone knows how to use Github CI/CD for a simple AOSP build.


## Disclaimer

**This project has no association with "Hetzner Online GmbH" in any form or manner, This project is purely Community work
and have no relationship with the company. This project purely exists for the community and will live if there is more contributions.**


## Hetzner Referral Program

Hetzner Online Gmbh has a referral program for loyal customers, if you signup using my referral link you will get free
20 euros cloud credit which you can build a ton of LineageOS builds for any device you like. The only problem is that 
Hetzner is pretty hard to register with but it is worth it. **I don't force you to use my referral link, it's totally upto
you.** [Hetzner Referral Link](https://hetzner.cloud/?ref=66oUbG2e4jXS)

Consider using the referral link as support towards this project. You can also star the project to make it more
credible.

## License

The BSD 3-Clause "New" or "Revised" License.

Copyright (C) 2022-present, D. Antony J.R.



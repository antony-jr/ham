<p align="center">
  <img src="artwork/logov1-export.png" height="180px" width=auto alt="HAM Logo">  <br>
</p>

# HAM [![GitHub issues](https://img.shields.io/github/issues/antony-jr/ham.svg?style=flat-square)](https://github.com/antony-jr/ham/issues) [![GitHub forks](https://img.shields.io/github/forks/antony-jr/ham.svg?style=flat-square)](https://github.com/antony-jr/ham/network) [![GitHub stars](https://img.shields.io/github/stars/antony-jr/ham.svg?style=flat-square)](https://github.com/antony-jr/ham/stargazers) [![GitHub license](https://img.shields.io/github/license/antony-jr/ham.svg?style=flat-square)](https://github.com/antony-jr/ham/blob/master/LICENSE) [![Deploy](https://github.com/antony-jr/ham/actions/workflows/deploy.yml/badge.svg)](https://github.com/antony-jr/ham/actions/workflows/deploy.yml)


HAM (Hetzner Android Make) is a Simple tool written in GO which can build LineageOS (or AOSP) from Source using Hetzner Cloud. 
**Build your Own Flavor of Android Under €1.** (Run Directly from your Android Phone too..)

<p align="center">
  <img src="artwork/pc-preview.gif" height=auto width=auto alt="Ham Preview PC">  <br>
</p>


<p align="center">
  <img src="artwork/preview.gif" height=auto width=auto alt="Ham Preview">  <br>
</p>

### Quickstart

Create your own recipes or browse recipes at https://github.com/ham-community, and then execute this in your terminal,

```
 # This is for building LineageOS 19.1 for OnePlus 6
 # Devices.

 # Only Once, to Initialize Hetzner Cloud API
 ./ham-linux-amd64 init

 ./ham-linux-amd64 get ~@gh/enchilada-los19.1:with_gapps

 # or without gapps and with F-Droid Priv Extensions
 ./ham-linux-amd64 get ~@gh/enchilada-los19.1
```

That's it, now your output should be uploaded by how the recipe describes. This recipe uploads the output to a github repo
given by the user. The repo can be private so you won't get any letter from Google for using gapps.

### What?

Ham is a simple tool which reads a recipe(git repo with a yaml file) and builds a Android build with the build instructions given in that 
recipe using Hetzner Cloud API. Everything is Automated, so you can start a build and just forget about it, the server destroys itself when
the job is done (i.e The build finishes and the assets are uploaded to some other server). The program makes sure that even if the build errors
out, it destroys itself. This makes each AOSP build economical and faster.

It also runs on **Termux** so you can build AOSP right from your Android Phone.

### Why?

Yes there are a lot of options and solutions for this problem, but none of it offers the lower cost per build like this tool, Thanks to Hetzner,
cloud is very cheap and powerful at the same time. Cloud's ultimate power comes to it's scaling powers, but in CI/CD build systems, there is no
way to scale down to zero (which is the ultimate scalability). We waste a lot of computing resources doing noting but waiting for a cron job to
actually do the work. 

With **HAM**, we can scale the cloud down to zero, ham creates a temporary server, reads a recipe and setups the environment and securely transfers
required files and variables over SSH, starts the build and tracks it. Even if the client program closes for some reason, the server is still running
and building Android. Server destroys itself when the work is finished without wasting costly computing resources.

**LineageOS 19.1 Signed Build for a Single Device(OnePlus 6) cost me about €0.50, the runtime of the build is around 5-6 hours.** I don't need to 
stay awake for the build, it just runs over night and the server destroys itself when the job is done. **And the best part is, the client program
can be run from Termux too, so I can just use my Android Phone to build a new Android OS for itself (remotely).**


### How?

Currently there are not a lot of recipes for Ham because there are no contributors, but if you want to support your device you are much welcomed.
All community recipes for Ham is at https://github.com/ham-community, each device recipe is represented by its own repo, it follows a naming scheme
something like this,

```
 <LineageOS Device Codename>-<Short Form of OS><Version>
```

Example: enchilada-los19.1 (for OnePlus 6)

Specific flavor of the build can be separated by branches, for example ```with_gapps_micro``` branch can include gapps right into the ROM build.


### Disclaimer

**This project has no association with "Hetzner Online GmbH" in any form or manner, This project is purely Community work
and have no relationship with the company. This project purely exists for the community and will live if there is more contributions.**


# License

The BSD 3-Clause "New" or "Revised" License.

Copyright (C) 2022-present, D. Antony J.R.



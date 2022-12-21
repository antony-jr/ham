---
slug: enchilada-los19-guide-relocking
title: Relocking your Bootloader with LineageOS 19.1 with MindTheGapps on OnePlus 6
authors: [antonyjr]
tags: [enchilada, guide, lineageos19.1, oneplus6]
---

The OnePlus 6 (codenamed _"enchilada"_) is a flagship smartphone from OnePlus.
It was released in May 2018.

| Basic                   | Spec Sheet                                                                                                                     |
| -----------------------:|:------------------------------------------------------------------------------------------------------------------------------ |
| CPU                     | Octa-core (4x2.8 GHz Kryo 385 Gold & 4x1.7 GHz Kryo 385 Silver)                                                                |
| Chipset                 | Qualcomm SDM845 Snapdragon 845                                                                                                 |
| GPU                     | Adreno 630                                                                                                                     |
| Memory                  | 6/8 GB RAM                                                                                                                     |
| Shipped Android Version | 8.1                                                                                                                            |
| Storage                 | 64/128/256 GB                                                                                                                  |
| Battery                 | Non-removable Li-Po 3300 mAh battery                                                                                           |
| Display                 | Optic AMOLED, 1080 x 2280 pixels, 19:9 ratio (~402 ppi density)                                                                |
| Camera (Back)           | Dual: 16 MP (f/1.7, 27mm, 1/2.6", 1.22µm, gyro-EIS, OIS) + 20 MP (16 MP effective, f/1.7, 1/2.8", 1.0µm), PDAF, dual-LED flash |
| Camera (Front)          | 16 MP (f/2.0, 25mm, 1/3", 1.0µm), gyro-EIS, Auto HDR, 1080p                                                                    |

![OnePlus 6](https://cdn2.gsmarena.com/vv/pics/oneplus/oneplus-6-5.jpg "OnePlus 6")

:::note

**This blog post and the entire HAM project was inspired from [this XDA forum post](https://forum.xda-developers.com/t/guide-re-locking-the-bootloader-on-the-oneplus-6t-with-a-self-signed-build-of-los.4113743/), so credits goes to the
OP, LineageOS developers and AOSP.**

:::


## Scope of this Guide

Creating an *unofficial build of LineageOS 19.1* suitable for using to re-lock the bootloader on a OnePlus 6 and
take you through the process of re-locking your bootloader after installing the above.

## Out of Scope

Remove *all* warning messages during boot (the yellow "Custom OS" message will be present though the orange 
"Unlocked bootloader" message will not) allow you to use official builds of LineageOS 19.1 on your device with a 
re-locked bootloader (more details near the end of the tutorial)

## Pre-requisites

* Basic knowledge of Terminal commands and features in your OS (Windows, Linux or MacOS)
* A OnePlus 6 Device
* A PC/Phone which supports you to run **[ham tool](https://antony-jr.github.io/ham/)**
* Finished **[setting up ham](https://antony-jr.github.io/ham/docs/get_started)** with "Getting Started" guide.
* A working USB cable
* Fastboot/Adb installed and functional
* **AndroidCerts.zip**, **Github API Key**, **Github Username** and **Github Repo** from "Getting Started" [guide](https://antony-jr.github.io/ham/docs/get_started)
 
:::danger

This process may brick your device. Do not proceed unless you are comfortable taking this risk. 

This process will delete all data on your phone! Do not proceed unless you have backed up your data!

Make sure you have read through this entire process at least once before attempting, if you are uncomfortable 
with any steps include in this guide, do not continue.

:::

**If you did not read the "Getting Started", please go and 
[read that first](https://antony-jr.github.io/ham/docs/get_started), 
then come back here.**


## Preparing ```answers.json```

We can pass a json file to HAM instead of answering everything during the prompt, you can **also** write these answers
during the prompt itself, that's your choice.

This file contains **important variables** that are needed from the user (i.e you), like a Github Repo to Upload the
output of the build, a Github API Key to authenticate to upload assets to the given Github repo. A Github Username and
Path to the AndroidCerts.zip file.

From "Getting Started" guide, you should have all these variables, if you already have **Android Certificates** then
zip all the certificates together, that's your AndroidCerts.zip file.

 

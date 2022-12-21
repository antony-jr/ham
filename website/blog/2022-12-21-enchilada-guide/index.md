---
slug: enchilada-los19-guide-relocking
title: Relocking your Bootloader with LineageOS 19.1 with MindTheGapps on OnePlus 6
authors: [antonyjr]
tags: [enchilada, guide, lineageos19.1, oneplus6]
---

:::info

This guide needs the [ham tool](https://antony-jr.github.io/ham/) to build, this guide cannot be used without
that tool. HAM(Hetzner Android Make) helps you build your own flavor of Android under one Euro using Hetzner Cloud.

:::

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

:::info

**This blog post and the entire HAM project was inspired from [this XDA forum post](https://forum.xda-developers.com/t/guide-re-locking-the-bootloader-on-the-oneplus-6t-with-a-self-signed-build-of-los.4113743/), so credits goes to the
OP, LineageOS developers and AOSP.**

:::


## Scope of this Guide

Creating an *unofficial build of LineageOS 19.1* suitable for using to re-lock the bootloader on a OnePlus 6 and
take you through the process of re-locking your bootloader after installing the above.

## Out of Scope

Remove *all* warning messages during boot (the yellow "Custom OS" message will be present though the orange 
"Unlocked bootloader" message will not), allow you to use official builds of LineageOS 19.1 on your device with a 
re-locked bootloader.

:::danger 

You **should not** flash the official LineageOS build after you have re-locked your phone, you can only build
unofficial LineageOS builds which are self signed by you, using HAM will make this process very easy. Do not lose
your AndroidCerts.zip file which contains your Android Certificates.

:::

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

:::danger

Keep the ```answers.json``` file safe since that contains a lot of sensitive information, ```~/.ham.json```
is the ham tool configuration file, this also contains a lot of sensitive information. Make sure that these
files are never shared with anyone.

::: 

:::tip

You may delete your ```answers.json``` after you finish the build successfully, this file is to make life easier
for you.

:::

Here is a sample ```answers.json``` file required for this build,

```json
{
   "github_token": "github_pat_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXx1",
   "github_user": "yourusername",
   "github_repo": "yourrepo",
   "android_certs": "/absolute/path/to/AndroidCerts.zip",
   "updater_url": ""
}
```

:::tip

You can give the ```updater_url``` a empty string, unless you know what you are doing. During the build, the updater
string will be updated to this url.

:::

:::note

```android_certs``` must be a file path to the AndroidCerts.zip file. In Windows this might be something like,
```C:\Users\Desktop\AndroidCerts.zip``` and so forth. If you are using Android to build, enable termux storage and
copy the certificates to your home.
:::


## Get LineageOS 19.1 with MindTheGapps (Self-Signed)

You should now have everything ready to start a remote build, start the build by entering the following command
in your operating system's Terminal. (In Windows, this is Command Prompt or Terminal App, In Android this is Termux and so forth).

```bash
 ham get -a answers.json ~@gh/enchilada-los19.1:gapps
```

Once you see the progress bar, and some progress like, **Installing Dependencies**, you may close this program using
**Ctrl + C** or **Esc** or **q**. You may now shutdown your device if you like to.. You can continue tracking the 
build or check the status of the build by running the same command on any device with the same ```~/.ham.json``` config
file which was used to start the build first.

Track using the following command,

```bash
 ham get ~@gh/enchilada-los19.1:gapps
```

If the build was successful, you will have the output in the given Github Repo releases, something like this,

![Demo Releases](releases.png)

After you see the build output, just to be safe run **```ham clean```** to destroy all build servers if exists, this
will be done by ham automatically on the server itself, but just to be safe. If build has failed, when you run 
the above get command, you will get a error saying there was a previous build which was **successful** or **failed**.

```bash
 ham clean # After you confirm the previous build was successful or failed.
 # or just run this if you want to stop every build asap to save cost.
```

## Download the LineageOS Build from Github Releases

After you see the new release on your repository, Download all files to a directory and and open a Terminal in that
directory.

Run the following commands, (**Make sure you have fastboot/adb installed in your system which you are going to use to
flash your newly built ROM**) 

```bash
 unzip lineage-19.1-*-recovery-enchilada.zip
 # If you are in Windows, use a Archive Extractor to extract 
 # the contents of the recovery zip 
```

You will now have a **pkmd.bin** file which is your public key that is going to be flashed into your Phone as the 
custom avb key, and a **```lineage-19.1-*-recovery-enchilada.img```** file when the archive is extracted.

**You must now unlock your phone, make backups and also if you want, use fastboot to backup your persist and 
other partition too.** See LineageOS docs for this [instruction](https://wiki.lineageos.org/devices/enchilada/install#unlocking-the-bootloader).
 
This is the usual you would do with flashing any ROM, except you should not patch or use any third party recovery,
we will only be using LineageOS recovery.

:::danger

You should not use TWRP or any Recovery other than LineageOS Recovery which you have downloaded from Github Releases,
similarly **YOU SHOULD ONLY SIDELOAD THE ROM AND NOTHING ELSE**. **DO NOT SIDELOAD GAPPS OR MAGISK, THIS WILL 
BRICK YOUR PHONE**.

:::

**Enable Developer Mode in your Phone by tapping the Build Number 7 Times, then enable ADB.** Now Connect your Phone
to your PC and do the following,

**Enable OEM unlock in the Developer options under device Settings**.

Connect the device to your PC via USB. On the computer, open a command prompt (on Windows) or 
terminal (on Linux or macOS) window, and type:

```bash
 adb reboot bootloader
```
Once the device is in fastboot mode, verify your PC finds it by typing:

```bash
 fastboot devices
```

To verify if fastboot connection is ok.

**Now let's flash the recovery we just extracted from the zip file,**

```
 fastboot flash boot lineage-19.1-*-recovery-enchilada.img 
```

Now reboot into recovery or enter the following command,
```
 fastboot reboot recovery
```

**You can also use the Volume Up/Down and navigate to Recovery and Press the Power Button.** You can also go into
recovery by pressing **Volume Down** + **Power**. Please follow LineageOS docs for more information.

Now Tap **Factory Reset**, then Format Data / Factory Reset and continue with the formatting process. This will remove encryption and delete all files stored in the internal storage, as well as format your cache partition 
(if you have one). Return to the main menu.

:::danger

You don't need to **Format Data** if you are already on a build from HAM and trying to update your ROM, this is 
only for a fresh install and should be done only if the ROM itself is different.

:::

Sideload the LineageOS .zip package:

On the device, select **“Apply Update”**, then **“Apply from ADB”** to begin sideload.
On the host machine, sideload the package using: **```adb sideload filename.zip.```**

In our case, the following command,
```bash
adb sideload lineage-19.1-*-release-enchilada-signed.zip
``` 
This is the file you downloaded from your releases.

Once you have installed everything successfully, click the back arrow in the top left of the screen, 
then **“Reboot system now”**.

## Check your Build Before Locking

After the reboot, you should have a working LineageOS 19.1 build with Google Apps embedded into it, check for bootloops
and other errors, don't setup anything or login into your google account. Just checkout the OS for bugs and fatal
errors.

Once checking is done,
Enable OEM unlocking and developer mode again, and reboot into fastboot again.

## Flashing your Public Key and Lock Bootloader

Run the following commands while you are in fastboot,

```bash
 fastboot flash avb_custom_key pkmd.bin
```

This **pkmd.bin** file is the file we extracted from the recovery zip we downloaded from our Github repo release which
is built by HAM.


Now we will lock the bootloader, this will wipe all data on your Phone, 

```bash
 fastboot reboot bootloader
 fastboot oem lock
```

**Now you will be rebooted and hopefully with no errors, with a locked bootloader. Now Setup your Phone and recover
your Backups from Google.**


## Conclusion

If you do end up with a bootloop after you lock your bootloader, please follow the unbrick guide for OnePlus 6 using
the MSM Tool. This will not happen in most cases, at least from build outputs of HAM. But there is no warranty of
any sorts. 

🎉 You now have a Custom ROM running in your OP6 with **Locked Bootloader** 🎉 

Special thanks to [@WhitbyGreg](https://forum.xda-developers.com/m/whitbygreg.1915770/) who made the XDA forum post,
and [Wunderment OS](https://github.com/Wunderment) from where I copied a lot of scripts to make the ham recipe.

import React from "react";
import Link from "@docusaurus/Link";
import clsx from "clsx";
import styles from "./styles.module.css";

import siteConfig from "@generated/docusaurus.config";

type PlatformItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<"svg">>;
  description: JSX.Element;
  arm_version_link: string;
  arm_dev_version_link: string;
  stable_version_link: string;
  development_version_link: string;
};

const PlatformList: PlatformItem[] = [
  {
    title: "GNU/Linux Distributions",
    Svg: require("@site/static/img/linux.svg").default,
    description: (
      <>
	 Runs on Distros newer than or old as <b>Ubuntu 18.04 LTS</b>.
	 Just set it as executable and use like a normal binary, make it executable with 
	  <code>chmod a+x ham-linux-amd64</code>.
      </>
    ),
    arm_version_link: "https://github.com/antony-jr/ham/releases/download/stable/ham-linux-arm64",
    arm_dev_version_link: "https://github.com/antony-jr/ham/releases/download/continuous/ham-linux-arm64",
    stable_version_link: "https://github.com/antony-jr/ham/releases/download/stable/ham-linux-amd64",
    development_version_link: "https://github.com/antony-jr/ham/releases/download/continuous/ham-linux-amd64",
  },
  {
    title: "Apple macOS",
    Svg: require("@site/static/img/macos.svg").default,
    description: (
      <>
	 Runs on macOS Catalina or later, same with linux, just set it as 
	 executable and use it from your Terminal. Just like a normal binary.
      </>
    ),
    arm_version_link: "https://github.com/antony-jr/ham/releases/download/stable/ham-macos-arm64",
    arm_dev_version_link: "https://github.com/antony-jr/ham/releases/download/continuous/ham-macos-arm64",
    stable_version_link: "https://github.com/antony-jr/ham/releases/download/stable/ham-macos-amd64",
    development_version_link: "https://github.com/antony-jr/ham/releases/download/continuous/ham-macos-amd64",
  },
  {
    title: "Microsoft Windows",
    Svg: require("@site/static/img/windows.svg").default,
    description: (
      <>
	 Download the EXE file and place it anywhere in your Windows PC, open a CMD window
	 or Powershell at that location and simply execute the EXE file.
      </>
    ),
    arm_version_link: "",
    arm_dev_version_link: "",
    stable_version_link: "https://github.com/antony-jr/ham/releases/download/stable/ham-windows-amd64.exe",
    development_version_link: "https://github.com/antony-jr/ham/releases/download/continuous/ham-windows-amd64.exe",
  },
  {
    title: "Android",
    Svg: require("@site/static/img/android.svg").default,
    description: (
      <>
	 Download the binary to your Phone, Now use Termux to set it as executable <code>chmod a+x ham-android-arm64</code>
	  and then simply start using it in Termux Terminal Emulator like you do in linux. It works without any issues.
      </>
    ),
    arm_version_link: "https://github.com/antony-jr/ham/releases/download/stable/ham-android-arm64",
    arm_dev_version_link: "https://github.com/antony-jr/ham/releases/download/continuous/ham-android-arm64",
    stable_version_link: "",
    development_version_link: "",
  },
];

function Platform({
  title,
  Svg,
  description,
  arm_version_link,
  arm_dev_version_link,
  stable_version_link,
  development_version_link,
}: PlatformItem) {
  return (
    <div className={clsx("col")}>
       <div class="card" style={{minHeight: "620px"}}>
        <div class="card__header">
          <h3>{title}</h3>
        </div>
        <div class="card__body">
          <div className="text--center">
            <Svg className={styles.platformSvg} role="img" />
          </div>

          <br />
          <p>{description}</p>
        </div>
        <div class="card__footer">
	  {stable_version_link != "" && <Link
            target="_self"
            to={stable_version_link}
            style={{ margin: "5px" }}
            class="button button--primary button--block"
          >
            Stable Version
          </Link>}
	  {development_version_link != "" && <Link
            target="_self"
            to={development_version_link}
            style={{ margin: "5px" }}
            class="button button--secondary button--block"
          >
            Development Version
          </Link>}
	  {arm_version_link != "" && <Link
            target="_self"
            to={arm_version_link}
            style={{ margin: "5px" }}
            class="button button--primary button--block"
          >
            Aarch64 Stable Version
          </Link>}
	  {arm_dev_version_link != "" && <Link
            target="_self"
            to={arm_dev_version_link}
            style={{ margin: "5px" }}
            class="button button--secondary button--block"
          >
	     Aarch64 Development Version
	  </Link>}
        </div>
      </div>
    </div>
  );
}

export default function HomepageDownloads(): JSX.Element {
  return (
    <div class="hero hero--secondary">
      <div class="container">
        <h1 class="hero__title">Download Now</h1>
        <p class="hero__subtitle">
          Supported on all major platforms, Select one of the options below to
          Download. You can verify the installer's checksum available at Github.
        </p>

        <div className="row">
          {PlatformList.map((props, idx) => (
            <Platform key={idx} {...props} />
          ))}
        </div>
        <br />
        <br />
      </div>
    </div>
  );
}

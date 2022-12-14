import React from "react";
import Link from "@docusaurus/Link";
import clsx from "clsx";
import styles from "./styles.module.css";

export default function HomepageAndroid(): JSX.Element {
  return (
    <div class="hero  hero--dark">
      <div class="container">
        <h1 class="hero__title">Build Android from Android</h1>

        <p class="hero__subtitle"> 
	   Since HAM can run on multiple architectures and platforms, it implements most of the
	   client code in pure golang instead of relying it on the operating system, this makes it
	   run on pretty much everything go is supported which include Android.
	</p>
	<p align="center">
	   <video style={{maxWidth: "550px", width: "80%", height: "auto"}} muted={true} autoPlay={true} loop={true}>
	   <source src={require("@site/static/vids/ham-phone.webm").default} type="video/webm" />
	</video>
	</p>
      </div>
    </div>
  );
}

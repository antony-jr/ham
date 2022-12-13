import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import HomepageDownloads from "@site/src/components/HomepageDownloads";
import HomepageAndroid from "@site/src/components/HomepageAndroid";

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <img height="200px" width="auto"
	 src={require("@site/static/img/logo.png").default} alt="Logo" /> 
	<h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
	    style={{margin: "10px"}}
            className="button button--secondary button--lg"
            to="/docs/intro"
          >
            Quick Tutorial ⏱️
          </Link>

          <Link
            style={{ margin: "10px" }}
            className="button button--secondary button--lg"
            to=""
            onClick={() => {
              document
                .getElementById("downloads")
                .scrollIntoView({ behavior: "smooth" });
            }}
          >
            Download Now 📥
          </Link>
        </div>
        <br />
	<video width="90%" height="" muted={true} autoPlay={true} loop={true}>
	   <source src={require("@site/static/vids/preview.webm").default} type="video/webm" />
	</video>
      </div>
    </header>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title="Build Android Under One Euro"
      description="Ham is simple tool written in GO which can build Andorid Under 1 Euro using Hetzner Cloud">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
	<HomepageAndroid />
	<div id="downloads">
          <HomepageDownloads />
        </div>
      </main>
    </Layout>
  );
}

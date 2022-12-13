import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Easy to Use',
    Svg: require('@site/static/img/easy.svg').default,
    description: (
      <>
	 Ham was made easy to use and user-friendly, it is very easy to use and hard to abuse,
	 it has the most sane defaults to prevent economical loss for the user.
      </>
    ),
  },
  {
    title: 'Cross Platform',
    Svg: require('@site/static/img/cross_platform.svg').default,
    description: (
      <>
	 Ham client program can on multiple platforms and architectures including <b>Android
	    itself using Termux Terminal Emulator</b>.
      </>
    ),
  },
  {
    title: 'Hackable',
    Svg: require('@site/static/img/hackable.svg').default,
    description: (
      <>
	 Ham uses a YAML recipe file held in a git repository, each device and ROM can have it's
	 own repo with different recipes, different branches can add or remove elements that are
	 built directly into the ROM.
      </>
    ),
  },
  {
    title: 'Modern',
    Svg: require('@site/static/img/modern.svg').default,
    description: (
      <>
	 Ham YAML files are losely based on Github Actions YAML file, and are designed to keep everything
	 modern, make things easy for the developers.
      </>
    ),
  },
  {
    title: 'Super Fast',
    Svg: require('@site/static/img/fast.svg').default,
    description: (
      <>
	 Unlike a CI/CD, We don't use any container like Docker or LXC and directly run the build script
	 on the VPS created by Hetzner Ubuntu 20.04 LTS image, this gives us extra edge on performance.
      </>
    ),
  },
  {
    title: 'Very Low Cost',
    Svg: require('@site/static/img/cheap.svg').default,
    description: (
      <>
	 Unlike other cloud companies, Hetzner gives the best performance for the cost,
	 LineageOS build for single device only cost <b>0.30 euros</b> which took only
	 <b> 3 hours</b> to finish. That's just <b>32 cents</b> in the US. 
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

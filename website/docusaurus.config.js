// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Ham',
  tagline: 'Build Your Own Flavor of Android under One Euro',
  url: 'https://antonyjr.in',
  baseUrl: '/ham',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'antony-jr', // Usually your GitHub org/user name.
  projectName: 'ham', // Usually your repo name.

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'Ham',
        logo: {
          alt: 'Ham Logo',
          src: 'img/logo.svg',
        },
        items: [
	  {
	    href: '/#downloads',
	    label: 'Downloads',
	    position: 'left',
	  },

          {
            type: 'doc',
            docId: 'intro',
            position: 'left',
            label: 'Documentation',
          },

          {to: '/blog', label: 'Blog', position: 'left'},
          {
            href: 'https://github.com/antony-jr/ham',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Documentation',
                to: '/docs/intro',
              },
	      {
		label: 'Get Started',
		to: '/docs/get_started',
	      },
	      {
		label: 'Ham Recipe Specification',
		to: '/docs/category/ham-recipe',
	      },
            ],
          },
          {
            title: 'Community',
            items: [
	      {
		label: 'Ham Recipes',
		href: 'https://github.com/ham-community',
	      },
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/ham',
              },
              {
                label: 'Twitter',
                href: 'https://twitter.com/antonyjr0',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Blog',
                to: '/blog',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/antony-jr/ham',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} HAM Documentation. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;

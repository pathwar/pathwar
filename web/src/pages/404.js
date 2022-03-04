/* eslint-disable react/display-name */
import React from "react";
import { Helmet } from "react-helmet";
import { graphql, withPrefix } from "gatsby";
import { Global, css } from "@emotion/core";
import { ThemeProvider, useTheme } from "emotion-theming";

import { FormattedMessage } from "react-intl";
import { lightTheme } from "../styles/themes";
import { globalStyle } from "../styles/globalStyle";
import darkBlueLogo from "../images/new-pathwar-logo-dark-blue.svg";
import colorsLogo from "../images/new_pathwar-logo.svg";
import footerLogo from "../images/new-pathwar-logo-grey.svg";
import footerLogoD from "../images/new-pathwar-logo-light-purple.svg";

import islandLeadingBg from "../images/island-light-mode-illustration.svg";
import islandLeadingBgDark from "../images/landing-island-darkmode-illustration.svg";

const browser = typeof window !== "undefined" && window;

const logoWrapper = () => `
  width: fit-content;

  img {
    width: 91px;
    height: 94px
  }

  @media (max-width: 991px) {
    text-align: center;
    width: 100%;

    img {
      width: 51px;
      height: 52px
    }
  }
`;

const leading = ({ type, colors }) => {
  const isDark = type === "dark";

  return `
    background-color: ${isDark ? colors.primary : colors.light};
    height: 813px;
    background-image: url(${isDark ? islandLeadingBgDark : islandLeadingBg});
    background-position: bottom right;
    background-repeat: no-repeat;
    background-size: contain;
    padding-top: 35px;

    @media (max-width: 991px) {
      height: 706px;
      padding-top: 25px;
    }
  `;
};

const leadingContent = () => `
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  margin-top: auto;
  margin-bottom: auto;
  width: 40%;
  height: 640px;

  .title-block,
  .sub-block {
    margin-bottom: 2.5rem;
  }

  @media (max-width: 991px) {
    height: auto;
    width: 100%;

    .title-block,
    .sub-block {
      margin-bottom: 1rem;
    }

    .title-block {
      margin-top: 2.812rem;
    }
  }
`;

const footer = ({ colors, type }) => `
  display: flex;
  align-items: self-end;
  justify-content: space-around;
  background-color: ${type === "dark" ? colors.dark : colors.light};
  color: ${type === "dark" ? colors.secondary : colors.tertiary};
  padding: 40px;

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  a {
    color: inherit;
  }

  .data-col {
    max-width: 150px;
  }

  @media (max-width: 991px) {
    flex-direction: column;
    align-items: center;
    text-align: center;

    .data-col {
      margin-bottom: 1.8rem;
    }
  }
`;

const NotFoundPage = ({ data }) => {
  const currentTheme = useTheme();
  const isDark = currentTheme.type === "dark";
  const title = data.site.siteMetadata.title;
  const description = data.site.siteMetadata.description;
  const logo = data.logo.childImageSharp.fixed;
  const featuredImage = `${data.site.siteMetadata.baseUrl}${logo.src}`;
  const headerLogo = isDark ? colorsLogo : darkBlueLogo;

  return (
    <>
      <Global styles={globalStyle(currentTheme)} />
      <Helmet>
        <title>{title} - 404 Not found</title>
        <meta name="description" content={description} />

        <meta
          name="go-import"
          content="pathwar.land git https://github.com/pathwar/pathwar"
        />
        <meta
          name="go-source"
          content="pathwar.land https://github.com/pathwar/pathwar https://github.com/pathwar/pathwar/tree/master{/dir} https://github.com/pathwar/pathwar/tree/master{/dir}/{file}#L{line}"
        />

        <meta property="og:description" content={description} />
        <meta property="og:url" content={description} />
        <meta property="og:site_name" content={title} />
        <meta property="og:type" content="website" />
        <meta property="og:image" content={featuredImage} />
        <meta property="og:image:width" content={logo.width} />
        <meta property="og:image:height" content={logo.height} />

        <meta property="twitter:card" content="summary" />
        <meta property="twitter:title" content={title} />
        <meta property="twitter:description" content={description} />
        <meta property="twitter:image" content={featuredImage} />
        <link
          href="https://fonts.googleapis.com/css2?family=Nunito:wght@400;700&display=swap"
          rel="stylesheet"
        ></link>
        <link
          href="https://fonts.googleapis.com/css2?family=Bungee&display=swap"
          rel="stylesheet"
        />
        <script async src={withPrefix("chat-init.js")} type="text/javascript" />
      </Helmet>

      <section
        className="siteContainer"
        css={theme =>
          css`
            ${leading(theme)}
          `
        }
      >
        <span
          css={theme => css`
            ${logoWrapper(theme)}
          `}
        >
          <img src={headerLogo} alt="Pathwar Logo" />
        </span>
        <div
          css={theme =>
            css`
              ${leadingContent(theme)}
            `
          }
        >
          <div className="title-block">
            <h1>
              <FormattedMessage id="404Page.title" />
            </h1>
          </div>
          <div className="sub-block">
            <h2>
              <FormattedMessage id="404Page.text" />
            </h2>
          </div>
        </div>
      </section>

      <footer
        css={theme =>
          css`
            ${footer(theme)}
          `
        }
      >
        <img src={isDark ? footerLogoD : footerLogo} />
        <p>
          © 2015-2022 Pathwar Staff Licensed under the Apache License, Version
          2.0
        </p>
      </footer>
    </>
  );
};

export default ({ data }) => {
  const themeToUse = lightTheme;

  return (
    browser && (
      <ThemeProvider theme={themeToUse}>
        <NotFoundPage data={data} />
      </ThemeProvider>
    )
  );
};

export const query = graphql`
  query {
    site {
      siteMetadata {
        title
        description
        baseUrl
      }
    }
    logo: file(relativePath: { eq: "images/new_pathwar-logo.png" }) {
      childImageSharp {
        fixed(width: 200, height: 200) {
          ...GatsbyImageSharpFixed
        }
      }
    }
  }
`;

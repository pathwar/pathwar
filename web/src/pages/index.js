/* eslint-disable react/display-name */
import React from "react";
import { Helmet } from "react-helmet";
import { graphql, withPrefix } from "gatsby";
import { Global, css } from "@emotion/core";
import { ThemeProvider, useTheme } from "emotion-theming";
import { lightTheme, darkTheme } from "../styles/themes";
import { globalStyle } from "../styles/globalStyle";

import hookIcon from "../images/hook-l-icon.svg";
import mapIcon from "../images/map-l-icon.svg";
import shipIcon from "../images/ship-l-icon.svg";
import hookIconD from "../images/hook-d-icon.svg";
import mapIconD from "../images/map-d-icon.svg";
import shipIconD from "../images/ship-d-icon.svg";
import darkBlueLogo from "../images/new-pathwar-logo-dark-blue.svg";
import colorsLogo from "../images/new_pathwar-logo.svg";
import footerLogo from "../images/new-pathwar-logo-grey.svg";
import footerLogoD from "../images/new-pathwar-logo-light-purple.svg";

import islandLeadingBg from "../images/island-light-mode-illustration.svg";
import islandLeadingBgDark from "../images/landing-island-darkmode-illustration.svg";
import { App } from "./app";

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

const leadingContent = ({ colors, type }) => `
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

  .cta-block {
    color: ${type === "dark" ? colors.secondary : colors.tertiary};
    text-align: center;

    a {
      font-weight: bold;
      color: ${type === "dark" ? colors.tertiary : colors.primary};
    }

    button {
      margin-bottom: 1rem;
    }
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

    .cta-block {
      text-align: left;
    }
  }
`;

const cardsArea = ({ shadows, colors, type }) => `
  flex-direction: row;
  align-items: center;
  justify-content: space-around;
  position: relative;
  top: -70px;

  .site-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    background-color: #fff;
    max-width: 358px;
    min-height: 389px;
    padding: 40px 40px 30px;
    box-shadow: ${shadows.card};
    color: ${type === "dark" ? colors.primary : "inherit"};
    border-radius: 6px;

    img {
      margin-bottom: 35px;
    }

    h3 {
      margin-bottom: 20px;
    }

    p {
      margin-bottom: 25px;
    }
  }

  @media (max-width: 991px) {
    flex-direction: column;

    .site-card {
      max-width: 270px;
      min-height: 264px;
      margin-bottom: 1.875rem;
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

const appIsRoot = process.env.APP_ROOT === "true";

const IndexPage = ({ data }) => {
  const currentTheme = useTheme();
  const isDark = currentTheme.type === "dark";
  const title = data.site.siteMetadata.title;
  const description = data.site.siteMetadata.description;
  const logo = data.logo.childImageSharp.fixed;
  const featuredImage = `${data.site.siteMetadata.baseUrl}${logo.src}`;
  const comingsoon = process.env.COMINGSOON === "true";
  const headerLogo = isDark ? colorsLogo : darkBlueLogo;

  return (
    <>
      <Global styles={globalStyle(currentTheme)} />
      <Helmet>
        <title>{title}</title>
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
        css={theme =>
          css`
            ${leading(theme)}
          `
        }
      >
        <div className="siteContainer">
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
              <h1>Learn, hack, challenge & more!</h1>
            </div>
            <div className="sub-block">
              <h2>
                Pathwar is an educational platform with a focus on security and
                cryptography.
              </h2>
            </div>
            {comingsoon && <h3>Coming soon...</h3>}
            {!comingsoon && (
              <div className="cta-block">
                <button className="custom-button">Join the adventure !</button>
                <p>
                  Already on board ? <a href="/app/challenges">Login</a>
                </p>
              </div>
            )}
          </div>
        </div>
      </section>

      {!comingsoon && (
        <>
          <section>
            <div
              css={theme =>
                css`
                  ${cardsArea(theme)}
                `
              }
              className="siteContainer"
            >
              <div className="site-card">
                <img src={isDark ? shipIconD : shipIcon} />
                <h3>Put your skills to the test</h3>
                <p>and improve them. Beat the challenges, learn new tricks.</p>
              </div>
              <div className="site-card">
                <img src={isDark ? mapIconD : mapIcon} />
                <h3>Participate in tournaments</h3>
                <p>
                  with your team and win prizes. Create or join a team and
                  compete with other players
                </p>
              </div>
              <div className="site-card">
                <img src={isDark ? hookIconD : hookIcon} />
                <h3>Hack everything</h3>
                <p>
                  Levels? Other players’ profiles? The platform itself?
                  Everything is fair game here!
                </p>
                <a
                  href="https://github.com/pathwar/pathwar/blob/master/CODE_OF_CONDUCT.md"
                  target="_blank"
                  rel="noreferrer noopener"
                  className="custom-button outline"
                >
                  Read our Code of Conduct
                </a>
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
            <div className="data-col">
              <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris
                eget est molestie est tincidunt varius. Suspendisse quis
                elementum odio, vitae euismod sem.
              </p>
            </div>
            <div className="data-col">
              <ul>
                <li>
                  <a href="#">CGU</a>
                </li>
                <li>
                  <a href="#">42</a>
                </li>
                <li>
                  <a href="#">Lorem lien 01</a>
                </li>
                <li>
                  <a href="#">Lorem lien 02</a>
                </li>
              </ul>
            </div>
            <div className="data-col">
              <ul>
                <li>
                  <a href="#">Lorem lien 03</a>
                </li>
                <li>
                  <a href="#">Lorem lien 04</a>
                </li>
                <li>
                  <a href="#">Lorem lien 05</a>
                </li>
                <li>
                  <a href="#">Lorem lien 06</a>
                </li>
              </ul>
            </div>
            <div className="data-col">
              <ul>
                <li>
                  <a href="#">Lorem lien 07</a>
                </li>
                <li>
                  <a href="#">Lorem lien 08</a>
                </li>
              </ul>
            </div>
            <div className="data-col">
              <p>Emplacement pour RS ou autre si besoin</p>
            </div>
          </footer>
        </>
      )}
    </>
  );
};

export default ({ data }) => {
  const browser = typeof window !== "undefined" && window;
  const browserInDarkMode =
    browser &&
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches;
  const themeToUse = browserInDarkMode ? darkTheme : lightTheme;

  return (
    <ThemeProvider theme={themeToUse}>
      {appIsRoot ? <App /> : <IndexPage data={data} />}
    </ThemeProvider>
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

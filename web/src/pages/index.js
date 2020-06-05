/* eslint-disable react/display-name */
import React, { useState } from "react";
import { Helmet } from "react-helmet";
import Img from "gatsby-image";
import { graphql } from "gatsby";
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
import footerLogo from "../images/new-pathwar-logo-grey.svg";
import footerLogoD from "../images/new-pathwar-logo-light-purple.svg";

import islandLeadingBg from "../images/island-light-mode-illustration.svg";
import islandLeadingBgDark from "../images/landing-island-darkmode-illustration.svg";

const logoLink = () => `
  cursor: pointer;

  @media (max-width: 991px) {
    text-align: center;
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

const IndexPage = ({ data, themeSwitch }) => {
  const currentTheme = useTheme();
  const isDark = currentTheme.type === "dark";
  const currentScreenWidth =
    typeof window !== "undefined" && window && window.screen.availWidth;
  const isMobile = currentScreenWidth <= 991;
  const title = data.site.siteMetadata.title;
  const description = data.site.siteMetadata.description;
  const logo = data.headerLogo.childImageSharp.fixed;
  const logoColors = data.headerLogoColors.childImageSharp.fixed;
  const logoMobile = data.headerLogoMobile.childImageSharp.fixed;
  const logoColorsMobile = data.headerLogoColorsMobile.childImageSharp.fixed;
  const featuredImage = `${data.site.siteMetadata.baseUrl}${logo.src}`;
  const comingsoon = process.env.COMINGSOON === "true";
  let logoToShow;

  if (isMobile) {
    logoToShow = isDark ? logoColorsMobile : logoMobile;
  } else {
    logoToShow = isDark ? logoColors : logo;
  }

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
      </Helmet>

      <section
        css={theme =>
          css`
            ${leading(theme)}
          `
        }
      >
        <div className="siteContainer">
          <a
            href="#"
            css={theme => css`
              ${logoLink(theme)}
            `}
            onClick={() => themeSwitch()}
          >
            <Img fixed={logoToShow} alt="Pathwar Logo" />
          </a>
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
                <button>Join the adventure !</button>
                <p>
                  Already on board ? <a href="#">Login</a>
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
                <button href="#" className="outline">
                  Read our Code of Conduct
                </button>
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
  const [darkMode, setDarkMode] = useState(false);
  const themeToUse = darkMode ? darkTheme : lightTheme;

  const switchTheme = () => {
    setDarkMode(mode => !mode);
  };

  return (
    <ThemeProvider theme={themeToUse}>
      <IndexPage data={data} themeSwitch={switchTheme} />
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
    headerLogo: file(
      relativePath: { eq: "images/new-pathwar-logo-dark-blue.png" }
    ) {
      childImageSharp {
        fixed(width: 91, height: 94) {
          ...GatsbyImageSharpFixed
        }
      }
    }
    headerLogoMobile: file(
      relativePath: { eq: "images/new-pathwar-logo-dark-blue.png" }
    ) {
      childImageSharp {
        fixed(width: 51, height: 52) {
          ...GatsbyImageSharpFixed
        }
      }
    }
    headerLogoColors: file(
      relativePath: { eq: "images/new_pathwar-logo.png" }
    ) {
      childImageSharp {
        fixed(width: 91, height: 94) {
          ...GatsbyImageSharpFixed
        }
      }
    }
    headerLogoColorsMobile: file(
      relativePath: { eq: "images/new_pathwar-logo.png" }
    ) {
      childImageSharp {
        fixed(width: 51, height: 52) {
          ...GatsbyImageSharpFixed
        }
      }
    }
  }
`;

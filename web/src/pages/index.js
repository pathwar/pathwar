/* eslint-disable react/display-name */
import React from "react";
import { Helmet } from "react-helmet";
import Img from "gatsby-image";
import { graphql } from "gatsby";
import { Global, css } from "@emotion/core";
import { useTheme } from "emotion-theming";
import { globalStyle } from "../styles/globalStyle";
import {
  leading,
  leadingContent,
  cardsArea,
  footer,
} from "./styles/indexStyle";

import hookIcon from "../images/hook-l-icon.svg";
import mapIcon from "../images/map-l-icon.svg";
import shipIcon from "../images/ship-l-icon.svg";
import hookIconD from "../images/hook-d-icon.svg";
import mapIconD from "../images/map-d-icon.svg";
import shipIconD from "../images/ship-d-icon.svg";
import footerLogo from "../images/new-pathwar-logo-grey.svg";
import footerLogoD from "../images/new-pathwar-logo-light-purple.svg";

export default ({ data }) => {
  const currentTheme = useTheme();
  const isDark = currentTheme.type === "dark";
  const title = data.site.siteMetadata.title;
  const description = data.site.siteMetadata.description;
  const logo = data.headerLogo.childImageSharp.fixed;
  const logoColors = data.headerLogoColors.childImageSharp.fixed;
  const featuredImage = `${data.site.siteMetadata.baseUrl}${logo.src}`;

  return (
    <>
      <Global
        styles={css`
          ${globalStyle}
        `}
      />
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

      <section css={theme => leading(theme)}>
        <div className="siteContainer">
          <Img fixed={isDark ? logoColors : logo} alt="Pathwar Logo" />

          <div css={theme => leadingContent(theme)}>
            <div className="title-block">
              <h1>Learn, hack, challenge & more!</h1>
            </div>
            <div className="sub-block">
              <h2>
                Pathwar is an educational platform with a focus on security and
                cryptography.
              </h2>
            </div>
            <div className="cta-block">
              <button>Join the adventure !</button>
              <p>
                Already on board ? <a href="#">Login</a>
              </p>
            </div>
          </div>
        </div>
      </section>

      <section>
        <div css={theme => cardsArea(theme)} className="siteContainer">
          <div className="site-card">
            <img src={isDark ? shipIconD : shipIcon} />
            <h3>Put your skills to the test</h3>
            <p>and improve them. Beat the challenges, learn new tricks.</p>
          </div>
          <div className="site-card">
            <img src={isDark ? mapIconD : mapIcon} />
            <h3>Participate in tournaments</h3>
            <p>
              with your team and win prizes. Create or join a team and compete
              with other players
            </p>
          </div>
          <div className="site-card">
            <img src={isDark ? hookIconD : hookIcon} />
            <h3>Hack everything</h3>
            <p>
              Levels? Other players’ profiles? The platform itself? Everything
              is fair game here!
            </p>
            <button href="#" className="outline">
              Read our Code of Conduct
            </button>
          </div>
        </div>
      </section>
      <footer css={theme => footer(theme)}>
        <img src={isDark ? footerLogoD : footerLogo} />
        <div className="data-col">
          <p>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris eget
            est molestie est tincidunt varius. Suspendisse quis elementum odio,
            vitae euismod sem.
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
    headerLogoColors: file(
      relativePath: { eq: "images/new_pathwar-logo.png" }
    ) {
      childImageSharp {
        fixed(width: 91, height: 94) {
          ...GatsbyImageSharpFixed
        }
      }
    }
  }
`;

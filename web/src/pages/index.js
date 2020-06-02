/* eslint-disable react/display-name */
import React from "react";
import { Helmet } from "react-helmet";
import Img from "gatsby-image";
import { graphql } from "gatsby";
import { Global, css } from "@emotion/core";

import islandLeadingBg from "../images/island-light-mode-illustration.svg";
import hookIcon from "../images/hook-l-icon.svg";
import mapIcon from "../images/map-l-icon.svg";
import shipIcon from "../images/ship-l-icon.svg";
import footerLogo from "../images/new-pathwar-logo-grey.svg";

const leading = css`
  background-color: #fff;
  height: 813px;
  background-image: url(${islandLeadingBg});
  background-position: bottom right;
  background-repeat: no-repeat;
  background-size: contain;
  padding-top: 35px;
`;

const leadingContent = css`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  margin-top: auto;
  margin-bottom: auto;
  width: 45%;
  height: 640px;

  .title-block,
  .sub-block {
    margin-bottom: 2.5rem;
  }

  .cta-block {
    color: #7493b0;
    text-align: center;

    a {
      font-weight: bold;
      color: #0081ff;
    }

    button {
      margin-bottom: 1rem;
    }
  }
`;

const cardsArea = css`
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
    box-shadow: 0 0 20px 0 rgba(56, 95, 200, 0.25);
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
`;

const footer = css`
  display: flex;
  align-items: self-end;
  justify-content: space-around;
  background-color: #fff;
  color: #7493b0;
  padding: 40px;

  ul {
    list-style: none;
  }

  a {
    color: inherit;
  }

  .data-col {
    max-width: 150px;
  }
`;

export default ({ data }) => {
  const title = data.site.siteMetadata.title;
  const description = data.site.siteMetadata.description;
  const logo = data.headerLogo.childImageSharp.fixed;
  const featuredImage = `${data.site.siteMetadata.baseUrl}${logo.src}`;

  return (
    <>
      <Global
        styles={css`
          body,
          html {
            background-color: #f3f9ff;
            font-family: "Nunito", sans-serif;
            font-size: 16px;
            color: #00376c;
          }

          h1 {
            color: #0071de;
            font-size: 3.125rem;
            font-family: "Bungee", cursive;
            margin: 0;
          }

          h2 {
            font-size: 1.25rem;
            margin: 0;
          }

          h3 {
            font-size: 1.125rem;
            font-weight: bold;
            margin: 0;
          }

          button {
            border: none;
            background: #0081ff;
            border-radius: 31px;
            font-size: 1.25rem;
            color: #ffffff;
            text-align: center;
            font-weight: bold;
            padding: 1rem 3rem;
            cursor: pointer;

            &.outline {
              color: #0081ff;
              background-color: transparent;
              font-size: 1rem;
              border: 2px solid #0081ff;
              border-radius: 31px;
              padding: 1rem 2rem;
            }

            &:hover {
              opacity: 0.7;
            }
          }

          .siteContainer {
            padding: 0 50px;
            display: flex;
            flex-direction: column;
            width: 100%;
          }
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

      <section css={leading}>
        <div className="siteContainer">
          <Img fixed={logo} alt="Pathwar Logo" />

          <div css={leadingContent}>
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
        <div css={cardsArea} className="siteContainer">
          <div className="site-card">
            <img src={shipIcon} />
            <h3>Put your skills to the test</h3>
            <p>and improve them. Beat the challenges, learn new tricks.</p>
          </div>
          <div className="site-card">
            <img src={mapIcon} />
            <h3>Participate in tournaments</h3>
            <p>
              with your team and win prizes. Create or join a team and compete
              with other players
            </p>
          </div>
          <div className="site-card">
            <img src={hookIcon} />
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
      <footer css={footer}>
        <img src={footerLogo} />
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
  }
`;

/* eslint-disable react/display-name */
import React from "react";
import { Helmet } from "react-helmet";
import Img from "gatsby-image";
import { graphql } from "gatsby";
import { Global, css } from "@emotion/core";

import islandLeadingBg from "../images/island-light-mode-illustration.svg";

const leading = css`
  background-color: #fff;
  height: 813px;
  background-image: url(${islandLeadingBg});
  background-position: bottom right;
  background-repeat: no-repeat;
  background-size: contain;
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

  .site-card {
    background-color: #fff;
    min-height: 389px;
    padding: 40px 40px 30px;
    box-shadow: 0 0 20px 0 rgba(56, 95, 200, 0.25);
    border-radius: 6px;
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

          button {
            border: none;
            background: #0081ff;
            border-radius: 31px;
            font-size: 1.25rem;
            color: #ffffff;
            text-align: center;
            font-weight: bold;
            padding: 1rem 3rem;
          }

          .siteContainer {
            padding: 35px 50px;
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
          <div className="site-card">Card 1</div>
          <div className="site-card">Card 2</div>
          <div className="site-card">Card 3</div>
        </div>
      </section>
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

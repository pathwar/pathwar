import React from "react"
import {Helmet} from "react-helmet"
import Img from "gatsby-image"

import { graphql } from "gatsby";

import styles from "./index.module.css";

export default ({ data }) => <div className={styles.page}>
  <Helmet>
    <title>{data.site.siteMetadata.title}</title>
    <meta name="description" content={data.site.siteMetadata.description} />

    <meta property="og:description" content={data.site.siteMetadata.description} />
    <meta property="og:url" content={data.site.siteMetadata.description} />
    <meta property="og:site_name" content={data.site.siteMetadata.title} />
    <meta property="og:type" content="website" />
    <meta property="og:image" content={data.file.childImageSharp.fixed} />
    <meta property="og:image:width" content="200" />
    <meta property="og:image:height" content="200" />

    <meta property="twitter:card" content="summary" />
    <meta property="twitter:title" content={data.site.siteMetadata.title} />
    <meta property="twitter:description" content={data.site.siteMetadata.description} />
    <meta property="twitter:image" content={data.file.childImageSharp.fixed} />
  </Helmet>

  <div className={styles.content}>
    <div><Img fixed={data.file.childImageSharp.fixed} alt="Pathwar Logo" /></div>
    <h1>{data.site.siteMetadata.title}</h1>
    <p>{data.site.siteMetadata.description}</p>
    <h2>Coming Soon..</h2>
  </div>
</div>

export const query = graphql`
  query {
    site {
      siteMetadata {
        title
        description
      }
    }
    file(relativePath: { eq: "images/pathwar-logo.png" }) {
      childImageSharp {
        fixed(width: 200, height: 200) {
          ...GatsbyImageSharpFixed
        }
      }
    }
  }
`

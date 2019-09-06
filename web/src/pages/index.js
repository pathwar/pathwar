import React from "react"
import {Helmet} from "react-helmet"
import Img from "gatsby-image"
import { graphql } from "gatsby";

import styles from "./index.module.css";

export default ({ data }) => {
  const title = data.site.siteMetadata.title;
  const description = data.site.siteMetadata.description;
  const logo = data.file.childImageSharp.fixed;
  const featuredImage = `${data.site.siteMetadata.baseUrl}${logo.src}`;

  return (
    <div className={styles.page}>
      <Helmet>
        <title>{title}</title>
        <meta name="description" content={description} />

        <meta name="go-import" content="pathwar.land git https://github.com/pathwar/pathwar" />
        <meta name="go-source" content="pathwar.land https://github.com/pathwar/pathwar https://github.com/pathwar/pathwar/tree/master{/dir} https://github.com/pathwar/pathwar/tree/master{/dir}/{file}#L{line}" />

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
      </Helmet>

      <div className={styles.content}>
        <div><Img fixed={logo} alt="Pathwar Logo" /></div>
        <p>{description}</p>
        <h2>COMING SOON..</h2>
      </div>
    </div>
  )
}

export const query = graphql`
  query {
    site {
      siteMetadata {
        title
        description
        baseUrl
      }
    }
    file(relativePath: { eq: "images/new_pathwar-logo.png" }) {
      childImageSharp {
        fixed(width: 200, height: 200) {
          ...GatsbyImageSharpFixed
        }
      }
    }
  }
`

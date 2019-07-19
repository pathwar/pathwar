import React from "react"
import {Helmet} from "react-helmet"
import Img from "gatsby-image"

import { graphql } from "gatsby";

import styles from "./index.module.css";

const title = "Pathwar.land";
const description = "Pathwar is an educational platform with a focus on security/cryptography, where you can go through levels (courses) to learn new things and hone your skills.";

export default ({ data }) => <div className={styles.page}>
  <Helmet defer={false}>
    <title>{title}</title>
    <meta name="description" content={description} />
  </Helmet>

  <div className={styles.content}>
    <div>
      <Img fixed={data.file.childImageSharp.fixed} alt="Pathwar Logo" />
    </div>
    <h1>Pathwar Land</h1>
    <p>{description}</p>
    <h2 className={styles.comingSoon}>Coming Soon..</h2>
  </div>
</div>

export const query = graphql`
  query {
    file(relativePath: { eq: "images/pathwar-logo.png" }) {
      childImageSharp {
        fixed(width: 200, height: 200) {
          ...GatsbyImageSharpFixed
        }
      }
    }
  }
`

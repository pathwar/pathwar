module.exports = {
  siteMetadata: {
    title: `Pathwar Land`,
    description: `Pathwar is an educational platform with a focus on security/cryptography, where you can go through levels (courses) to learn new things and hone your skills.`
  },
  plugins: [
    `gatsby-plugin-react-helmet`,
    `gatsby-transformer-sharp`,
    `gatsby-plugin-sharp`,
    {
      resolve: `gatsby-source-filesystem`,
      options: {
        path: `${__dirname}/src/`,
      },
    },
  ],
}

require("dotenv").config({
  path: `.env`,
});

const fs = require("fs");

let keycloakConfig = {
  realm: `${process.env.KEYCLOAK_REALM}`,
  "auth-server-url": `${process.env.KEYCLOAK_BASE_URL}auth`,
  "ssl-required": "external",
  resource: "platform-front",
  "public-client": true,
  "confidential-port": 0,
};

let keycloakFileData = JSON.stringify(keycloakConfig);
fs.writeFileSync("static/keycloak.json", keycloakFileData);

module.exports = {
  siteMetadata: {
    title: `☠️ Pathwar Land ☠️ `,
    description: `Pathwar is an educational platform with a focus on security/cryptography, where you can go through challenges (courses) to learn new things and hone your skills.`,
    baseUrl: `https://www.pathwar.land`,
  },
  plugins: [
    `gatsby-plugin-react-helmet`,
    `gatsby-transformer-sharp`,
    `gatsby-plugin-sharp`,
    `gatsby-plugin-emotion`,
    `gatsby-plugin-use-query-params`,
    {
      resolve: `gatsby-source-filesystem`,
      options: {
        path: `${__dirname}/src/`,
      },
    },
    {
      resolve: `gatsby-plugin-favicon`,
      options: {
        logo: `${__dirname}/src/images/pathwar-favicon.png`,

        dir: "auto",
        lang: "en-US",
        background: "#fff",
        theme_color: "#fff",
        display: "standalone",
        orientation: "any",
        start_url: "/?homescreen=1",
        version: "1.0",

        icons: {
          android: true,
          appleIcon: true,
          appleStartup: true,
          coast: true,
          favicons: true,
          firefox: true,
          yandex: true,
          windows: true,
        },
      },
    },
    {
      resolve: `gatsby-plugin-google-analytics`,
      options: {
        trackingId: "UA-47629346-5",
      },
    },
  ],
};

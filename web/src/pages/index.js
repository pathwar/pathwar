import React from "react";
import { ThemeProvider } from "emotion-theming";
import { lightTheme } from "../styles/themes";
import { Auth0Provider } from "@auth0/auth0-react";

import { App } from "./app";

const Index = () => {
  const themeToUse = lightTheme;

  return (
    <Auth0Provider
    domain={process.env.GATSBY_AUTH0_REALM}
    clientId={process.env.GATSBY_AUTH0_CLIENT_ID}
    authorizationParams={{
      redirect_uri: window.location.origin,
      audience: process.env.GATSBY_AUTH0_AUDIENCE
    }}
  >
      <ThemeProvider theme={themeToUse}>
        <App />
      </ThemeProvider>
    </Auth0Provider>
  );
};

export default Index;

import React from "react";
import { ThemeProvider } from "emotion-theming";
import { lightTheme } from "../styles/themes";
import { Auth0Provider } from "@auth0/auth0-react";

import { App } from "./app";

const Index = () => {
  const themeToUse = lightTheme;

  return (
    <Auth0Provider
    domain="dev-5ccwzy8qtcsjsnpf.us.auth0.com"
    clientId="bJpLWOLTRseEVfM9kvFhKfi9wUBmm8Gh"
    authorizationParams={{
      redirect_uri: window.location.origin,
      audience: "https://pathwar.net/"
    }}
  >
      <ThemeProvider theme={themeToUse}>
        <App />
      </ThemeProvider>
    </Auth0Provider>
  );
};

export default Index;

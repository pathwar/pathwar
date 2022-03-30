import React from "react";
import { ThemeProvider } from "emotion-theming";
import { Helmet } from "react-helmet";
import { lightTheme } from "../styles/themes";

import { App } from "./app";

const Index = () => {
  const themeToUse = lightTheme;

  return (
    <>
      <Helmet>
        <script async defer src="https://sa.moul.io/latest.js"></script>
        <noscript>
          {`
          <img
            src="https://queue.simpleanalyticscdn.com/noscript.gif"
            alt=""
            referrerpolicy="no-referrer-when-downgrade"
          />
          `}
        </noscript>
      </Helmet>
      <ThemeProvider theme={themeToUse}>
        <App />
      </ThemeProvider>
    </>
  );
};

export default Index;

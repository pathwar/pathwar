/* eslint-disable react/display-name */
import React from "react";
import { ThemeProvider } from "emotion-theming";
import { lightTheme } from "../styles/themes";

import { App } from "./app";

export default () => {
  const themeToUse = lightTheme;

  return (
    <ThemeProvider theme={themeToUse}>
      <App />
    </ThemeProvider>
  );
};

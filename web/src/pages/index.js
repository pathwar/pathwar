import React from "react";
import { ThemeProvider } from "emotion-theming";
import { lightTheme } from "../styles/themes";

import { App } from "./app";

const Index = () => {
  const themeToUse = lightTheme;

  return (
    <ThemeProvider theme={themeToUse}>
      <App />
    </ThemeProvider>
  );
};

export default Index;

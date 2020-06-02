import { lightTheme, darkTheme } from "./themes";

const activeTheme = window.activeTheme;
const isDark = activeTheme === "dark";
const themeToUse = isDark ? darkTheme : lightTheme;

export const globalStyle = `
body,
html {
  background-color: ${themeToUse.colors.body};
  font-family: ${themeToUse.font.family.body};
  font-size: ${themeToUse.font.size.base};
  color: ${themeToUse.colors.secondary};
}

h1 {
  color: ${isDark ? themeToUse.colors.light : themeToUse.colors.primary};
  font-size: 3.125rem;
  font-family: ${themeToUse.font.family.h1};
  margin: 0;
}

h2 {
  font-size: 1.25rem;
  font-weight: normal;
  margin: 0;
}

h3 {
  font-size: 1.125rem;
  font-weight: bold;
  margin: 0;
}

button {
  border: none;
  background:  ${
    isDark ? themeToUse.colors.tertiary : themeToUse.colors.primary
  };
  border-radius: 31px;
  font-size: 1.25rem;
  color: #ffffff;
  text-align: center;
  font-weight: bold;
  padding: 1rem 3rem;
  cursor: pointer;
  transition: all .3s linear;

  &.outline {
    color:  ${themeToUse.colors.primary};
    background-color: transparent;
    font-size: 1rem;
    border: 2px solid  ${themeToUse.colors.primary};
    border-radius: 31px;
    padding: 1rem 2rem;
  }

  &:hover {
    opacity: 0.7;
  }
}

.siteContainer {
  padding: 0 50px;
  display: flex;
  flex-direction: column;
  width: 100%;
}
`;

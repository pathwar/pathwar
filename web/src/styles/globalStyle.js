import { css } from "@emotion/core";

export const globalStyle = ({ colors, font, type }) => {
  const isDark = type === "dark";

  return css`
    body,
    html {
      background-color: ${colors.body};
      font-family: ${font.family.body};
      font-size: ${font.size.base};
      color: ${colors.secondary};
      margin: 0;
      padding: 0;
    }

    h1 {
      color: ${isDark ? colors.light : colors.primary};
      font-size: 3.125rem;
      font-family: ${font.family.h1};
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
      background: ${isDark ? colors.tertiary : colors.primary};
      border-radius: 31px;
      font-size: 1.25rem;
      color: #ffffff;
      text-align: center;
      font-weight: bold;
      padding: 1rem 3rem;
      cursor: pointer;
      transition: all 0.3s linear;

      &.outline {
        color: ${colors.primary};
        background-color: transparent;
        font-size: 1rem;
        border: 2px solid ${colors.primary};
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

    @media (max-width: 991px) {
      h1 {
        font-size: 1.562rem;
      }

      h2 {
        font-size: 1.125rem;
        font-weight: normal;
        margin: 0;
      }

      button {
        font-size: 1.125rem;
        padding: 0.4rem 1rem;
      }

      .siteContainer {
        padding: 0 22px;
        display: flex;
        flex-direction: column;
        width: 100%;
      }
    }
  `;
};

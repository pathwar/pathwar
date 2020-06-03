import { css } from "@emotion/core";
import islandLeadingBg from "../../images/island-light-mode-illustration.svg";
import islandLeadingBgDark from "../../images/landing-island-darkmode-illustration.svg";

export const logoLink = () => css`
  @media (max-width: 991px) {
    text-align: center;
  }
`;

export const leading = ({ type, colors }) => {
  const isDark = type === "dark";

  return css`
    background-color: ${isDark ? colors.primary : colors.light};
    height: 813px;
    background-image: url(${isDark ? islandLeadingBgDark : islandLeadingBg});
    background-position: bottom right;
    background-repeat: no-repeat;
    background-size: contain;
    padding-top: 35px;

    @media (max-width: 991px) {
      height: 706px;
      padding-top: 25px;
    }
  `;
};

export const leadingContent = ({ colors, type }) => css`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  margin-top: auto;
  margin-bottom: auto;
  width: 45%;
  height: 640px;

  .title-block,
  .sub-block {
    margin-bottom: 2.5rem;
  }

  .cta-block {
    color: ${type === "dark" ? colors.secondary : colors.tertiary};
    text-align: center;

    a {
      font-weight: bold;
      color: ${type === "dark" ? colors.tertiary : colors.primary};
    }

    button {
      margin-bottom: 1rem;
    }
  }

  @media (max-width: 991px) {
    height: auto;
    width: 100%;

    .title-block,
    .sub-block {
      margin-bottom: 1rem;
    }

    .title-block {
      margin-top: 2.812rem;
    }

    .cta-block {
      text-align: left;
    }
  }
`;

export const cardsArea = ({ shadows, colors, type }) => css`
  flex-direction: row;
  align-items: center;
  justify-content: space-around;
  position: relative;
  top: -70px;

  .site-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    background-color: #fff;
    max-width: 358px;
    min-height: 389px;
    padding: 40px 40px 30px;
    box-shadow: ${shadows.card};
    color: ${type === "dark" ? colors.primary : "inherit"};
    border-radius: 6px;

    img {
      margin-bottom: 35px;
    }

    h3 {
      margin-bottom: 20px;
    }

    p {
      margin-bottom: 25px;
    }
  }

  @media (max-width: 991px) {
    flex-direction: column;

    .site-card {
      max-width: 270px;
      min-height: 264px;
      margin-bottom: 1.875rem;
    }
  }
`;

export const footer = ({ colors, type }) => css`
  display: flex;
  align-items: self-end;
  justify-content: space-around;
  background-color: ${type === "dark" ? colors.dark : colors.light};
  color: ${type === "dark" ? colors.secondary : colors.tertiary};
  padding: 40px;

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  a {
    color: inherit;
  }

  .data-col {
    max-width: 150px;
  }

  @media (max-width: 991px) {
    flex-direction: column;
    align-items: center;
    text-align: center;

    .data-col {
      margin-bottom: 1.8rem;
    }
  }
`;

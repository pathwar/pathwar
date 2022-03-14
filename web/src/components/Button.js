/* eslint-disable react/prop-types */
import React from "react";
import { css } from "@emotion/core";

const buttonStyle = css`
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  text-align: center;
  font-weight: 900;
  width: 100%;
  padding: 1rem 0;
  border-radius: 5px;
  border: none;
  cursor: pointer;

  &:hover {
    opacity: 0.8;
  }
`;

const Button = ({ children, color, onClick, ...rest }) => {
  return (
    <button
      onClick={onClick}
      css={theme => [
        buttonStyle,
        css`
          background-color: ${theme.colors[color] || theme.colors.primary};
        `,
      ]}
      {...rest}
    >
      {children}
    </button>
  );
};

export default Button;

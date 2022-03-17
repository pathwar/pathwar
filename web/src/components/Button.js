/* eslint-disable react/prop-types */
import React from "react";
import { css } from "@emotion/core";

const buttonStyle = css`
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  text-align: center;
  font-weight: 900;
  padding: 0.7rem 3rem;
  border-radius: 5px;
  border: none;
  cursor: pointer;

  &:hover {
    opacity: 0.8;
  }
`;

const Button = ({
  children,
  color,
  onClick,
  textColor,
  emotionStyle,
  disabled,
  ...rest
}) => {
  return (
    <button
      onClick={onClick}
      css={theme => [
        buttonStyle,
        css`
          background-color: ${theme.colors[color] || theme.colors.primary};
          color: ${theme.colors[textColor] || theme.colors.light};
          ${disabled && `opacity: 0.7;`}
          ${emotionStyle};
        `,
      ]}
      disabled={disabled}
      {...rest}
    >
      {children}
    </button>
  );
};

export default Button;

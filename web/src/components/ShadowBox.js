import React from "react";
import { css } from "@emotion/core";

const wrapper = css`
  background-color: #fff;
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  margin-bottom: 0.5rem;
  padding: 1rem 1rem;
  border-radius: 8px;
  min-height: 200px;
  width: 100%;
`;

const ShadowBox = ({ children }) => {
  return <div css={wrapper}>{children}</div>;
};

export default ShadowBox;

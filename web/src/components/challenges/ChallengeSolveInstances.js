import React from "react";
import { css } from "@emotion/core";
import { Icon } from "tabler-react";
import { Terminal } from "react-window-ui";

const terminal = css`
  margin-bottom: 1rem;

  a {
    color: #fff;
    margin-right: 1.5rem;
  }

  span {
    color: #16b279;
  }
`;

const ChallengeSolveInstances = ({ instances, purchased }) => {
  return (
    <Terminal
      minHeight="2rem"
      boxShadow="0px 2px 15px -8px rgba(0,0,0,0.41)"
      css={terminal}
    >
      {!purchased && <p>Purchase the challenge to get resolution links...</p>}
      {purchased &&
        instances.map(item => (
          <div key={item.id}>
            <a href={item.nginx_url} target="_blank" rel="noreferrer">
              {item.nginx_url}
            </a>
            <span>
              <Icon name="check-circle" />
              {item.status}
            </span>
          </div>
        ))}
    </Terminal>
  );
};

export default React.memo(ChallengeSolveInstances);

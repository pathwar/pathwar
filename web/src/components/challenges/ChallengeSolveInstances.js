import React from "react";
import { css } from "@emotion/core";
import { Icon } from "tabler-react";
import { Terminal } from "react-window-ui";
import { FormattedMessage } from "react-intl";

const terminal = css`
  margin-bottom: 1rem;

  a {
    color: #fff;
    margin-right: 1.5rem;
  }
`;

const statusStyle = status => css`
  color: ${status === "Available" ? "#16b279" : "orange"};
`;

const ChallengeSolveInstances = ({ instances, purchased }) => {
  return (
    <Terminal
      minHeight="2rem"
      boxShadow="0px 2px 15px -8px rgba(0,0,0,0.41)"
      css={terminal}
    >
      {!purchased && (
        <p>
          <FormattedMessage id="ChallengeSolveInstances.purchasePhrase" />
        </p>
      )}
      {purchased &&
        instances.map(item => {
          const isAvailable = item.status === "Available";

          return (
            <div key={item.id}>
              <a href={item.nginx_url} target="_blank" rel="noreferrer">
                {item.nginx_url}
              </a>
              <span css={statusStyle(item.status)}>
                <Icon name={isAvailable ? "check-circle" : "x-circle"} />
                {isAvailable ? (
                  item.status
                ) : (
                  <FormattedMessage id="ChallengeSolveInstances.availableSoonPhrase" />
                )}
              </span>
            </div>
          );
        })}
    </Terminal>
  );
};

export default React.memo(ChallengeSolveInstances);

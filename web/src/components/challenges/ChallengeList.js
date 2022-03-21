/* eslint-disable react/prop-types */
import React from "react";
import { useSelector } from "react-redux";
import { Dimmer, Grid } from "tabler-react";
import { css } from "@emotion/core";
import { FormattedMessage } from "react-intl";
import ChallengeCard from "./ChallengeCard";

const container = theme => css`
  .sectionTitle {
    font-weight: 700;
    border-bottom: 1px solid ${theme.colors.gray};
    margin-bottom: 0.5rem;
    padding-bottom: 0.5rem;
    display: inline-block;
  }
`;

const ChallengeList = props => {
  const activeTeam = useSelector(state => state.seasons.activeTeam);

  const { challenges } = props;

  const purchasedNotSolved =
    challenges &&
    challenges.filter(challenge => {
      const { subscriptions } = challenge;
      const isClosed = subscriptions && subscriptions[0].status === "Closed";

      if (subscriptions && !isClosed) {
        return challenge;
      }
    });

  const notPurchased =
    challenges &&
    challenges.filter(challenge => {
      const { subscriptions } = challenge;

      if (!subscriptions) {
        return challenge;
      }
    });

  const closed =
    challenges &&
    challenges.filter(challenge => {
      const { subscriptions } = challenge;
      const isClosed = subscriptions && subscriptions[0].status === "Closed";

      if (subscriptions && isClosed) {
        return challenge;
      }
    });

  return !challenges ? (
    <Dimmer active loader />
  ) : (
    <div css={theme => container(theme)}>
      {purchasedNotSolved.length > 0 && (
        <>
          <h2 className="sectionTitle">
            <FormattedMessage id="ChallengeList.purchasedTitle" />
          </h2>

          <Grid.Row>
            {purchasedNotSolved.map(challenge => (
              <Grid.Col lg={6} sm={12} xs={12} key={challenge.id}>
                <ChallengeCard challenge={challenge} teamID={activeTeam.id} />
              </Grid.Col>
            ))}
          </Grid.Row>
        </>
      )}
      {notPurchased.length > 0 && (
        <>
          <h2 className="sectionTitle">
            <FormattedMessage id="ChallengeList.openTitle" />
          </h2>

          <Grid.Row>
            {notPurchased.map(challenge => (
              <Grid.Col lg={6} sm={12} xs={12} key={challenge.id}>
                <ChallengeCard challenge={challenge} teamID={activeTeam.id} />
              </Grid.Col>
            ))}
          </Grid.Row>
        </>
      )}
      {closed.length > 0 && (
        <>
          <h2 className="sectionTitle">
            <FormattedMessage id="ChallengeList.closedTitle" />
          </h2>

          <Grid.Row>
            {closed.map(challenge => (
              <Grid.Col lg={6} sm={12} xs={12} key={challenge.id}>
                <ChallengeCard challenge={challenge} teamID={activeTeam.id} />
              </Grid.Col>
            ))}
          </Grid.Row>
        </>
      )}
    </div>
  );
};

export default React.memo(ChallengeList);

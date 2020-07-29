import React, { memo } from "react";
import { Button } from "tabler-react";
import { css } from "@emotion/core";

const rewardText = css`
  font-size: 0.75rem;
  font-weight: 800;
  margin-top: 0.5rem;
`;

const ChallengeBuyButton = ({ challenge, buyChallenge, isClosed, ...rest }) => {
  const { subscriptions, flavor } = challenge;
  const hasSubscriptions = subscriptions;
  const { purchase_price: price, validation_reward: reward } = flavor;

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.flavor_id, challenge.season_id);
  };

  return (
    <>
      <Button
        icon={hasSubscriptions ? "check" : "dollar-sign"}
        color="indigo"
        disabled={hasSubscriptions || isClosed}
        onClick={handleBuyChallenge}
        {...rest}
      >
        {hasSubscriptions ? "Purchased" : `${price || 0} Buy`}
      </Button>
      <p css={rewardText}>Reward: {reward}</p>
    </>
  );
};

export default memo(ChallengeBuyButton);

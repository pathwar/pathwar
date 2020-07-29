import React, { memo } from "react";
import { Button } from "tabler-react";

const ChallengeBuyButton = ({ challenge, buyChallenge, ...rest }) => {
  const { subscriptions, flavor } = challenge;
  const hasSubscriptions = subscriptions;
  const { purchase_price: price } = flavor;

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.flavor_id, challenge.season_id);
  };

  return (
    <>
      <Button
        icon={hasSubscriptions ? "check" : "dollar-sign"}
        color="indigo"
        disabled={hasSubscriptions}
        onClick={handleBuyChallenge}
        {...rest}
      >
        {hasSubscriptions ? "Purchased" : `${price || 0} Buy`}
      </Button>
    </>
  );
};

export default memo(ChallengeBuyButton);

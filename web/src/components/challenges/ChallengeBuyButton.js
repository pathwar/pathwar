import React, { memo } from "react";
import { FormattedMessage } from "react-intl";
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
        {hasSubscriptions ? (
          <FormattedMessage id="ChallengeBuyButton.buy" />
        ) : (
          <>
            {price || 0} <FormattedMessage id="ChallengeBuyButton.buy" />
          </>
        )}
      </Button>
    </>
  );
};

export default memo(ChallengeBuyButton);

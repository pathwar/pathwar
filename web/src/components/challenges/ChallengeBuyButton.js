import React, { memo } from "react";
import { FormattedMessage } from "react-intl";
import { useDispatch } from "react-redux";
import { buyChallenge as buyChallengeAction } from "../../actions/seasons";

import Button from "../Button";

const ChallengeBuyButton = ({ challenge, ...rest }) => {
  const dispatch = useDispatch();

  const buyChallenge = (flavorChallengeID, seasonID) =>
    dispatch(buyChallengeAction(flavorChallengeID, seasonID));

  const { subscriptions } = challenge;
  const hasSubscriptions = subscriptions;

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.flavor_id, challenge.season_id);
  };

  return (
    <>
      <Button
        color="yellow"
        textColor="secondary"
        emotionStyle={`
          width: 100%;
        `}
        disabled={hasSubscriptions}
        onClick={handleBuyChallenge}
        {...rest}
      >
        {hasSubscriptions ? (
          <FormattedMessage id="ChallengeBuyButton.purchased" />
        ) : (
          <FormattedMessage id="ChallengeBuyButton.buy" />
        )}
      </Button>
    </>
  );
};

export default memo(ChallengeBuyButton);

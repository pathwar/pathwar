import React from "react"
import { Button } from "tabler-react"

const ChallengeBuyButton = ({
  challenge,
  teamID,
  buyChallenge,
  isClosed,
  ...rest
}) => {
  const hasSubscriptions = challenge.subscriptions

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.id, teamID, true);
  }

  return (
    <Button icon={hasSubscriptions ? "check" : "dollar-sign"} color="success" disabled={hasSubscriptions || isClosed} onClick={handleBuyChallenge} {...rest}>
      {hasSubscriptions ? "Purchased" : "Buy"}
    </Button>
  )
}

export default ChallengeBuyButton

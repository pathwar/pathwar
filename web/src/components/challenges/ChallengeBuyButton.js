import React from "react"
import { Button } from "tabler-react"

const ChallengeBuyButton = ({
  challenge,
  teamID,
  buyChallenge
}) => {
  const hasSubscriptions = challenge.subscriptions

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.id, teamID);
  }

  return (
    <Button icon={hasSubscriptions ? "check" : "dollar-sign"} color="success" disabled={hasSubscriptions} onClick={handleBuyChallenge}>
      {hasSubscriptions ? "Purchased" : "Buy"}
    </Button>
  )
}

export default ChallengeBuyButton

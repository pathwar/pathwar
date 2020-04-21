import React from "react"
import { Link } from "gatsby"
import { StampCard } from "tabler-react"

const ChallengeBuyStampCard = ({
  challenge,
  teamID,
  buyChallenge
}) => {
  const hasSubscriptions = challenge.subscriptions

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.id, teamID);
  }

  const CardHeader = () => hasSubscriptions ? (
      <small>Purchased</small>
    ) : (
      <Link to="/" onClick={handleBuyChallenge}>
        <small>Buy</small>
      </Link>
    )

  const cardFooterText = hasSubscriptions ? "You have this challenge" : "Buy challenge";

  return (
      <StampCard
        color="success"
        icon={hasSubscriptions ? "check" : "dollar-sign"}
        header={<CardHeader />}
        footer={cardFooterText}
      />
  )
}

export default ChallengeBuyStampCard

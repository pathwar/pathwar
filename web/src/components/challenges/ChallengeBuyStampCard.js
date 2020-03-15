import React, { useState, useEffect } from "react"
import { Link } from "gatsby"
import { isEmpty } from "ramda"
import { StampCard, Form, Button } from "tabler-react"

const ChallengeBuyStampCard = ({
  challenge,
  teamID,
  buyChallenge
}) => {
  const hasSubscriptions = challenge.subscriptions

  const handleBuyChallenge = async event => {
    event.preventDefault();
    await buyChallenge(challenge.id, teamID, challenge.season_id);
  }

  const CardHeader = () => hasSubscriptions ? (
      <small>Purchased</small>
    ) : (
      <Link to="/" onClick={handleBuyChallenge}>
        <small>Buy</small>
      </Link>
    )

  const cardFooterText = hasSubscriptions ? "You have this challenge" : "Buy challenge";

  console.log(challenge)

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

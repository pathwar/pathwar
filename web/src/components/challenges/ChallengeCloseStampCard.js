import React from "react"
import { Link } from "gatsby"
import { StampCard } from "tabler-react"

const ChallengeCloseStampCard = ({
  challenge,
  closeChallenge
}) => {
  const hasSubscriptions = challenge.subscriptions
  const subscription =
  hasSubscriptions &&
  challenge.subscriptions.find(item => item.status === "Active")
  const subscriptionHasValidations = subscription && subscription.validations

  const handleCloseChallenge = async event => {
    event.preventDefault();
    await closeChallenge(subscription.id);
  }

  const CardHeader = () => hasSubscriptions && subscriptionHasValidations ? (
      <Link to="/" onClick={handleCloseChallenge}>
        <small>Close</small>
      </Link>
    ) : (
        <small>Can't Close</small>
    )

  return (
      <StampCard
        color="danger"
        icon="x-circle"
        header={<CardHeader />}
        footer={""}
      />
  )
}

export default ChallengeCloseStampCard

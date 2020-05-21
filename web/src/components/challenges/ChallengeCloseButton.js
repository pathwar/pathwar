import React from "react"
import { Button } from "tabler-react"

const ChallengeCloseButton = ({
  challenge,
  closeChallenge,
  isClosed,
  ...rest
}) => {
  const hasSubscriptions = challenge.subscriptions
  const subscription = hasSubscriptions && challenge.subscriptions[0]

  const handleCloseChallenge = async event => {
    event.preventDefault()
    await closeChallenge(subscription.id)
  }

  return (
    <Button
      icon={"x-circle"}
      color="danger"
      onClick={handleCloseChallenge}
      disabled={isClosed}
      {...rest}
    >
      {isClosed ? "Closed" : "Close"}
    </Button>
  )
}

export default ChallengeCloseButton

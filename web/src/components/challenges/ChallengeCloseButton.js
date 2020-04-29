import React from "react"
import { Button } from "tabler-react"

const ChallengeCloseButton = ({ challenge, closeChallenge }) => {
  const hasSubscriptions = challenge.subscriptions
  const subscription =
    hasSubscriptions &&
    challenge.subscriptions.find(item => item.status === "Active")

  const handleCloseChallenge = async event => {
    event.preventDefault()
    await closeChallenge(subscription.id)
  }

  return (
    <Button icon={"x-circle"} color="danger" onClick={handleCloseChallenge}>
      Close
    </Button>
  )
}

export default ChallengeCloseButton

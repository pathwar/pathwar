import React, { useState } from "react"
import { Link } from "gatsby"
import { StampCard, Form, Button, Badge } from "tabler-react"

const ChallengeValidateStampCard = ({
  challenge,
  validateChallenge
}) => {
  const [isValidateOpen, setValidateOpen] = useState(false)
  const [formData, setFormData] = useState({ passphrase: "", comment: "" })
  const hasSubscriptions = challenge.subscriptions
  const subscription =
    hasSubscriptions &&
    challenge.subscriptions.find(item => item.status === "Active");

  const subscriptionHasValidations = subscription && subscription.validations


  const submitValidate = event => {
    event.preventDefault();
    const validateDataSet = {
      ...formData,
      subscriptionID: subscription.id,
    }
    validateChallenge(validateDataSet, challenge.season_id);
  }

  const handleChange = event => {
    setFormData({
      ...formData,
      [event.target.name]: event.target.value,
    });
  }

  const handleFormOpen = event => {
    event.preventDefault();
    setValidateOpen(prev => !prev);
  }

  const CardHeader = () => hasSubscriptions ? (
      <Link to="/" onClick={handleFormOpen}>
        <small>Validate</small>
      </Link>
    ) : ( <small>Validate</small> )

  const cardFooterText = hasSubscriptions && subscriptionHasValidations
  ? <p>Validations: <Badge color="primary" className="mr-1">{subscription.validations.length}</Badge></p>
  : hasSubscriptions ? "Validate with a passphrase"
  : "You need to purchase to validate";

  console.log(challenge)

  return (
    <>
      <StampCard
        color="warning"
        icon={hasSubscriptions ? "check-circle" : "circle"}
        header={<CardHeader />}
        footer={cardFooterText}
      />
      {isValidateOpen && (
        <form onSubmit={submitValidate}>
          <Form.FieldSet>
            <Form.Group isRequired label="Passphrase">
              <Form.Input name="passphrase" onChange={handleChange} />
            </Form.Group>
            <Form.Group isRequired label="Comment">
              <Form.Textarea
                name="comment"
                onChange={handleChange}
                placeholder="Leave a comment..."
                rows={6}
              />
            </Form.Group>
            <Form.Group>
              <Button type="submit" color="primary" className="ml-auto">
                Send
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </>
  )
}

export default ChallengeValidateStampCard

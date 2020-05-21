import React, { useState } from "react"
import { Form, Button } from "tabler-react"
import styles from "../../styles/layout/button.module.css"

const ChallengeValidateButton = ({ challenge, validateChallenge, ...rest }) => {
  const [isValidateOpen, setValidateOpen] = useState(false)
  const [isFetching, setFetching] = useState(false)
  const [formData, setFormData] = useState({ passphrase: "", comment: "" })
  const hasSubscriptions = challenge.subscriptions
  const subscription =
    hasSubscriptions &&
    challenge.subscriptions.find(item => item.status === "Active")

  const submitValidate = event => {
    event.preventDefault()
    const validateDataSet = {
      ...formData,
      subscriptionID: subscription.id,
    }

    setFetching(true)
    validateChallenge(validateDataSet, challenge.season_id).then(() => {
      setValidateOpen(false)
      setFetching(false)
    })
  }

  const handleChange = event => {
    setFormData({
      ...formData,
      [event.target.name]: event.target.value,
    })
  }

  const handleFormOpen = event => {
    event.preventDefault()
    setValidateOpen(prev => !prev)
  }

  return (
    <>
      <Button
        icon={"check-circle"}
        color="warning"
        className={styles.btn}
        onClick={handleFormOpen}
        {...rest}
      >
        Validate
      </Button>
      {isValidateOpen && (
        <form onSubmit={submitValidate}>
          <Form.FieldSet>
            <Form.Group isRequired label="Passphrase">
              <Form.Input
                name="passphrase"
                onChange={handleChange}
                placeholder="Insert passphrase here"
              />
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
              <Button
                type="submit"
                color="primary"
                className="ml-auto"
                disabled={isFetching}
              >
                Send
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </>
  )
}

export default ChallengeValidateButton

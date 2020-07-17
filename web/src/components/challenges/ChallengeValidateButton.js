import React, { useState, useEffect } from "react";
import { Form, Button } from "tabler-react";
import { isEmpty } from "ramda";
import styles from "../../styles/layout/button.module.css";

const initialErrorObj = { withError: false, fieldsWithError: [] };

const ChallengeValidateButton = ({ challenge, validateChallenge, ...rest }) => {
  const [isValidateOpen, setValidateOpen] = useState(false);
  const [isFetching, setFetching] = useState(false);
  const [formData, setFormData] = useState({ passphrases: "", comment: "" });
  const [error, setError] = useState(initialErrorObj);

  const hasSubscriptions = challenge.subscriptions;
  const subscription =
    hasSubscriptions &&
    challenge.subscriptions.find(item => item.status === "Active");

  useEffect(() => {
    if (!isEmpty(formData.passphrase) && !isEmpty(formData.comment)) {
      setError(initialErrorObj);
    }
  }, [formData, error]);

  const submitValidate = event => {
    event.preventDefault();

    if (isEmpty(formData.passphrases) || isEmpty(formData.comment)) {
      let fields = [];

      isEmpty(formData.passphrases) && fields.push("passphrases");
      isEmpty(formData.comment) && fields.push("comment");

      setError({ withError: true, fieldsWithError: fields });

      return;
    } else {
      const validateDataSet = {
        ...formData,
        passphrases: formData.passphrases.split(/[, ]+/),
        subscriptionID: subscription.id,
      };

      setFetching(true);
      validateChallenge(validateDataSet, challenge.season_id).then(() => {
        setValidateOpen(false);
        setFetching(false);
      });
    }
  };

  const handleChange = event => {
    setFormData({
      ...formData,
      [event.target.name]: event.target.value,
    });
  };

  const handleFormOpen = event => {
    event.preventDefault();
    setValidateOpen(prev => !prev);
  };

  const passphraseWithError =
    error.withError && error.fieldsWithError.includes("passphrases");

  const commentWithError =
    error.withError && error.fieldsWithError.includes("comment");

  return (
    <>
      <Button
        icon={"check-circle"}
        color="indigo"
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
                name="passphrases"
                onChange={handleChange}
                placeholder="Insert passphrases here separated by ','"
                invalid={passphraseWithError}
                cross={passphraseWithError}
                feedback={
                  passphraseWithError && "Please, insert a least one passphrase"
                }
              />
            </Form.Group>
            <Form.Group isRequired label="Comment">
              <Form.Textarea
                name="comment"
                onChange={handleChange}
                placeholder="Leave a comment..."
                rows={3}
                invalid={commentWithError}
                cross={commentWithError}
                feedback={commentWithError && "Please, insert a comment"}
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
  );
};

export default ChallengeValidateButton;

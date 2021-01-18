import React, { useState, useEffect } from "react";
import { Form, Button } from "tabler-react";
import { css } from "@emotion/core";
import { isEmpty } from "ramda";
import { useIntl, FormattedMessage } from "react-intl";

const initialErrorObj = { withError: false, fieldsWithError: [] };

const formStyle = css`
  margin-top: 1rem;
  text-align: left;
`;

const ChallengeValidateForm = ({ challenge, validateChallenge, ...rest }) => {
  const intl = useIntl();
  const [isValidateOpen, setValidateOpen] = useState(false);
  const [isFetching, setFetching] = useState(false);
  const [formData, setFormData] = useState({ passphrases: "", comment: "" });
  const [error, setError] = useState(initialErrorObj);

  const hasSubscriptions = challenge.subscriptions;
  const subscription =
    hasSubscriptions &&
    challenge.subscriptions.find(item => item.status === "Active");

  useEffect(() => {
    if (!isEmpty(formData.passphrase)) {
      setError(initialErrorObj);
    }
  }, [formData]);

  const submitValidate = event => {
    event.preventDefault();

    if (isEmpty(formData.passphrases)) {
      let fields = [];

      isEmpty(formData.passphrases) && fields.push("passphrases");
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

  const passphrasePlaceholderIntl = intl.formatMessage({
    id: "ChallengeValidateForm.passphrasePlaceholder",
  });

  const commentPlaceholderIntl = intl.formatMessage({
    id: "ChallengeValidateForm.commentPlaceholder",
  });

  return (
    <>
      <Button icon={"check-circle"} color="indigo" onClick={handleFormOpen}>
        <FormattedMessage id="ChallengeValidateForm.validate" />
      </Button>
      {isValidateOpen && (
        <form onSubmit={submitValidate} css={formStyle} {...rest}>
          <Form.FieldSet>
            <Form.Group
              isRequired
              label={
                <FormattedMessage id="ChallengeValidateForm.passphraseLabel" />
              }
            >
              <Form.Input
                name="passphrases"
                onChange={handleChange}
                placeholder={passphrasePlaceholderIntl}
                invalid={passphraseWithError}
                cross={passphraseWithError}
                feedback={
                  passphraseWithError && "Please, insert a least one passphrase"
                }
                autoFocus={true}
              />
            </Form.Group>
            <Form.Group
              label={
                <FormattedMessage id="ChallengeValidateForm.commentLabel" />
              }
            >
              <Form.Textarea
                name="comment"
                onChange={handleChange}
                placeholder={commentPlaceholderIntl}
                rows={3}
              />
            </Form.Group>
            <Form.Group>
              <Button
                type="submit"
                color="primary"
                className="ml-auto"
                disabled={isFetching}
              >
                <FormattedMessage id="ChallengeValidateForm.send" />
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </>
  );
};

export default ChallengeValidateForm;

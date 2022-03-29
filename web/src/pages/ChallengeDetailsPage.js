import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Helmet } from "react-helmet";
import { css } from "@emotion/core";
import PropTypes from "prop-types";
import { Page, Grid, Dimmer, Tag } from "tabler-react";
import moment from "moment";
import siteMetaData from "../constants/metadata";
import {
  fetchChallengeDetail as fetchChallengeDetailAction,
  validateChallenge as validateChallengeAction,
  // closeChallenge as closeChallengeAction,
} from "../actions/seasons";
// import ChallengeCloseButton from "../components/challenges/ChallengeCloseButton";
import ChallengeValidateForm from "../components/challenges/ChallengeValidateForm";
// import ValidationsList from "../components/challenges/ValidationsList";
import ChallengeSolveInstances from "../components/challenges/ChallengeSolveInstances";
import { CLEAN_CHALLENGE_DETAIL } from "../constants/actionTypes";
import { FormattedMessage } from "react-intl";
import ChallengeCard from "../components/challenges/ChallengeCard";

const statusTag = css`
  font-size: 0.75rem;
  opacity: 0.7;
`;

const whiteWrapper = css`
  background-color: #fff;
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  margin-bottom: 0.5rem;
  padding: 1rem 1rem;
  border-radius: 8px;
  min-height: 200px;
  width: 100%;
`;

const ChallengeDetailsPage = props => {
  const dispatch = useDispatch();
  const challenge = useSelector(state => state.seasons.challengeInDetail);

  const validateChallenge = (validationData, seasonId) =>
    dispatch(validateChallengeAction(validationData, seasonId));
  // const closeChallenge = subscriptionID =>
  //   dispatch(closeChallengeAction(subscriptionID));
  const fetchChallengeDetail = challengeID =>
    dispatch(fetchChallengeDetailAction(challengeID));

  useEffect(() => {
    const { uri, challengeID: challengeIDFromProps } = props;
    const challengeIDFromURI = uri && uri.split("/")[2];
    const challhengeID = challengeIDFromURI || challengeIDFromProps;

    fetchChallengeDetail(challhengeID);

    return () => dispatch({ type: CLEAN_CHALLENGE_DETAIL });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (!challenge) {
    return <Dimmer active loader />;
  }

  const { flavor, subscriptions } = challenge || {};

  const subscription = subscriptions && subscriptions[0];
  const validations = subscription && subscription.validations;
  const validation = validations && validations[0];
  const isClosed = subscription && subscription.status === "Closed";
  const purchased = subscriptions;
  const { title, description } = siteMetaData;

  const validationStatusColor =
    validation && validation.status === "NeedReview"
      ? "orange"
      : validation && validation.status === "Rejected"
      ? "red"
      : "green";

  return (
    <>
      <Helmet>
        <title>{`${title} - ${flavor.challenge.name} Challenge`}</title>
        <meta name="description" content={description} />
      </Helmet>
      <Page.Content
        title={flavor.challenge.name}
        subTitle={`Author: ${flavor.challenge.author}`}
      >
        <Grid.Row className="mb-3">
          {challenge && (
            <Grid.Col sm={12} md={7}>
              <ChallengeCard challenge={challenge} withModal={false} />
            </Grid.Col>
          )}
          <Grid.Col sm={12} md={5}>
            <div css={whiteWrapper}>
              <h3>
                <FormattedMessage id="ChallengeDetailsPage.solve" />
              </h3>
              <ChallengeSolveInstances
                instances={flavor.instances}
                purchased={purchased}
              />
              {purchased && !validations && (
                <Tag css={statusTag}>
                  <FormattedMessage id="ChallengeDetailsPage.purchased" />{" "}
                  {moment(subscription.created_at).calendar()}
                </Tag>
              )}
              {purchased && validations && (
                <Tag color={validationStatusColor}>
                  <FormattedMessage id="ChallengeDetailsPage.validated" />{" "}
                  {moment(validation.created_at).calendar()}
                </Tag>
              )}
            </div>
          </Grid.Col>
          {/* <Grid.Col width={12} sm={12} md={7}>
             {subscriptions && (
                <ChallengeCloseButton
                  challenge={challenge}
                  closeChallenge={closeChallenge}
                  isClosed={isClosed}
                />
              )}
          </Grid.Col> */}
        </Grid.Row>
        <Grid.Row>
          {purchased && !validations && (
            <Grid.Col width={12} sm={12} md={12}>
              <ChallengeValidateForm
                challenge={challenge}
                validateChallenge={validateChallenge}
                disabled={isClosed}
              />
            </Grid.Col>
          )}
        </Grid.Row>
      </Page.Content>
    </>
  );
};

ChallengeDetailsPage.propTypes = {
  fetchChallengeDetailAction: PropTypes.func,
};

export default React.memo(ChallengeDetailsPage);

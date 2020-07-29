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
  buyChallenge as buyChallengeAction,
  validateChallenge as validateChallengeAction,
  // closeChallenge as closeChallengeAction,
} from "../actions/seasons";
import ChallengeBuyButton from "../components/challenges/ChallengeBuyButton";
// import ChallengeCloseButton from "../components/challenges/ChallengeCloseButton";
import ChallengeValidateForm from "../components/challenges/ChallengeValidateForm";
import ValidationsList from "../components/challenges/ValidationsList";
import ChallengeSolveInstances from "../components/challenges/ChallengeSolveInstances";
import { CLEAN_CHALLENGE_DETAIL } from "../constants/actionTypes";

const paragraph = css`
  margin-top: 0.5rem;
`;

const statusStamp = css`
  font-size: 0.75rem;
  opacity: 0.7;
`;

const rewardText = css`
  font-size: 0.75rem;
  font-weight: 800;
  margin-top: 0.5rem;
`;

const ChallengeDetailsPage = props => {
  const dispatch = useDispatch();
  const challenge = useSelector(state => state.seasons.challengeInDetail);

  const buyChallenge = (flavorChallengeID, seasonID) =>
    dispatch(buyChallengeAction(flavorChallengeID, seasonID));
  const validateChallenge = (validationData, seasonId) =>
    dispatch(validateChallengeAction(validationData, seasonId));
  // const closeChallenge = subscriptionID =>
  //   dispatch(closeChallengeAction(subscriptionID));
  const fetchChallengeDetail = challengeID =>
    dispatch(fetchChallengeDetailAction(challengeID));

  useEffect(() => {
    const { uri, challengeID: challengeIDFromProps } = props;
    const challengeIDFromURI = uri && uri.split("/")[3];
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
  const isClosed = subscription && subscription.status === "Closed";
  const purchased = subscriptions;
  const { title, description } = siteMetaData;

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
        <Grid.Row className="mb-6">
          <Grid.Col width={12} sm={12} md={7}>
            <p css={paragraph}>
              Hello Ol&apos;salt! Try to beat the {flavor.challenge.name}{" "}
              challenge. Heave ho!
            </p>
          </Grid.Col>
          <Grid.Col md={5} sm={12} width={12} className="text-right">
            {purchased && (
              <Tag css={statusStamp}>
                purchased {moment(subscription.created_at).calendar()}
              </Tag>
            )}
            {!purchased && (
              <ChallengeBuyButton
                challenge={challenge}
                buyChallenge={buyChallenge}
              />
            )}
            <p css={rewardText}>Reward: {flavor.validation_reward}</p>
            {/* {subscriptions && (
                <ChallengeCloseButton
                  challenge={challenge}
                  closeChallenge={closeChallenge}
                  isClosed={isClosed}
                />
              )} */}
          </Grid.Col>
        </Grid.Row>
        {/* <hr /> */}
        <Grid.Row>
          <Grid.Col width={12} sm={12} md={12}>
            <h3>Solve challenge</h3>
            <ChallengeSolveInstances
              instances={flavor.instances}
              purchased={purchased}
            />
          </Grid.Col>
          {purchased && (
            <Grid.Col width={12} sm={12} md={12} className="text-right">
              <ChallengeValidateForm
                challenge={challenge}
                validateChallenge={validateChallenge}
                disabled={isClosed}
              />
              {validations && <ValidationsList validations={validations} />}
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

import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid, Dimmer, Button } from "tabler-react";
import {
  fetchChallengeDetail as fetchChallengeDetailAction,
  buyChallenge as buyChallengeAction,
  validateChallenge as validateChallengeAction,
  closeChallenge as closeChallengeAction,
} from "../actions/seasons";
import ChallengeBuyButton from "../components/challenges/ChallengeBuyButton";
import ChallengeValidateButton from "../components/challenges/ChallengeValidateButton";
import ChallengeCloseButton from "../components/challenges/ChallengeCloseButton";
import ValidationsList from "../components/challenges/ValidationsList";
import { CLEAN_CHALLENGE_DETAIL } from "../constants/actionTypes";

import styles from "./styles/ChallengeDetailsPage.module.css";

const ChallengeDetailsPage = props => {
  const dispatch = useDispatch();
  const challenge = useSelector(state => state.seasons.challengeInDetail);
  const activeTeam = useSelector(state => state.seasons.activeTeam);

  const buyChallenge = (challengeID, teamID, seasonId) =>
    dispatch(buyChallengeAction(challengeID, teamID, seasonId));
  const validateChallenge = (validationData, seasonId) =>
    dispatch(validateChallengeAction(validationData, seasonId));
  const closeChallenge = subscriptionID =>
    dispatch(closeChallengeAction(subscriptionID));
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

  const {
    flavor: { challenge: flavorChallenge, instances } = {
      challenge: "no challenge",
    },
    subscriptions,
  } = challenge || {};

  const teamID = activeTeam && activeTeam.id;

  const subscription = challenge.subscriptions && challenge.subscriptions[0];
  const validations = subscription && subscription.validations;
  const isClosed = subscription && subscription.status === "Closed";

  return (
    <Page.Content title={flavorChallenge.name}>
      <Grid.Row>
        <Grid.Col lg={4} md={4} sm={4} xs={4}>
          <h4>Author</h4>
          <p className={styles.p}>{flavorChallenge.author}</p>
          <Button
            href={flavorChallenge.homepage}
            target="_blank"
            RootComponent="a"
            social="github"
            size="sm"
          >
            Visit page
          </Button>
        </Grid.Col>
        <Grid.Col lg={4} md={4} sm={4} xs={4}>
          <h4>Actions</h4>
          <Button.List>
            <ChallengeBuyButton
              challenge={challenge}
              buyChallenge={buyChallenge}
              teamID={teamID}
              isClosed={isClosed}
            />
            <Button
              RootComponent="a"
              target="_blank"
              href={instances[0].nginx_url}
              color="gray-dark"
              icon="terminal"
              disabled={isClosed || !subscription}
            >
              Solve
            </Button>
            <ChallengeCloseButton
              challenge={challenge}
              closeChallenge={closeChallenge}
              isClosed={isClosed}
            />
          </Button.List>
        </Grid.Col>
      </Grid.Row>
      <hr />
      <Grid.Row>
        {subscriptions && (
          <Grid.Col lg={12} md={12} sm={12} xs={12}>
            <div style={{ marginBottom: "1rem" }}>
              <h3>Validations</h3>
              <ChallengeValidateButton
                challenge={challenge}
                validateChallenge={validateChallenge}
                disabled={isClosed}
              />
            </div>
            {validations && <ValidationsList validations={validations} />}
          </Grid.Col>
        )}
      </Grid.Row>
    </Page.Content>
  );
};

ChallengeDetailsPage.propTypes = {
  fetchChallengeDetailAction: PropTypes.func,
};

export default React.memo(ChallengeDetailsPage);

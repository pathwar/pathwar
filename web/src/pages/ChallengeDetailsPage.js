import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Helmet } from "react-helmet";
import { css } from "@emotion/core";
import PropTypes from "prop-types";
import { Page, Grid, Dimmer, Button, Icon } from "tabler-react";
import { Terminal } from "react-window-ui";
import siteMetaData from "../constants/metadata";
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

const paragraph = css`
  margin-top: 0.5rem;
`;

const ChallengeSolveInstances = ({ instances }) => {
  const terminal = css`
    margin-bottom: 1rem;

    a {
      color: #fff;
      margin-right: 1.5rem;
    }

    span {
      color: #16b279;
    }
  `;

  return (
    <Terminal
      minHeight="2rem"
      boxShadow="0px 2px 15px -8px rgba(0,0,0,0.41)"
      css={terminal}
    >
      {instances.map(item => (
        <div key={item.id}>
          <a href={item.nginx_url} target="_blank" rel="noreferrer">
            {item.nginx_url}
          </a>
          <span>
            <Icon name="check-circle" />
            {item.status}
          </span>
        </div>
      ))}
    </Terminal>
  );
};

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
  const { title, description } = siteMetaData;

  return (
    <>
      <Helmet>
        <title>{`${title} - ${flavorChallenge.name} Challenge`}</title>
        <meta name="description" content={description} />
      </Helmet>
      <Page.Content
        title={flavorChallenge.name}
        subTitle={`Author: ${flavorChallenge.author}`}
      >
        <Grid.Row className="mb-6">
          <Grid.Col width={12} sm={12} md={8}>
            <p css={paragraph}>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin sem
              arcu, tristique id elementum quis, pulvinar ac lorem. Integer id
              sem condimentum, aliquam erat in, lobortis nisi. Aliquam pretium
              mi purus. Donec sit amet neque nulla. Pellentesque mollis egestas
              nisl a placerat.
            </p>
            <Button.List>
              <Button
                href={flavorChallenge.homepage}
                target="_blank"
                RootComponent="a"
                social="github"
                size="sm"
              >
                Visit page
              </Button>
            </Button.List>
          </Grid.Col>
          <Grid.Col md={4} sm={12} width={12} className="text-right">
            <Button.List>
              <ChallengeBuyButton
                challenge={challenge}
                buyChallenge={buyChallenge}
                teamID={teamID}
                isClosed={isClosed}
              />
              {/* <Button
                RootComponent="a"
                target="_blank"
                href={instances[0].nginx_url}
                color="gray-dark"
                icon="terminal"
                disabled={isClosed || !subscription}
              >
                Solve
              </Button> */}
              {subscriptions && (
                <ChallengeCloseButton
                  challenge={challenge}
                  closeChallenge={closeChallenge}
                  isClosed={isClosed}
                />
              )}
            </Button.List>
          </Grid.Col>
        </Grid.Row>
        {/* <hr /> */}
        <Grid.Row>
          <Grid.Col width={12} sm={12} md={8}>
            <h3>Solve challenge</h3>
            <ChallengeSolveInstances instances={instances} />
          </Grid.Col>
          {subscriptions && (
            <Grid.Col width={12} sm={12} md={4}>
              <Button.List className="text-right">
                <ChallengeValidateButton
                  challenge={challenge}
                  validateChallenge={validateChallenge}
                  disabled={isClosed}
                />
              </Button.List>
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

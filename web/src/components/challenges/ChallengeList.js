/* eslint-disable react/prop-types */
import React from "react";
import { useSelector } from "react-redux";
import { Link } from "gatsby";
import { Button, Dimmer, Card, Grid } from "tabler-react";
import styles from "../../styles/layout/loader.module.css";

const ChallengeCard = ({ challenge, buyChallenge, teamID }) => {
  const { flavor, subscriptions, id: challengeID } = challenge;
  const isClosed = subscriptions && subscriptions[0].status === "Closed";

  return (
    <Card
      title={flavor.challenge.name}
      body={
        <Button.List align="center">
          <Button
            onClick={() => buyChallenge(challengeID, teamID, false)}
            size="sm"
            color={isClosed ? "red" : "success"}
            disabled={subscriptions || isClosed}
            icon={
              subscriptions ? (isClosed ? "x-circle" : "check") : "dollar-sign"
            }
          >
            {isClosed ? "Closed" : "Buy"}
          </Button>
          <Button
            RootComponent={Link}
            to={`/app/challenge/${challengeID}`}
            target="_blank"
            color="info"
            size="sm"
            icon="eye"
          >
            Open
          </Button>
          {subscriptions && flavor.instances && (
            <Button
              RootComponent="a"
              target="_blank"
              href={flavor.instances[0].nginx_url}
              size="sm"
              color="gray-dark"
              icon="terminal"
              disabled={isClosed}
            >
              Solve
            </Button>
          )}
        </Button.List>
      }
    />
  );
};

const ChallengeList = props => {
  const activeUserSession = useSelector(
    state => state.userSession.activeUserSession
  );
  const activeTeam = useSelector(state => state.seasons.activeTeam);

  const { challenges, buyChallenge } = props;

  return !challenges || !activeUserSession ? (
    <Dimmer className={styles.dimmer} active loader />
  ) : (
    <>
      <Grid.Row>
        {challenges.map(challenge => (
          <Grid.Col lg={4} sm={4} md={4} xs={4} key={challenge.id}>
            <ChallengeCard
              challenge={challenge}
              buyChallenge={buyChallenge}
              teamID={activeTeam.id}
            />
          </Grid.Col>
        ))}
      </Grid.Row>
    </>
  );
};

export default ChallengeList;

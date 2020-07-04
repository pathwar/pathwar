/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { useSelector } from "react-redux";
import { Link, navigate } from "gatsby";
import { Button, Dimmer, Card, Grid } from "tabler-react";
import { Modal } from "react-responsive-modal";
import styles from "../../styles/layout/loader.module.css";
import { func } from "prop-types";

const ChallengeCard = ({ challenge }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const { flavor, id: challengeID } = challenge;
  // const isClosed = subscriptions && subscriptions[0].status === "Closed";

  const onCloseModal = function() {
    setModalOpen(false);
  };

  const openModal = function() {
    setModalOpen(true);
  };

  const openChallengePage = function() {
    setModalOpen(false);
    navigate(`/app/challenge/${challengeID}`);
  };

  return (
    <Card
      title={flavor.challenge.name}
      statusColor="green"
      statusSide
      body={
        <>
          <Grid.Row>
            <Grid.Col auto>
              <p>Author: {flavor.challenge.author}</p>
            </Grid.Col>
          </Grid.Row>
          <Grid.Row>
            <Grid.Col auto>
              <p>Version: {flavor.version}</p>
            </Grid.Col>
          </Grid.Row>
          <Grid.Row>
            <Grid.Col auto>
              <Button.List align="center">
                <Button
                  color="indigo"
                  size="sm"
                  icon="compass"
                  onClick={openModal}
                >
                  View
                </Button>
                {/* <Button
                onClick={() => buyChallenge(challengeID, teamID, false)}
                size="sm"
                color={isClosed ? "red" : "success"}
                disabled={subscriptions || isClosed}
                icon={
                  subscriptions
                    ? isClosed
                      ? "x-circle"
                      : "check"
                    : "dollar-sign"
                }
              >
                {isClosed ? "Closed" : "Buy"}
              </Button> */}
                {/* {subscriptions && flavor.instances && (
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
                )}*/}
              </Button.List>
            </Grid.Col>
          </Grid.Row>
          <Modal open={modalOpen} onClose={onCloseModal}>
            <h2>Modal</h2>
            <Button
              color="info"
              size="sm"
              icon="eye"
              onClick={openChallengePage}
            >
              Challenge page
            </Button>
          </Modal>
        </>
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

/* eslint-disable react/prop-types */
import React from "react";
import { useSelector } from "react-redux";
import { Dimmer, Grid } from "tabler-react";
import ChallengeCard from "./ChallengeCard";
import styles from "../../styles/layout/loader.module.css";

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

export default React.memo(ChallengeList);

/* eslint-disable react/prop-types */
import React from "react";
import { useSelector } from "react-redux";
import { Dimmer, Grid } from "tabler-react";
import ChallengeCard from "./ChallengeCard";

const ChallengeList = props => {
  const activeTeam = useSelector(state => state.seasons.activeTeam);

  const { challenges } = props;

  return !challenges ? (
    <Dimmer active loader />
  ) : (
    <>
      <Grid.Row>
        {challenges.map(challenge => (
          <Grid.Col lg={4} sm={4} md={4} xs={4} key={challenge.id}>
            <ChallengeCard challenge={challenge} teamID={activeTeam.id} />
          </Grid.Col>
        ))}
      </Grid.Row>
    </>
  );
};

export default React.memo(ChallengeList);

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
          <Grid.Col lg={6} sm={6} md={6} xs={12} key={challenge.id}>
            <ChallengeCard challenge={challenge} teamID={activeTeam.id} />
          </Grid.Col>
        ))}
      </Grid.Row>
    </>
  );
};

export default React.memo(ChallengeList);

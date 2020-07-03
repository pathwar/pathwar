/* eslint-disable react/prop-types */
import React from "react";
import { useSelector } from "react-redux";
import { Grid, Page } from "tabler-react";
import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";

const StatisticsPage = props => {
  const activeSeason = useSelector(state => state.seasons.activeSeason);
  const allTeamsOnSeason = useSelector(state => state.seasons.allTeamsOnSeason);

  return (
    <Page.Content title={`Statistics`}>
      <Grid.Row cards={true}>
        <Grid.Col xs={12} sm={12} lg={6}>
          <AllTeamsOnSeasonList
            activeSeason={activeSeason}
            allTeamsOnSeason={allTeamsOnSeason}
          />
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default StatisticsPage;

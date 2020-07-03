/* eslint-disable react/prop-types */
import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Grid, Page } from "tabler-react";
import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";
import { fetchAllSeasonTeams as fetchAllSeasonTeamsAction } from "../actions/seasons";

const StatisticsPage = () => {
  const dispatch = useDispatch();
  const activeSeason = useSelector(state => state.seasons.activeSeason);
  const allTeamsOnSeason = useSelector(state => state.seasons.allTeamsOnSeason);

  useEffect(
    state => {
      console.log(state);
      if (!allTeamsOnSeason && activeSeason) {
        dispatch(fetchAllSeasonTeamsAction(activeSeason.id));
      }
    },
    [activeSeason, allTeamsOnSeason]
  );

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

/* eslint-disable react/prop-types */
import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Grid, Page } from "tabler-react";
import siteMetaData from "../constants/metadata";
import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";
import CreateTeamButton from "../components/team/CreateTeamButton";
import {
  fetchAllSeasonTeams as fetchAllSeasonTeamsAction,
  createTeam as createTeamAction,
} from "../actions/seasons";

const StatisticsPage = () => {
  const dispatch = useDispatch();
  const activeSeason = useSelector(state => state.seasons.activeSeason);
  const activeTeamInSeason = useSelector(
    state => state.seasons.activeTeamInSeason
  );
  const allTeamsOnSeason = useSelector(state => state.seasons.allTeamsOnSeason);
  const dispatchCreateTeamAction = dispatch(createTeamAction);

  useEffect(
    state => {
      console.log(state);
      if (!allTeamsOnSeason && activeSeason) {
        dispatch(fetchAllSeasonTeamsAction(activeSeason.id));
      }
    },
    [activeSeason, allTeamsOnSeason]
  );

  const { title, description } = siteMetaData;

  return (
    <>
      <Helmet>
        <title>{title} - Statistics</title>
        <meta name="description" content={description} />
      </Helmet>
      <Page.Content title={`Statistics`}>
        <Grid.Row cards={true}>
          <Grid.Col auto>
            <CreateTeamButton
              activeSeason={activeSeason}
              createTeam={dispatchCreateTeamAction}
              activeTeamInSeason={activeTeamInSeason}
            />
          </Grid.Col>
          <Grid.Col xs={12} sm={12} lg={9}>
            <AllTeamsOnSeasonList
              activeSeason={activeSeason}
              allTeamsOnSeason={allTeamsOnSeason}
            />
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    </>
  );
};

export default StatisticsPage;

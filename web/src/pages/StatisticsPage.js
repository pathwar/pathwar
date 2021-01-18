import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Helmet } from "react-helmet";
import { Grid, Page } from "tabler-react";
import { useIntl } from "react-intl";
import siteMetaData from "../constants/metadata";
import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";
import CreateTeamButton from "../components/team/CreateTeamButton";
import {
  fetchAllSeasonTeams as fetchAllSeasonTeamsAction,
  createTeam as createTeamAction,
} from "../actions/seasons";

const StatisticsPage = () => {
  const intl = useIntl();
  const dispatch = useDispatch();
  const activeSeason = useSelector(state => state.seasons.activeSeason);
  const activeTeamInSeason = useSelector(
    state => state.seasons.activeTeamInSeason
  );
  const allTeamsOnSeason = useSelector(state => state.seasons.allTeamsOnSeason);
  const dispatchCreateTeamAction = dispatch(createTeamAction);

  useEffect(() => {
    if (!allTeamsOnSeason && activeSeason) {
      dispatch(fetchAllSeasonTeamsAction(activeSeason.id));
    }
  }, [activeSeason, allTeamsOnSeason, dispatch]);

  const { title, description } = siteMetaData;
  const statisticsIntl = intl.formatMessage({ id: "nav.statistics" });

  return (
    <>
      <Helmet>
        <title>
          {title} - {statisticsIntl}
        </title>
        <meta name="description" content={description} />
      </Helmet>
      <Page.Content title={statisticsIntl}>
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={12} md={3}>
            <CreateTeamButton
              activeSeason={activeSeason}
              createTeam={dispatchCreateTeamAction}
              activeTeamInSeason={activeTeamInSeason}
            />
          </Grid.Col>
          <Grid.Col xs={12} sm={12} md={9}>
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

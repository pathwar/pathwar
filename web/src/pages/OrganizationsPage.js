import {useIntl} from "react-intl";
import {Grid, Page} from "tabler-react";
import React, {useEffect} from "react";
import {Helmet} from "react-helmet";
import siteMetaData from "../constants/metadata";
import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";
import {useDispatch, useSelector} from "react-redux";
import {createTeam, fetchAllSeasonTeams as fetchAllSeasonTeamsAction} from "../actions/seasons";
import CreateTeamButton from "../components/team/CreateTeamButton";

//TODO: Lister les organisations de l'utilisateur dans un tableau
//TODO: Créer un boutton permettant de créer une organisation pour l'utilisateur
//TODO: Le button devrait ouvrier un modal permettant de créer une organisation
//TODO: Lister les invitations de l'utilisateur dans un tableau
const OrganizationsPage = () => {
  const intl = useIntl();
  const organizationsIntl = intl.formatMessage({ id: "nav.organizations" });
  const pageTitleIntl = intl.formatMessage({ id: "OrganizationsPage.title" });
  const { title, description } = siteMetaData;
  const dispatch = useDispatch();

  const activeSeason = useSelector(state => state.seasons.activeSeason);
  // const activeTeamInSeason = useSelector(
  //   state => state.seasons.activeTeamInSeason
  // );
  const allTeamsOnSeason = useSelector(state => state.seasons.allTeamsOnSeason);
  // const dispatchCreateTeamAction = dispatch(createTeamAction);

  useEffect(() => {
    if (!allTeamsOnSeason && activeSeason) {
      dispatch(fetchAllSeasonTeamsAction(activeSeason.id));
    }
  }, [activeSeason, allTeamsOnSeason, dispatch]);

  const activeTeamInSeason = useSelector(
     state => state.seasons.activeTeamInSeason
  );

  const dispatchCreateTeamAction = dispatch(createTeam);

  return (
    <>
    <Helmet>
      <title>
        {title} - {organizationsIntl}
      </title>
      <meta name="description" content={description} />
    </Helmet>
    <Page.Content title={pageTitleIntl}>
      <Grid.Row>
        <Grid.Col offset={10}>
          <CreateTeamButton
            activeSeason={activeSeason}
            createTeam={dispatchCreateTeamAction}
            activeTeamInSeason={activeTeamInSeason}
          />
        </Grid.Col>
      </Grid.Row>
      <Grid.Row cards={true}>
        <Grid.Col xs={12} sm={12} md={6}>
          <AllTeamsOnSeasonList
            activeSeason={activeSeason}
            allTeamsOnSeason={allTeamsOnSeason}
          />
        </Grid.Col>
        <Grid.Col xs={12} sm={12} md={6}>
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

export default React.memo(OrganizationsPage);

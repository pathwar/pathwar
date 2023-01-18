import {useIntl} from "react-intl";
import {Grid, Page} from "tabler-react";
import React, {useEffect} from "react";
import {Helmet} from "react-helmet";
import siteMetaData from "../constants/metadata";
import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";
import {useDispatch, useSelector} from "react-redux";
import CreateOrganizationButton from "../components/organization/CreateOrganizationButton";
import {fetchOrganizationsList} from "../actions/organizations";
import UserOrganizationsList from "../components/organization/UserOrganizationsList";

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
  const userOrganizations = useSelector(state => state.organizations.userOrganizationsList);
  // const dispatchCreateTeamAction = dispatch(createTeamAction);

  useEffect(() => {
    if (!userOrganizations) {
      dispatch(fetchOrganizationsList());
    }
  }, [userOrganizations, dispatch]);


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
          <CreateOrganizationButton />
        </Grid.Col>
      </Grid.Row>
      <Grid.Row cards={true}>
        <Grid.Col xs={12} sm={12} md={6}>
          <UserOrganizationsList
            userOrganizationsList={userOrganizations}
          />
        </Grid.Col>
        <Grid.Col xs={12} sm={12} md={6}>
          <UserOrganizationsList
            userOrganizationsList={userOrganizations}
          />
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
    </>
  );
};

export default React.memo(OrganizationsPage);

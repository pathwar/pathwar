import * as React from "react";

import { Page, Grid } from "tabler-react";

import AllOrganizationsList from "../components/organizations/AllOrganizationsList";
// import UserTeamsCard from "../components/teams/UserTeamsList";

const HomePage = () => {
  return (
    <Page.Content title="Dashboard">
      <Grid.Row cards={true}>
        <Grid.Col xs={12} sm={12} lg={6}>
          <AllOrganizationsList />
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default HomePage;

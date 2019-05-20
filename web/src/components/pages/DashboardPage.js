import * as React from "react";

import {
  Page,
  Grid
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import AllTeamsCard from "../teams/AllTeamsList";
import UserTeamsCard from "../teams/UserTeamsList";
import TeamStatsStampCard from "../teams/TeamsStatsStampCard";

const Dashboard = () => {
  return (
    <SiteWrapper>
      <Page.Content title="Dashboard">
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={12} lg={6}>
            <UserTeamsCard />
            <TeamStatsStampCard />
          </Grid.Col>

          <Grid.Col xs={12} sm={12} lg={6}>
            <AllTeamsCard />
          </Grid.Col>

        </Grid.Row>
      </Page.Content>
    </SiteWrapper>
  );
}

export default Dashboard;
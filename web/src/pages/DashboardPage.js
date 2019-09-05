import * as React from "react";

import {
  Page,
  Grid
} from "tabler-react";

import AllTeamsCard from "../components/teams/AllTeamsList";
import UserTeamsCard from "../components/teams/UserTeamsList";
import TeamStatsStampCard from "../components/teams/TeamsStatsStampCard";

const Dashboard = () => {
  return (
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
  );
}

export default Dashboard;

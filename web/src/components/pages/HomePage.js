import * as React from "react";

import {
  Page,
  Grid
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import TeamsCard from "../teams/TeamsCards";
import TeamStatsStampCard from "../teams/TeamsStatsStampCard";

function Home() {
  return (
    <SiteWrapper>
      <Page.Content title="Dashboard">
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={8} lg={6}>
            <TeamsCard />
          </Grid.Col>
          
          <Grid.Col xs={12} sm={4} lg={3}>
            <TeamStatsStampCard />
          </Grid.Col>

        </Grid.Row>
      </Page.Content>
    </SiteWrapper>
  );
}

export default Home;
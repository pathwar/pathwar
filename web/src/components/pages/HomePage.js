import * as React from "react";

import {
  Page,
  Grid,
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import TeamsCard from "../teams/TeamsCards";

function Home() {
  return (
    <SiteWrapper>
      <Page.Content title="Dashboard">
        <Grid.Row cards={true}>
          <Grid.Col lg={6}>
            <TeamsCard />
          </Grid.Col>
        </Grid.Row>
      </Page.Content></SiteWrapper>
  );
}

export default Home;
import * as React from "react";

import {
  Page,
  Grid
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";

const CompetitionsPage = () => {
  return (
    <SiteWrapper>
      <Page.Content title="Competitions">
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={8} lg={6}>
          </Grid.Col>
          
          <Grid.Col xs={12} sm={4} lg={3}>
          </Grid.Col>

        </Grid.Row>
      </Page.Content>
    </SiteWrapper>
  );
}

export default CompetitionsPage;
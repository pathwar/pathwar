import * as React from "react";

import {
  Page,
  Grid
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import LevelsCardPreview from "../levels/LevelCardPreview";

const CompetitionsPage = () => {
  return (
    <SiteWrapper>
      <Page.Content title="Competitions">
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={8} lg={6}>
            <h3>Levels</h3>
            <LevelsCardPreview />
          </Grid.Col>

        </Grid.Row>
      </Page.Content>
    </SiteWrapper>
  );
}

export default CompetitionsPage;
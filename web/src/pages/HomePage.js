import React from "react";
import { Page, Grid } from "tabler-react";
import { useSelector } from "react-redux";
import moment from "moment";

const HomePage = () => {
  const activeUserSession = useSelector(
    state => state.userSession.activeUserSession
  );
  const {
    user: { active_team_member: activeTeam },
  } = activeUserSession;

  console.log(activeTeam);

  return (
    <Page.Content title="Home">
      <Grid.Row>
        <Grid.Col xs={12} sm={12} lg={6}>
          <h3>created at:</h3>
          <p>{moment(activeTeam.created_at).format("ll")}</p>
          <h3>role:</h3>
          <p>{activeTeam.role}</p>
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(HomePage);

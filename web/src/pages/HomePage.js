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
  } = activeUserSession || { user: {} };

  const createDate = activeTeam && activeTeam.created_at;
  const role = activeTeam && activeTeam.role;

  return (
    <Page.Content title="Home">
      <Grid.Row>
        <Grid.Col xs={12} sm={12} lg={6}>
          <h3>created at:</h3>
          <p>{moment(createDate).format("ll")}</p>
          <h3>role:</h3>
          <p>{role}</p>
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(HomePage);

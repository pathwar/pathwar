import React from "react";
import { Page, Grid } from "tabler-react";
import { useSelector } from "react-redux";
import { useIntl, FormattedMessage } from "react-intl";
import moment from "moment";

const HomePage = () => {
  const intl = useIntl();

  const activeUserSession = useSelector(
    state => state.userSession.activeUserSession
  );
  const {
    user: { active_team_member: activeTeam },
  } = activeUserSession || { user: {} };

  const createDate = activeTeam && activeTeam.created_at;
  const role = activeTeam && activeTeam.role;

  const pageTitleIntl = intl.formatMessage({ id: "HomePage.title" });

  return (
    <Page.Content title={pageTitleIntl}>
      <Grid.Row>
        <Grid.Col xs={12} sm={12} lg={6}>
          <h3>
            <FormattedMessage id="HomePage.createdAt" />
          </h3>
          <p>{moment(createDate).format("ll")}</p>
          <h3>role:</h3>
          <p>{role}</p>
        </Grid.Col>
      </Grid.Row>
    </Page.Content>
  );
};

export default React.memo(HomePage);

import React, { useEffect } from "react";
import { Page, Grid, Avatar } from "tabler-react";
import { useSelector, useDispatch } from "react-redux";
import { useIntl, FormattedMessage } from "react-intl";
import moment from "moment";
import ShadowBox from "../components/ShadowBox";

import { fetchTeamDetails as fetchTeamDetailsAction } from "../actions/seasons";

const HomePage = () => {
  const intl = useIntl();
  const dispatch = useDispatch();

  const activeUserSession = useSelector(
    state => state.userSession.activeUserSession
  );
  const {
    user: {
      active_team_member: activeTeam,
      gravatar_url,
      username,
      email,
      created_at,
    },
  } = activeUserSession || { user: {} };

  const role = activeTeam && activeTeam.role;

  const pageTitleIntl = intl.formatMessage({ id: "HomePage.title" });

  useEffect(() => {
    if (activeTeam) {
      dispatch(fetchTeamDetailsAction(activeTeam.team_id));
    }
  }, [activeTeam, dispatch]);

  return (
    <Page.Content title={pageTitleIntl}>
      <div>
        <Grid.Row>
          <Grid.Col width={12} lg={4}>
            <ShadowBox>
              <div
                css={{
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                }}
              >
                {gravatar_url ? (
                  <Avatar size="xxl" imageURL={`${gravatar_url}?d=identicon`} />
                ) : (
                  <Avatar size="xxl" icon="users" />
                )}
                <h2 className="mb-0 mt-2">{username}</h2>
                <p>{email}</p>
                <h3>
                  <FormattedMessage id="HomePage.createdAt" />
                </h3>
                <p>{moment(created_at).format("ll")}</p>
              </div>
            </ShadowBox>
          </Grid.Col>
          <Grid.Col width={12} lg={8}>
            <ShadowBox></ShadowBox>
          </Grid.Col>
        </Grid.Row>
      </div>
    </Page.Content>
  );
};

export default React.memo(HomePage);

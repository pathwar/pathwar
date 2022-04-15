import React, { useEffect } from "react";
import { Page, Grid, Avatar, Dimmer, ProgressCard } from "tabler-react";
import { useSelector, useDispatch } from "react-redux";
import { useIntl, FormattedMessage } from "react-intl";
import moment from "moment";
import ShadowBox from "../components/ShadowBox";
import iconPwn from "../images/icon-pwn-small.svg";

import {
  fetchChallenges,
  fetchTeamDetails as fetchTeamDetailsAction,
} from "../actions/seasons";
import { isNil } from "ramda";
import UserChallengesView from "../components/home/UserChallengesView";

const cardStyle = {
  margin: "1rem",
  ".progress": {
    display: "none",
  },
  ".display-4": {
    margin: "0 !important",
  },
};

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

  const pageTitleIntl = intl.formatMessage({ id: "HomePage.title" });

  const activeChallenges = useSelector(state => state.seasons.activeChallenges);
  const activeSeason = useSelector(state => state.seasons.activeSeason);
  const teamDetails = useSelector(state => state.seasons.teamInDetail);

  useEffect(() => {
    if (activeTeam && !teamDetails) {
      dispatch(fetchTeamDetailsAction(activeTeam.team_id));
    }
  }, [activeTeam, dispatch, teamDetails]);

  useEffect(() => {
    if (isNil(activeChallenges) && activeSeason) {
      dispatch(fetchChallenges(activeSeason.id));
    }
  }, [activeChallenges, activeSeason, dispatch]);

  if (!activeUserSession) {
    return <Dimmer active loader />;
  }

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
            <ShadowBox>
              <h2>Season stats</h2>
              {teamDetails && (
                <div css={{ display: "flex", justifyContent: "space-around" }}>
                  {/* <span
                    css={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "center",
                      fontWeight: "bold",
                    }}
                  >
                    Score:
                    <span>{teamDetails.score}</span>
                  </span> */}
                  <ProgressCard
                    css={cardStyle}
                    content={teamDetails.score}
                    header="Score"
                  />
                  <ProgressCard
                    css={cardStyle}
                    content={`$${teamDetails.cash}`}
                    header={
                      <>
                        Cash{" "}
                        <img
                          css={{ display: "inline-block" }}
                          src={iconPwn}
                          className="img-responsive"
                        />
                      </>
                    }
                  />
                  {/* <span
                    css={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "center",
                      fontWeight: "bold",
                    }}
                  >
                    Cash:
                    <span>${teamDetails.cash}</span>
                    <span>
                      <img src={iconPwn} className="img-responsive" />
                    </span>
                  </span> */}
                </div>
              )}
            </ShadowBox>

            <ShadowBox>
              <UserChallengesView challenges={activeChallenges} />
            </ShadowBox>
          </Grid.Col>
        </Grid.Row>
      </div>
    </Page.Content>
  );
};

export default React.memo(HomePage);

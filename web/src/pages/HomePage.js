import React, { useEffect, useState } from "react";
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
import getTeamRank from "../utils/getTeamRank";
import SwitchSeasonInput from "../components/season/SwitchSeasonInput";
import UserOrganizationBadges from "../components/organization/UserOrganizationBadges";

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
  const userOrganizations = useSelector(state => state.organizations.userOrganizationsList);

  const [rank, setRank] = useState();

  useEffect(() => {
    if (activeTeam && activeSeason && !rank) {
      const fetchRank = async () => {
        const calculatedRank = await getTeamRank(
          activeTeam.team_id,
          activeSeason.id
        );
        setRank(calculatedRank);
      };

      fetchRank();
    }
  }, [activeSeason, activeTeam, rank]);

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
                <h3 className="mb-0 mt-2">
                  <FormattedMessage id="HomePage.createdAt" />
                </h3>
                <p>{moment(created_at).format("ll")}</p>
                <h3 className="mb-0 mt-2">
                  <FormattedMessage id="HomePage.activeSeason" />
                </h3>
                <p>{activeSeason.slug}</p>
                <h3 className="mb-2 mt-2">
                  <FormattedMessage id="HomePage.organizations" />
                </h3>
                <UserOrganizationBadges organizations={userOrganizations}/>
              </div>
            </ShadowBox>
          </Grid.Col>
          <Grid.Col width={12} lg={8}>
            <ShadowBox>
              <h2>
                <FormattedMessage id="HomePage.statsTitle" />
              </h2>
              {teamDetails && (
                <div css={{ display: "flex", justifyContent: "space-around" }}>
                  {/* <ProgressCard
                    css={cardStyle}
                    content={`# ${rank || "-"}`}
                    header={intl.formatMessage({ id: "HomePage.rank" })}
                  />
                  <ProgressCard
                    css={cardStyle}
                    content={teamDetails.score ? teamDetails.score : 0}
                    header={intl.formatMessage({ id: "HomePage.score" })}
                  /> */}
                  <ProgressCard
                    css={cardStyle}
                    content={`$${teamDetails.cash ? teamDetails.cash : 0}`}
                    header={
                      <>
                        <FormattedMessage id="HomePage.cash" />{" "}
                        <img
                          css={{ display: "inline-block" }}
                          src={iconPwn}
                          className="img-responsive"
                        />
                      </>
                    }
                  />
                </div>
              )}
            </ShadowBox>
            <UserChallengesView challenges={activeChallenges} />
            <ShadowBox>
              <h2>
                <FormattedMessage id="HomePage.switchSeason" />
              </h2>
                <>
                  <div
                    css={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "center",
                    }}
                  >
                  <SwitchSeasonInput />
                  </div>
                </>
            </ShadowBox>
          </Grid.Col>
        </Grid.Row>
      </div>
    </Page.Content>
  );
};

export default React.memo(HomePage);

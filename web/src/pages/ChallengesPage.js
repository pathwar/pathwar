import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid, Button } from "tabler-react";
import { isNil } from "ramda";

import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList";
import ChallengeList from "../components/challenges/ChallengeList";
import ValidateCouponButton from "../components/coupon/ValidateCouponButton";
import CreateTeamButton from "../components/team/CreateTeamButton";

import {
  fetchChallenges as fetchChallengesAction,
  fetchAllSeasonTeams as fetchAllSeasonTeamsAction,
  buyChallenge as buyChallengeAction,
  createTeam as createTeamAction,
} from "../actions/seasons";

class SeasonPage extends React.Component {
  componentDidUpdate(prevProps) {
    const {
      fetchAllSeasonTeamsAction,
      fetchChallengesAction,
      activeSeason,
      allTeamsOnSeason,
      activeChallenges,
    } = this.props;

    const { activeSeason: prevActiveSeason } = prevProps;

    if (
      (isNil(prevActiveSeason) && activeSeason) ||
      prevActiveSeason.id === activeSeason.id
    ) {
      if (isNil(allTeamsOnSeason)) {
        fetchAllSeasonTeamsAction(activeSeason.id);
      }

      if (isNil(activeChallenges)) {
        fetchChallengesAction(activeSeason.id);
      }
    }
  }

  render() {
    const {
      buyChallengeAction,
      createTeamAction,
      activeSeason,
      activeChallenges,
      allTeamsOnSeason,
      activeTeamInSeason,
    } = this.props;
    const name = activeSeason ? activeSeason.name : undefined;

    return (
      <Page.Content title="Season" subTitle={name}>
        <Grid.Row>
          <Grid.Col lg={4} md={4} sm={4} xs={4}>
            <Button.List>
              <ValidateCouponButton />
            </Button.List>
          </Grid.Col>
        </Grid.Row>
        <hr />
        <Grid.Row>
          <Grid.Col xs={12} sm={3} lg={3}>
            <h4>Teams</h4>
            <CreateTeamButton
              activeSeason={activeSeason}
              createTeam={createTeamAction}
              activeTeamInSeason={activeTeamInSeason}
            />
            <AllTeamsOnSeasonList
              activeSeason={activeSeason}
              allTeamsOnSeason={allTeamsOnSeason}
            />
          </Grid.Col>
          <Grid.Col xs={12} sm={9} lg={9}>
            <h4>Challenges</h4>
            <ChallengeList
              challenges={activeChallenges}
              buyChallenge={buyChallengeAction}
            />
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    );
  }
}

SeasonPage.propTypes = {
  seasons: PropTypes.object,
  fetchChallengesAction: PropTypes.func,
};

const mapStateToProps = state => ({
  activeSeason: state.seasons.activeSeason,
  activeChallenges: state.seasons.activeChallenges,
  allTeamsOnSeason: state.seasons.allTeamsOnSeason,
  activeTeamInSeason: state.seasons.activeTeamInSeason,
});

const mapDispatchToProps = {
  fetchChallengesAction: seasonID => fetchChallengesAction(seasonID),
  fetchAllSeasonTeamsAction: seasonID => fetchAllSeasonTeamsAction(seasonID),
  buyChallengeAction: (seasonID, teamID) =>
    buyChallengeAction(seasonID, teamID),
  createTeamAction: (seasonID, name) => createTeamAction(seasonID, name),
};

export default connect(mapStateToProps, mapDispatchToProps)(SeasonPage);

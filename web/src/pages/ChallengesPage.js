import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid, Button } from "tabler-react";
import { isNil } from "ramda";

import ChallengeList from "../components/challenges/ChallengeList";
import ValidateCouponButton from "../components/coupon/ValidateCouponButton";

import {
  fetchChallenges as fetchChallengesAction,
  fetchAllSeasonTeams as fetchAllSeasonTeamsAction,
  buyChallenge as buyChallengeAction,
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
    const { buyChallengeAction, activeSeason, activeChallenges } = this.props;
    const name = activeSeason ? activeSeason.name : undefined;

    return (
      <Page.Content title="Challenges" subTitle={name}>
        {/* <Grid.Row>
          <Grid.Col lg={4} md={4} sm={4} xs={4}>
            <Button.List>
              <ValidateCouponButton />
            </Button.List>
          </Grid.Col>
        </Grid.Row>
        <hr /> */}
        <Grid.Row>
          <Grid.Col xs={12} sm={12} lg={12}>
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
});

const mapDispatchToProps = {
  fetchChallengesAction: seasonID => fetchChallengesAction(seasonID),
  fetchAllSeasonTeamsAction: seasonID => fetchAllSeasonTeamsAction(seasonID),
  buyChallengeAction: (seasonID, teamID) =>
    buyChallengeAction(seasonID, teamID),
};

export default connect(mapStateToProps, mapDispatchToProps)(SeasonPage);

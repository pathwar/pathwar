import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid } from "tabler-react";
import { isNil } from "ramda";

import AllTeamsOnTournamentList from "../components/tournament/AllTeamsOnTournamentList"
import ChallengesCardPreview from "../components/challenges/ChallengeCardPreview";
import ValidationCouponStamp from "../components/coupon/ValidateCouponStampCard";
import {
  fetchChallenges as fetchChallengesAction ,
  fetchAllTournamentTeams as fetchAllTournamentTeamsAction
} from "../actions/tournaments";

class TournamentPage extends React.Component {

    componentDidUpdate(prevProps) {
      const { fetchChallengesAction, fetchAllTournamentTeamsAction, tournaments: { activeTournament } } = this.props;
      const { tournaments: { activeTournament: prevActiveTournament } } = prevProps;

      if (isNil(prevActiveTournament) && activeTournament ) {
        fetchChallengesAction(activeTournament.id);
        fetchAllTournamentTeamsAction(activeTournament.id)
      }
    }

    render() {
        const { tournaments: { activeTournament, activeChallenges, allTeamsOnTournament } } = this.props;
        const name = activeTournament ? activeTournament.name : undefined;

        return (
            <Page.Content title="Tournament" subTitle={name}>
                <Grid.Row>
                  <Grid.Col xs={12} sm={3} lg={3}>
                    <h3>Teams</h3>
                    <AllTeamsOnTournamentList activeTournament={activeTournament} allTeamsOnTournament={allTeamsOnTournament} />
                  </Grid.Col>
                  <Grid.Col xs={12} sm={6} lg={6}>
                    <h3>Challenges</h3>
                    <ChallengesCardPreview challenges={activeChallenges} />
                  </Grid.Col>
                  <Grid.Col xs={12} sm={3} lg={3}>
                    <h3>Actions</h3>
                    <ValidationCouponStamp />
                  </Grid.Col>
                </Grid.Row>
              </Page.Content>
          );
    }
}

TournamentPage.propTypes = {
    tournaments: PropTypes.object,
    activeTeam: PropTypes.object,
    fetchChallengesAction: PropTypes.func
};

const mapStateToProps = state => ({
    tournaments: state.tournaments
});

const mapDispatchToProps = {
  fetchChallengesAction: (tournamentID) => fetchChallengesAction(tournamentID),
  fetchAllTournamentTeamsAction: (tournamentID) => fetchAllTournamentTeamsAction(tournamentID)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(TournamentPage);

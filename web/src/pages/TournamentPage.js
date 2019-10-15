import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid } from "tabler-react";

import ChallengesCardPreview from "../components/challenges/ChallengeCardPreview";
import ValidationCouponStamp from "../components/coupon/ValidateCouponStampCard";
import { fetchChallenges as fetchChallengesAction } from "../actions/tournaments";

class TournamentPage extends React.Component {

    componentDidMount() {
      const { fetchChallengesAction, tournaments: { activeTournament } } = this.props;
      fetchChallengesAction(activeTournament.id);
    }

    render() {
        const { tournaments: { activeTournament, activeChallenges } } = this.props;
        const name = activeTournament ? activeTournament.name : undefined;

        return (
            <Page.Content title="Tournament" subTitle={name}>
                <Grid.Row>
                </Grid.Row>
                <Grid.Row>
                  <Grid.Col xs={12} sm={8} lg={6}>
                    <h3>Challenges</h3>
                    {activeChallenges && <ChallengesCardPreview challenges={activeChallenges} />}
                  </Grid.Col>
                  <Grid.Col xs={12} sm={4} lg={3}>
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
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(TournamentPage);

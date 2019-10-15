import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid } from "tabler-react";

import AllTeamsOnTournamentList from "../components/tournament/AllTeamsOnTournamentList"
import LevelsCardPreview from "../components/levels/LevelCardPreview";
import ValidationCouponStamp from "../components/coupon/ValidateCouponStampCard";
import { fetchLevels as fetchLevelsAction } from "../actions/tournaments";

class TournamentPage extends React.Component {

    componentDidMount() {
      const { fetchLevelsAction, tournaments: { activeTournament } } = this.props;
      // fetchLevelsAction(activeTournament.id);
    }

    render() {
        const { tournaments: { activeTournament, activeLevels } } = this.props;
        const name = activeTournament ? activeTournament.name : undefined;

        return (
            <Page.Content title="Tournament" subTitle={name}>
                <Grid.Row>
                  <Grid.Col xs={12} sm={6} lg={6}>
                    <h3>Teams</h3>
                    {activeTournament && <AllTeamsOnTournamentList />}
                  </Grid.Col>
                  <Grid.Col xs={12} sm={6} lg={6}>
                    <h3>Levels</h3>
                    {activeLevels && <LevelsCardPreview levels={activeLevels} />}
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
    fetchLevelsAction: PropTypes.func
};

const mapStateToProps = state => ({
    tournaments: state.tournaments
});

const mapDispatchToProps = {
  fetchLevelsAction: (tournamentID) => fetchLevelsAction(tournamentID),
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(TournamentPage);

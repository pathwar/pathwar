import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {
  Page,
  Grid
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import AllTournamentsList from "../tournament/AllTournamentsList";
import AllTeamTournamentsList from "../tournament/AllTeamTournamentsList";
import { 
  fetchAllTournaments as fetchAllTournamentsAction,
  setActiveTournament as setActiveTournamentAction
} from "../../actions/tournaments"

class AllTournamentsPage extends React.PureComponent {

    componentDidMount() {
        const { fetchAllTournamentsAction } = this.props;
        fetchAllTournamentsAction();
    }

    render() {
        const { 
          tournaments: { allTournaments, allTeamTournaments, activeTournament },
          activeTeam,
          setActiveTournamentAction
        } = this.props;
        return (
            <SiteWrapper>
              <Page.Content title="All Tournaments">
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    { allTournaments && <AllTournamentsList tournaments={allTournaments} /> }
                  </Grid.Col>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    { allTeamTournaments && 
                      <AllTeamTournamentsList 
                        teamTournaments={allTeamTournaments}
                        activeTournament={activeTournament}
                        setActive={setActiveTournamentAction}
                        activeTeam={activeTeam}
                      /> 
                    }
                  </Grid.Col>
                </Grid.Row>
              </Page.Content>
            </SiteWrapper>
          );
    }
}

AllTournamentsPage.propTypes = {
    tournaments: PropTypes.object,
    activeTeam: PropTypes.object,
    fetchAllTournamentsAction: PropTypes.func,
    setActiveTournamentAction: PropTypes.func
};

const mapStateToProps = state => ({
    tournaments: state.tournaments,
    activeTeam: state.teams.activeTeam
});

const mapDispatchToProps = {
    fetchAllTournamentsAction: () => fetchAllTournamentsAction(),
    setActiveTournamentAction: (teamID, tournament) => setActiveTournamentAction(teamID, tournament)
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(AllTournamentsPage);
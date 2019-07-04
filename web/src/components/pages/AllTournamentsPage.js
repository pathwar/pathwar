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
  fetchTeamTournaments as fetchTeamTournamentsAction
} from "../../actions/tournaments"

class AllTournamentsPage extends React.PureComponent {

    componentDidMount() {
        const { fetchAllTournamentsAction } = this.props;
        fetchAllTournamentsAction();
    }

    render() {
        const { 
          tournaments: { allTournaments },
          activeTeam,
        } = this.props;
        return (
            <SiteWrapper>
              <Page.Content title="All Tournaments">
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    {activeTeam &&
                      <AllTeamTournamentsList /> 
                    }
                  </Grid.Col>
                  <Grid.Col xs={12} sm={12} lg={6}>
                     { allTournaments && <AllTournamentsList tournaments={allTournaments} /> }
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
    fetchAllTournamentsAction: PropTypes.func
};

const mapStateToProps = state => ({
    tournaments: state.tournaments,
    activeTeam: state.teams.activeTeam
});

const mapDispatchToProps = {
    fetchAllTournamentsAction: () => fetchAllTournamentsAction(),
    fetchTeamTournamentsAction: (teamID) => fetchTeamTournamentsAction(teamID)
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(AllTournamentsPage);
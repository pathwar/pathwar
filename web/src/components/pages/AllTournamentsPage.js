import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {
  Page,
  Grid
} from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import AllTournamentsList from "../tournament/AllTournamentsList";
import { fetchAllTournaments as fetchAllTournamentsAction } from "../../actions/tournaments"

class AllTournamentsPage extends React.PureComponent {

    componentDidMount() {
        const { fetchAllTournamentsAction } = this.props;
        fetchAllTournamentsAction();
    }

    render() {
        const { tournaments: { allTournaments } } = this.props;
        return (
            <SiteWrapper>
              <Page.Content title="All Tournaments">
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    { allTournaments && <AllTournamentsList tournaments={allTournaments} /> }
                  </Grid.Col>
        
                  <Grid.Col xs={12} sm={12} lg={6}>
        
                  </Grid.Col>
        
                </Grid.Row>
              </Page.Content>
            </SiteWrapper>
          );
    }
}

AllTournamentsPage.propTypes = {
    tournaments: PropTypes.object,
    fetchAllTournamentsAction: PropTypes.func
};

const mapStateToProps = state => ({
    tournaments: state.tournaments
});

const mapDispatchToProps = {
    fetchAllTournamentsAction: () => fetchAllTournamentsAction()
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(AllTournamentsPage);
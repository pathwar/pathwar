import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {
  Page,
  Grid,
  Dimmer
} from "tabler-react";

import {
  fetchTeamDetails as fetchTeamDetailsAction
} from "../actions/seasons"

class TeamDetailsPage extends React.PureComponent {

    componentDidMount() {
        const { fetchTeamDetailsAction, uri } = this.props;
        const teamID = uri.split("/")[3];
        fetchTeamDetailsAction(teamID);
    }

    render() {
        const { teamInDetail } = this.props;
        return ( !teamInDetail ? <Dimmer active loader /> :
            <Page.Content title={`Team ${teamInDetail.organization.name}`}>
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>

                  </Grid.Col>
                </Grid.Row>
              </Page.Content>
          );
    }
}

TeamDetailsPage.propTypes = {
  teamInDetail: PropTypes.object,
  fetchTeamDetailsAction: PropTypes.func
};

const mapStateToProps = state => ({
  teamInDetail: state.seasons.teamInDetail
});

const mapDispatchToProps = {
  fetchTeamDetailsAction: (teamID) => fetchTeamDetailsAction(teamID)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(TeamDetailsPage);

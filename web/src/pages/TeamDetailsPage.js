import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {
  Page,
  Grid
} from "tabler-react";

import {
  fetchTeamDetails as fetchTeamDetailsAction
} from "../actions/teams"

class TeamDetailsPage extends React.PureComponent {

    componentDidMount() {
        const { fetchTeamDetailsAction, uri } = this.props;
        const teamID = uri.split("/")[3];
        fetchTeamDetailsAction(teamID);
    }

    render() {
        return (
            <Page.Content title="TEAM XYZ">
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>

                  </Grid.Col>
                </Grid.Row>
              </Page.Content>
          );
    }
}

TeamDetailsPage.propTypes = {
    fetchTeamDetailsAction: PropTypes.func
};

const mapStateToProps = state => ({
});

const mapDispatchToProps = {
  fetchTeamDetailsAction: (teamID) => fetchTeamDetailsAction(teamID)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(TeamDetailsPage);

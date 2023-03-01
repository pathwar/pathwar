import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {Page, Grid, Dimmer, Avatar} from "tabler-react";

import { fetchTeamDetails as fetchTeamDetailsAction } from "../actions/seasons";
import ShadowBox from "../components/ShadowBox";
import {FormattedMessage} from "react-intl";
import moment from "moment";
import TeamsOnOrganizationList from "../components/organization/AllTeamsOnOrganization";

class TeamDetailsPage extends React.PureComponent {
  componentDidMount() {
    const { fetchTeamDetailsAction, uri } = this.props;
    const teamID = uri.split("/")[2];
    fetchTeamDetailsAction(teamID);
  }

  render() {
    const { teamInDetail } = this.props;
    return !teamInDetail ? (
      <Dimmer active loader />
    ) : (
      <Page.Content title={`Team ${teamInDetail.organization.name} - ${teamInDetail.season.name}`}>
        <Grid.Row css={{
          "margin-bottom": "15px",
        }}>
        </Grid.Row>
        <Grid.Row>
          <Grid.Col width={12} lg={5} >
            <ShadowBox>
              <div
                css={{
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                }}
              >
                <h2 className="mb-4 mt-2" style={{fontSize: '1.60rem'}}>{teamInDetail.organization.name}</h2>
                <h2>{teamInDetail.season.name}</h2>
                {teamInDetail.organization.gravatar_url ? (
                  <Avatar size="xxl" imageURL={`${teamInDetail.organization.gravatar_url}?d=identicon`} />
                ) : (
                  <Avatar size="xxl" icon="users" />
                )}
                <h3 className="mb-0 mt-6">
                  <FormattedMessage id="TeamDetails.stats" />
                </h3>
                <p>Score: {teamInDetail.score ? teamInDetail.score : 0} Cash: ${teamInDetail.cash ? teamInDetail.cash : 0}</p>
              </div>
            </ShadowBox>
          </Grid.Col>
          <Grid.Col xs={12} sm={12} md={6}>
            <TeamsOnOrganizationList
            />
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    );
  }
}

TeamDetailsPage.propTypes = {
  teamInDetail: PropTypes.object,
  fetchTeamDetailsAction: PropTypes.func,
};

const mapStateToProps = state => ({
  teamInDetail: state.seasons.teamInDetail,
});

const mapDispatchToProps = {
  fetchTeamDetailsAction: teamID => fetchTeamDetailsAction(teamID),
};

export default connect(mapStateToProps, mapDispatchToProps)(TeamDetailsPage);

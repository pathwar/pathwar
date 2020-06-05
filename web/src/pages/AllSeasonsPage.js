import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import { Page, Grid } from "tabler-react";

import AllSeasonsList from "../components/season/AllSeasonsList";
import { fetchAllSeasons as fetchAllSeasonsAction } from "../actions/seasons";

class AllSeasonsPage extends React.PureComponent {
  componentDidMount() {
    const { fetchAllSeasonsAction } = this.props;
    fetchAllSeasonsAction();
  }

  render() {
    const {
      seasons: { allSeasons },
    } = this.props;
    return (
      <Page.Content title="All Seasons">
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={12} lg={6}></Grid.Col>
          <Grid.Col xs={12} sm={12} lg={6}>
            {allSeasons && <AllSeasonsList seasons={allSeasons} />}
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    );
  }
}

AllSeasonsPage.propTypes = {
  seasons: PropTypes.object,
  fetchAllSeasonsAction: PropTypes.func,
};

const mapStateToProps = state => ({
  seasons: state.seasons,
});

const mapDispatchToProps = {
  fetchAllSeasonsAction: () => fetchAllSeasonsAction(),
};

export default connect(mapStateToProps, mapDispatchToProps)(AllSeasonsPage);

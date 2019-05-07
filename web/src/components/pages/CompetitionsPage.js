import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Page, Grid } from "tabler-react";

import SiteWrapper from "../SiteWrapper";
import LevelsCardPreview from "../levels/LevelCardPreview";

import { fetchLevels as fetchLevelsAction } from "../../actions/levels";

class CompetitionsPage extends React.Component {

    componentDidMount() {
        const { fetchLevelsAction } = this.props;
        fetchLevelsAction();
    }
  
    render() {
        const { competitions } = this.props;

        return (
            <SiteWrapper>
              <Page.Content title="Competitions">
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={8} lg={6}>
                    <h3>Levels</h3>
                    <LevelsCardPreview levels={competitions.levels} />
                  </Grid.Col>
        
                </Grid.Row>
              </Page.Content>
            </SiteWrapper>
          );
    }
}

CompetitionsPage.propTypes = {
    competitions: PropTypes.object,
    fetchLevelsAction: PropTypes.func
};

const mapStateToProps = state => ({
    competitions: state.competitions
});

const mapDispatchToProps = {
    fetchLevelsAction: () => fetchLevelsAction()
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(CompetitionsPage);


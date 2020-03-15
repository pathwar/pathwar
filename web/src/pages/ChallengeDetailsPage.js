import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import { fetchChallengeDetail as fetchChallengeDetailAction } from "../actions/seasons"
import styles from "./styles/ChallengeDetailsPage.module.css"

import {
  Page,
  Grid,
  Dimmer
} from "tabler-react";

class ChallengeDetailsPage extends React.PureComponent {

    componentDidMount() {
      const { fetchChallengeDetailAction, uri } = this.props;
      const challengeID = uri.split("/")[3];
      fetchChallengeDetailAction(challengeID);
    }

    render() {
        const { challenge: { flavor: { challenge } } = { flavor: "no flavor" } } = this.props;

        if(!challenge) {
          return <Dimmer active />
        }

        return (
            <Page.Content title={challenge.name}>
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    <h4>Name</h4>
                    <p className={styles.p}>{challenge.name}</p>

                    <h4>Author</h4>
                    <p className={styles.p}>{challenge.author}</p>

                    <h4>Page</h4>
                    <p className={styles.p}>{challenge.homepage}</p>
                  </Grid.Col>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    <h3>Actions</h3>
                    <p className={styles.p}>{challenge.name}</p>
                  </Grid.Col>
                </Grid.Row>
              </Page.Content>
        );
    }
}

ChallengeDetailsPage.propTypes = {
  fetchChallengeDetailAction: PropTypes.func
};

const mapStateToProps = state => ({
  challenge: state.seasons.challengeInDetail
});

const mapDispatchToProps = {
  fetchChallengeDetailAction: (challengeID) => fetchChallengeDetailAction(challengeID)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ChallengeDetailsPage);

import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {
  fetchChallengeDetail as fetchChallengeDetailAction,
  buyChallenge as buyChallengeAction,
  validateChallenge as validateChallengeAction,
  closeChallenge as closeChallengeAction
} from "../actions/seasons";
import ChallengeBuyStampCard from "../components/challenges/ChallengeBuyStampCard";
import ChallengeValidateStampCard from "../components/challenges/ChallengeValidateStampCard";
import ChallengeCloseStampCard from "../components/challenges/ChallengeCloseStampCard";

import styles from "./styles/ChallengeDetailsPage.module.css";

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
        const {
          challenge,
          activeTeam: { id: teamID } = { id: "no id" },
          buyChallengeAction,
          validateChallengeAction,
          closeChallengeAction
        } = this.props;

        const { flavor: { challenge: flavorChallenge } = { challenge: "no challenge" } } = challenge || {};

        if(!challenge) {
          return <Dimmer active />
        }

        return (
            <Page.Content title={flavorChallenge.name}>
                <Grid.Row cards={true}>
                  <Grid.Col lg={6} md={6} sm={12} xs={12}>
                    <h3>Info</h3>
                    <h4>Name</h4>
                    <p className={styles.p}>{flavorChallenge.name}</p>

                    <h4>Author</h4>
                    <p className={styles.p}>{flavorChallenge.author}</p>

                    <h4>Page</h4>
                    <p className={styles.p}>{flavorChallenge.homepage}</p>
                  </Grid.Col>
                  <Grid.Col lg={6} md={6} sm={12} xs={12}>
                    <h3>Actions</h3>
                    <ChallengeBuyStampCard
                      challenge={challenge}
                      buyChallenge={buyChallengeAction}
                      teamID={teamID}
                    />
                    <ChallengeValidateStampCard
                      challenge={challenge}
                      validateChallenge={validateChallengeAction}
                    />
                    <ChallengeCloseStampCard
                      challenge={challenge}
                      closeChallenge={closeChallengeAction}
                    />
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
  challenge: state.seasons.challengeInDetail,
  activeTeam: state.seasons.activeTeam
});

const mapDispatchToProps = {
  buyChallengeAction: (challengeID, teamID, seasonId) => buyChallengeAction(challengeID, teamID, seasonId),
  validateChallengeAction: (validationData, seasonId) => validateChallengeAction(validationData, seasonId),
  closeChallengeAction: subscriptionID => closeChallengeAction(subscriptionID),
  fetchChallengeDetailAction: (challengeID) => fetchChallengeDetailAction(challengeID)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ChallengeDetailsPage);

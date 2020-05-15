/* eslint-disable react/prop-types */
import React from "react"
import { connect } from "react-redux"
import PropTypes from "prop-types"
import { Link } from "gatsby"
import {
  Button,
  Dimmer,
  ProgressCard,
  Grid,
} from "tabler-react"
import styles from "../../styles/layout/loader.module.css"

import { buyChallenge as buyChallengeAction } from "../../actions/seasons"

const ChallengeCard = (challenge, buyChallenge) => {
  const {
    challenge: { flavor, subscriptions, id: challengeID },
    teamId,
  } = challenge

  return (
    <ProgressCard
      header={flavor.challenge.name}
      content={
        <Button.List>
          <Button
            RootComponent={Link}
            to={`/app/challenge/${challengeID}`}
            target="_blank"
            color="info"
            size="sm"
            icon="eye"
          >
            Open
          </Button>
          <Button
            onClick={() => buyChallenge(challengeID, teamId)}
            size="sm"
            color="success"
            disabled={subscriptions}
            icon={subscriptions ? "check" : "dollar-sign"}
          >
            Buy
          </Button>
        </Button.List>
      }
      progressColor="green"
      progressWidth={84}
    />
  )
}

const ChallengeList = props => {
  const { challenges, activeUserSession, buyChallengeAction } = props

  return !challenges || !activeUserSession ? (
    <Dimmer className={styles.dimmer} active loader />
  ) : (
    <>
      <Grid.Row>
        {challenges.map(challenge => (
          <Grid.Col lg={4} sm={4} md={4} xs={4}>
            <ChallengeCard
              challenge={challenge}
              buyChallenge={buyChallengeAction}
            />
          </Grid.Col>
        ))}
      </Grid.Row>
    </>
  )
}

ChallengeList.propTypes = {
  buyChallengeAction: PropTypes.func,
  activeTeamId: PropTypes.string,
}

const mapStateToProps = state => ({
  activeUserSession: state.userSession.activeUserSession,
})

const mapDispatchToProps = {
  buyChallengeAction: (challengeID, teamID, seasonId) =>
    buyChallengeAction(challengeID, teamID, seasonId),
}

export default connect(mapStateToProps, mapDispatchToProps)(ChallengeList)

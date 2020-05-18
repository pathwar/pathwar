/* eslint-disable react/prop-types */
import React from "react"
import { useSelector } from "react-redux"
import { Link } from "gatsby"
import { Button, Dimmer, ProgressCard, Grid } from "tabler-react"
import styles from "../../styles/layout/loader.module.css"

const ChallengeCard = ({ challenge, buyChallenge, teamID }) => {
  const { flavor, subscriptions, id: challengeID } = challenge

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
            onClick={() => buyChallenge(challengeID, teamID, false)}
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
  const activeUserSession = useSelector(state => state.userSession.activeUserSession)
  const activeTeam = useSelector(state => state.seasons.activeTeam)

  const { challenges, buyChallenge } = props

  return !challenges || !activeUserSession ? (
    <Dimmer className={styles.dimmer} active loader />
  ) : (
    <>
      <Grid.Row>
        {challenges.map(challenge => (
          <Grid.Col lg={4} sm={4} md={4} xs={4}>
            <ChallengeCard
              challenge={challenge}
              buyChallenge={buyChallenge}
              teamID={activeTeam.id}
            />
          </Grid.Col>
        ))}
      </Grid.Row>
    </>
  )
}

export default ChallengeList;

/* eslint-disable react/prop-types */
import React, { useState } from "react"
import { connect } from "react-redux"
import PropTypes from "prop-types"
import { Link } from "gatsby"
import { Card, Button, Dimmer, Table, Progress, Form } from "tabler-react"
import styles from "../../styles/layout/loader.module.css"

import { buyChallenge as buyChallengeAction } from "../../actions/seasons"

const ChallengeRow = ({ challenge, teamId, seasonId, buyChallenge }) => {
  const [isValidateOpen, setValidateOpen] = useState(false)
  const { flavor } = challenge
  const hasSubscriptions = challenge.subscriptions

  return (
    <>
      <Table.Row>
        <Table.Col>
          <strong>{flavor.challenge.name}</strong>
        </Table.Col>
        <Table.Col className="text-nowrap">{flavor.challenge.author}</Table.Col>
        <Table.Col>
          <div className="clearfix">
            <div className="float-left">
              <strong>42%</strong>
            </div>
          </div>
          <Progress size="sm">
            <Progress.Bar color="yellow" width={42} />
          </Progress>
        </Table.Col>
        <Table.Col className="w-1">
          <Button
            RootComponent={Link}
            to={`/app/challenge/${challenge.id}`}
            color="info"
            size="sm"
            icon="eye"
          />
        </Table.Col>
        <Table.Col className="w-1">
          <Button
            onClick={() => buyChallenge(challenge.id, teamId, seasonId)}
            value="Buy"
            size="sm"
            color="success"
            icon={hasSubscriptions ? "check" : "dollar-sign"}
          />
        </Table.Col>
        <Table.Col className="w-1">
          <Button
            value="Validate"
            size="sm"
            color="warning"
            icon="check"
            onClick={() => setValidateOpen(!isValidateOpen)}
          >
            Validate
          </Button>
        </Table.Col>
        <Table.Col className="w-1">
          <Button value="Github page" social="github" size="sm" />
        </Table.Col>
        <Table.Col className="w-1">
          <Button value="Close" size="sm" color="danger" icon="x-circle" />
        </Table.Col>
      </Table.Row>
      {isValidateOpen && (
        <Table.Row>
          <Table.Col colSpan="8">
          <Form.FieldSet>
            <Form.Group isRequired label="Passphrase">
              <Form.Input name="passphrase" />
            </Form.Group>
            <Form.Group isRequired label="Comment">
              <Form.Textarea
                defaultValue=" Oh! Come and see the violence inherent in the system! Help, help, I'm being repressed! We shall say 'Ni' again to you, if you do not appease us. I'm not a witch. I'm not a witch. Camelot!"
                name="example-textarea"
                placeholder="Content.."
                rows={6}
              />
            </Form.Group>
          </Form.FieldSet>
          </Table.Col>
        </Table.Row>
      )}
    </>
  )
}

const ChallengeTable = ({ challenges, teamId, seasonId, buyChallenge }) => {
  return (
    <Table
      cards={true}
      striped={true}
      responsive={true}
      className="table-vcenter"
    >
      <Table.Header>
        <Table.Row>
          <Table.ColHeader>Flavor</Table.ColHeader>
          <Table.ColHeader>Author</Table.ColHeader>
          <Table.ColHeader>Progress</Table.ColHeader>
          <Table.ColHeader>View</Table.ColHeader>
          <Table.ColHeader>Buy</Table.ColHeader>
          <Table.ColHeader />
          <Table.ColHeader>Page</Table.ColHeader>
          <Table.ColHeader>Close</Table.ColHeader>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {challenges.map(challenge => {
          return (
            <ChallengeRow
              challenge={challenge}
              buyChallenge={buyChallenge}
              teamId={teamId}
              seasonId={seasonId}
            />
          )
        })}
      </Table.Body>
    </Table>
  )
}

const ChallengeCardPreview = props => {
  const { challenges, activeUserSession, buyChallengeAction } = props
  const { active_team_member, active_season_id } =
    (activeUserSession && activeUserSession.user) || {}

  debugger

  return !challenges || !activeUserSession ? (
    <Dimmer className={styles.dimmer} active loader />
  ) : (
    <Card>
      <ChallengeTable
        challenges={challenges}
        buyChallenge={buyChallengeAction}
        teamId={active_team_member.team_id}
        seasonId={active_season_id}
      />
    </Card>
  )
}

ChallengeCardPreview.propTypes = {
  fetchChallengesAction: PropTypes.func,
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

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ChallengeCardPreview)

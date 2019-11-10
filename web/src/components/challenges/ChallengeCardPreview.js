/* eslint-disable react/prop-types */
import * as React from "react"
import { Link } from "gatsby"
import { Card, Button, Dimmer, Table, Progress } from "tabler-react"
import styles from "../../styles/layout/loader.module.css"

const ChallengeTable = ({ challenges }) => {
  return (
    <Table
      cards={true}
      striped={true}
      responsive={true}
      className="table-vcenter"
    >
      <Table.Header>
        <Table.Row>
          <Table.ColHeader>Name</Table.ColHeader>
          <Table.ColHeader>Author</Table.ColHeader>
          <Table.ColHeader>Progress</Table.ColHeader>
          <Table.ColHeader />
          <Table.ColHeader />
          <Table.ColHeader />
          <Table.ColHeader />
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {challenges.map(challenge => {
          const { flavor } = challenge
          return (
            <Table.Row>
              <Table.Col>{flavor.challenge.name}</Table.Col>
              <Table.Col className="text-nowrap">
                {flavor.challenge.author}
              </Table.Col>
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
                >
                  View
                </Button>
              </Table.Col>
              <Table.Col className="w-1">
                <Button RootComponent={Link} to="/" color="success" size="sm">
                  Buy
                </Button>
              </Table.Col>
              <Table.Col className="w-1">
                <Button RootComponent={Link} to="/" color="warning" size="sm">
                  Validate
                </Button>
              </Table.Col>
              <Table.Col className="w-1">
                <Button social="github" />
              </Table.Col>
            </Table.Row>
          )
        })}
      </Table.Body>
    </Table>
  )
}

const ChallengeCardPreview = props => {
  const { challenges } = props

  return !challenges ? (
    <Dimmer className={styles.dimmer} active loader />
  ) : (
    <Card>
      <ChallengeTable challenges={challenges} />
    </Card>
  )
}

export default ChallengeCardPreview

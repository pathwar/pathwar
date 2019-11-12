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
          <Table.ColHeader>View</Table.ColHeader>
          <Table.ColHeader>Buy</Table.ColHeader>
          <Table.ColHeader />
          <Table.ColHeader>Close</Table.ColHeader>
          <Table.ColHeader>Page</Table.ColHeader>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {challenges.map(challenge => {
          const { flavor } = challenge
          return (
            <Table.Row>
              <Table.Col><strong>{flavor.challenge.name}</strong></Table.Col>
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
                  icon="eye"
                />
              </Table.Col>
              <Table.Col className="w-1">
                <Button value="Buy" size="sm" color="success" icon="dollar-sign" />
              </Table.Col>
              <Table.Col className="w-1">
                <Button value="Validate" size="sm" color="warning" icon="check">
                  Validate
                </Button>
              </Table.Col>
              <Table.Col className="w-1">
                <Button value="Close" size="sm" color="danger" icon="x-circle" />
              </Table.Col>
              <Table.Col className="w-1">
                <Button value="Github page" social="github" size="sm"/>
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

import * as React from "react";

import {
  Page,
  Grid,
  Card,
  Table,
  Avatar,
} from "tabler-react";


import SiteWrapper from "../SiteWrapper";

function Home() {
  return (
    <SiteWrapper>
      <Page.Content title="Dashboard">
        <Grid.Row cards={true}>
          <Grid.Col lg={6}>
            <Card>
              <Card.Header>
                <Card.Title>Teams</Card.Title>
              </Card.Header>
              <Table
                cards={true}
                striped={true}
                responsive={true}
                className="table-vcenter"
              >
                <Table.Header>
                  <Table.Row>
                    <Table.ColHeader colSpan={2}>Name</Table.ColHeader>
                    <Table.ColHeader>Locale</Table.ColHeader>
                    <Table.ColHeader />
                  </Table.Row>
                </Table.Header>
                <Table.Body>
                  <Table.Row>
                    <Table.Col className="w-1">
                      <Avatar imageURL="./demo/faces/male/9.jpg" />
                    </Table.Col>
                    <Table.Col>Team 1</Table.Col>
                    <Table.Col>fr_FR</Table.Col>
                  </Table.Row>
                </Table.Body>
              </Table>
            </Card>
          </Grid.Col>

        </Grid.Row>
      </Page.Content></SiteWrapper>
  );
}

export default Home;
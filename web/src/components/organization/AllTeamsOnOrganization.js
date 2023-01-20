import React from "react";
import { Card, Table, Dimmer } from "tabler-react";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";

const TeamsOnOrganizationRow = ({ teams }) => {
  return teams.map((item, idx) => {
    return (
      <Table.Row key={item.id}>
        <Table.Col alignContent="center">
          {item.season.name}
        </Table.Col>
        <Table.Col alignContent="center">
          {item.score}
        </Table.Col>
        <Table.Col alignContent="center">
          {item.cash}
        </Table.Col>
      </Table.Row>
    );
  });
}

const TeamsOnOrganizationList = ({ teams }) => {
  return !teams ? (
    <Dimmer active loader />
  ) : (
    <Card>
      <Table
        striped={true}
        responsive={true}
        verticalAlign="center"
        className="mb-0"
      >
        <Table.Header>
          <Table.Row>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnOrganization.season" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnOrganization.score" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnOrganization.cash" />
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {teams && (
            <TeamsOnOrganizationRow teams={teams} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

export default TeamsOnOrganizationList;

import React from "react";
import { Card, Table, Dimmer } from "tabler-react";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";

const MembersOnOrganizationRow = ({ members }) => {
  return members.map((item, idx) => {
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

const MembersOnOrganizationList = ({ members }) => {
  return (
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
              Index
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              Name
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              Role
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              Joined at
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {members && (
            <MembersOnOrganizationRow members={members} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

export default MembersOnOrganizationList;

import React from "react";
import {Card, Table, Avatar} from "tabler-react";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";

const MembersOnTeamsRow = ({ members }) => {
  return members.map((item, idx) => {
    return (
      <Table.Row key={item.id}>
        <Table.Col alignContent="center">
          {idx+1}
        </Table.Col>
        <Table.Col alignContent="center">
          <span>{item.user.username}</span>
        </Table.Col>
        <Table.Col alignContent="center">
          {item.role ? item.role : "Member"}
        </Table.Col>
      </Table.Row>
    );
  });
}

const MembersOnTeamsList = ({ members }) => {
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
              <FormattedMessage id="AllMembersOnOrganization.index" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllMembersOnOrganization.name" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllMembersOnOrganization.role" />
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {members && (
            <MembersOnTeamsRow members={members} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

export default MembersOnTeamsList;

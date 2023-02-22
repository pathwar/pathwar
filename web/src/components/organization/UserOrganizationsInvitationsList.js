import React from "react";
import { Card, Table, Dimmer, Avatar } from "tabler-react";
import PropTypes from "prop-types";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";

const OrganizationsInvitationsRows = ({ organizationsInvitations }) => {
  return organizationsInvitations.map((item, idx) => {
    return (
      <Table.Row key={item.id}>
        <Table.Col alignContent="center">{idx + 1}</Table.Col>
      </Table.Row>
    );
  });
}

const UserOrganizationsInvitationsList = ({ userOrganizationsInvitationsList }) => {
  return !UserOrganizationsInvitationsList ? (
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
              <FormattedMessage id="UserOrganizationsList.rank" />
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {userOrganizationsInvitationsList && (
            <OrganizationsInvitationsRows organizationsInvitations={userOrganizationsInvitationsList} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

UserOrganizationsInvitationsList.propTypes = {
  seasons: PropTypes.object,
};

export default UserOrganizationsInvitationsList;

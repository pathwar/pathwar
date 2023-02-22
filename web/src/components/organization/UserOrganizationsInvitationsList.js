import React from "react";
import { Card, Table, Dimmer, Avatar } from "tabler-react";
import PropTypes from "prop-types";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";
import {Link} from "gatsby";

const OrganizationsInvitationsRows = ({ organizationsInvitations }) => {
  return organizationsInvitations.map((item, idx) => {
    return (
      <Table.Row key={item.id}>
        <Table.Col alignContent="center">{idx + 1}</Table.Col>

        <Table.Col alignContent="center"
        >
          <Avatar
            className="mr-2"
            imageURL={`${item.gravatar_url}?d=identicon`}
          />
          <Link
            className="link"
            to={"/organization/" + item.id}
            activeStyle={{
              fontWeight: "bold",
              color: currentTheme.colors.primary,
            }}
          >
            {item.name}
          </Link>
        </Table.Col>
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
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserOrganizationsList.organization" />
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

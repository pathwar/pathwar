import React from "react";
import {Card, Table, Dimmer, Avatar, Button} from "tabler-react";
import PropTypes from "prop-types";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";
import {Link} from "gatsby";
import {useTheme} from "emotion-theming";
import AcceptOrganizationInviteButton from "../organization/AcceptOrganizationInviteButton";

const TeamsInvitationsRows = ({ teamsInvitations }) => {
  const currentTheme = useTheme();

  return teamsInvitations.map((item, idx) => {
    return (
      <Table.Row key={item.id}>
        <Table.Col alignContent="center">{idx + 1}</Table.Col>

        <Table.Col alignContent="center"
        >
          <Avatar
            className="mr-2"
            imageURL={`${item.organization.gravatar_url}?d=identicon`}
          />
          <Link
            className="link"
            to={"/organization/" + item.organization.id}
            activeStyle={{
              fontWeight: "bold",
              color: currentTheme.colors.primary,
            }}
          >
            {item.organization.name}
          </Link>
        </Table.Col>

        <Table.Col alignContent="center">{item.user.slug}</Table.Col>
        <Table.Col alignContent="center">
          <AcceptOrganizationInviteButton organizationInvite={item} />
        </Table.Col>
      </Table.Row>
    );
  });
}

const UserTeamsInvitationsList = ({ userTeamsInvitationsList }) => {
  return !userTeamsInvitationsList ? (
    <h3><FormattedMessage id="UserTeamsInvitationsList.NoInvitations" /></h3>
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
              <FormattedMessage id="UserTeamsInvitationsList.season" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserTeamsInvitationsList.organization" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserTeamsInvitationsList.invitedBy" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserTeamsInvitationsList.accept" />
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {userTeamsInvitationsList && (
            <TeamsInvitationsRows teamsInvitations={userTeamsInvitationsList} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

UserTeamsInvitationsList.propTypes = {
  seasons: PropTypes.object,
};

export default UserTeamsInvitationsList;

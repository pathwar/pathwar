import React from "react";
import {Card, Table, Dimmer, Avatar, Button} from "tabler-react";
import PropTypes from "prop-types";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";
import {Link} from "gatsby";
import {useTheme} from "emotion-theming";
import {
  acceptOrganizationInvite as acceptOrganizationInviteAction,
  rejectOrganizationInvite as rejectOrganizationInviteAction,
} from "../../actions/organizations";
import {useDispatch} from "react-redux";

const OrganizationsInvitationsRows = ({ organizationsInvitations }) => {
  const currentTheme = useTheme();
  const dispatch = useDispatch();
  const acceptOrganizationInvite = (organizationInviteID) =>
    dispatch(acceptOrganizationInviteAction(organizationInviteID));
  const rejectOrganizationInvite = (organizationInviteID) =>
    dispatch(rejectOrganizationInviteAction(organizationInviteID));

  return organizationsInvitations.map((item, idx) => {
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
          <Button.List>
            <Button color="success" className="mx-lg-auto" onClick={() => acceptOrganizationInvite(item.id)}>✔</Button>
            <Button color="red" className="ml-auto" onClick={() => rejectOrganizationInvite(item.id)}>✖</Button>
          </Button.List>
        </Table.Col>
      </Table.Row>
    );
  });
}

const UserOrganizationsInvitationsList = ({ userOrganizationsInvitationsList }) => {
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
              <FormattedMessage id="UserOrganizationsInvitationsList.rank" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserOrganizationsInvitationsList.organization" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserOrganizationsInvitationsList.invitedBy" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="UserOrganizationsInvitationsList.accept" />
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

import React, { useMemo } from "react";
import {Card, Table, Dimmer, Avatar, Button} from "tabler-react";
import PropTypes from "prop-types";
import { css } from "@emotion/core";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";
import {Link} from "gatsby";
import {useTheme} from "emotion-theming";

const OrganizationsRows = ({ organizations }) => {
  const currentTheme = useTheme();

  return organizations.map((item, idx) => {
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
          <Table.Col alignContent="center">
            <Button.List>
              <Button color="primary" className="mx-lg-auto">See Details</Button>
            </Button.List>
          </Table.Col>
        </Table.Row>
      );
  });
}

const UserOrganizationsList = ({ userOrganizationsList }) => {
  return !userOrganizationsList ? (
    <h3>User don't have any organization for now</h3>
  ) : (
    <Card>
      <div css={{maxHeight: "435px", overflow: "auto"}}>
        <Table
          striped={true}
          responsive={true}
          verticalAlign="center"
          className="mb-0"
        >
          <Table.Header>
            <Table.Row>
              <Table.ColHeader alignContent="center">
                <FormattedMessage id="UserOrganizationsList.index" />
              </Table.ColHeader>
              <Table.ColHeader alignContent="center">
                <FormattedMessage id="UserOrganizationsList.organization" />
              </Table.ColHeader>
              <Table.ColHeader alignContent="center">
                <FormattedMessage id="UserOrganizationsList.details" />
              </Table.ColHeader>
            </Table.Row>
          </Table.Header>
            <Table.Body>
                {userOrganizationsList && (
                  <OrganizationsRows organizations={userOrganizationsList} />
                )}
            </Table.Body>
        </Table>
      </div>
    </Card>
  );
};

UserOrganizationsList.propTypes = {
  seasons: PropTypes.object,
};

export default UserOrganizationsList;

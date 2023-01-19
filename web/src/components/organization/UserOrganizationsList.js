import React, { useMemo } from "react";
import { Card, Table, Dimmer, Avatar } from "tabler-react";
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

          <Table.Col
            css={css`
              display: flex;
              align-items: center;
            `}
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
        </Table.Row>
      );
  });
}

const UserOrganizationsList = ({ userOrganizationsList }) => {
  return !userOrganizationsList ? (
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
            <Table.ColHeader>
              <FormattedMessage id="UserOrganizationsList.organization" />
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {userOrganizationsList && (
            <OrganizationsRows organizations={userOrganizationsList} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

UserOrganizationsList.propTypes = {
  seasons: PropTypes.object,
};

export default UserOrganizationsList;

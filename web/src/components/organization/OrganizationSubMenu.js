import React from "react";
import {Grid} from "tabler-react";
import {Link} from "gatsby";
import {FormattedMessage} from "react-intl";
import {useTheme} from "emotion-theming";

export const OrganizationSubMenu = ({organization}) => {
  const currentTheme = useTheme();

  return (
    <>
      <Grid.Col xs={12} sm={12} md={3} offsetMd={1}>
        <Link
          className="link"
          to={"/organization/" + organization.id}
          activeStyle={{
            fontWeight: "bold",
            color: currentTheme.colors.primary,
          }}
        >
          {organization.name}
        </Link>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={3} offsetMd={1}>
        <Link
          className="link"
          to={"/organization/" + organization.id + "/members"}
          activeStyle={{
            fontWeight: "bold",
            color: currentTheme.colors.primary,
          }}
        >
          <FormattedMessage id="OrganizationDetailsSubmenu.members" />
        </Link>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={3}>
        <Link
          className="link"
          to={"/organization/" + organization.id + "/teams"}
          activeStyle={{
            fontWeight: "bold",
            color: currentTheme.colors.primary,
          }}
        >
          <FormattedMessage id="OrganizationDetailsSubmenu.teams" />
        </Link>
      </Grid.Col>
    </>
);
}

export default OrganizationSubMenu;

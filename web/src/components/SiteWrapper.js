import * as React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import {
  Site,
  Nav,
  Dropdown,
  Card,
  Grid,
  Avatar,
  Text,
  Tag,
} from "tabler-react";
import { Link } from "gatsby";

import logo from "../images/pathwar-favicon.png";

const navBarItems = [
  {
    value: "Season",
    to: "/app/season",
    icon: "flag",
    LinkComponent: Link,
    useExact: "false",
  },
  {
    value: "Dashboard",
    to: "/app/dashboard",
    icon: "home",
    LinkComponent: Link,
    useExact: "false",
  },
];

const accountDropdownProps = (
  { activeUserSession, activeKeycloakSession },
  activeSeason
) => {
  const { user, claims } = activeUserSession || {};

  const username =
    claims && claims.preferred_username ? claims.preferred_username : "Account";
  const avatar = user && user.gravatar_url ? user.gravatar_url : logo;
  const description = activeSeason && activeSeason.name;
  const options = [];

  if (!activeUserSession && !activeKeycloakSession) {
    options.push({ icon: "log-in", value: "Log in", to: "/app/login" });
  }

  if (activeUserSession && activeKeycloakSession) {
    options.push("profile");
    options.push({
      icon: "edit",
      value: "Edit account",
      to: activeKeycloakSession.tokenParsed.iss + "/account",
    });
    options.push("divider");
    options.push({
      value: "App settings",
      to: "/app/settings",
      icon: "settings",
    });
  }

  options.push({
    icon: "help-circle",
    value: "FAQ",
    to: "https://github.com/pathwar/pathwar/wiki/FAQ",
    target: "_blank",
  });

  if (activeUserSession && activeKeycloakSession) {
    options.push({
      icon: "log-out",
      value: "Log out",
      to: "/app/logout",
    });
  }

  return {
    avatarURL: avatar,
    name: `${username}`,
    description: description,
    options: options,
    optionsRootComponent: Link,
  };
};

const navItemsProps = ({ activeUserSession }, activeSeason) => {
  const clicked = e => {
    e.preventDefault();
    alert(activeSeason.name);
  };

  const items =
    activeUserSession &&
    activeUserSession.seasons.map(dataSet => {
      const { season } = dataSet;
      const isActive = activeSeason && season.id === activeSeason.id;

      return (
        <Dropdown.Item
          className={isActive && "active bold"}
          key={season.id}
          to="#"
          onClick={e => clicked(e)}
        >
          <div style={{ fontWeight: isActive ? "bold" : "initial" }}>
            {season.name}
          </div>
          <div>
            <Tag.List>
              <Tag addOn={season.status} addOnColor="indigo">
                Status
              </Tag>
              <Tag addOn={season.visibility} addOnColor="indigo">
                Visibility
              </Tag>
            </Tag.List>
          </div>
        </Dropdown.Item>
      );
    });

  return (
    <Nav.Item type="div" className="d-none d-md-flex">
      <Dropdown
        triggerContent={activeSeason && activeSeason.name}
        type="button"
        color="primary"
        icon="flag"
        items={items}
      />
    </Nav.Item>
  );
};

class SiteWrapper extends React.Component {
  render() {
    const { userSession, activeSeason } = this.props;
    return (
      <Site.Wrapper
        headerProps={{
          href: "/",
          alt: "Pathwar Project",
          imageURL: logo,
          accountDropdown: accountDropdownProps(userSession, activeSeason),
          navItems: navItemsProps(userSession, activeSeason),
        }}
        navProps={{ itemsObjects: navBarItems }}
      >
        {this.props.children}
      </Site.Wrapper>
    );
  }
}

SiteWrapper.propTypes = {
  children: PropTypes.node,
  userSession: PropTypes.object,
  lastActiveTeam: PropTypes.object,
};

const mapStateToProps = state => ({
  userSession: state.userSession,
  activeSeason: state.seasons.activeSeason,
});

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(SiteWrapper);

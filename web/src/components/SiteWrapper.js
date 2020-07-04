import * as React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Site, Nav, Button } from "tabler-react";
import { navigate, Link } from "gatsby";

import logo from "../images/new-pathwar-logo-dark-blue.png";

const navBarItems = [
  {
    value: "Home",
    to: "/app/home",
    icon: "home",
    LinkComponent: Link,
    useExact: "false",
  },
  {
    value: "Challenges",
    to: "/app/challenges",
    icon: "anchor",
    LinkComponent: Link,
    useExact: "false",
  },
  {
    value: "Statistics",
    to: "/app/statistics",
    icon: "activity",
    LinkComponent: Link,
    useExact: "false",
  },
];

const accountDropdownProps = ({ activeUserSession, activeKeycloakSession }) => {
  const { user, claims } = activeUserSession || {};

  const username =
    claims && claims.preferred_username ? claims.preferred_username : "Account";
  const avatar = user && user.gravatar_url ? user.gravatar_url : logo;
  const description = claims && claims.email ? claims.email : "Log in?";
  const options = [];
  if (activeUserSession) {
    options.push("profile");
  }
  if (activeUserSession) {
    options.push("divider");
  }
  options.push("help");
  if (!activeUserSession && !activeKeycloakSession) {
    options.push({ icon: "log-in", value: "Log in", to: "/app/login" });
  }
  if (activeUserSession && activeKeycloakSession) {
    options.push({
      icon: "edit",
      value: "Edit account",
      to: activeKeycloakSession.tokenParsed.iss + "/account",
    });
  }
  return {
    avatarURL: avatar,
    name: `${username}`,
    description: description,
    options: options,
  };
};

class SiteWrapper extends React.Component {
  render() {
    const { userSession } = this.props;
    return (
      <Site.Wrapper
        headerProps={{
          href: "/",
          alt: "Pathwar Project",
          imageURL: logo,
          accountDropdown: accountDropdownProps(userSession),
          navItems: (
            <Nav.Item type="div" className="d-none d-md-flex">
              {userSession.activeKeycloakSession && (
                <Button link onClick={() => navigate("/app/logout")}>
                  Log out
                </Button>
              )}
            </Nav.Item>
          ),
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
});

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(SiteWrapper);

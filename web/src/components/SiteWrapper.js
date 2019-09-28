import * as React from "react";
import PropTypes from "prop-types";
import { connect } from 'react-redux';
import {
  Site,
  Nav,
  Button
} from "tabler-react";
import { Link } from "@reach/router";
import { navigate } from "gatsby";

import logo from "../images/pathwar-logo.png";

const navBarItems = [
  {
    value: "Dashboard",
    to: "/app/dashboard",
    icon: "clipboard",
    LinkComponent: Link
  },
  {
    value: "Statistics",
    to: "/app/statistics",
    icon: "bar-chart-2",
    LinkComponent: Link
  },
  {
    value: "Tournament",
    to: "/app/tournament",
    icon: "flag",
    LinkComponent: Link
  }
];

const accountDropdownProps = ({activeUserSession, activeKeycloakSession}) => {
    const { user, claims } =  activeUserSession || {};

    const username = claims && claims.preferred_username ? claims.preferred_username : "Account"
    const avatar = user && user.gravatar_url ? user.gravatar_url : logo;
    const description = claims && claims.email ? claims.email : "Log in?";
    const options = [];
    if (activeUserSession) { options.push("profile"); }
    if (activeUserSession) { options.push("divider"); }
    options.push("help");
    if (!activeUserSession && !activeKeycloakSession) { options.push({icon: "settings", value: "Log in", to: "/app/login"}); }
    if (activeUserSession && activeKeycloakSession) { options.push({icon: "settings", value: "Settings", to: activeKeycloakSession.tokenParsed.iss+"/account"}); }
    if (activeUserSession && activeKeycloakSession) { options.push({icon: "settings", value: "Log out", to: "/app/logout"}); }
    return {
        avatarURL: avatar,
        name: `${username}`,
        description: description,
        options: options,
    }
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
              {userSession.activeSession && (
                <Button
                  link
                  onClick={() => navigate("/app/logout")}
                >
                Log out
              </Button>
              )}
            </Nav.Item>
          )
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
    userSession: state.userSession
});

const mapDispatchToProps = {};

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SiteWrapper);

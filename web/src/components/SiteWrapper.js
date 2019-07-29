import * as React from "react";
import PropTypes from "prop-types";
import { connect } from 'react-redux';
import {
  Site
} from "tabler-react";
import { Link } from "@reach/router";

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

const accountDropdownProps = ({activeUser}) => {
    const username = activeUser ? activeUser.username : "Log In?";

    return {
        avatarURL: logo,
        name: `${username}`,
        options: [
            { icon: "user", value: "Profile" },
            { icon: "settings", value: "Settings" },
            { isDivider: true },
            { icon: "help-circle", value: "Need help?" },
            { icon: "log-out", value: "Sign out", to: "/logout" },
        ],
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
          accountDropdown: accountDropdownProps(userSession)
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
    performLoginAction: PropTypes.func
};

const mapStateToProps = state => ({
  userSession: state.userSession
});

const mapDispatchToProps = {};

export default connect(
mapStateToProps,
mapDispatchToProps
)(SiteWrapper);
import * as React from "react";
import { withRouter, NavLink } from "react-router-dom";
import { connect } from 'react-redux';
import PropTypes from "prop-types";
import {
  Site,
  RouterContextProvider,
} from "tabler-react";

const navBarItems = [
  {
    value: "Dashboard",
    to: "/dashboard",
    icon: "clipboard",
    useExact: true,
    LinkComponent: withRouter(NavLink)
  },
  {
    value: "Statistics",
    to: "/statistics",
    icon: "bar-chart-2",
    LinkComponent: withRouter(NavLink)
  },
  {
    value: "Competition",
    to: "/competition",
    icon: "flag",
    LinkComponent: withRouter(NavLink)
  }
];

const accountDropdownProps = ({activeUser}) => {
    const username = activeUser ? activeUser.username : "Log In?";
    return {
        avatarURL: `"${require('../images/pathwar-logo.png')}"`,
        name: `${username}`,
        description: "Description",
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
          imageURL: "/pathwar-logo.png",
          accountDropdown: accountDropdownProps(userSession)
        }}
        navProps={{ itemsObjects: navBarItems }}
        routerContextComponentType={withRouter(RouterContextProvider)}
      >
        {this.props.children}
      </Site.Wrapper>
    );
  }
}

SiteWrapper.propTypes = {
    children: PropTypes.node,
    userSession: PropTypes.object,
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
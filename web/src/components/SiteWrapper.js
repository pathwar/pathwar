import * as React from "react";
import { withRouter, NavLink } from "react-router-dom";
import { connect } from 'react-redux';
import PropTypes from "prop-types";
import {
  Site,
  RouterContextProvider,
} from "tabler-react";

import { fetchUserSession as fetchUserSessionAction } from "../actions/userSession";

const navBarItems = [
  {
    value: "Dashboard",
    to: "/",
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
            { icon: "log-out", value: "Sign out" },
        ],
    }
};

class SiteWrapper extends React.Component {

  componentDidMount() {
      const { fetchUserSessionAction } = this.props;
      fetchUserSessionAction();
  }

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
    fetchUserSessionAction: PropTypes.func
};

const mapStateToProps = state => ({
    userSession: state.userSession
});

const mapDispatchToProps = {
    fetchUserSessionAction: () => fetchUserSessionAction()
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(SiteWrapper);
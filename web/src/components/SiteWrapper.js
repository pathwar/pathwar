import * as React from "react";
import { withRouter, NavLink } from "react-router-dom";
import { connect } from 'react-redux';
import PropTypes from "prop-types";
import {
  Site,
  RouterContextProvider,
} from "tabler-react";

import { fetchUserSession as fetchUserSessionAction } from "../actions/session";

const navBarItems = [
  {
    value: "Dashboard",
    to: "/",
    icon: "clipboard",
    useExact: true,
  },
  {
    value: "Statistics",
    to: "/statistics",
    icon: "bar-chart-2",
    LinkComponent: withRouter(NavLink)
  },
  {
    value: "Competitions",
    icon: "flag",
  }
];

const accountDropdownProps = ({activeSession}) => {
    const username = activeSession ? activeSession.username : "Log In?";
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
    const { session } = this.props;

    return (
      <Site.Wrapper
        headerProps={{
          href: "/",
          alt: "Pathwar Project",
          imageURL: "/pathwar-logo.png",
          accountDropdown: accountDropdownProps(session)
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
    session: PropTypes.object,
    fetchUserSessionAction: PropTypes.func
};

const mapStateToProps = state => ({
    session: state.session
});

const mapDispatchToProps = {
    fetchUserSessionAction: () => fetchUserSessionAction()
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(SiteWrapper);
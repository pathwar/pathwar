import * as React from "react";
import { NavLink, withRouter } from "react-router-dom";
import {
  Site,
  RouterContextProvider,
} from "tabler-react";

const navBarItems = [
  {
    value: "Dashboard",
    to: "/",
    icon: "clipboard",
    LinkComponent: withRouter(NavLink),
    useExact: true,
  },
  {
    value: "Statistics",
    icon: "bar-chart-2",
  },
  {
    value: "Competitions",
    icon: "flag",
  }
];

const accountDropdownProps = {
  avatarURL: `"${require('../images/pathwar-logo.png')}"`,
  name: "Logged User",
  description: "Active session: session",
  options: [
    { icon: "user", value: "Profile" },
    { icon: "settings", value: "Settings" },
    { isDivider: true },
    { icon: "help-circle", value: "Need help?" },
    { icon: "log-out", value: "Sign out" },
  ],
};

class SiteWrapper extends React.Component {

  render() {
    return (
      <Site.Wrapper
        headerProps={{
          href: "/",
          alt: "Pathwar Project",
          imageURL: "/pathwar-logo.png",
          accountDropdown: accountDropdownProps
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
    children: null,
};
  

export default SiteWrapper;
import * as React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Site, Nav, Dropdown, Tag, Grid } from "tabler-react";
import { Link } from "gatsby";
import ValidateCouponForm from "../components/coupon/ValidateCouponForm";

import logo from "../images/new-pathwar-logo-dark-blue.png";

const navBarItems = [
  // {
  //   value: "Home",
  //   to: "/app/home",
  //   icon: "home",
  //   LinkComponent: Link,
  //   useExact: "false",
  // },
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
    claims && claims.preferred_username ? claims.preferred_username : "Log in";
  const avatar = user && user.gravatar_url ? user.gravatar_url : logo;
  const options = [];

  if (!activeUserSession && !activeKeycloakSession) {
    options.push({ icon: "log-in", value: "Log in", to: "/app/challenges" });
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
    description: undefined,
    options: options,
    optionsRootComponent: Link,
  };
};

const SeasonDropdownSelector = ({ userSession, activeSeason }) => {
  const { activeUserSession } = userSession || {};
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
    <Dropdown
      triggerContent={(activeSeason && activeSeason.name) || "Loading.."}
      type="button"
      color="primary"
      icon="flag"
      items={items}
    />
  );
};

const NavBar = ({ userSession }) => {
  const {
    activeUserSession: {
      user: {
        active_team_member: { team },
      },
    },
  } = userSession;

  return (
    <Grid.Row className="align-items-center">
      <Grid.Col width={6} className="ml-auto text-right" ignoreCol={true}>
        <ValidateCouponForm />
        <Tag color="lime" addOn={team.cash || "$0"} addOnColor="green">
          Cash
        </Tag>
      </Grid.Col>
      <Grid.Col className="col-lg order-lg-first">
        <Nav
          tabbed
          className="border-0 flex-column flex-lg-row"
          itemsObjects={navBarItems}
        />
      </Grid.Col>
    </Grid.Row>
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
          accountDropdown: accountDropdownProps(userSession),
          navItems: (
            <Nav.Item type="div" className="d-none d-md-flex">
              <SeasonDropdownSelector
                userSession={userSession}
                activeSeason={activeSeason}
              />
            </Nav.Item>
          ),
        }}
        navProps={{
          children: <NavBar userSession={userSession} />,
        }}
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

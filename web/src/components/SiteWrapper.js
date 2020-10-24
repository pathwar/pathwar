import * as React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Site, Nav, Dropdown, Tag, Grid } from "tabler-react";
import { Link } from "gatsby";
import { css } from "@emotion/core";

import ValidateCouponForm from "../components/coupon/ValidateCouponForm";

import logo from "../images/new-pathwar-logo-dark-blue.svg";
import iconProfile from "../images/icon-profile.svg";
import iconMail from "../images/icon-mail.svg";
import iconNotifications from "../images/icon-notifications.svg";
import iconPwn from "../images/icon-pwn-small.svg";
import iconClose from "../images/icon-close-filled.svg";

const wrapper = css`
  font-family: "Barlow", sans-serif;
  font-weight: 500;
  width: 100%;
  z-index: 1;
  box-sizing: border-box;
  background-color: transparent;
  top: 0px;
  padding: 1rem 4rem;
  margin-top: 1rem;
  display: flex;
  align-items: center;

  @media (max-width: 700px) {
    height: 54px;
  }
  @media (min-width: 701px) and (min-height: 600px) {
    height: 72px;
  }
  .headerMenu {
    list-style: none;
    margin: 0;
    padding-left: 10px;
    display: flex;
    align-items: center;
  }

  .link {
    display: block;
    text-decoration: none;
    padding: 1rem;
    color: #919aa3;

    &:hover {
      opacity: 0.8;
    }
  }
  @media (max-width: 700px) {
    .link {
      font-size: 16px;
    }
  }
  @media (max-width: 360px) {
    .link {
      font-size: 13px;
    }
  }
`;

const dropdown = css`
  position: relative;
  display: inline-block;
  margin-left: auto;

  &:hover {
    .dropdown-content {
      display: block;
    }

    .dropbtn {
      border-bottom-left-radius: 0;
      border-bottom-right-radius: 0;
    }
  }

  .dropbtn {
    background-color: white;
    color: #072a44;
    box-shadow: 2px 2px 24px -12px rgba(0, 0, 0, 0.75);
    border-radius: 10px;
    font-weight: 500;
    padding: 16px;
    font-size: 21px;
    border: none;
    cursor: pointer;
    min-width: 160px;
  }

  .dropdown-content {
    display: none;
    position: absolute;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
    background-color: white;
    min-width: 160px;
    box-shadow: -2px 16px 24px -12px rgba(0, 0, 0, 0.75);
    z-index: 1;

    ul {
      list-style: none;
      display: flex;
      flex-direction: column;
      margin: 0;
      padding: 0 1rem;

      li {
        display: flex;
        flex-direction: row;

        a {
          padding: 12px 16px;
          text-decoration: none;
          display: block;
        }
      }
    }
  }
`;

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
  let items;

  const multipleItems =
    activeUserSession &&
    activeUserSession.seasons &&
    activeUserSession.seasons.length >= 2;

  const clicked = e => {
    e.preventDefault();
    alert(activeSeason.name);
  };

  if (multipleItems) {
    items =
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
  }

  return (
    <Dropdown
      triggerContent={(activeSeason && activeSeason.name) || "Loading.."}
      toggle={multipleItems}
      color="primary"
      icon="flag"
      items={items}
    />
  );
};

const NavBar = ({ userSession }) => {
  const { cash } = userSession;

  return (
    <Grid.Row className="align-items-center">
      <Grid.Col width={6} className="ml-auto text-right" ignoreCol={true}>
        <ValidateCouponForm />
        <Tag
          color="lime"
          addOn={(cash && `$${cash}`) || "$0"}
          addOnColor="green"
        >
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

const listItems = [
  { link: "/app/challenges", name: "Challenges" },
  { link: "/app/missions", name: "Missions" },
  { link: "/app/events", name: "Tournaments" },
  { link: "/app/community", name: "Community" },
  { link: "/blog", name: "Blog" },
];

class SiteWrapper extends React.Component {
  render() {
    const { userSession, activeSeason } = this.props;

    const { activeUserSession: { claims } = {} } = userSession;

    const username =
      claims && claims.preferred_username
        ? claims.preferred_username
        : "Log in";

    return (
      <>
        <header css={wrapper}>
          <img
            src={logo}
            className="img-responsive"
            style={{ width: "59px" }}
          />
          <ul className="headerMenu">
            {listItems.map(item => (
              <li key={listItems.name}>
                <Link className="link" to={item.link}>
                  {item.name}
                </Link>
              </li>
            ))}
          </ul>
          <div css={dropdown}>
            <button className="dropbtn">{`@${username}`}</button>
            <div className="dropdown-content">
              <ul>
                <li>
                  <img src={iconProfile} className="img-responsive" />
                  <a href="#" className="link">
                    Profile
                  </a>
                </li>
                <li>
                  <img src={iconMail} className="img-responsive" />
                  <a href="#" className="link">
                    Messages
                  </a>
                </li>
                <li>
                  <img src={iconPwn} className="img-responsive" />
                  <a href="#" className="link">
                    Wallet
                  </a>
                </li>
                <li>
                  <img src={iconNotifications} className="img-responsive" />
                  <a href="#" className="link">
                    Notifications
                  </a>
                </li>
                <li>
                  <img src={iconClose} className="img-responsive" />
                  <a href="#" className="link">
                    Disconnect
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </header>
        <body>{this.props.children}</body>
      </>
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

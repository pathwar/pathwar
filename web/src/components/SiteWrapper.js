import * as React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Link } from "gatsby";
import { css } from "@emotion/core";
import { FormattedMessage } from "react-intl";
import ValidateCouponForm from "../components/coupon/ValidateCouponForm";
import logo from "../images/new-pathwar-logo-light-blue.svg";
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
  padding: 1rem 4rem 0;
  margin-top: 1rem;
  display: flex;
  flex-wrap: wrap;
  align-items: center;

  .headerMenu {
    list-style: none;
    margin: 0;
    padding-left: 10px;
    display: flex;
    align-items: center;
  }

  .subHeader {
    flex-basis: 100%;
    display: flex;
    align-items: center;
    justify-content: flex-end;

    .cash {
      display: flex;
      margin-right: 1.5rem;

      .value {
        font-weight: bold;
        margin: 0;
        padding-left: 0.5rem;
      }
    }
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
    height: 54px;
  }
  @media (min-width: 701px) and (min-height: 600px) {
    height: 140px;
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

const listItems = [
  { link: "/app/challenges", name: <FormattedMessage id="nav.challenges" /> },
  { link: "/app/missions", name: <FormattedMessage id="nav.missions" /> },
  { link: "/app/events", name: <FormattedMessage id="nav.events" /> },
  { link: "/app/community", name: <FormattedMessage id="nav.community" /> },
  { link: "/blog", name: <FormattedMessage id="nav.blog" /> },
];

class SiteWrapper extends React.Component {
  render() {
    const { userSession } = this.props;

    const {
      cash,
      activeKeycloakSession,
      activeUserSession: { claims } = {},
    } = userSession;

    const username =
      claims && claims.preferred_username
        ? claims.preferred_username
        : "Log in";

    return (
      <>
        <header css={wrapper}>
          <img src={logo} className="img-responsive" />
          <ul className="headerMenu">
            {listItems.map(item => (
              <li key={item.link}>
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
                  <a
                    href={
                      activeKeycloakSession &&
                      activeKeycloakSession.tokenParsed.iss + "/account"
                    }
                    className="link"
                  >
                    <FormattedMessage id="userNav.profile" />
                  </a>
                </li>
                <li>
                  <img src={iconMail} className="img-responsive" />
                  <a href="#" className="link">
                    <FormattedMessage id="userNav.messages" />
                  </a>
                </li>
                <li>
                  <img src={iconPwn} className="img-responsive" />
                  <a href="#" className="link">
                    <FormattedMessage id="userNav.wallet" />
                  </a>
                </li>
                <li>
                  <img src={iconNotifications} className="img-responsive" />
                  <a href="#" className="link">
                    <FormattedMessage id="userNav.notifications" />
                  </a>
                </li>
                <li>
                  <img src={iconClose} className="img-responsive" />
                  <Link className="link" to="/app/logout">
                    <FormattedMessage id="userNav.disconnect" />
                  </Link>
                </li>
              </ul>
            </div>
          </div>
          <div className="subHeader">
            <div className="cash">
              <img src={iconPwn} className="img-responsive" />
              <p className="value">{(cash && `$${cash}`) || "$0"}</p>
            </div>
            <ValidateCouponForm />
          </div>
        </header>
        {this.props.children}
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

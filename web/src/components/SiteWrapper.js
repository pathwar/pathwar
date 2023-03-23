import * as React from "react";
import {useDispatch, useSelector} from "react-redux";
import { Link } from "gatsby";
import { css } from "@emotion/core";
import { FormattedMessage } from "react-intl";
import { useTheme } from "emotion-theming";
import ValidateCouponForm from "../components/coupon/ValidateCouponForm";
import logo from "../images/new-pathwar-logo-light-blue.svg";
import iconProfile from "../images/icon-profile.svg";
// import iconMail from "../images/icon-mail.svg";
// import iconNotifications from "../images/icon-notifications.svg";
import iconPwn from "../images/icon-pwn-small.svg";
import iconClose from "../images/icon-close-filled.svg";
import {Avatar, Table} from "tabler-react";
import moment from "moment/moment";
import {switchSeason} from "../actions/seasons";

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
  margin-left: 1rem;

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
    display: flex;
    align-items: center;
    justify-content: center;
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
    width: 100%;

    ul {
      list-style: none;
      display: flex;
      flex-direction: column;
      margin: 0;
      padding: 0 1rem;
      max-height: 200px;
      overflow-y: auto;

      li {
        display: flex;
        flex-direction: row;

        a {
          text-decoration: none;
          display: block;
        }

      }
    }
  }
`;

const linkSpan = css`
    color: #919aa3;
    cursor: pointer;

    &:hover {
      opacity: 0.8;
    }
`;

const langSwitcher = css`
  margin-left: 1.5rem;

  span {
    cursor: pointer;

    &:hover {
      font-weight: bold;
    }
  }
`;

const listItems = [
  {
    link: `/`,
    name: <FormattedMessage id="nav.home" />,
  },
  {
    link: `/challenges`,
    name: <FormattedMessage id="nav.challenges" />,
  },
  {
    link: `/statistics`,
    name: <FormattedMessage id="nav.statistics" />,
  },
  {
    link: `/organizations`,
    name: <FormattedMessage id="nav.organizations" />,
  },
  // { link: `${appPrefix}/events`, name: <FormattedMessage id="nav.events" /> },
  // {
  //   link: `${appPrefix}/community`,
  //   name: <FormattedMessage id="nav.community" />,
  // },
  // { link: `${appPrefix}/blog`, name: <FormattedMessage id="nav.blog" /> },
];

const SiteWrapper = ({ children }) => {
  const userSession = useSelector(state => state.userSession);
  const activeTeam = useSelector(state => state.seasons.activeTeam);
  const activeSeason = useSelector(state => state.seasons.activeSeason);
  const userSeasons = useSelector(state => state.seasons.userSeasons);
  const currentTheme = useTheme();

  const {
    cash,
    activeKeycloakSession,
    activeUserSession: { claims, user } = {},
  } = userSession;

  const username =
    claims && claims.preferred_username ? claims.preferred_username : "Log in";

  const switchLanguage = lang => {
    const browser = typeof window !== "undefined" && window;
    if (browser) {
      window.localStorage.setItem("pw.lang", lang);
      window.location.reload();
    }
  };

  return (
    <>
      <header css={wrapper}>
        <Link className="link" to={"/"}>
          <img src={logo} className="img-responsive" />
        </Link>
        <ul className="headerMenu">
          {listItems.map(item => (
            <li key={item.link}>
              <Link
                className="link"
                to={item.link}
                activeStyle={{
                  fontWeight: "bold",
                  color: currentTheme.colors.primary,
                }}
              >
                {item.name}
              </Link>
            </li>
          ))}
        </ul>
        <div css={{marginLeft: "auto"}}>
          <div css={dropdown}>
            <button className="dropbtn">
              <span className="mr-2">
                {activeTeam && activeTeam.organization && activeTeam.organization.gravatar_url ? (
                  <Avatar imageURL={`${activeTeam.organization.gravatar_url}?d=identicon`} />
                ) : (
                  <Avatar icon="users" />
                )}
              </span>
              <span>
                <FormattedMessage id="userNav.activeSeason" />
                {`: ${activeSeason && activeSeason.name ? activeSeason.name : 'Loading'}`}</span>
            </button>
            <div className="dropdown-content">
              <ul>
                {userSeasons && (
                  <UserSeasonsItems userSeasons={userSeasons} />
                )}
              </ul>
            </div>
          </div>
          <div css={dropdown}>
            <button className="dropbtn">
              <span className="mr-2">
                {user && user.gravatar_url ? (
                  <Avatar imageURL={`${user.gravatar_url}?d=identicon`} />
                ) : (
                  <Avatar icon="users" />
                )}
              </span>
              <span>{`${username}`}</span>
            </button>
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
                    target="_blank"
                    rel="noreferrer"
                  >
                    <FormattedMessage id="userNav.profile" />
                  </a>
                </li>
                {/* <li>
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
                </li> */}
                <li>
                  <img src={iconClose} className="img-responsive" />
                  <Link className="link" to="/logout">
                    <FormattedMessage id="userNav.disconnect" />
                  </Link>
                </li>
              </ul>
            </div>

          </div>
        </div>
        <div css={langSwitcher}>
          <span onClick={() => switchLanguage("en")}>EN</span> •{" "}
          <span onClick={() => switchLanguage("fr")}>FR</span>
        </div>
        <div className="subHeader">
          <div className="cash">
            <img src={iconPwn} className="img-responsive" />
            <p className="value">{(cash && `$${cash}`) || "$0"}</p>
          </div>
          <ValidateCouponForm />
        </div>
      </header>
      {children}
    </>
  );
};

const UserSeasonsItems = ({ userSeasons }) => {
  const dispatch = useDispatch();
  const setActiveSeasonDispatch = season => dispatch(switchSeason(season));
  const SwitchSeason = async season => {
    const id = season && season.slug ? season.slug : 0;
    await setActiveSeasonDispatch(id);
    window.location.reload();
  };

  return userSeasons.map((item, idx) => {
    return (
      <li css={{padding: "12px 16px"}}>
        <span className="mr-2">
          {item && item.team && item.team.organization && item.team.organization.gravatar_url ? (
            <Avatar size="sm" imageURL={`${item.team.organization.gravatar_url}?d=identicon`} />
          ) : (
            <Avatar size="sm" icon="users" />
          )}
        </span>
        <span css={linkSpan} onClick={() => SwitchSeason(item.season)}>{`${item && item.season && item.season.name ? item.season.name : 'Loading'}`}</span>
      </li>
    );
  });
}

export default React.memo(SiteWrapper);

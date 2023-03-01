import React from "react";
import { Router, Location } from "@reach/router";
import { Helmet } from "react-helmet";
import { withPrefix } from "gatsby";
import loadable from "@loadable/component";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Global } from "@emotion/core";
import { ThemeProvider, useTheme } from "emotion-theming";
import siteMetaData from "../constants/metadata";
import { lightTheme } from "../styles/themes";
import { globalStyle } from "../styles/globalStyle";
import HomePage from "./HomePage";
import LogoutPage from "./LogoutPage";
import ChallengesPage from "./ChallengesPage";
import ChallengeDetailsPage from "./ChallengeDetailsPage";
import StatisticsPage from "./StatisticsPage";
import SiteWrapper from "../components/SiteWrapper";
import TeamDetailsPage from "./TeamDetailsPage";
import SettingsPage from "./SettingsPage";
import NotAvailablePage from "./NotAvailablePage";
import * as Sentry from "@sentry/browser";

import logo from "../images/new_pathwar-logo.svg";
const ProtectedRoute = loadable(() => import("../components/ProtectedRoute"));

//Third part libs global styles
import "tabler-react/dist/Tabler.css";
import "react-responsive-modal/styles.css";
import OrganizationsPage from "./organization/UserOrganizationsPage";
import OrganizationDetailsPage from "./organization/OrganizationDetailsPage";
import OrganizationMembersPage from "./organization/OrganizationMembersPage";
import OrganizationTeamsAndSeasonsPage from "./organization/OrganizationTeamsAndSeasons";

Sentry.init({
  dsn:
    "https://8605d8e8fa21419d9a0e3f36a54df5cb@o406102.ingest.sentry.io/5272916",
});
toast.configure();

export const App = () => {
  const currentTheme = useTheme();
  const { title, description } = siteMetaData;

  return (
    <>
      <Helmet>
        <title>{title}</title>
        <meta name="description" content={description} />

        <meta
          name="go-import"
          content="pathwar.land git https://github.com/pathwar/pathwar"
        />
        <meta
          name="go-source"
          content="pathwar.land https://github.com/pathwar/pathwar https://github.com/pathwar/pathwar/tree/master{/dir} https://github.com/pathwar/pathwar/tree/master{/dir}/{file}#L{line}"
        />

        <meta property="og:description" content={description} />
        <meta property="og:url" content={description} />
        <meta property="og:site_name" content={title} />
        <meta property="og:type" content="website" />
        <meta property="og:image" content={logo} />
        <meta property="og:image:width" content={200} />
        <meta property="og:image:height" content={200} />

        <meta property="twitter:card" content="summary" />
        <meta property="twitter:title" content={title} />
        <meta property="twitter:description" content={description} />
        <meta property="twitter:image" content={logo} />
        <link
          rel="stylesheet"
          href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Nunito:wght@400;500;700;900&display=swap"
          rel="stylesheet"
        ></link>
        <link
          href="https://fonts.googleapis.com/css2?family=Bungee&display=swap"
          rel="stylesheet"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Barlow:wght@300;500;700&display=swap"
          rel="stylesheet"
        ></link>
        <script async src={withPrefix("chat-init.js")} type="text/javascript" />
      </Helmet>

      <Global styles={globalStyle(currentTheme, true)} />
      <SiteWrapper />
      <Location>
        {({ location }) => (
          <Router location={location}>
            <ProtectedRoute path={`/challenges`} component={ChallengesPage} />
            <ProtectedRoute path={`/statistics`} component={StatisticsPage} />
            <ProtectedRoute path={`/organizations`} component={OrganizationsPage} />
            <ProtectedRoute
              path={`/team/:teamId`}
              component={TeamDetailsPage}
            />
            <ProtectedRoute
              path={`/organization/:organizationId`}
              component={OrganizationDetailsPage}
            />
            <ProtectedRoute
              path={`/organization/:organizationId/members`}
              component={OrganizationMembersPage}
            />
            <ProtectedRoute
              path={`/organization/:organizationId/teams`}
              component={OrganizationTeamsAndSeasonsPage}
            />
            <ProtectedRoute
              path={`/team/:teamId`}
              component={TeamDetailsPage}
            />
            <ProtectedRoute
              path={`/challenges/:challengeId`}
              component={ChallengeDetailsPage}
            />
            <ProtectedRoute path={`/settings`} component={SettingsPage} />
            <ProtectedRoute path={`/logout`} component={LogoutPage} />
            <ProtectedRoute path="/" component={HomePage} />
            <ProtectedRoute default component={NotAvailablePage} />
          </Router>
        )}
      </Location>
    </>
  );
};

const ThemedApp = () => (
  <ThemeProvider theme={lightTheme}>
    <App />
  </ThemeProvider>
);

export default ThemedApp;

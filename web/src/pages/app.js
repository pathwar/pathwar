import React from "react";
import { Router, Location } from "@reach/router";
import { Helmet } from "react-helmet";
import loadable from "@loadable/component";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Global } from "@emotion/core";
import { ThemeProvider, useTheme } from "emotion-theming";
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
import * as Sentry from "@sentry/browser";
const ProtectedRoute = loadable(() => import("../components/ProtectedRoute"));

import "tabler-react/dist/Tabler.css";

Sentry.init({
  dsn:
    "https://8605d8e8fa21419d9a0e3f36a54df5cb@o406102.ingest.sentry.io/5272916",
});
toast.configure();

const App = () => {
  const currentTheme = useTheme();
  return (
    <>
      <Helmet>
        <link
          rel="stylesheet"
          href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Nunito:wght@400;700&display=swap"
          rel="stylesheet"
        ></link>
        <link
          href="https://fonts.googleapis.com/css2?family=Bungee&display=swap"
          rel="stylesheet"
        />
      </Helmet>
      <Global styles={globalStyle(currentTheme)} />
      <SiteWrapper />
      <Location>
        {({ location }) => (
          <Router location={location}>
            <ProtectedRoute path="/app/home" component={HomePage} />
            <ProtectedRoute path="/app/challenges" component={ChallengesPage} />
            <ProtectedRoute path="/app/statistics" component={StatisticsPage} />
            <ProtectedRoute
              path="/app/team/:teamId"
              component={TeamDetailsPage}
            />
            <ProtectedRoute
              path="/app/challenge/:challengeId"
              component={ChallengeDetailsPage}
            />
            <ProtectedRoute path="/app/settings" component={SettingsPage} />
            <ProtectedRoute path="/app/logout" component={LogoutPage} />
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

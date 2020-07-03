import React from "react";
import { Router, Location } from "@reach/router";
import { Helmet } from "react-helmet";
import loadable from "@loadable/component";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import HomePage from "./HomePage";
import LogoutPage from "./LogoutPage";
import ChallengesPage from "./ChallengesPage";
import AllSeasonsPage from "./AllSeasonsPage";
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

const App = () => (
  <div>
    <Helmet>
      <link
        rel="stylesheet"
        href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"
      />
    </Helmet>
    <SiteWrapper />
    <Location>
      {({ location }) => (
        <Router location={location}>
          <ProtectedRoute path="/app/home" component={HomePage} />
          <ProtectedRoute path="/app/challenges" component={ChallengesPage} />
          <ProtectedRoute path="/app/statistics" component={StatisticsPage} />
          <ProtectedRoute path="/app/all-seasons" component={AllSeasonsPage} />
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
  </div>
);

export default App;

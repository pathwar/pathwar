/* eslint-disable import/first */
import React from "react"
import { Router, Location } from "@reach/router"
import { Helmet } from "react-helmet"
import loadable from "@loadable/component"
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import DashboardPage from "./DashboardPage"
import LogoutPage from "./LogoutPage"
import SeasonPage from "./SeasonPage"
import AllSeasonsPage from "./AllSeasonsPage"
import ChallengeDetailsPage from "./ChallengeDetailsPage"
import SiteWrapper from "../components/SiteWrapper"
import TeamDetailsPage from "./TeamDetailsPage"
import SettingsPage from "./SettingsPage"
const ProtectedRoute = loadable(() => import("../components/ProtectedRoute"))

import "tabler-react/dist/Tabler.css"

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
          <ProtectedRoute path="/app/dashboard" component={DashboardPage} />
          <ProtectedRoute path="/app/season" component={SeasonPage} />
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
          <LogoutPage path="/app/logout" component={LogoutPage} />
          <ProtectedRoute path="/app/" redirect to="/app/season" component={SeasonPage} />
        </Router>
      )}
    </Location>
  </div>
)

export default App

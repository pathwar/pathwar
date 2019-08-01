import React from "react"
import { Router, Location } from "@reach/router"
import ProtectedRoute from "../components/ProtectedRoute"
import DashboardPage from "./DashboardPage"
import LogoutPage from "./LogoutPage";
import TournamentPage from "./TournamentPage";
import AllTournamentsPage from "./AllTournamentsPage";
import SiteWrapper from "../components/SiteWrapper";

import "tabler-react/dist/Tabler.css";


const App = () => (
    <div>
      <SiteWrapper />
      <Location>
        {({ location }) => (
          <Router location={location}>
            <ProtectedRoute path="/app/dashboard" component={DashboardPage} />
            <ProtectedRoute path="/app/tournament" component={TournamentPage} />
            <ProtectedRoute path="/app/all-tournaments" component={AllTournamentsPage} />
            <LogoutPage path="/app/logout" component={LogoutPage} />
          </Router>
        )}
      </Location>
    </div>
)

export default App
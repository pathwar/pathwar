import React from "react"
import { Router } from "@reach/router"
import ProtectedRoute from "../components/ProtectedRoute"
import DashboardPage from "./DashboardPage"
import LoginPage from "./LoginPage"
import LogoutPage from "./LogoutPage";
import TournamentPage from "./TournamentPage";
import AllTournamentsPage from "./AllTournamentsPage";

import "tabler-react/dist/Tabler.css";

const App = () => (
    <Router>
      <ProtectedRoute path="/app/dashboard" component={DashboardPage} />
      <LoginPage path="/app/login" />
      <ProtectedRoute path="/app/tournament" component={TournamentPage} />
      <ProtectedRoute path="/app/all-tournaments" component={AllTournamentsPage} />
      <LogoutPage path="/app/logout" component={LogoutPage} />
    </Router>
)

export default App
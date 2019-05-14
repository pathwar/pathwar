/* eslint-disable import/no-named-as-default */
import { Route, Switch } from "react-router-dom";
import PropTypes from "prop-types";
import React from "react";
import { hot } from "react-hot-loader";
import DashboardPage from "./pages/DashboardPage";
import LoginPage from "./pages/LoginPage";
import NotFoundPage from "./pages/NotFoundPage";
import CompetitionPage from "./pages/CompetitionPage";

class App extends React.Component {
  render() {
    return (
      <div>
        <Switch>
          <Route exact path="/" component={LoginPage} />
          <Route exact path="/login" component={LoginPage} />
          <Route exact path="/dashboard" component={DashboardPage} />
          <Route exact path="/competition" component={CompetitionPage} />
          <Route component={NotFoundPage} />
        </Switch>
      </div>
    );
  }
}

App.propTypes = {
  children: PropTypes.element
};

export default hot(module)(App);

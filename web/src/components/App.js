/* eslint-disable import/no-named-as-default */
import { Route, Switch } from "react-router-dom";
import PropTypes from "prop-types";
import React from "react";
import { hot } from "react-hot-loader";
import HomePage from "./pages/HomePage";
import NotFoundPage from "./pages/NotFoundPage";
import CompetitionsPage from "./pages/CompetitionsPage";

class App extends React.Component {
  render() {
    return (
      <div>
        <Switch>
          <Route exact path="/" component={HomePage} />
          <Route exact path="/competitions" component={CompetitionsPage} />
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

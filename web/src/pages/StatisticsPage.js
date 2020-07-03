/* eslint-disable react/prop-types */
import React from "react";
import { useSelector } from "react-redux";
import { Link } from "gatsby";
import { Button, Dimmer, Card, Grid } from "tabler-react";
import styles from "../../styles/layout/loader.module.css";

const StatisticsPage = props => {
  const activeUserSession = useSelector(
    state => state.userSession.activeUserSession
  );
  const activeTeam = useSelector(state => state.seasons.activeTeam);

  return <h1>StatisticsPage</h1>;
};

export default StatisticsPage;

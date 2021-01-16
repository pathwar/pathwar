import * as React from "react";
// import { Link } from "@reach/router";
import { Card, Table, Dimmer } from "tabler-react";
import PropTypes from "prop-types";
import moment from "moment";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";

const TeamsRows = ({ teams }) => {
  return teams.map(item => {
    return (
      <Table.Row key={item.organization.id}>
        <Table.Col>{item.organization.name}</Table.Col>
        <Table.Col alignContent="center">
          {(item.cash && `$${item.cash}`) || "$0"}
        </Table.Col>
        <Table.Col alignContent="center">
          {moment(item.created_at).calendar()}
        </Table.Col>
        {/* <Table.Col colSpan={1} alignContent="center">
          {item.score}
        </Table.Col>
        <Table.Col colSpan={1} alignContent="center">
          {item.gold_medals || 0}
        </Table.Col>
        <Table.Col colSpan={1} alignContent="center">
          {item.silver_medals || 0}
        </Table.Col>
        <Table.Col colSpan={1} alignContent="center">
          {item.bronze_medals || 0}
        </Table.Col>
        <Table.Col colSpan={1} alignContent="center">
          {item.nb_achievements || 0}
        </Table.Col> */}
      </Table.Row>
    );
  });
};

const AllTeamsOnSeasonList = ({ activeSeason, allTeamsOnSeason }) => {
  return !activeSeason || !allTeamsOnSeason ? (
    <Dimmer active loader />
  ) : (
    <Card>
      <Table
        striped={true}
        responsive={true}
        verticalAlign="center"
        className="mb-0"
      >
        <Table.Header>
          <Table.Row>
            <Table.ColHeader>
              <FormattedMessage id="AllTeamsOnSeasonList.team" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnSeasonList.cash" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnSeasonList.joined" />
            </Table.ColHeader>
            {/* <Table.ColHeader colSpan={1} alignContent="center">
              Score
            </Table.ColHeader>
            <Table.ColHeader colSpan={1}>
              <Icon name="award" className={styles.goldMedal} />
            </Table.ColHeader>
            <Table.ColHeader colSpan={1} alignContent="center">
              <Icon name="award" className={styles.sivlerMedal} />
            </Table.ColHeader>
            <Table.ColHeader colSpan={1} alignContent="center">
              <Icon name="award" className={styles.bronzeMedal} />
            </Table.ColHeader>
            <Table.ColHeader colSpan={1} alignContent="center">
              Achievements
            </Table.ColHeader> */}
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {activeSeason && allTeamsOnSeason && (
            <TeamsRows teams={allTeamsOnSeason} />
          )}
        </Table.Body>
      </Table>
    </Card>
  );
};

AllTeamsOnSeasonList.propTypes = {
  seasons: PropTypes.object,
};

export default AllTeamsOnSeasonList;

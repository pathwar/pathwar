import * as React from "react";
import { Link } from "@reach/router";
import { Card, Table, Dimmer, Icon } from "tabler-react";
import PropTypes from "prop-types";

import styles from "./style.module.css";

const TeamsRows = ({ teams }) => {
  return teams.map(item => {
    return (
      <Table.Row key={item.organization.id}>
        <Table.Col>{item.organization.name}</Table.Col>
        <Table.Col alignContent="center">{item.cash || "$0"}</Table.Col>
        <Table.Col alignContent="center">{item.created_at}</Table.Col>
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
            <Table.ColHeader>Team</Table.ColHeader>
            <Table.ColHeader alignContent="center">Cash</Table.ColHeader>
            <Table.ColHeader alignContent="center">Joined at</Table.ColHeader>
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

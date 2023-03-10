import React, { useMemo } from "react";
// import { Link } from "@reach/router";
import { Card, Table, Dimmer, Avatar } from "tabler-react";
import PropTypes from "prop-types";
import { css } from "@emotion/core";
import * as R from "ramda";

// import styles from "./style.module.css";
import { FormattedMessage } from "react-intl";

const TeamsRows = ({ teams }) => {
  const scoreCashSort = R.sortWith([
    R.descend(R.prop("score")),
    R.descend(R.prop("cash")),
  ]);

  const parsedTeams = teams.map(item => ({
    ...item,
    score: item.score ? parseInt(item.score, 10) : 0,
    cash: item.cash ? parseInt(item.cash, 10) : 0,
  }));

  const sortedTeamsByScoreAndCash = useMemo(() => scoreCashSort(parsedTeams), [
    scoreCashSort,
    parsedTeams,
  ]);

  return sortedTeamsByScoreAndCash.map((item, idx) => {
    if (item.score && item.cash) {
      return (
        <Table.Row key={item.organization.id}>
          <Table.Col alignContent="center">{idx + 1}</Table.Col>

          <Table.Col
            css={css`
              display: flex;
              align-items: center;
            `}
          >
            <Avatar
              className="mr-2"
              imageURL={`${item.organization.gravatar_url}?d=identicon`}
            />
            <span>{item.organization.name}</span>
          </Table.Col>
          <Table.Col alignContent="center">{item.score}</Table.Col>
          {/* <Table.Col colSpan={1} alignContent="center">
          {item.nb_achievements}
        </Table.Col> */}
          <Table.Col alignContent="center">
            {(item.cash && `$${item.cash}`) || "$0"}
          </Table.Col>

          {/*<Table.Col colSpan={1} alignContent="center">
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
    }
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
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnSeasonList.rank" />
            </Table.ColHeader>
            <Table.ColHeader>
              <FormattedMessage id="AllTeamsOnSeasonList.team" />
            </Table.ColHeader>
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnSeasonList.score" />
            </Table.ColHeader>
            {/* <Table.ColHeader colSpan={1} alignContent="center">
              <FormattedMessage id="AllTeamsOnSeasonList.achievements" />
            </Table.ColHeader> */}
            <Table.ColHeader alignContent="center">
              <FormattedMessage id="AllTeamsOnSeasonList.cash" />
            </Table.ColHeader>
            {/* <Table.ColHeader colSpan={1}>
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

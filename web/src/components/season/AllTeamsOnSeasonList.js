import * as React from "react";
import { Link } from "@reach/router";
import { Card, Table, Dimmer, Icon } from "tabler-react";
import PropTypes from "prop-types";

import loaderStyles from "../../styles/layout/loader.module.css"
import styles from "./style.module.css"

const TeamsRows = ({ teams }) => {
    return teams.map((item) => {
        return (
            <Table.Row key={item.organization.id}>
                <Table.Col colSpan={2}><Link to={`/app/team/${item.id}`}>{item.organization.name}</Link></Table.Col>
                <Table.Col colSpan={1}>{item.score}</Table.Col>
                <Table.Col colSpan={1}>{item.gold_medals || 0}</Table.Col>
                <Table.Col colSpan={1}>{item.silver_medals || 0}</Table.Col>
                <Table.Col colSpan={1}>{item.bronze_medals || 0}</Table.Col>
            </Table.Row>
        )
    })
}

const AllTeamsOnSeasonList = ({ activeSeason, allTeamsOnSeason }) => {
    return (
        !activeSeason || !allTeamsOnSeason ? <Dimmer className={loaderStyles.dimmer} active loader /> :
        <Card>
            <Table
                cards={true}
                striped={true}
                responsive={true}
                className="table-vcenter"
            >
                <Table.Header>
                    <Table.Row>
                        <Table.ColHeader colSpan={2}>Name</Table.ColHeader>
                        <Table.ColHeader colSpan={1}>Score</Table.ColHeader>
                        <Table.ColHeader colSpan={1}><Icon name="award" className={styles.goldMedal} /></Table.ColHeader>
                        <Table.ColHeader colSpan={1}><Icon name="award" className={styles.sivlerMedal} /></Table.ColHeader>
                        <Table.ColHeader colSpan={1}><Icon name="award" className={styles.bronzeMedal} /></Table.ColHeader>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    { activeSeason && allTeamsOnSeason &&
                        <TeamsRows
                            teams={allTeamsOnSeason}
                        />
                    }
                </Table.Body>
            </Table>
        </Card>
    )
}

AllTeamsOnSeasonList.propTypes = {
    seasons: PropTypes.object,
};

export default AllTeamsOnSeasonList;

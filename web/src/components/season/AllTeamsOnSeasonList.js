import * as React from "react";
import { Link } from "@reach/router";
import { Card, Table, Dimmer } from "tabler-react";
import PropTypes from "prop-types";

import styles from "../../styles/layout/loader.module.css"

const TeamsRows = ({ teams }) => {
  debugger
    return teams.map((item) => {
        return (
            <Table.Row key={item.organization.id}>
                <Table.Col colSpan={2}><Link to={`/app/team/${item.id}`}>{item.organization.name}</Link></Table.Col>
            </Table.Row>
        )
    })
}

const AllTeamsOnSeasonList = ({ activeSeason, allTeamsOnSeason }) => {
    return (
        !activeSeason || !allTeamsOnSeason ? <Dimmer className={styles.dimmer} active loader /> :
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

import * as React from "react";
import { Card, Table, Button } from "tabler-react";
import PropTypes from "prop-types";

const TeamsRows = ({allTournaments}) => {
    return allTournaments.map((tournament) => {
        return (
            <Table.Row key={tournament.metadata.id}>
                <Table.Col colSpan={2}>{tournament.name}</Table.Col>
                <Table.Col>{tournament.status}</Table.Col>
                <Table.Col>{tournament.visibility}</Table.Col>
                <Table.Col>
                    <Button color="info" size="sm">Set active</Button>
                </Table.Col>
            </Table.Row>
        )
    })
}

const AllTournamentsList = props => {
    const { tournaments } = props;
    return (
        <Card>
            <Card.Header>
                <Card.Title>All Tournaments</Card.Title>
            </Card.Header>
            <Table
                cards={true}
                striped={true}
                responsive={true}
                className="table-vcenter"
            >
                <Table.Header>
                    <Table.Row>
                        <Table.ColHeader colSpan={2}>Name</Table.ColHeader>
                        <Table.ColHeader>Status</Table.ColHeader>
                        <Table.ColHeader>Visibility</Table.ColHeader>
                        <Table.ColHeader></Table.ColHeader>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {tournaments && <TeamsRows allTournaments={tournaments}  />}
                </Table.Body>
            </Table>
        </Card>
    )
}

AllTournamentsList.propTypes = {
    tournaments: PropTypes.array,
};

export default AllTournamentsList
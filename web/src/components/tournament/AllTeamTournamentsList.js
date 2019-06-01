import * as React from "react";
import { Card, Table, Button } from "tabler-react";
import PropTypes from "prop-types";

const TournamentsRows = ({ teamTournaments, setActive, activeTeam, activeTournament }) => {
    return teamTournaments.map((tournament) => {
        const isActive = activeTournament ? tournament.metadata.id === activeTournament.metadata.id : false;
        return (
            <Table.Row key={tournament.metadata.id}>
                <Table.Col colSpan={2}>{tournament.name}</Table.Col>
                <Table.Col>{tournament.status}</Table.Col>
                <Table.Col>{tournament.visibility}</Table.Col>
                {isActive && <Table.Col>
                    Active
                </Table.Col>}
                {!isActive && <Table.Col>
                    <Button color="info" size="sm" onClick={() => setActive(activeTeam.metadata.id, tournament)}>Set active</Button>
                </Table.Col>}
            </Table.Row>
        )
    })
}

const AllTeamTournamentsList = props => {
    const { teamTournaments, setActive, activeTeam, activeTournament } = props;
    return (
        <Card>
            <Card.Header>
                <Card.Title>Team Tournaments</Card.Title>
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
                    { teamTournaments && 
                        <TournamentsRows 
                            teamTournaments={teamTournaments}
                            activeTournament={activeTournament}
                            setActive={setActive} 
                            activeTeam={activeTeam} 
                        />
                    }
                </Table.Body>
            </Table>
        </Card>
    )
}

AllTeamTournamentsList.propTypes = {
    teamTournaments: PropTypes.array,
    setActive: PropTypes.func.metadata,
    activeTeam: PropTypes.object,
    activeTournament: PropTypes.object
};

export default AllTeamTournamentsList
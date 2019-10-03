import * as React from "react";
import { connect } from "react-redux"
import { Card, Table, Button } from "tabler-react";
import PropTypes from "prop-types";
import {
    setActiveTournament as setActiveTournamentAction,
    fetchTeamTournaments as fetchTeamTournamentsAction
} from "../../actions/tournaments"


const TournamentsRows = ({ teamTournaments, setActive, activeTournament }) => {
    return teamTournaments.map((tournament) => {
        const isActive = activeTournament ? tournament.id === activeTournament.id : false;
        return (
            <Table.Row key={tournament.id}>
                <Table.Col colSpan={2}>{tournament.name}</Table.Col>
                <Table.Col>{tournament.status}</Table.Col>
                <Table.Col>{tournament.visibility}</Table.Col>
                {isActive && <Table.Col>
                    Active
                </Table.Col>}
                {!isActive && <Table.Col>
                    <Button color="info" size="sm" onClick={() => setActive(tournament)}>Set active</Button>
                </Table.Col>}
            </Table.Row>
        )
    })
}

class AllTeamTournamentsList extends React.PureComponent {

    componentDidMount() {
        const { fetchTeamTournamentsAction, activeTeam } = this.props;
        fetchTeamTournamentsAction(activeTeam.id)
    }

    render() {
        const { setActiveTournamentAction, tournaments: { activeTournament, allTeamTournaments } } = this.props;
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
                        { allTeamTournaments && activeTournament &&
                            <TournamentsRows
                                teamTournaments={allTeamTournaments}
                                activeTournament={activeTournament}
                                setActive={setActiveTournamentAction}
                            />
                        }
                    </Table.Body>
                </Table>
            </Card>
        )
    }
}

AllTeamTournamentsList.propTypes = {
    activeTeam: PropTypes.object,
    fetchTeamTournamentsAction: PropTypes.func,
    tournaments: PropTypes.object,
    setActiveTournamentAction: PropTypes.func
};

const mapStateToProps = state => ({
    tournaments: state.tournaments,
    activeTeam: state.teams.activeTeam,
});

const mapDispatchToProps = {
    fetchTeamTournamentsAction: (teamID) => fetchTeamTournamentsAction(teamID),
    setActiveTournamentAction: (tournamentData) => setActiveTournamentAction(tournamentData)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AllTeamTournamentsList);

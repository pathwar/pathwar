import * as React from "react";
import { connect } from "react-redux"
import { Card, Table, Avatar, Button } from "tabler-react";
import PropTypes from "prop-types";
import {
    fetchTeamsList as fetchTeamsListAction,
    joinTeam as joinTeamAction
} from "../../actions/teams"

const TeamsRows = ({allTeamsList, userTeamsList, joinTeam, userID}) => {
    return allTeamsList.map((team) => {

        const isUserOnTem = userTeamsList ? userTeamsList.find((userTeam) => team.metadata.id === userTeam.metadata.id) : undefined;

        return (
            <Table.Row key={team.metadata.id}>
            <Table.Col className="w-1">
                <Avatar imageURL={team.gravatar_url} />
            </Table.Col>
            <Table.Col>{team.name}</Table.Col>
            <Table.Col>{team.locale}</Table.Col>
            {isUserOnTem && <Table.Col>
                Joined
            </Table.Col>}
            {!isUserOnTem && <Table.Col>
                <Button color="info" size="sm" onClick={() => joinTeam(userID, team.metadata.id)}>Join</Button>
            </Table.Col>}
        </Table.Row>
        )
    })

}

class AllTeamsCard extends React.PureComponent {

    componentDidMount() {
        const { fetchTeamsListAction } = this.props;
        fetchTeamsListAction();
    }

    render() {
        const { teams, joinTeamAction, activeSession } = this.props;
        return (
            <Card>
                  <Card.Header>
                    <Card.Title>All Teams</Card.Title>
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
                        <Table.ColHeader>Locale</Table.ColHeader>
                        <Table.ColHeader></Table.ColHeader>
                      </Table.Row>
                    </Table.Header>
                    <Table.Body>
                        {teams.allTeamsList && <TeamsRows
                            allTeamsList={teams.allTeamsList}
                            userTeamsList={teams.userTeamsList}
                            joinTeam={joinTeamAction}
                            userID={activeSession.sessionId}
                        />}
                    </Table.Body>
                  </Table>
                </Card>
        )
    }
}

AllTeamsCard.propTypes = {
    teams: PropTypes.object,
    activeSession: PropTypes.object,
    fetchTeamsListAction: PropTypes.func,
    joinTeamAction: PropTypes.func
};

const mapStateToProps = state => ({
    teams: state.teams,
    activeSession: state.userSession.activeSession
});

const mapDispatchToProps = {
    fetchTeamsListAction: () => fetchTeamsListAction(),
    joinTeamAction: (userID, teamID) => joinTeamAction(userID, teamID)
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AllTeamsCard);

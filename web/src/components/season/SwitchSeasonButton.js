import { FormattedMessage } from "react-intl";
import * as React from "react";
import { Card, Table } from "tabler-react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { Button } from "tabler-react";
import {setSwitchSeason} from "../../actions/seasons";

const SeasonsRows = ({ allSeasons }) => {

  const [isFetching, setFetching] = React.useState(false);

  const dispatch = useDispatch();
  const setActiveSeasonDispatch = seasonID => dispatch(setSwitchSeason(seasonID));

  const SwitchSeason = async seasonID => {
    setFetching(true);
    setActiveSeasonDispatch(seasonID).then(response => {
      setFetching(false);
      return response;
    });
  };

  return allSeasons.map(season => {

    return (
      <Table.Row key={season.id}>
        <Table.Col colSpan={2}>
          <Button.List>
            <Button
              onClick={() => SwitchSeason(season.id)}
              loading={isFetching}
              color="primary"
            >{season.name}
            </Button>
          </Button.List>
        </Table.Col>
        <Table.Col>{season.status}</Table.Col>
        <Table.Col>{season.visibility}</Table.Col>
      </Table.Row>
    );
  });
};

const SwitchSeasonButton = props => {
  const { seasons } = props;
  return (
    <Card>
      <Card.Header>
        <Card.Title>
          <FormattedMessage id="AllSeasonsList.title" />
        </Card.Title>
      </Card.Header>
      <Table
        cards={true}
        striped={true}
        responsive={true}
        className="table-vcenter"
      >
        <Table.Header>
          <Table.Row>
            <Table.ColHeader colSpan={2}>
              <FormattedMessage id="AllSeasonsList.name" />
            </Table.ColHeader>
            <Table.ColHeader>
              <FormattedMessage id="AllSeasonsList.status" />
            </Table.ColHeader>
            <Table.ColHeader>
              <FormattedMessage id="AllSeasonsList.visibility" />
            </Table.ColHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {seasons && <SeasonsRows allSeasons={seasons} />}
        </Table.Body>
      </Table>
    </Card>
  );
};

SwitchSeasonButton.propTypes = {
  seasons: PropTypes.array,
};

export default SwitchSeasonButton;

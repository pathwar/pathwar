import { FormattedMessage } from "react-intl";
import * as React from "react";
import { Card } from "tabler-react";;
import { useDispatch } from "react-redux";
import { Button } from "tabler-react";
import {setSwitchSeason} from "../../actions/seasons";

const SwitchSeasonInput = () => {

  const [season, setSeason] = React.useState('');

  const dispatch = useDispatch();
  const setActiveSeasonDispatch = seasonID => dispatch(setSwitchSeason(seasonID));

  const SwitchSeason = async seasonID => {
    setActiveSeasonDispatch(seasonID).then(response => {
      return response;
    });
  };

  return (
    <Card>
      <Card.Header>
        <Card.Title>
          <FormattedMessage id="SwitchSeasonInput.title"/>
        </Card.Title>
      </Card.Header>
        <input value={season} onChange={e => setSeason(e.target.value)}/>
        <Button.List>
          <Button onClick={() => SwitchSeason(season)} color="primary">
            <FormattedMessage id="SwitchSeasonInput.action"/>
          </Button>
        </Button.List>
    </Card>
  );
};

export default SwitchSeasonInput;

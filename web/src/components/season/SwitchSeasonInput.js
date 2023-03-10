import { FormattedMessage } from "react-intl";
import * as React from "react";
import { Card } from "tabler-react";
import {useDispatch} from "react-redux";
import {Autocomplete} from "@material-ui/lab";
import {TextField} from "@material-ui/core";
import {switchSeason} from "../../actions/seasons";
import {Button} from "tabler-react";

const SwitchSeasonInput = ({seasons}) => {
  const dispatch = useDispatch();
  const setActiveSeasonDispatch = season => dispatch(switchSeason(season));
  const [season, setSeason] = React.useState('');

  const SwitchSeason = async seasonSlug => {
    await setActiveSeasonDispatch(seasonSlug);
    window.location.reload();
  }

  return (
    <Card>
      <Card.Header css={{ display: "flex", justifyContent: "space-around" }}>
        <Card.Title>
          <FormattedMessage id="SwitchSeasonInput.title"/>
        </Card.Title>
      </Card.Header>
      <Autocomplete
        freeSolo
        autoComplete
        autoHighlight
        onInputChange={(event, newInputValue) => {setSeason(newInputValue)}}
        options={seasons ? seasons.filter(season => season.team).map(season => season.season.slug) : []}
        renderInput={(params) => (
          <TextField {...params}
                     variant="outlined"
                     label="Search Box"
          />
        )}
      />
      <Button.List css={{ display: "flex", justifyContent: "space-around" }}>
        <Button onClick={() => SwitchSeason(season)} color="primary">
          <FormattedMessage id="SwitchSeasonInput.action"/>
        </Button>
      </Button.List>
    </Card>
  );
};

export default SwitchSeasonInput;

import { FormattedMessage } from "react-intl";
import * as React from "react";
import {useDispatch} from "react-redux";
import {switchSeason} from "../../actions/seasons";
import {Button} from "tabler-react";

const SwitchSeasonButton = ({seasonID}) => {
  const dispatch = useDispatch();
  const setActiveSeasonDispatch = season => dispatch(switchSeason(season));

  const SwitchSeason = async () => {
    await setActiveSeasonDispatch(seasonID);
    window.location.reload();
  }

  return (
        <Button onClick={SwitchSeason} color="warning">
          <FormattedMessage id="SwitchSeasonButton.title"/>
        </Button>
  );
};

export default SwitchSeasonButton;
